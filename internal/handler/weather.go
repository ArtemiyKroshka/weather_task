package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"weather_task/internal/model"
	"weather_task/pkg/weather"
)

// weatherGetter is the subset of weather.Client used by WeatherHandler.
// Accepting an interface makes the handler easy to test without a real HTTP server.
type weatherGetter interface {
	GetWeather(ctx context.Context, city string) (model.Weather, error)
}

// WeatherHandler handles GET /api/weather requests.
type WeatherHandler struct {
	client weatherGetter
}

// NewWeatherHandler creates a WeatherHandler with the given weather client.
func NewWeatherHandler(client weatherGetter) *WeatherHandler {
	return &WeatherHandler{client: client}
}

// GetWeather godoc
//
//	@Summary		Get current weather for a city
//	@Tags			weather
//	@Produce		json
//	@Param			city	query		string			true	"City name"
//	@Success		200		{object}	model.Weather
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/weather [get]
func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		writeError(w, http.StatusBadRequest, "city is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := h.client.GetWeather(ctx, city)
	if err != nil {
		switch {
		case errors.Is(err, weather.ErrCityNotFound):
			writeError(w, http.StatusNotFound, "city not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}
