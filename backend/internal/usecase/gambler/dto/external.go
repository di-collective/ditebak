package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Topic database object
type Topic struct {
	ID        primitive.ObjectID `json:"id"`
	CreatedAt *time.Time         `json:"created_at"`
	ClosingAt *time.Time         `json:"closing_at"`
	Banner    string             `json:"banner"`
	Question  string             `json:"question"`
	Answer    string             `json:"answer"`
	Context   string             `json:"context"`
	State     string             `json:"state"`
}

// User dto
type User struct {
	ID          string
	CreatedAt   *time.Time `json:"created_at"`
	VerifiedAt  *time.Time `json:"verified_at"`
	Provider    string     `json:"provider"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Photo       string     `json:"photo"`
	Reputation  int64      `json:"reputation"`
}

// Bet dto
type Bet struct {
	ID         string     `json:"id,omitempty"`
	TopicID    string     `json:"topic_id"`
	CreatedAt  *time.Time `json:"created_at"`
	Owner      string     `json:"owner"`      // who made the bet (email)
	Prediction string     `json:"prediction"` // whats his/her prediction
	Reputation int        `json:"reputation"` // how many reputation at stake
	State      string     `json:"state"`
}

// Wrapper to data
type Wrapper struct {
	Data interface{} `json:"data"`
}
