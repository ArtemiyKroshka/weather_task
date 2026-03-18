package repository

import (
	"context"
	"weather_task/internal/model"

	"gorm.io/gorm"
)

// SubscriptionRepo defines the data-access contract for subscriptions.
type SubscriptionRepo interface {
	Create(ctx context.Context, sub *model.Subscription) error
	FindByEmail(ctx context.Context, email string) (*model.Subscription, error)
	FindByToken(ctx context.Context, token string) (*model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, sub *model.Subscription) error
	FindConfirmedByFrequency(ctx context.Context, freq string) ([]model.Subscription, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

// NewSubscriptionRepo creates a GORM-backed SubscriptionRepo.
func NewSubscriptionRepo(db *gorm.DB) SubscriptionRepo {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *model.Subscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *subscriptionRepo) FindByEmail(ctx context.Context, email string) (*model.Subscription, error) {
	var sub model.Subscription
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepo) FindByToken(ctx context.Context, token string) (*model.Subscription, error) {
	var sub model.Subscription
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepo) Update(ctx context.Context, sub *model.Subscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *subscriptionRepo) Delete(ctx context.Context, sub *model.Subscription) error {
	return r.db.WithContext(ctx).Delete(sub).Error
}

func (r *subscriptionRepo) FindConfirmedByFrequency(ctx context.Context, freq string) ([]model.Subscription, error) {
	var subs []model.Subscription
	if err := r.db.WithContext(ctx).Where("frequency = ? AND confirmed = ?", freq, true).Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}
