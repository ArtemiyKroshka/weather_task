package server

import (
	"fmt"
	"net/http"

	weatherapi "weather_task/api"
	"weather_task/internal/handler"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Weather API – Swagger UI</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
  window.onload = function() {
    SwaggerUIBundle({
      url: "/api/openapi.yaml",
      dom_id: "#swagger-ui",
      presets: [SwaggerUIBundle.presets.apis],
      layout: "BaseLayout",
      deepLinking: true
    });
  };
</script>
</body>
</html>`

// New registers all HTTP routes and returns a ready-to-use ServeMux.
func New(weatherH *handler.WeatherHandler, subH *handler.SubscriptionHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// API routes (RESTful)
	mux.HandleFunc("GET /api/weather", weatherH.GetWeather)
	mux.HandleFunc("POST /api/subscriptions", subH.Subscribe)
	mux.HandleFunc("POST /api/subscriptions/confirm/{token}", subH.Confirm)
	mux.HandleFunc("DELETE /api/subscriptions/{token}", subH.Unsubscribe)

	// Interactive API documentation at /docs
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, swaggerUIHTML)
	})

	// OpenAPI spec for Swagger UI
	mux.HandleFunc("GET /api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(weatherapi.OpenAPISpec)
	})

	return mux
}
