package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Topic states
var (
	enum               = state("")
	TopicStates states = &enum
)

type state string

type states interface {
	Draft() state
	Published() state
	Closed() state
	Answered() state
}

func (t *state) Draft() state {
	return state("draft")
}

func (t *state) Published() state {
	return state("published")
}

func (t *state) Closed() state {
	return state("closed")
}

func (t *state) Answered() state {
	return state("answered")
}

// Topic database object
type Topic struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt *time.Time         `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at,omitempty"`
	ClosingAt *time.Time         `json:"closing_at" bson:"closing_at"`
	Banner    string             `json:"banner" bson:"banner"`
	Question  string             `json:"question" bson:"question"`
	Answer    string             `json:"answer" bson:"answer"`
	Context   string             `json:"context" bson:"context"`
	State     state              `json:"state" bson:"state"`
}
