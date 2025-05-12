package server

import (
	"GmailManagement/internal/auth"
	"GmailManagement/internal/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"
	//"github.com/gorilla/pat"
	"net/http"

	//"github.com/gorilla/pat"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Get("/auth/{provider}/callback", s.CallbackHandler)

	r.Get("/health", s.healthHandler)

	r.Get("/home", func(writer http.ResponseWriter, request *http.Request) {
	})

	r.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		q := r.URL.Query()
		q.Add("provider", provider)
		r.URL.RawQuery = q.Encode()
		//gothic.StoreInSession("provider", provider, r, w)
		gothic.GetContextWithProvider(r, provider)
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			fmt.Println(gothUser)
			http.Redirect(w, r, "http://localhost:8080/home", http.StatusTemporaryRedirect)
		}
		gothic.BeginAuthHandler(w, r)
	})

	r.Get("/emails", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		q := r.URL.Query()
		q.Add("provider", provider)
		r.URL.RawQuery = q.Encode()
		//gothic.StoreInSession("provider", provider, r, w)
		gothic.GetContextWithProvider(r, provider)
		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			fmt.Println(gothUser)
			http.Redirect(w, r, "http://localhost:8080/home", http.StatusTemporaryRedirect)
		}
		gothic.BeginAuthHandler(w, r)
	})

	r.Post("/sendEmails", s.sendEmails)
	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		session, _ := auth.Store.Get(r, "session")
		session.Values["jwt"] = nil // need to set path to the frontend
		err := gothic.Logout(w, r)
		if err != nil {
			return
		}
		if err := session.Save(r, w); err != nil {
			fmt.Println("Error saving session:", err)
		}
		msg := map[string]string{"msg": "successfully logged out"}
		json.NewEncoder(w).Encode(msg)
	})

	r.Get("/getemails", s.handleEmails)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request to /") // Add debug logging
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	r.Get("/whoamI", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if checkCookie(w, r) {
			cookie, err := auth.Store.Get(r, "session")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Cookies not found"))
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if cookie.Values["jwt"] != nil {
				fmt.Println(cookie)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("logged in"))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("unauthorized"))
			}
		}
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081") // Allow frontend
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow cookies & sessions
}

func (s *Server) sendEmails(w http.ResponseWriter, r *http.Request) {
	//enableCORS(w, r)
	if checkCookie(w, r) {
		user1, err := s.db.GetUserByEmail(os.Getenv("EMAIL"))
		if err != nil {
			fmt.Println("Error while fetching user")
			return
		}
		fmt.Println(r.Body)
		var email models.EmailData
		err = json.NewDecoder(r.Body).Decode(&email)
		fmt.Println(email)
		if err != nil {
			fmt.Println("Error while decoding email body: ", err)
			//json.NewEncoder(w).Encode()
			return
		}
		if user1.TokenExpiry.Unix() < time.Now().Unix() {
			provider, err := goth.GetProvider("google")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Could not find out provider"))
				return
			}
			newToken, err := provider.RefreshToken(user1.RefreshToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Could not generate new token"))
				return
			}
			//update user with new access token and expiry
			user1.AccessToken = newToken.AccessToken
			user1.TokenExpiry = time.Now().Add(time.Hour * 1) //add expiry time for 1 hour
			err = s.db.UpdateUser(*user1)
			if err != nil {
				return
			}
		}
		srv, err := getEmailService(user1.AccessToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Could not find out email service"))
			return
		}
		var message gmail.Message
		emailTo := mail.Address{"", email.To}
		emailFrom := mail.Address{"", "me"}
		header := make(map[string]string)
		header["From"] = emailFrom.String()
		header["To"] = emailTo.String()
		header["Subject"] = email.Subject
		header["Content-Type"] = "text/html; charset=UTF-8"

		var msg strings.Builder
		for k, v := range header {
			fmt.Fprintf(&msg, "%s: %s\r\n", k, v)
		}
		fmt.Fprintf(&msg, "\r\n %s", email.Body)
		message.Raw = base64.URLEncoding.EncodeToString([]byte(msg.String()))
		_, err = srv.Users.Messages.Send("me", &message).Do()
		if err != nil {
			fmt.Println("unable to send email: %v", err)
			return
		}
		return
	}
	return
}

func checkCookie(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := auth.Store.Get(r, "session")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Cookies not found"))
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unauthorized"))
		//w.Write("unauthorized")
		return false
	}
	jwts, ok := cookie.Values["jwt"].(string)
	if !ok || jwts == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return false
	}
	claims, err := validateToken(jwts)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return false
	}
	fmt.Println(len(claims))
	return true
}

func (s *Server) handleEmails(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	emails, err := s.getEmailsFromDB()
	//emails, err := s.getEmails(*user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emails)
}

func (s *Server) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	//if provider == r.URL.Query().Get("provider") {
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()
	sessions.Save(r, w)
	//q.Set("google", provider)
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	//generate session id
	sessionid, err := generateSessionid()
	if err != nil {
		http.Error(w, "Failed to generate session ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user1, err := s.db.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println("Failed to get user by email")
		if user1 == nil {
			u := models.User{GoogleID: user.UserID, Email: user.Email, Name: user.Name, FamilyName: user.LastName, ProfilePicture: user.AvatarURL, AccessToken: user.AccessToken, RefreshToken: user.RefreshToken, TokenExpiry: user.ExpiresAt, CreatedAt: time.Now(), ExpiresAt: time.Now().Add(time.Hour * 1)}
			err = s.db.StoreUser(u)
			if err != nil {
				return
			}
		}
	}

	jwtstring, err := createToken(sessionid)
	if err != nil {
		fmt.Println(err)
	}

	//create a cookie sessionls
	session, _ := auth.Store.Get(r, "session")
	session.Values["jwt"] = jwtstring // need to set path to the frontend
	if err := session.Save(r, w); err != nil {
		fmt.Println("Error saving session:", err)
	}
	http.Redirect(w, r, "http://localhost:8081/", http.StatusFound)
}

func (s *Server) getEmailsFromDB() ([]models.CategorizedEmail, error) {
	//var emails []models.RawEmails
	emails, err := s.db.GetAllCategorizedEmails()
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func getEmailService(accessToken string) (*gmail.Service, error) {
	client := &http.Client{}
	client.Transport = &oauth2.Transport{Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}),
		Base: http.DefaultTransport}
	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func GetEmailDetails(srv *gmail.Service, id string) (id1 string, from, to, subject, body, date string, err error) {
	// Fetch the email message with format="full"
	gmailMessage, err := srv.Users.Messages.Get("me", id).Format("full").Do()
	if err != nil {
		return "", "", "", "", "", "", err
	}

	if gmailMessage == nil || gmailMessage.Payload == nil {
		return "", "", "", "", "", "", fmt.Errorf("message or payload is nil")
	}

	// Extract headers
	headers := gmailMessage.Payload.Headers
	for _, header := range headers {
		switch header.Name {
		case "From":
			from = header.Value
		case "To":
			to = header.Value
		case "Subject":
			subject = header.Value
		case "Date":
			date = header.Value
		}
	}

	// Extract body content
	body = extractEmailBody(gmailMessage.Payload)

	return id, from, to, subject, body, date, nil
}

// Helper function to extract the email body from the payload
func extractEmailBody(payload *gmail.MessagePart) string {
	if payload.Body != nil && payload.Body.Data != "" {
		data, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err != nil {
			log.Println("Error decoding body data:", err)
			return ""
		}
		return string(data)
	}

	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" || part.MimeType == "text/html" {
			data, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err != nil {
				log.Println("Error decoding part data:", err)
				continue
			}
			return string(data)
		}
	}

	return ""
}

func generateSessionid() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sessionId": username,
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(os.Getenv("UID_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("UID_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Invalid token")
}
