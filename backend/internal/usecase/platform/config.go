package platform

import (
	"context"
	"net/url"
	"path"
	"time"

	"firebase.google.com/go/auth"
)

// AuthClient masks firebase auth
type AuthClient interface {
	VerifyIDTokenAndCheckRevoked(ctx context.Context, token string) (*auth.Token, error)
	SessionCookie(ctx context.Context, token string, dur time.Duration) (string, error)
}

// ConfigConst ...
type ConfigConst struct {
	Prod            bool
	SessionDuration time.Duration
	AuthClient      AuthClient
}

// ConfigURL ...
type ConfigURL struct {
	Cred  string
	User  string
	Topic string
	Bet   string
}

// Config ...
type Config struct {
	Const *ConfigConst
	URL   *ConfigURL
}

// GetUserURL based on user id
func (conf *Config) GetUserURL(id string) string {
	uri, _ := url.Parse(conf.URL.User)
	uri.Path = path.Join(uri.Path, id)
	return uri.String()
}

// GetTopicURL based on topic id
func (conf *Config) GetTopicURL(id string) string {
	uri, _ := url.Parse(conf.URL.Topic)
	uri.Path = path.Join(uri.Path, id)
	return uri.String()
}

// GetBetURL based on topic and owner
func (conf *Config) GetBetURL(topic, owner string) string {
	uri, _ := url.Parse(conf.URL.Bet)

	q := uri.Query()
	q.Set("topic", topic)
	q.Set("owner", owner)

	uri.RawQuery = q.Encode()
	return uri.String()
}

// GetMyBetURL based on owner
func (conf *Config) GetMyBetURL(owner string) string {
	uri, _ := url.Parse(conf.URL.Bet)

	q := uri.Query()
	q.Set("owner", owner)

	uri.RawQuery = q.Encode()
	return uri.String()
}
