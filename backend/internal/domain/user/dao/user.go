package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Providers variety
var (
	enum      = provider("")
	Providers = &enum
)

type provider string

type providers interface {
	Google() provider
	Facebook() provider
	Email() provider
	Bot() provider
}

func (t *provider) Google() provider {
	return provider("google")
}

func (t *provider) Facebook() provider {
	return provider("facebook")
}

func (t *provider) Email() provider {
	return provider("email")
}

func (t *provider) Bot() provider {
	return provider("bot")
}

// User database object
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt   *time.Time         `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at" bson:"updated_at,omitempty"`
	VerifiedAt  *time.Time         `json:"verified_at" bson:"verified_at,omitempty"`
	Provider    provider           `json:"provider" bson:"provider,omitempty"`
	Email       string             `json:"email" bson:"email"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	Photo       string             `json:"photo" bson:"photo"`
	Reputation  int64              `json:"reputation" bson:"reputation"`
}
