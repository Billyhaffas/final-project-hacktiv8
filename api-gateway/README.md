# api-gateway

HTTP entry point for the Climate Action backend. Accepts REST requests from clients, validates JWTs, and fans out to the appropriate downstream microservice.

## Architecture

```
Client (HTTP/REST)
    │
    ▼
api-gateway  :8080
    ├── HTTP reverse proxy ──▶ auth-service      :8081
    └── gRPC client        ──▶ count-emission-service :50051
                                  (GetUserPreferences / SetUserPreferences /
                                   CreateUserEmission / GetUserDailyEmission /
                                   GetUserMonthlyEmission / GetUserYearlyEmission)
```

Stubs (return 501 until downstream services are built):
- `GET /api/v1/emissions/alert` — waiting on notification-service
- `GET /api/v1/emissions/convert` — waiting on convert-emission-service

## Project Structure

```
api-gateway/
├── cmd/
│   └── main.go                 # Entry point — wires dependencies, starts Echo
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go     # Reverse-proxies all /auth/* to auth-service
│   │   ├── emission_handler.go # gRPC calls to count-emission-service
│   │   └── preference_handler.go # gRPC calls for user preferences
│   ├── middleware/
│   │   └── jwt.go              # Bearer token validation, injects user_id + email
│   └── router/
│       └── router.go           # Route definitions
├── helper/
│   └── response.go             # Standard { success, data } / { success, error } envelope
├── proto/
│   └── emission/               # Generated gRPC stubs (copied from count-emission-service)
└── .env.example
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `GATEWAY_PORT` | `8080` | Port this service listens on |
| `AUTH_SERVICE_URL` | — | Base URL of auth-service (e.g. `http://localhost:8081`) |
| `COUNT_EMISSION_SERVICE_ADDR` | — | gRPC address of count-emission-service (e.g. `localhost:50051`) |
| `NOTIFICATION_SERVICE_ADDR` | — | gRPC address of notification-service (future) |
| `CONVERT_EMISSION_SERVICE_ADDR` | — | gRPC address of convert-emission-service (future) |
| `JWT_SECRET` | — | Secret used to validate HS256 JWTs issued by auth-service |

Copy `.env.example` to `.env` and fill in real values before running.

## Running Locally

```bash
cd api-gateway
cp .env.example .env   # then edit .env
go run cmd/main.go
```

Requires `auth-service` and `count-emission-service` to be running.

## API Endpoints

### Auth (proxied to auth-service, no JWT required)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login → returns JWT + refresh token |
| `POST` | `/api/v1/auth/refresh` | Exchange refresh token for new JWT |
| `POST` | `/api/v1/auth/forgot-password` | Request password reset |
| `POST` | `/api/v1/auth/reset-password` | Complete password reset |
| `POST` | `/api/v1/auth/logout` | **JWT required** — revoke refresh token |

### Emissions (JWT required)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/emissions` | Log a commute trip |
| `GET` | `/api/v1/emissions/today` | Today's total CO₂ |
| `GET` | `/api/v1/emissions/report` | Monthly breakdown |
| `GET` | `/api/v1/emissions/alert` | Daily limit alert *(stub — 501)* |
| `GET` | `/api/v1/emissions/convert` | Convert kg CO₂ to IDR *(stub — 501)* |

**POST `/api/v1/emissions` request body:**
```json
{
  "vehicle_type": "Car-Size-Medium",
  "fuel_type": "Petrol",
  "distance_km": 12.5
}
```
`fuel_type` is optional (defaults to `"Unknown"`). `distance_km` must be > 0.
`vehicle_type` values: `Car-*`, `Motorbike-*`, `Bus-*`, `Train-*`, `Taxi-Local`, `Bicycle`, `Walk`.

### Preferences (JWT required)

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/v1/preferences` | Get user's emission preferences |
| `PUT` | `/api/v1/preferences` | Update preferences |

**PUT `/api/v1/preferences` request body:**
```json
{
  "country_code": "IDN",
  "custom_daily_limit_kg_co2": 5.0
}
```
Set `custom_daily_limit_kg_co2` to `0` to clear a custom limit (falls back to country default).

## Authentication

Protected routes require `Authorization: Bearer <token>` header. The middleware:
1. Parses the HS256 JWT using `JWT_SECRET`
2. Returns `401 TOKEN_EXPIRED` if the token is expired
3. Returns `401 TOKEN_INVALID` for any other validation failure
4. On success, sets `user_id` and `email` in the Echo context
5. Passes `user-id` as gRPC metadata to downstream services

## Response Format

All responses use a standard envelope:

```json
// Success
{ "success": true, "data": { ... } }

// Error
{ "success": false, "error": { "code": "VALIDATION_ERROR", "message": "..." } }
```

## Running Tests

```bash
cd api-gateway
go test ./...
```

Tests cover JWT middleware (6 cases) and all handler methods using hand-rolled gRPC client mocks — no network required.
