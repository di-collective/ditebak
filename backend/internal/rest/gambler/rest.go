package gambler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	firebase "firebase.google.com/go"
	log "github.com/sirupsen/logrus"

	"github.com/di-collective/ditebak/backend/internal/usecase/gambler"
	"github.com/di-collective/ditebak/backend/internal/usecase/gambler/command"
	"github.com/di-collective/ditebak/backend/internal/usecase/gambler/dto"
	"github.com/di-collective/ditebak/backend/pkg/exception"
	"github.com/di-collective/ditebak/backend/pkg/global"
	"github.com/di-collective/ditebak/backend/pkg/rest"

	"github.com/julienschmidt/httprouter"
)

// rest API for gambler
type restapi struct {
	conf *gambler.Config
	ggw  *gambler.Gateway
}

// New gambler micro API gateway
func New() rest.REST {
	ctx := context.Background()
	fap, _ := firebase.NewApp(ctx, nil)
	fac, _ := fap.Auth(ctx)

	conf := &gambler.Config{
		Const: &gambler.ConfigConst{
			MaxStake:   10,
			AuthClient: fac,
		},
		URL: &gambler.ConfigURL{
			User:  defaultOnEmptyEnv(os.Getenv("URL_USER"), "http://localhost:8080/users"),
			Topic: defaultOnEmptyEnv(os.Getenv("URL_TOPIC"), "http://localhost:8080/topics"),
			Bet:   defaultOnEmptyEnv(os.Getenv("URL_BET"), "http://localhost:8080/bets"),
		},
	}
	api := &restapi{
		conf: conf,
		ggw:  gambler.New(conf),
	}

	return api
}

// WithRouter initialize routes using julienschmidth httprouter
func (api *restapi) WithRouter(router *httprouter.Router) {
	router.Handle("GET", "/ggw/profile", api.guard(api.MyProfile))

	router.Handle("GET", "/ggw/topics", api.TopicList)                    // list of topics
	router.Handle("GET", "/ggw/topics/:id", api.Topic)                    // one topic
	router.Handle("GET", "/ggw/topics/:id/bet/", api.guard(api.TopicBet)) // list of bet on a topic but only expects 1

	router.Handle("GET", "/ggw/bets", api.guard(api.MyBets))    // list of bets
	router.Handle("POST", "/ggw/bets", api.guard(api.PlaceBet)) // place a bet
}

func (api *restapi) guard(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		res := rest.NewAPIResponse(w, r)
		session, err := r.Cookie("fa-session")
		if err != nil || session == nil || session.Value == "" {
			res.Error("You are not authenticated", err).Respond(http.StatusUnauthorized)
			return
		}

		jwtok, err := api.conf.Const.AuthClient.VerifySessionCookie(ctx, session.Value)
		if err != nil {
			res.Error("You are not authenticated", err).Respond(http.StatusUnauthorized)
			return
		}

		log.Tracef("TOKEN: %+v\n", jwtok)
		log.Tracef("CLAIMS: %+v\n", jwtok.Claims)

		email, _ := jwtok.Claims["email"]
		ctx = context.WithValue(ctx, global.Context.Email(), email)
		r = r.WithContext(ctx)
		next(w, r, p)
	}
}

// MyProfile ...
func (api *restapi) MyProfile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	res := rest.NewAPIResponse(w, r)
	email := ctx.Value(global.Context.Email()).(string)

	user, err := api.ggw.MyProfile(ctx, email)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed get your profile"), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Payload(user).Respond(http.StatusOK)
}

// TopicList find all topic where
// @state != draft
func (api *restapi) TopicList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := rest.NewAPIResponse(w, r)
	ctx := r.Context()

	topicURL, _ := url.Parse(api.conf.URL.Topic)
	rq := topicURL.Query()
	rq.Set("state", "published,closed,answered")

	//TODO: Don't use redirect!
	// http.Redirect(w, r, topicURL.String(), http.StatusMovedPermanently)
	topicURL.RawQuery = rq.Encode()
	result, err := api.ggw.Forward(ctx, topicURL.String())
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed get topics"), err).
			Respond(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Topics one
func (api *restapi) Topic(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	//TODO: Don't use redirect!
	http.Redirect(w, r, api.conf.GetTopicURL(id), http.StatusPermanentRedirect)
}

// TopicBet bet on a topic
func (api *restapi) TopicBet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	topic := p.ByName("id")
	owner := p.ByName("owner")

	rq := url.Values{
		"id":    []string{topic},
		"owner": []string{owner},
		"page":  []string{"1"},
		"size":  []string{"1"},
	}
	betURL, _ := url.Parse(api.conf.URL.Bet)
	betURL.RawQuery = rq.Encode()

	log.Trace("Redirect to:", betURL.String())

	//TODO: Don't use redirect!
	http.Redirect(w, r, betURL.String(), http.StatusPermanentRedirect)
}

// MyBets list
func (api *restapi) MyBets(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := rest.NewAPIResponse(w, r)
	ctx := r.Context()

	// optimistic coding ...
	// .MyBets is guarded endpoint so email should already put by middleware
	email := ctx.Value(global.Context.Email()).(string)
	result, err := api.ggw.MyBets(ctx, email)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed get topics"), err).
			Respond(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// PlaceBet ...
func (api *restapi) PlaceBet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := rest.NewAPIResponse(w, r)
	ctx := r.Context()
	pbt := &command.PlaceBet{}
	if err := defaultRequestUnwrapper(pbt)(r.Body); err != nil {
		res.Error("Failed to parse request body", err).Respond(http.StatusBadRequest)
		return
	}

	// optimistic coding ...
	// .PlaceBet is guarded endpoint so email should already put by middleware
	email := ctx.Value(global.Context.Email()).(string)

	bet, err := api.ggw.PlaceBet(ctx, email, pbt)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to place a bet"), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Payload(bet).Respond(http.StatusCreated)
}

// @obj: please send a pointer to a struct
func defaultRequestUnwrapper(obj interface{}) func(body io.ReadCloser) error {
	return func(body io.ReadCloser) error {
		data, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		wrapper := &dto.Wrapper{
			Data: obj,
		}
		return json.Unmarshal(data, wrapper)
	}
}

func defaultOnEmptyEnv(env, def string) string {
	obj := os.Getenv(env)
	if obj == "" {
		obj = def
	}

	return obj
}
