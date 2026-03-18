package mailer

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"weather_task/internal/config"
	"weather_task/internal/model"
)

// Mailer sends transactional emails.
type Mailer interface {
	SendConfirmation(to, confirmURL string) error
	SendWeatherUpdate(to, city string, weather model.Weather, unsubscribeURL string) error
}

// SMTPMailer sends emails via SMTP using stdlib net/smtp.
type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
}

// New creates an SMTPMailer from application config.
func New(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUser,
		password: cfg.SMTPPassword,
	}
}

// SendConfirmation sends a subscription confirmation email.
func (m *SMTPMailer) SendConfirmation(to, confirmURL string) error {
	subject := "Please Confirm Your Weather Update Subscription"
	body := fmt.Sprintf(`<p>Good day!</p>
<p>To confirm your subscription to weather updates, <a href="%s">click here</a>.</p>
<p>If you did not sign up, simply ignore this message.</p>`, confirmURL)
	return m.send(to, subject, body)
}

// SendWeatherUpdate sends a weather update email with an unsubscribe link.
func (m *SMTPMailer) SendWeatherUpdate(to, city string, weather model.Weather, unsubscribeURL string) error {
	subject := fmt.Sprintf("Weather Update for %s", city)
	body := fmt.Sprintf(`<p>Current weather in <strong>%s</strong>:</p>
<ul>
  <li><strong>Temperature:</strong> %.1f &deg;C</li>
  <li><strong>Humidity:</strong> %.0f %%</li>
  <li><strong>Condition:</strong> %s</li>
</ul>
<p style="font-size:0.9em;color:#555;">
  <a href="%s">Unsubscribe</a>
</p>`,
		city,
		weather.Temperature,
		weather.Humidity,
		weather.Description,
		unsubscribeURL,
	)
	return m.send(to, subject, body)
}

func (m *SMTPMailer) send(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	msg := buildMessage(m.username, to, subject, htmlBody)
	return smtp.SendMail(addr, auth, m.username, []string{to}, msg)
}

func buildMessage(from, to, subject, htmlBody string) []byte {
	var sb strings.Builder
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + to + "\r\n")
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(htmlBody)
	return []byte(sb.String())
}

// LogMailer is a no-op mailer that logs emails instead of sending them.
// Use it when SMTP credentials are not configured (e.g. local development).
type LogMailer struct{}

func (LogMailer) SendConfirmation(to, confirmURL string) error {
	slog.Info("confirmation email (not sent)", "to", to, "confirm_url", confirmURL)
	return nil
}

func (LogMailer) SendWeatherUpdate(to, city string, w model.Weather, unsubscribeURL string) error {
	slog.Info("weather update email (not sent)",
		"to", to,
		"city", city,
		"temp", w.Temperature,
		"unsubscribe_url", unsubscribeURL,
	)
	return nil
}

// NewFromConfig returns an SMTPMailer when credentials are present,
// or a LogMailer when SMTP_USER / SMTP_PASSWORD are not set.
func NewFromConfig(cfg *config.Config) Mailer {
	if cfg.SMTPUser == "" || cfg.SMTPPassword == "" {
		slog.Warn("SMTP credentials not configured — emails will be logged only")
		return LogMailer{}
	}
	return New(cfg)
}
