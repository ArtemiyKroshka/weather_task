package service

import (
	"context"
	"errors"
	"fmt"

	"weather_task/internal/model"
	"weather_task/internal/repository"
	"weather_task/pkg/mailer"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrAlreadySubscribed is returned when an email is already registered.
var ErrAlreadySubscribed = errors.New("email already subscribed")

// ErrTokenNotFound is returned when no subscription matches the given token.
var ErrTokenNotFound = errors.New("token not found")

// SubscriptionService handles subscription lifecycle operations.
type SubscriptionService interface {
	Subscribe(ctx context.Context, email, city, frequency string) error
	Confirm(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
}

type subscriptionService struct {
	repo    repository.SubscriptionRepo
	mailer  mailer.Mailer
	baseURL string
}

// NewSubscriptionService wires up the subscription service with its dependencies.
func NewSubscriptionService(
	repo repository.SubscriptionRepo,
	m mailer.Mailer,
	baseURL string,
) SubscriptionService {
	return &subscriptionService{
		repo:    repo,
		mailer:  m,
		baseURL: baseURL,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, email, city, frequency string) error {
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return ErrAlreadySubscribed
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check existing subscription: %w", err)
	}

	sub := &model.Subscription{
		Email:     email,
		City:      city,
		Frequency: frequency,
		Token:     uuid.NewString(),
	}
	if err := s.repo.Create(ctx, sub); err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}

	confirmURL := fmt.Sprintf("%s/api/subscriptions/confirm/%s", s.baseURL, sub.Token)
	if err := s.mailer.SendConfirmation(email, confirmURL); err != nil {
		return fmt.Errorf("send confirmation email: %w", err)
	}

	return nil
}

func (s *subscriptionService) Confirm(ctx context.Context, token string) error {
	sub, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTokenNotFound
		}
		return fmt.Errorf("find subscription: %w", err)
	}

	sub.Confirmed = true
	sub.Token = uuid.NewString()
	if err := s.repo.Update(ctx, sub); err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}

	return nil
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, token string) error {
	sub, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTokenNotFound
		}
		return fmt.Errorf("find subscription: %w", err)
	}

	if err := s.repo.Delete(ctx, sub); err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	return nil
}
