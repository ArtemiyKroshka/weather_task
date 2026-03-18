### Weather Task

A weather subscription API that allows users to:

- Get current weather for any city.
- Subscribe an email address to receive hourly or daily weather updates.
- Confirm and unsubscribe via tokens sent by email.

## Prerequisites

- Docker Engine (v20.10+)
- Docker Compose v2+
- Make
- Go 1.24+ (for local development)

---

## Getting started

### 1. Initialize environment file

```
make init-env
```

Open `.env` and set the required variables:

- `WEATHER_API_KEY` — your [WeatherAPI.com](https://www.weatherapi.com/) key
- `SMTP_USER` — SMTP username (e.g. Gmail address)
- `SMTP_PASSWORD` — SMTP password or app password

### 2. Start with Docker Compose

```
make up
```

This builds the Go image and starts the API together with PostgreSQL.

- API: `http://localhost:8080`
- Interactive docs: `http://localhost:8080/docs`

### 3. Stop

```
make down
```

---

## API Endpoints

Interactive Swagger UI is available at `http://localhost:8080/docs`.

| Method   | Path                                        | Description                                      |
|----------|---------------------------------------------|--------------------------------------------------|
| `GET`    | `/api/weather?city={city}`                  | Get current weather (temperature, humidity, ...) |
| `POST`   | `/api/subscriptions`                        | Subscribe an email for updates                   |
| `POST`   | `/api/subscriptions/confirm/{token}`        | Confirm email subscription                       |
| `DELETE` | `/api/subscriptions/{token}`                | Unsubscribe from updates                         |

### Subscribe (form-encoded body)

```
POST /api/subscriptions
Content-Type: application/x-www-form-urlencoded

email=user@example.com&city=Kyiv&frequency=daily
```

---

## Development

```bash
# Run tests
make test

# Build binary
make build
```
