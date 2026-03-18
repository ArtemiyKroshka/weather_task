package db

import (
	"fmt"
	"weather_task/internal/config"
	"weather_task/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New opens a PostgreSQL connection using the provided configuration.
func New(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return db, nil
}

// Migrate runs auto-migration for all application models.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Subscription{}); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}
	return nil
}
