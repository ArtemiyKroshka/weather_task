package db

import (
	"fmt"
	"log"
	"os"
	"weather_task/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "db"),
		getEnv("DB_USER", "weather_user"),
		getEnv("DB_PASSWORD", "weather_pass"),
		getEnv("DB_NAME", "weather_db"),
		getEnv("DB_PORT", "5432"),
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Migrate() {
	if err := DB.AutoMigrate(&models.Subscription{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func GetUser(email string) (models.Subscription, error) {
	var user models.Subscription
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.Subscription{}, err
	}
	return user, nil
}
