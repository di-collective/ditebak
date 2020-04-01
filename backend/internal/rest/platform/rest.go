package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"

	"github.com/di-collective/ditebak/backend/internal/usecase/platform"
	"github.com/di-collective/ditebak/backend/internal/usecase/platform/command"
	"github.com/di-collective/ditebak/backend/internal/usecase/platform/dto"
	"github.com/di-collective/ditebak/backend/pkg/exception"
	"github.com/di-collective/ditebak/backend/pkg/rest"

	"github.com/julienschmidt/httprouter"
)

// rest API for platform
type restapi struct {
	conf *platform.Config
	pgw  *platform.Gateway
}

// New platform micro API gateway
func New() rest.REST {
	ctx := context.Background()
	fap, _ := firebase.NewApp(ctx, nil)
	fac, _ := fap.Auth(ctx)

	conf := &platform.Config{
		Const: &platform.ConfigConst{
			SessionDuration: 5 * 24 * time.Hour,
			AuthClient:      fac,
			Prod:            os.Getenv("PROD") == "true",
		},
		URL: &platform.ConfigURL{
			Cred:  defaultOnEmptyEnv(os.Getenv("URL_CRED"), "http://localhost:8080/credentials"),
			User:  defaultOnEmptyEnv(os.Getenv("URL_USER"), "http://localhost:8080/users"),
			Topic: defaultOnEmptyEnv(os.Getenv("URL_TOPIC"), "http://localhost:8080/topics"),
			Bet:   defaultOnEmptyEnv(os.Getenv("URL_BET"), "http://localhost:8080/bets"),
		},
	}
	api := &restapi{
		conf: conf,
		pgw:  platform.New(conf),
	}

	return api
}

// WithRouter initialize routes using julienschmidth httprouter
func (api *restapi) WithRouter(router *httprouter.Router) {
	router.Handle("POST", "/pgw/login", api.Login)
	router.Handle("GET", "/pgw/logout", api.Logout)

	router.Handle("POST", "/pgw/answers", api.Answer) // answer a topic
}

// Login to platform
func (api *restapi) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := rest.NewAPIResponse(w, r)
	ctx := r.Context()
	lgn := &command.Login{}
	if err := defaultRequestUnwrapper(lgn)(r.Body); err != nil {
		res.Error("Failed to parse request body", err).Respond(http.StatusBadRequest)
		return
	}

	session, user, err := api.pgw.Login(ctx, lgn)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to login"), err).
			Respond(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Path:     "/",
		Name:     "fa-session",
		Value:    session,
		MaxAge:   int(api.conf.Const.SessionDuration.Seconds()),
		HttpOnly: true,
	}

	if api.conf.Const.Prod {
		cookie.Domain = ".ditebak.com"
		cookie.Secure = true
	}

	w.Header().Set("Set-Cookie", cookie.String())
	res.Payload(user).Respond(http.StatusOK)
}

// Login from platform
func (api *restapi) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cookie := http.Cookie{
		Path:     "/",
		Domain:   ".ditebak.com",
		Name:     "fa-session",
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
	}

	w.Header().Set("Set-Cookie", cookie.String())
	rest.NewAPIResponse(w, r).Respond(http.StatusOK)
	return
}

// Answer a topic
func (api *restapi) Answer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := rest.NewAPIResponse(w, r)
	ctx := r.Context()
	ans := &command.Answer{}
	if err := defaultRequestUnwrapper(ans)(r.Body); err != nil {
		res.Error("Failed to parse request body", err).Respond(http.StatusBadRequest)
		return
	}

	// go func() {
	// 	stat, err := api.pgw.Answer(ctx, ans)
	// 	log.Tracef("Answer stat: %+v, err: %v", stat, err)
	// }()

	stat, err := api.pgw.Answer(ctx, ans)
	exc, throw := exception.IsException(err)
	if throw {
		res.Error(exc.Message(), err).Respond(exc.Code())
		return
	} else if err != nil {
		res.Error(fmt.Sprintf("Failed to get answer a topic"), err).
			Respond(http.StatusInternalServerError)
		return
	}

	res.Payload(stat).Respond(http.StatusOK)
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
