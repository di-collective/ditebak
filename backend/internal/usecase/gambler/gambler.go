package gambler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"github.com/di-collective/ditebak/backend/internal/usecase/gambler/command"
	"github.com/di-collective/ditebak/backend/internal/usecase/gambler/dto"
	"github.com/di-collective/ditebak/backend/pkg/exception"
)

type req struct {
	mtd   string
	res   string
	url   string
	pay   interface{}
	err   func(res *resty.Response) error
	parse func([]byte) error
}

// Gateway ...
type Gateway struct {
	conf *Config
	rc   *resty.Client
	ac   AuthClient
}

// New instance of gambler's gateway
func New(conf *Config) *Gateway {
	gw := &Gateway{
		conf: conf,
		rc:   resty.New().EnableTrace(),
		ac:   conf.Const.AuthClient,
	}

	return gw
}

// Forward request to API
func (gw *Gateway) Forward(ctx context.Context, uri string) ([]byte, error) {
	var result []byte
	err := gw.doReq(ctx, &req{
		mtd: "GET",
		url: uri,
		err: defaultResponseHandler,
		parse: func(b []byte) error {
			result = b
			return nil
		},
	})
	return result, err
}

// MyBets forwarded from API
func (gw *Gateway) MyBets(ctx context.Context, email string) ([]byte, error) {
	users := []*dto.User{}
	gw.doReq(ctx, &req{
		res: "users",
		mtd: "GET",
		url: gw.conf.FindUserURL(url.Values{
			"email": []string{email},
			"page":  []string{"1"},
			"size":  []string{"1"},
		}),
		err:   defaultResponseHandler,
		parse: defaultResponseUnwrapper(&users),
	})

	// if user is not found
	if len(users) <= 0 {
		return nil, exception.New(http.StatusUnauthorized, "You are not authenticated. Please login before placing bet")
	}

	var result []byte
	user := users[0]
	return result, gw.doReq(ctx, &req{
		mtd: "GET",
		res: "bets",
		url: gw.conf.GetBetURL("", user.ID),
		err: defaultResponseHandler,
		parse: func(b []byte) error {
			result = b
			return nil
		},
	})
}

// PlaceBet ...
// Verification:
// 1. Reputation at stake must be more than 0 and can't be more than configured maximum
// 2. User must exists
// 3. Topic must exists
// 	  - state is published
//    - not expired at closing time
// 4. Cannot bet more than once
// Then:
// create / place the bet!
func (gw *Gateway) PlaceBet(ctx context.Context, email string, pb *command.PlaceBet) (*dto.Bet, error) {
	// verify PB command
	if pb.Stake < 1 || pb.Stake > gw.conf.Const.MaxStake {
		return nil, exception.New(http.StatusBadRequest, "Reputation at stake must be between 1 and %d", gw.conf.Const.MaxStake)
	}
	if pb.Topic == "" {
		return nil, exception.New(http.StatusBadRequest, "Topic can't be empty")
	}

	users := []*dto.User{}
	topic := &dto.Topic{}
	bets := []*dto.Bet{}

	reqs := []*req{{
		res: "users",
		mtd: "GET",
		url: gw.conf.FindUserURL(url.Values{
			"email": []string{email},
			"page":  []string{"1"},
			"size":  []string{"1"},
		}),
		err:   defaultResponseHandler,
		parse: defaultResponseUnwrapper(&users),
	}, {
		res:   "topics",
		mtd:   "GET",
		url:   gw.conf.GetTopicURL(pb.Topic),
		err:   defaultResponseHandler,
		parse: defaultResponseUnwrapper(topic),
	}, {
		res:   "bets",
		mtd:   "GET",
		url:   gw.conf.GetBetURL(pb.Topic, email),
		err:   defaultResponseHandler,
		parse: defaultResponseUnwrapper(&bets),
	}}

	// fetch user, topic, and bet
	for _, req := range reqs {
		if err := gw.doReq(ctx, req); err != nil {
			return nil, err
		}
	}

	log.Tracef("PlaceBet User: %+v, Topic: %+v, Bets: %+v\n", users, topic, bets)

	// if user is not found
	if len(users) <= 0 {
		return nil, exception.New(http.StatusUnauthorized, "You are not authenticated. Please login before placing bet")
	}

	// verify topic
	if topic.State != "published" {
		switch topic.State {
		case "closed", "answered":
			return nil, exception.New(http.StatusBadRequest, "Topic is already closed")
		case "draft":
			return nil, exception.New(http.StatusBadRequest, "Topic is not published yet")
		default:
			return nil, exception.New(http.StatusBadRequest, "Topic state is unknown")
		}
	}

	// now exceed closing time
	if time.Now().Unix() > topic.ClosingAt.Unix() {
		return nil, exception.New(http.StatusBadRequest, "Topic is already closed")
	}

	// verify bet
	if len(bets) > 0 {
		return nil, exception.New(http.StatusConflict, "Bet already exists")
	}

	// create a bet
	user := users[0]
	return gw.createBet(ctx, &dto.Bet{
		Owner:      user.ID,
		TopicID:    pb.Topic,
		Prediction: pb.Prediction,
		Reputation: pb.Stake,
	})
}

// MyProfile fetch currently logged in user profile based on email
func (gw *Gateway) MyProfile(ctx context.Context, email string) (*dto.User, error) {
	users := []*dto.User{}
	if err := gw.doReq(ctx, &req{
		mtd: "GET",
		res: "users",
		url: gw.conf.FindUserURL(url.Values{
			"page":  []string{"1"},
			"size":  []string{"1"},
			"email": []string{email},
		}),
		err:   defaultResponseHandler,
		parse: defaultResponseUnwrapper(&users),
	}); err != nil {
		return nil, err
	}

	if len(users) <= 0 {
		return nil, exception.New(http.StatusNotFound, "User not found")
	}

	return users[0], nil
}

func (gw *Gateway) doReq(ctx context.Context, req *req) error {
	var call func(string) (*resty.Response, error)

	api := gw.rc.R().SetContext(ctx)
	switch req.mtd {
	case http.MethodPost:
		call = api.Post
	case http.MethodPut:
		call = api.Put
	case http.MethodPatch:
		call = api.Patch
	default:
		call = api.Get
	}

	if req.pay != nil {
		api.SetBody(req.pay)
	}

	res, err := call(req.url)
	if err != nil {
		return exception.New(http.StatusBadGateway, "Failed to [%s] to url: %s, err: %v", req.mtd, req.url, err)
	}

	// delegated error condition
	if err = req.err(res); err != nil {
		return err
	}

	// delegate response parsing
	if err = req.parse(res.Body()); err != nil {
		return exception.New(http.StatusBadGateway, "Failed to parse response from url: %s, err: %s", req.url, err.Error())
	}

	return nil
}

func (gw *Gateway) createBet(ctx context.Context, bet *dto.Bet) (*dto.Bet, error) {
	payload := &dto.Wrapper{Data: bet}
	return bet, gw.doReq(ctx, &req{
		res:   "bets",
		mtd:   "POST",
		url:   gw.conf.URL.Bet,
		pay:   payload,
		err:   defaultResponseHandler,
		parse: defaultResponseUnwrapper(bet),
	})
}

func defaultResponseHandler(res *resty.Response) error {
	if res.IsError() {
		switch res.StatusCode() {
		case http.StatusRequestTimeout:
			return exception.New(http.StatusGatewayTimeout, "Request timed out to [%s] url: %s", res.Request.Method, res.Request.URL)
		case http.StatusBadRequest:
			return exception.New(http.StatusBadRequest, "Invalid request to [%s] url: %s", res.Request.Method, res.Request.URL)
		case http.StatusNotFound:
			return exception.New(http.StatusNotFound, "Resource with [%s] url: %s, is not found", res.Request.Method, res.Request.URL)
		}

		return exception.New(res.StatusCode(), "Failed to [%s] to url: %s", res.Request.Method, res.Request.URL)
	}

	return nil
}

// @obj: please send a pointer to a struct
func defaultResponseUnwrapper(obj interface{}) func(body []byte) error {
	return func(body []byte) error {
		wrapper := &dto.Wrapper{
			Data: obj,
		}
		return json.Unmarshal(body, wrapper)
	}
}
