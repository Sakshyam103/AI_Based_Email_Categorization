package agent

import (
	database2 "GmailManagement/internal/database"
	"GmailManagement/internal/models"
	"encoding/base64"
	"fmt"
	"github.com/markbates/goth"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
	"os"
	"time"
)

var database = database2.New()

func main() {
	ticker := time.NewTicker(10 * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Fetching emails")
				user, err := database.GetUserByEmail(os.Getenv("EMAILS"))
				if err != nil {
					fmt.Println("Errors fetching emails: ", err.Error())
					return
				}
				if user.TokenExpiry.Unix() < time.Now().Unix() {
					provider, err := goth.GetProvider("google")
					if err != nil {
						fmt.Println("Could not find provider: ", err.Error())
						//w.WriteHeader(http.StatusBadRequest)
						//w.Write([]byte("Could not find out provider"))
						return
					}
					newToken, err := provider.RefreshToken(user.RefreshToken)
					if err != nil {
						fmt.Println("Could not generate new token: ", err.Error())
						//w.WriteHeader(http.StatusUnauthorized)
						//w.Write([]byte("Could not generate new token"))
						return
					}
					//update user with new access token and expiry
					user.AccessToken = newToken.AccessToken
					user.TokenExpiry = time.Now().Add(time.Hour * 1) //add expiry time for 1 hour
					err = database.UpdateUser(*user)
					if err != nil {
						fmt.Println("Error while updating user: ", err.Error())
						return
					}
				}
				//emails, err := s.getEmailsFromDB()
				emails, err := getEmails(user)
				if err != nil {
					fmt.Println("Error while fetching email: ", err.Error())
					return
					//w.WriteHeader(http.StatusUnauthorized)
				}
				fmt.Println("Fetched emails: ", len(emails))

				//w.Header().Set("Content-Type", "application/json")
				//json.NewEncoder(w).Encode(emails)
				//getEmailService(user.AccessToken)
			case <-done:
				fmt.Println("Fetching emails done")
				return
			}
		}
	}()
	select {} //running this forever
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

func getEmails(user *models.User) ([]models.RawEmails, error) {
	srv, err := getEmailService(user.AccessToken)
	if err != nil {
		return nil, err
	}
	userId := "me"
	pageToken := ""
	q := "newer_than:10m"
	var emails = []models.RawEmails{}
	for {
		r := srv.Users.Messages.List(userId)
		r.MaxResults(20)
		r.PageToken(pageToken)
		r.Q(q)
		resp, err := r.Do()
		if err != nil {
			return nil, err
		}
		for _, m := range resp.Messages {
			id, from, to, subject, body, date, historyId, err := GetEmailDetails(srv, m.Id)
			if err != nil {
				return nil, err
			}
			//after getting from, to, subject, body, er
			rawEmail := models.RawEmails{EmailId: id, From: from, To: to, Subject: subject, HistoryId: historyId, Body: body, Date: date}

			emails = append(emails, rawEmail)
			err = database.StoreEmail(rawEmail)
			fmt.Printf("From: %s, To: %s, Subject: %s, Body: %s\n", from, to, subject, body)
		}
		pageToken = resp.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return emails, nil
}

func GetEmailDetails(srv *gmail.Service, id string) (id1 string, from, to, subject, body, date string, historyid uint64, err error) {
	// Fetch the email message with format="full"
	gmailMessage, err := srv.Users.Messages.Get("me", id).Format("full").Do()
	if err != nil {
		return "", "", "", "", "", "", 0, err
	}

	if gmailMessage == nil || gmailMessage.Payload == nil {
		return "", "", "", "", "", "", 0, fmt.Errorf("message or payload is nil")
	}

	// Extract headers
	headers := gmailMessage.Payload.Headers
	historyid = gmailMessage.HistoryId
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

	return id, from, to, subject, body, date, historyid, nil
}

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
