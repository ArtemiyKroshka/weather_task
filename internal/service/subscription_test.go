package service

import (
	"context"
	"testing"

	"weather_task/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- mock repository ---

type mockRepo struct {
	sub          *model.Subscription
	findByEmail  error
	findByToken  error
	createErr    error
	updateErr    error
	deleteErr    error
	updated      *model.Subscription
	deleted      *model.Subscription
}

func (m *mockRepo) Create(ctx context.Context, sub *model.Subscription) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.sub = sub
	return nil
}

func (m *mockRepo) FindByEmail(_ context.Context, _ string) (*model.Subscription, error) {
	if m.findByEmail != nil {
		return nil, m.findByEmail
	}
	return m.sub, nil
}

func (m *mockRepo) FindByToken(_ context.Context, _ string) (*model.Subscription, error) {
	if m.findByToken != nil {
		return nil, m.findByToken
	}
	return m.sub, nil
}

func (m *mockRepo) Update(_ context.Context, sub *model.Subscription) error {
	m.updated = sub
	return m.updateErr
}

func (m *mockRepo) Delete(_ context.Context, sub *model.Subscription) error {
	m.deleted = sub
	return m.deleteErr
}

func (m *mockRepo) FindConfirmedByFrequency(_ context.Context, _ string) ([]model.Subscription, error) {
	return nil, nil
}

// --- mock mailer ---

type mockMailer struct {
	confirmCalled bool
	confirmErr    error
}

func (m *mockMailer) SendConfirmation(_, _ string) error {
	m.confirmCalled = true
	return m.confirmErr
}

func (m *mockMailer) SendWeatherUpdate(_, _ string, _ model.Weather, _ string) error {
	return nil
}

// --- tests ---

func TestSubscribe_Success(t *testing.T) {
	repo := &mockRepo{findByEmail: gorm.ErrRecordNotFound}
	ml := &mockMailer{}
	svc := NewSubscriptionService(repo, ml, "http://localhost:8080")

	err := svc.Subscribe(context.Background(), "user@example.com", "Kyiv", "daily")

	require.NoError(t, err)
	require.NotNil(t, repo.sub)
	assert.Equal(t, "user@example.com", repo.sub.Email)
	assert.Equal(t, "Kyiv", repo.sub.City)
	assert.Equal(t, "daily", repo.sub.Frequency)
	assert.NotEmpty(t, repo.sub.Token)
	assert.True(t, ml.confirmCalled)
}

func TestSubscribe_AlreadySubscribed(t *testing.T) {
	repo := &mockRepo{sub: &model.Subscription{Email: "user@example.com"}}
	svc := NewSubscriptionService(repo, &mockMailer{}, "http://localhost:8080")

	err := svc.Subscribe(context.Background(), "user@example.com", "Kyiv", "daily")

	assert.ErrorIs(t, err, ErrAlreadySubscribed)
}

func TestSubscribe_MailerError(t *testing.T) {
	repo := &mockRepo{findByEmail: gorm.ErrRecordNotFound}
	ml := &mockMailer{confirmErr: assert.AnError}
	svc := NewSubscriptionService(repo, ml, "http://localhost:8080")

	err := svc.Subscribe(context.Background(), "user@example.com", "Kyiv", "daily")

	assert.ErrorContains(t, err, "send confirmation email")
}

func TestConfirm_Success(t *testing.T) {
	sub := &model.Subscription{Token: "old-token", Confirmed: false}
	repo := &mockRepo{sub: sub}
	svc := NewSubscriptionService(repo, &mockMailer{}, "http://localhost:8080")

	err := svc.Confirm(context.Background(), "old-token")

	require.NoError(t, err)
	require.NotNil(t, repo.updated)
	assert.True(t, repo.updated.Confirmed)
	assert.NotEqual(t, "old-token", repo.updated.Token)
}

func TestConfirm_TokenNotFound(t *testing.T) {
	repo := &mockRepo{findByToken: gorm.ErrRecordNotFound}
	svc := NewSubscriptionService(repo, &mockMailer{}, "http://localhost:8080")

	err := svc.Confirm(context.Background(), "invalid-token")

	assert.ErrorIs(t, err, ErrTokenNotFound)
}

func TestUnsubscribe_Success(t *testing.T) {
	sub := &model.Subscription{Email: "user@example.com", Token: "some-token"}
	repo := &mockRepo{sub: sub}
	svc := NewSubscriptionService(repo, &mockMailer{}, "http://localhost:8080")

	err := svc.Unsubscribe(context.Background(), "some-token")

	require.NoError(t, err)
	assert.Equal(t, sub, repo.deleted)
}

func TestUnsubscribe_TokenNotFound(t *testing.T) {
	repo := &mockRepo{findByToken: gorm.ErrRecordNotFound}
	svc := NewSubscriptionService(repo, &mockMailer{}, "http://localhost:8080")

	err := svc.Unsubscribe(context.Background(), "invalid-token")

	assert.ErrorIs(t, err, ErrTokenNotFound)
}
