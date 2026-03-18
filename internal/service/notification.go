package service

import (
	"context"
	"fmt"
	"log/slog"

	"weather_task/internal/model"
	"weather_task/internal/repository"
	"weather_task/pkg/mailer"
)

// WeatherClient fetches current weather for a city.
// Implemented by *weather.Client.
type WeatherClient interface {
	GetWeather(ctx context.Context, city string) (model.Weather, error)
}

// NotificationService sends scheduled weather update emails.
type NotificationService interface {
	SendWeatherUpdates(ctx context.Context, frequency string) error
}

type notificationService struct {
	repo          repository.SubscriptionRepo
	weatherClient WeatherClient
	mailer        mailer.Mailer
	baseURL       string
}

// NewNotificationService wires up the notification service.
func NewNotificationService(
	repo repository.SubscriptionRepo,
	weatherClient WeatherClient,
	m mailer.Mailer,
	baseURL string,
) NotificationService {
	return &notificationService{
		repo:          repo,
		weatherClient: weatherClient,
		mailer:        m,
		baseURL:       baseURL,
	}
}

// SendWeatherUpdates fetches and emails the current weather to all confirmed
// subscribers with the given frequency. Per-subscriber errors are logged and
// skipped so that one bad address does not block the rest.
func (s *notificationService) SendWeatherUpdates(ctx context.Context, frequency string) error {
	subs, err := s.repo.FindConfirmedByFrequency(ctx, frequency)
	if err != nil {
		return fmt.Errorf("load %s subscriptions: %w", frequency, err)
	}

	for _, sub := range subs {
		if err := s.sendOne(ctx, sub); err != nil {
			slog.Error("send weather update", "email", sub.Email, "err", err)
		}
	}

	return nil
}

func (s *notificationService) sendOne(ctx context.Context, sub model.Subscription) error {
	w, err := s.weatherClient.GetWeather(ctx, sub.City)
	if err != nil {
		return fmt.Errorf("get weather for %s: %w", sub.City, err)
	}

	unsubURL := fmt.Sprintf("%s/api/subscriptions/%s", s.baseURL, sub.Token)
	return s.mailer.SendWeatherUpdate(sub.Email, sub.City, w, unsubURL)
}
