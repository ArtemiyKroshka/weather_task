package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"weather_task/internal/model"
)

// cityNotFound is the WeatherAPI error code for an unknown city.
// See: https://www.weatherapi.com/docs/#intro-error-codes
const cityNotFound = 1006

// ErrCityNotFound is returned when the requested city does not exist in the weather API.
var ErrCityNotFound = errors.New("city not found")

// ErrInvalidRequest is returned for any other non-OK response from the weather API.
var ErrInvalidRequest = errors.New("invalid request")

// Client calls the WeatherAPI to retrieve current weather conditions.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a Client. httpClient may be nil to use http.DefaultClient.
func NewClient(apiKey string, httpClient *http.Client) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("weather API key is required")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    "https://api.weatherapi.com/v1/current.json",
		apiKey:     apiKey,
		httpClient: httpClient,
	}, nil
}

// GetWeather fetches current weather for the given city.
// The caller is responsible for setting an appropriate deadline on ctx.
func (c *Client) GetWeather(ctx context.Context, city string) (model.Weather, error) {
	reqURL := fmt.Sprintf("%s?key=%s&q=%s&lang=uk", c.baseURL, c.apiKey, city)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return model.Weather{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.Weather{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return model.Weather{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
		}
		if errResp.Error.Code == cityNotFound {
			return model.Weather{}, ErrCityNotFound
		}
		return model.Weather{}, ErrInvalidRequest
	}

	var success successResponse
	if err := json.NewDecoder(resp.Body).Decode(&success); err != nil {
		return model.Weather{}, fmt.Errorf("decode response: %w", err)
	}

	return model.Weather{
		Temperature: success.Current.TempC,
		Humidity:    float64(success.Current.Humidity),
		Description: success.Current.Condition.Text,
	}, nil
}

type errorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type successResponse struct {
	Current struct {
		TempC    float64 `json:"temp_c"`
		Humidity int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}
