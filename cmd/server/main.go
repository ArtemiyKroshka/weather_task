package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"weather_task/internal/config"
	"weather_task/internal/db"
	"weather_task/internal/handler"
	"weather_task/internal/repository"
	"weather_task/internal/scheduler"
	"weather_task/internal/server"
	"weather_task/internal/service"
	"weather_task/pkg/mailer"
	"weather_task/pkg/weather"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	database, err := db.New(cfg)
	if err != nil {
		slog.Error("connect database", "err", err)
		os.Exit(1)
	}
	if err := db.Migrate(database); err != nil {
		slog.Error("migrate database", "err", err)
		os.Exit(1)
	}

	weatherClient, err := weather.NewClient(cfg.WeatherAPIKey, nil)
	if err != nil {
		slog.Error("create weather client", "err", err)
		os.Exit(1)
	}

	m := mailer.NewFromConfig(cfg)

	subRepo := repository.NewSubscriptionRepo(database)
	subSvc := service.NewSubscriptionService(subRepo, m, cfg.BaseURL)
	notifSvc := service.NewNotificationService(subRepo, weatherClient, m, cfg.BaseURL)

	sched, err := scheduler.Start(notifSvc)
	if err != nil {
		slog.Error("start scheduler", "err", err)
		os.Exit(1)
	}
	defer sched.Shutdown()

	weatherH := handler.NewWeatherHandler(weatherClient)
	subH := handler.NewSubscriptionHandler(subSvc)
	mux := server.New(weatherH, subH)

	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	}
}
