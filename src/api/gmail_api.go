package api

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/joho/godotenv"
	"google.golang.org/api/gmail/v1"
)

// GmailContent represents content of email
type GmailContent struct {
	Email   entity.Email
	Subject string
	Body    string
}

// Send send gmail
func (e GmailContent) Send() (err error) {
	if err = godotenv.Load(); err != nil {
		panic("failed to load .env file")
	}

	srv, err := gmail.New(googleOauth2Client())
	if err != nil {
		log.Print("Failed to connect gmail client")
	}

	temp := []byte("From: 'Share office'\r\n" +
		"reply-to: " + util.EMAIL + "\r\n" +
		"To: " + string(e.Email) + "\r\n" +
		"Subject: " + util.ConvertUtf8ToISOHelper(e.Subject) + "\r\n" +
		"\r\n" + e.Body)

	message := gmail.Message{}
	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send(util.EMAIL, &message).Do()
	if err != nil {
		return err
	}
	return nil
}
