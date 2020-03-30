package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/di-collective/ditebak/backend/internal/domain/bet/dao"
	betDao "github.com/di-collective/ditebak/backend/internal/domain/bet/dao"
	topicDao "github.com/di-collective/ditebak/backend/internal/domain/topic/dao"
	"github.com/di-collective/ditebak/backend/internal/usecase/platform/command"
	"github.com/di-collective/ditebak/backend/internal/usecase/platform/dto"
	"github.com/di-collective/ditebak/backend/pkg/exception"
	"github.com/go-resty/resty"
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

// New instance of platform API gateway
func New(conf *Config) *Gateway {
	gw := &Gateway{
		conf: conf,
		rc:   resty.New().EnableTrace().SetRetryCount(3),
		ac:   conf.Const.AuthClient,
	}

	return gw
}

// Login to platform
// Process:
// 1. Create user & credential if new,
// 2. Else, update credential
// 3. Verify the token
// 4. Creat, store and return session string
func (gw *Gateway) Login(ctx context.Context, login *command.Login) (string, *dto.User, error) {
	log.Tracef("Login %+v\n", login)
	method := "POST"
	user := &dto.User{}
	cred := &dto.Credential{}
	cred.FromLogin(login)
	credURL, _ := url.Parse(gw.conf.URL.Cred)
	if login.AdditionalUserInfo.IsNewUser {
		user.FromLogin(login)

		// 1. Create user and credentials
		if err := gw.doReq(ctx, &req{
			mtd:   "POST",
			res:   "users",
			url:   gw.conf.URL.User,
			pay:   &dto.Wrapper{Data: &user},
			err:   defaultResponseHandler,
			parse: nil,
		}); err != nil {
			return "", nil, err
		}
	} else {
		// 1. Get existing user's credential
		emailQuery := url.Values{
			"page":  []string{"1"},
			"size":  []string{"1"},
			"email": []string{login.User.Email},
		}.Encode()
		creds := []*dto.Credential{}
		credURL.RawQuery = emailQuery
		if err := gw.doReq(ctx, &req{
			mtd:   "GET",
			res:   "credentials",
			url:   credURL.String(),
			err:   defaultResponseHandler,
			parse: defaultResponseUnwrapper(&creds),
		}); err != nil {
			return "", nil, err
		}

		if len(creds) < 1 {
			return "", nil, exception.New(http.StatusInternalServerError, "Failed to authenticate credential")
		}

		cred = creds[0]
		method = "PATCH"
		credURL.RawQuery = ""
		credURL.Path = path.Join(credURL.Path, cred.ID)

		// Get existing user's
		users := []*dto.User{}
		userURL, _ := url.Parse(gw.conf.URL.User)
		userURL.RawQuery = emailQuery
		if err := gw.doReq(ctx, &req{
			mtd:   "GET",
			res:   "users",
			url:   userURL.String(),
			err:   defaultResponseHandler,
			parse: defaultResponseUnwrapper(&users),
		}); err != nil {
			return "", nil, err
		}

		if len(creds) < 1 {
			return "", nil, exception.New(http.StatusInternalServerError, "Failed to authenticate credential")
		}
		user = users[0]
	}

	// 2. Create / Update user's credential
	gw.doReq(ctx, &req{
		mtd:   method,
		res:   "credentials",
		url:   credURL.String(),
		pay:   &dto.Wrapper{Data: cred},
		err:   defaultResponseHandler,
		parse: defaultResponseUnwrapper(cred),
	})

	// 3. Verify session
	token := login.User.StsTokenManager.AccessToken
	if _, err := gw.ac.VerifyIDTokenAndCheckRevoked(ctx, token); err != nil {
		return "", nil, err
	}

	// TODO: store session string for in-mem cache
	// 4. Create and return session
	duration := gw.conf.Const.SessionDuration
	session, err := gw.ac.SessionCookie(ctx, token, duration)
	return session, user, err
}

// Answer an existing topic
// Verification:
// 1. Topic must exists
//    - already past closing time
// Then:
// 1. Update topic as answered
// 2. Find all bets with state = placed
// 3. Flag wrong answer with "lost"
// 4. Flag correct answer with "won"
// 5. Reward karma!
func (gw *Gateway) Answer(ctx context.Context, ans *command.Answer) (*dto.Answered, error) {
	if ans.Topic == "" {
		return nil, exception.New(http.StatusBadRequest, "Topic can't be empty")
	}
	if ans.Answer == "" {
		return nil, exception.New(http.StatusBadRequest, "Answer can't be empty")
	}

	stat := &dto.Answered{
		Lost:  0,
		Won:   0,
		Total: 0,
	}

	// 1. Update the topic as answered
	if err := gw.doReq(ctx, &req{
		mtd: "PATCH",
		res: "topics",
		url: gw.conf.GetTopicURL(ans.Topic),
		pay: &dto.Wrapper{Data: &map[string]interface{}{
			"answer": ans.Answer,
			"state":  topicDao.TopicStates.Answered(),
		}},
		err:   defaultResponseHandler,
		parse: nil,
	}); err != nil {
		log.Errorln("Failed to answer topic:", err)
		return nil, err
	}

	// 2. Find all bets on a topic
	bets := []*dto.Bet{}
	betURL, _ := url.Parse(gw.conf.URL.Bet)
	page, rpp := 0, "100"
	for {
		// partial paging
		page++
		rq := url.Values{
			"topic": []string{ans.Topic},
			"state": []string{string(betDao.BetStates.Placed())},
			"page":  []string{fmt.Sprintf("%d", page)},
			"size":  []string{rpp},
		}
		betURL.RawQuery = rq.Encode()

		partial := []*dto.Bet{}
		if err := gw.doReq(ctx, &req{
			mtd:   "GET",
			res:   "bets",
			url:   betURL.String(),
			err:   defaultResponseHandler,
			parse: defaultResponseUnwrapper(&partial),
		}); err != nil {
			log.Errorln("Failed to collect all bets:", err)
			return nil, err
		}

		// no more bets, break from loop
		if len(partial) <= 0 {
			break
		}

		bets = append(bets, partial...)
	}

	// - Nobody bets
	if len(bets) <= 0 {
		return stat, nil
	}

	// 3. Find winner
	betURL.RawQuery = "" // clean betURL query
	stat.Total = len(bets)
	betRequests := make([]*req, stat.Total)
	for i, bet := range bets {
		if ans.IsTrue(bet.Prediction) {
			bet.State = string(betDao.BetStates.Won())
			stat.Won++
		} else {
			bet.State = string(dao.BetStates.Lost())
			stat.Lost++
		}

		betURL.Path = path.Join(betURL.Path, bet.ID)
		betRequests[i] = &req{
			mtd:   "PATCH",
			res:   "bets",
			url:   betURL.String(),
			pay:   &dto.Wrapper{Data: bet},
			err:   defaultResponseHandler,
			parse: nil,
		}
	}

	// 4. Update all bets to database
	for _, req := range betRequests {
		if err := gw.doReq(ctx, req); err != nil {
			log.Errorf("Failed to update bet: %+v", req.pay)
			continue
		}
	}

	// 5. Reward karma!
	for _, bet := range bets {
		user := &dto.User{}
		userURL := gw.conf.GetUserURL(bet.Owner)
		if err := gw.doReq(ctx, &req{
			mtd:   "GET",
			res:   "users",
			url:   userURL,
			err:   defaultResponseHandler,
			parse: defaultResponseUnwrapper(user),
		}); err != nil {
			payload, _ := json.Marshal(bet)
			log.Errorf("Failed to update user's karma, err: %s, payload: %s", err, string(payload))
			continue
		}

		rep := user.Reputation
		switch bet.State {
		case string(betDao.BetStates.Lost()):
			rep -= bet.Reputation
		case string(betDao.BetStates.Won()):
			rep += bet.Reputation
		default:
			continue // unknown state
		}

		if err := gw.doReq(ctx, &req{
			mtd: "PATCH",
			res: "users",
			url: userURL,
			pay: &dto.Wrapper{Data: &map[string]interface{}{
				"reputation": rep,
			}},
			err:   defaultResponseHandler,
			parse: nil,
		}); err != nil {
			payload, _ := json.Marshal(bet)
			log.Errorf("Failed to update user's karma, err: %s, payload: %s", err, string(payload))
			continue
		}
	}

	return stat, nil
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
	if req.parse != nil {
		if err = req.parse(res.Body()); err != nil {
			return exception.New(http.StatusBadGateway, "Failed to parse response from url: %s, err: %s", req.url, err.Error())
		}
	}

	return nil
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
