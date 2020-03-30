package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Credential database object
type Credential struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Firebase struct {
		ID             string `json:"id" bson:"id"`
		IDToken        string `json:"id_token" bson:"id_token"`
		AccessToken    string `json:"access_token" bson:"access_token"`
		RefreshToken   string `json:"refresh_token" bson:"refresh_token"`
		ExpirationTime int64  `json:"expiration_time" bson:"expiration_time"`
	} `json:"firebase" bson:"firebase"`
	Google struct {
		ID           string `json:"id" bson:"id"`
		IDToken      string `json:"id_token" bson:"id_token"`
		AccessToken  string `json:"access_token" bson:"access_token"`
		RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	} `json:"google" bson:"google"`
}
