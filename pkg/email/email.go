package email

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

func init() {
	_ = godotenv.Load()
}

var Receiver string

func NewSend(to string) {
	Receiver = to
}

func ClearSenrder() {
	Receiver = ""
}

func SendEmail(subject, body string) error {
	if Receiver == "" {
		return fmt.Errorf("receiver email is not set")
	}

	message := gomail.NewMessage()

	message.SetHeader("From", os.Getenv("FROM_EMAIL"))
	message.SetHeader("To", Receiver)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", body)

	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

	return nil
}
