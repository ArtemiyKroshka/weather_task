package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather_task/internal/model"
	"weather_task/pkg/weather"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWeatherGetter struct {
	result model.Weather
	err    error
}

func (m *mockWeatherGetter) GetWeather(_ context.Context, _ string) (model.Weather, error) {
	return m.result, m.err
}

func TestGetWeather_Success(t *testing.T) {
	mock := &mockWeatherGetter{result: model.Weather{Temperature: 22.5, Humidity: 65, Description: "Sunny"}}
	h := NewWeatherHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/weather?city=Kyiv", nil)
	w := httptest.NewRecorder()
	h.GetWeather(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got model.Weather
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, 22.5, got.Temperature)
	assert.Equal(t, float64(65), got.Humidity)
	assert.Equal(t, "Sunny", got.Description)
}

func TestGetWeather_MissingCity(t *testing.T) {
	h := NewWeatherHandler(&mockWeatherGetter{})

	req := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	w := httptest.NewRecorder()
	h.GetWeather(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assertErrorBody(t, w, "city is required")
}

func TestGetWeather_CityNotFound(t *testing.T) {
	mock := &mockWeatherGetter{err: weather.ErrCityNotFound}
	h := NewWeatherHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/weather?city=NoSuchPlace", nil)
	w := httptest.NewRecorder()
	h.GetWeather(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assertErrorBody(t, w, "city not found")
}

func TestGetWeather_InternalError(t *testing.T) {
	mock := &mockWeatherGetter{err: assert.AnError}
	h := NewWeatherHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/weather?city=Kyiv", nil)
	w := httptest.NewRecorder()
	h.GetWeather(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// assertErrorBody checks that the response JSON has the expected "error" field.
func assertErrorBody(t *testing.T, w *httptest.ResponseRecorder, contains string) {
	t.Helper()
	var body map[string]string
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Contains(t, body["error"], contains)
}
