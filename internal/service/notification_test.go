package service

import (
	"context"
	"testing"

	"weather_task/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockWeatherClient implements WeatherClient for tests.
type mockWeatherClient struct {
	weather model.Weather
	err     error
}

func (m *mockWeatherClient) GetWeather(_ context.Context, _ string) (model.Weather, error) {
	return m.weather, m.err
}

// mockNotifRepo embeds mockRepo to satisfy SubscriptionRepo for notification tests.
type mockNotifRepo struct {
	subs    []model.Subscription
	findErr error
}

func (m *mockNotifRepo) Create(_ context.Context, _ *model.Subscription) error { return nil }
func (m *mockNotifRepo) FindByEmail(_ context.Context, _ string) (*model.Subscription, error) {
	return nil, nil
}
func (m *mockNotifRepo) FindByToken(_ context.Context, _ string) (*model.Subscription, error) {
	return nil, nil
}
func (m *mockNotifRepo) Update(_ context.Context, _ *model.Subscription) error { return nil }
func (m *mockNotifRepo) Delete(_ context.Context, _ *model.Subscription) error { return nil }
func (m *mockNotifRepo) FindConfirmedByFrequency(_ context.Context, _ string) ([]model.Subscription, error) {
	return m.subs, m.findErr
}

// mockNotifMailer records SendWeatherUpdate calls.
type mockNotifMailer struct {
	calls []string
	err   error
}

func (m *mockNotifMailer) SendConfirmation(_, _ string) error { return nil }
func (m *mockNotifMailer) SendWeatherUpdate(to, _ string, _ model.Weather, _ string) error {
	m.calls = append(m.calls, to)
	return m.err
}

func TestSendWeatherUpdates_Success(t *testing.T) {
	repo := &mockNotifRepo{subs: []model.Subscription{
		{Email: "a@example.com", City: "Kyiv", Token: "tok1"},
		{Email: "b@example.com", City: "Lviv", Token: "tok2"},
	}}
	wc := &mockWeatherClient{weather: model.Weather{Temperature: 20, Humidity: 60, Description: "Clear"}}
	ml := &mockNotifMailer{}

	svc := NewNotificationService(repo, wc, ml, "http://localhost:8080")
	err := svc.SendWeatherUpdates(context.Background(), "hourly")

	require.NoError(t, err)
	assert.Equal(t, []string{"a@example.com", "b@example.com"}, ml.calls)
}

func TestSendWeatherUpdates_WeatherError_ContinuesOtherSubs(t *testing.T) {
	repo := &mockNotifRepo{subs: []model.Subscription{
		{Email: "a@example.com", City: "BadCity", Token: "tok1"},
		{Email: "b@example.com", City: "Kyiv", Token: "tok2"},
	}}

	callCount := 0
	// Use a per-city mock via a simple function-based client.
	wc := &funcWeatherClient{fn: func(city string) (model.Weather, error) {
		callCount++
		if city == "BadCity" {
			return model.Weather{}, assert.AnError
		}
		return model.Weather{Temperature: 20}, nil
	}}
	ml := &mockNotifMailer{}

	svc := NewNotificationService(repo, wc, ml, "http://localhost:8080")
	err := svc.SendWeatherUpdates(context.Background(), "hourly")

	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
	assert.Equal(t, []string{"b@example.com"}, ml.calls)
}

func TestSendWeatherUpdates_RepoError(t *testing.T) {
	repo := &mockNotifRepo{findErr: assert.AnError}
	svc := NewNotificationService(repo, &mockWeatherClient{}, &mockNotifMailer{}, "http://localhost:8080")

	err := svc.SendWeatherUpdates(context.Background(), "daily")
	assert.ErrorContains(t, err, "load daily subscriptions")
}

type funcWeatherClient struct {
	fn func(city string) (model.Weather, error)
}

func (f *funcWeatherClient) GetWeather(_ context.Context, city string) (model.Weather, error) {
	return f.fn(city)
}
