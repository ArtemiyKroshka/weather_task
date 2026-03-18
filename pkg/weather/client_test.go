package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestClient creates a Client that points to the given test server URL.
func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	return &Client{
		baseURL:    baseURL,
		apiKey:     "test-key",
		httpClient: http.DefaultClient,
	}
}

func TestNewClient_MissingAPIKey(t *testing.T) {
	_, err := NewClient("", nil)
	assert.ErrorContains(t, err, "weather API key is required")
}

func TestGetWeather_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.URL.Query().Get("key"))
		assert.Equal(t, "Kyiv", r.URL.Query().Get("q"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"current": {
				"temp_c": 18.5,
				"humidity": 70,
				"condition": {"text": "Partly cloudy"}
			}
		}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	got, err := client.GetWeather(context.Background(), "Kyiv")

	require.NoError(t, err)
	assert.Equal(t, 18.5, got.Temperature)
	assert.Equal(t, float64(70), got.Humidity)
	assert.Equal(t, "Partly cloudy", got.Description)
}

func TestGetWeather_CityNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":1006,"message":"No matching location found."}}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.GetWeather(context.Background(), "NonExistentXYZ")

	assert.ErrorIs(t, err, ErrCityNotFound)
}

func TestGetWeather_InvalidRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":1002,"message":"API key is invalid."}}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.GetWeather(context.Background(), "Kyiv")

	assert.ErrorIs(t, err, ErrInvalidRequest)
}

func TestGetWeather_UnexpectedStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.GetWeather(context.Background(), "Kyiv")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status 500")
}

func TestGetWeather_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not valid json`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)
	_, err := client.GetWeather(context.Background(), "Kyiv")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode response")
}
