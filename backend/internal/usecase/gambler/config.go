package gambler

import (
	"context"
	"net/url"
	"path"

	"firebase.google.com/go/auth"
)

// AuthClient masks firebase auth
type AuthClient interface {
	VerifySessionCookie(ctx context.Context, cookie string) (*auth.Token, error)
}

// ConfigConst ...
type ConfigConst struct {
	MaxStake   int
	AuthClient AuthClient
}

// ConfigURL ...
type ConfigURL struct {
	User  string
	Topic string
	Bet   string

	Login  string
	Logout string
}

// Config ...
type Config struct {
	Const *ConfigConst
	URL   *ConfigURL
}

// FindUserURL ...
func (conf *Config) FindUserURL(query url.Values) string {
	uri, _ := url.Parse(conf.URL.User)
	uri.RawQuery = query.Encode()

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
