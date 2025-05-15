package models

import "time"

type SubscriptionRequest struct {
	Email     string `form:"email" binding:"required,email"`
	City      string `form:"city" binding:"required"`
	Frequency string `form:"frequency" binding:"required,oneof=hourly daily"`
}

type Subscription struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	City      string `gorm:"not null"`
	Frequency string `gorm:"not null"` // "hourly" or "daily"
	Confirmed bool   `gorm:"default:false"`
	Token     string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
