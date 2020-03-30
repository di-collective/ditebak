package dto

import (
	"time"

	"github.com/di-collective/ditebak/backend/internal/usecase/platform/command"
)

// Topic database object
type Topic struct {
	ID        string     `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	ClosingAt *time.Time `json:"closing_at"`
	Banner    string     `json:"banner"`
	Question  string     `json:"question"`
	Answer    string     `json:"answer"`
	Context   string     `json:"context"`
	State     string     `json:"state"`
}

// User dto
type User struct {
	ID          string     `json:"id,omitempty"`
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

// FromLogin to Credential
func (u *User) FromLogin(login *command.Login) {
	var verified *time.Time
	if login.User.EmailVerified {
		now := time.Now()
		verified = &now
	}

	provider := "email"
	if login.Credential.ProviderID == "google.com" {
		provider = "google"
	}

	u.VerifiedAt = verified
	u.Provider = provider
	u.Email = login.User.Email
	u.FirstName = login.AdditionalUserInfo.Profile.GivenName
	u.LastName = login.AdditionalUserInfo.Profile.FamilyName
	u.DisplayName = login.User.DisplayName
	u.Photo = login.User.PhotoURL
	u.Reputation = 0
}

// Bet dto
type Bet struct {
	ID         string     `json:"id,omitempty"`
	TopicID    string     `json:"topic_id"`
	CreatedAt  *time.Time `json:"created_at"`
	Owner      string     `json:"owner"`      // who made the bet (email)
	Prediction string     `json:"prediction"` // whats his/her prediction
	Reputation int64      `json:"reputation"` // how many reputation at stake
	State      string     `json:"state"`
}

// Credential database object
type Credential struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email"`
	Firebase struct {
		ID             string `json:"id"`
		IDToken        string `json:"id_token"`
		AccessToken    string `json:"access_token"`
		RefreshToken   string `json:"refresh_token"`
		ExpirationTime int64  `json:"expiration_time"`
	} `json:"firebase"`
	Google struct {
		ID           string `json:"id"`
		IDToken      string `json:"id_token"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"google"`
}

// FromLogin to Credential
func (cred *Credential) FromLogin(login *command.Login) {
	cred.Email = login.User.Email
	cred.Firebase.ID = login.User.UID
	cred.Firebase.IDToken = ""
	cred.Firebase.AccessToken = login.User.StsTokenManager.AccessToken
	cred.Firebase.RefreshToken = login.User.StsTokenManager.AccessToken
	cred.Firebase.ExpirationTime = login.User.StsTokenManager.ExpirationTime

	cred.Google.ID = login.AdditionalUserInfo.Profile.ID
	cred.Google.IDToken = login.Credential.OauthIDToken
	cred.Google.AccessToken = login.Credential.OauthAccessToken
	cred.Google.RefreshToken = ""
}
