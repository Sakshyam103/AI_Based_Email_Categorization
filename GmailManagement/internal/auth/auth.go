package auth

import (
	"GmailManagement/internal/models"
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	MaxAge = 86400 // 1 day
	IsProd = false // for local use, use true for production
	// Update callback URL to point to your Vite frontend
	callbackURL = "http://localhost:8080/auth/google/callback" // Base frontend URL
)

var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Get session secret from environment
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET must be set in environment")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientId == "" {
		log.Fatal("GOOGLE_CLIENT_ID must be set in environment")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET must be set in environment")
	}

	//store := sessions.NewCookieStore([]byte(sessionSecret))
	Store.MaxAge(MaxAge)
	Store.Options.Path = "/"
	Store.Options.Domain = "localhost"
	Store.Options.HttpOnly = true
	Store.Options.Secure = IsProd
	Store.Options.SameSite = http.SameSiteLaxMode

	//gothic.Store = store
	//register our user model so that session can store user model
	gob.Register(models.User{})

	provider := google.New(googleClientId,
		googleClientSecret,
		callbackURL,
		"email",
		"profile", "https://www.googleapis.com/auth/gmail.readonly", "https://www.googleapis.com/auth/gmail.send",
	)
	provider.SetPrompt("consent")
	provider.SetAccessType("offline")
	//provider.Client()
	goth.UseProviders(
		provider,
	)

}
