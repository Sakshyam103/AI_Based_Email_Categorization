package models

import (
	"time"
)

type User struct {
	//ID             string    `json:"id"`
	GoogleID       string    `json:"google_id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	GivenName      string    `json:"given_name"`
	FamilyName     string    `json:"family_name"`
	ProfilePicture string    `json:"profile_picture"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenExpiry    time.Time `json:"token_expiry"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}
