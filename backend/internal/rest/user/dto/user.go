package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User dto
type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	VerifiedAt  *time.Time         `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	Provider    string             `json:"provider,omitempty" bson:"provider,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	DisplayName string             `json:"display_name,omitempty" bson:"display_name,omitempty"`
	FirstName   string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName    string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Photo       *string            `json:"photo,omitempty" bson:"photo,omitempty"`
	Reputation  *int64             `json:"reputation,omitempty" bson:"reputation,omitempty"`
}
