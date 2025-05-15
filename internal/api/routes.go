package api

import (
	"weather_task/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")

	api.GET("/weather", handlers.GetWeatherHandler)
	api.POST("/subscribe", handlers.SubscriptionHandler)
	api.GET("/confirm/:token", handlers.ConfirmSubscriptionHandler)
	api.GET("/unsubscribe/:token", handlers.UnsubscribeHandler)
}
