package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Topic database object
type Topic struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	ClosingAt *time.Time         `json:"closing_at,omitempty" bson:"closing_at,omitempty"`
	Banner    string             `json:"banner,omitempty" bson:"banner,omitempty"`
	Question  string             `json:"question,omitempty" bson:"question,omitempty"`
	Answer    string             `json:"answer,omitempty" bson:"answer,omitempty"`
	Context   string             `json:"context,omitempty" bson:"context,omitempty"`
	State     string             `json:"state,omitempty" bson:"state,omitempty"`
}
