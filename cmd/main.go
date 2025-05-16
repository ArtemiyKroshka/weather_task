package main

import (
	"log"
	"weather_task/internal/api"
	"weather_task/internal/db"
	"weather_task/internal/scheduler"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	db.Migrate()

	scheduler.Init()

	routes := gin.Default()
	api.RegisterRoutes(routes)

	if err := routes.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
