package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"weather_task/internal/service"

	"github.com/go-co-op/gocron/v2"
)

// Start registers all scheduled jobs and starts the scheduler.
// Returns the scheduler so the caller can shut it down gracefully.
func Start(svc service.NotificationService) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	_, err = s.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			if err := svc.SendWeatherUpdates(ctx, "hourly"); err != nil {
				slog.Error("hourly weather updates failed", "err", err)
			}
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("register hourly job: %w", err)
	}

	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(9, 0, 0))),
		gocron.NewTask(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			if err := svc.SendWeatherUpdates(ctx, "daily"); err != nil {
				slog.Error("daily weather updates failed", "err", err)
			}
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("register daily job: %w", err)
	}

	s.Start()
	return s, nil
}
