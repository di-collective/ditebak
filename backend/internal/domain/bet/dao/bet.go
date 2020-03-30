package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Bet states
var (
	enum             = state("")
	BetStates states = &enum
)

type state string

type states interface {
	Placed() state
	Lost() state
	Won() state
}

func (t *state) Placed() state {
	return state("placed")
}

func (t *state) Lost() state {
	return state("lost")
}

func (t *state) Won() state {
	return state("won")
}

// Bet database object
type Bet struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt  *time.Time         `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at,omitempty"`
	TopicID    string             `json:"topic_id" bson:"topic_id,omitempty"`
	Owner      string             `json:"owner" bson:"owner,omitempty"` // who made the bet (email)
	Prediction string             `json:"prediction" bson:"prediction"` // whats his/her prediction
	Reputation int                `json:"reputation" bson:"reputation"` // how many reputation at stake
	State      state              `json:"state" bson:"state"`
}
