package api

import (
	"encoding/base64"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Yukio0315/reservation-backend/src/entity"
	"github.com/Yukio0315/reservation-backend/src/util"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// SendGmail send gmail
func SendGmail(email entity.EmailContent) (err error) {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env file")
	}

	config := oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{"https://mail.google.com/"},
	}

	expiry, _ := time.Parse("2006-01-02 03:04:05", "2020-04-14 23:15:00")
	token := oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN_GMAIL"),
		TokenType:    "Bearer",
		RefreshToken: os.Getenv("REFRESH_TOKEN_GMAIL"),
		Expiry:       expiry,
	}

	client := config.Client(oauth2.NoContext, &token)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatal("Failed to connect gmail client")
	}

	temp := []byte("From: 'Share office'\r\n" +
		"reply-to: kurosunotai@gmail.com\r\n" +
		"To: " + string(email.Email) + "\r\n" +
		"Subject: " + util.ConvertUtf8ToISOHelper(email.Subject) + "\r\n" +
		"\r\n" + email.Body)

	var message gmail.Message
	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send("kurosunotai@gmail.com", &message).Do()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
