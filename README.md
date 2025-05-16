### Weather Task

A simple weather API application that allows users to:

- Get the current weather for a chosen city.
- Subscribe an email address to receive weather updates (hourly or daily).
- Confirm and unsubscribe via email tokens.

## Prerequisites

- Docker Engine (v20.10+)
- Docker Compose v2+ (or `docker-compose` for v1)
- Make

---

## Steps to start the app:

> **Note:** All commands below use the `docker compose` syntax. If your Docker installation requires the standalone `docker-compose`, replace `docker compose` with `docker-compose` accordingly.

### 1. Initialize environment file

```
make init-env
```

This will create a `.env` file from `.env.example`. Open `.env` and set the following:

- `WEATHER_API_KEY` — your WeatherAPI.com API key
- `SMTP_USER` — SMTP username
- `SMTP_PASSWORD` — SMTP password

2. Run 'up' command

```
make up
```

This will:

1. Build the Go API image
2. Start application and PostgreSQL

After startup:

- API is available at `http://localhost:8080`

## API Endpoints

All endpoints are prefixed with `/api`:

| Method | Path                       | Description                                        |
| ------ | -------------------------- | -------------------------------------------------- |
| GET    | `/api/weather?city={city}` | Get current weather (temperature, humidity, text). |
| POST   | `/api/subscribe`           | Subscribe an email for updates.                    |
| GET    | `/api/confirm/{token}`     | Confirm email subscription.                        |
| GET    | `/api/unsubscribe/{token}` | Unsubscribe from updates.                          |
