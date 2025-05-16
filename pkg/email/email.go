package email

import (
	"fmt"
	"os"
	"time"

	"weather_task/internal/db"
	"weather_task/internal/service"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

func init() {
	_ = godotenv.Load()
}

func sendEmail(to, subject, body string) error {
	if to == "" {
		return fmt.Errorf("receiver email is not set")
	}
	if subject == "" {
		return fmt.Errorf("subject is not set")
	}
	if body == "" {
		return fmt.Errorf("body is not set")
	}

	message := gomail.NewMessage()

	message.SetHeader("From", os.Getenv("SMTP_USER"))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", body)

	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))

	if err := dialer.DialAndSend(message); err != nil {
		return err
	} else {
		fmt.Println("Email sent successfully!")
	}

	return nil
}

func SendConfirmationEMail(to, token string) error {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	confirmURL := fmt.Sprintf("%s/api/confirm/%s", baseURL, token)

	subject := "Please Confirm Your Weather Update Subscription"
	htmlBody := fmt.Sprintf(
		`<p>Good day!</p>
			 <p>To confirm your subscription to weather updates, <a href="%s">click here</a>.</p>
			 <p>If you did not sign up, simply ignore this message.</p>
			 <hr>`,
		confirmURL)

	return sendEmail(to, subject, htmlBody)
}

func SendWeatherEmail(to string) error {
	user, err := db.GetUser(to)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	client, err := service.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	weather, err := client.GetWeather(user.City)
	if err != nil {
		return fmt.Errorf("failed to get weather: %w", err)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	unsubscribeLink := fmt.Sprintf("%s/api/unsubscribe/%s", baseURL, user.Token)

	date := time.Now().Format("2006-01-02 15:04:05")

	subject := "Weather Update"
	htmlBody := fmt.Sprintf(
		`<p>Good day! Below is the detailed weather forecast for <strong>%s</strong> in %s:</p>
			 <ul style="list-style: none; padding: 0;">
				 <li><strong>🌡️ Temperature:</strong> %.1f °C</li>
				 <li><strong>💧 Humidity:</strong> %.0f %%</li>
				 <li><strong>☁️ Description:</strong> %s</li>
			 </ul>
			 <hr>
			 <p style="font-size:0.9em; color: #555;">
				 If you wish to unsubscribe, please click here:
				 <a href="%s">Unsubscribe</a>
			 </p>`,
		date,
		user.City,
		weather.Temperature,
		weather.Humidity,
		weather.Description,
		unsubscribeLink,
	)

	return sendEmail(to, subject, htmlBody)
}
