package scheduler

import (
	"log"
	"time"
	"weather_task/internal/db"
	"weather_task/internal/models"
	"weather_task/pkg/email"

	"github.com/go-co-op/gocron"
)

func Init() {
	sched := gocron.NewScheduler(time.UTC)

	sched.Every(1).Hour().Do(SendPending, "hourly")
	sched.Every(1).Day().At("09:00").Do(SendPending, "daily")

	sched.StartAsync()
}

func SendPending(freq string) {
	var subs []models.Subscription

	if err := db.DB.Where("frequency = ? AND confirmed = ?", freq, true).Find(&subs).Error; err != nil {
		log.Printf("[Scheduler] Error loading %s subscriptions: %v", freq, err)
		return
	}

	for _, sub := range subs {
		if err := email.SendWeatherEmail(sub.Email); err != nil {
			log.Printf("[Scheduler] Failed to send email to %s: %v", sub.Email, err)
		}
	}
}
