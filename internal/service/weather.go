package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"weather_task/internal/models"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type rawErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type rawSuccessResponse struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

var (
	ErrMissingAPIKey  = errors.New("missing WEATHER_API_KEY")
	ErrInvalidRequest = errors.New("invalid request")
	ErrCityNotFound   = errors.New("city not found")
)

var CITY_NOT_FOUND = 1006

func NewClient() (*Client, error) {
	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		return nil, ErrMissingAPIKey
	}
	return &Client{
		baseURL:    "http://api.weatherapi.com/v1/current.json",
		apiKey:     key,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *Client) GetWeather(city string) (models.Weather, error) {
	reqURL := fmt.Sprintf("%s?key=%s&q=%s&lang=uk", c.baseURL, c.apiKey, city)

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return models.Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp rawErrorResponse
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errResp); decodeErr != nil {
			return models.Weather{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		switch errResp.Error.Code {
		case CITY_NOT_FOUND:
			return models.Weather{}, ErrCityNotFound
		default:
			return models.Weather{}, ErrInvalidRequest
		}
	}

	var success rawSuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&success); err != nil {
		return models.Weather{}, err
	}

	return models.Weather{
		Temperature: success.Current.TempC,
		Humidity:    float64(success.Current.Humidity),
		Description: success.Current.Condition.Text,
	}, nil
}
