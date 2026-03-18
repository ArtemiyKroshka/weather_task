package model

import "time"

// Subscription represents a user's weather update subscription stored in the database.
type Subscription struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	City      string    `gorm:"not null"`
	Frequency string    `gorm:"not null"` // "hourly" or "daily"
	Confirmed bool      `gorm:"default:false"`
	Token     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
