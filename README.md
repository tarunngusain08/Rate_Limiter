# Rate Limiter

A high-performance, extensible rate limiting service written in Go. This service provides distributed and in-memory rate limiting for APIs, with support for custom strategies, metrics, and flexible configuration.

<p align="center">
  <img src="https://github.com/user-attachments/assets/bd635fc0-39f8-4bf2-9370-985e66c65f8a" width="600"/>
</p>

---

## Features

- **Distributed Rate Limiting**: Uses Redis for distributed state management.
- **In-Memory Rate Limiting**: Fast, local fallback for high performance.
- **Multiple Limiting Strategies**:
  - Request rate limiting (RPS, burst, concurrent)
  - Worker utilization limiting
  - Fleet usage limiting
- **Configurable per-user limits**
- **Metrics aggregation**
- **REST API with Gin**
- **Pluggable middleware**

---

## Architecture

- **cmd/server**: Entry point for the HTTP server.
- **internal/limiter**: Core rate limiting logic and manager.
- **internal/config**: Configuration loading and validation.
- **internal/constants**: Shared constants.
- **internal/storage**: Redis and in-memory storage backends.
- **internal/middleware**: Gin middleware for rate limiting.

---

## Configuration

Configuration is loaded from environment variables or `.env` files. Example (`internal/config/config.local.env`):

```env
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

RATE_LIMIT_DEFAULT_RPS=30
RATE_LIMIT_DEFAULT_BURST=30
RATE_LIMIT_DEFAULT_CONCURRENT=10

FLEET_USAGE_CRITICAL_RESERVATION=0.2
FLEET_USAGE_MONITORING_WINDOW=60s

WORKER_THRESHOLD_EMERGENCY=0.95
WORKER_THRESHOLD_HIGH=0.85
WORKER_THRESHOLD_MEDIUM=0.75
WORKER_THRESHOLD_LOW=0.60
```

You can override these by setting environment variables.

---

## Usage

### Running the Server

```bash
go run cmd/server/main.go
```

The server listens on port `8080` by default.

### API Endpoints

- `GET /health`  
  Returns service health status.

- `GET /test`  
  Triggers a worker release (for testing).

All other endpoints are protected by the rate limiting middleware.

---

## Summary of Rate Limiting Results

When sending a burst of requests to the server (e.g., using `curl` in a loop), you may observe the following HTTP status codes:

- **200**: Request allowed (within rate limits).
- **503**: Service unavailable, typically due to worker or fleet utilization limits being exceeded.
- **429**: Too many requests, standard rate limiting response when request rate or concurrency limits are hit.

Example output for 50 rapid requests:

```
200
200
...
503
503
...
429
429
...
```

This demonstrates that the rate limiter is enforcing multiple layers of protection, including request rate, worker utilization, and fleet usage.

---

## Extending
<img width="923" height="634" alt="Screenshot 2025-08-04 at 7 16 24 AM" src="https://github.com/user-attachments/assets/c517ce8a-3aa4-41eb-af67-81929f31965a" />
<img width="672" height="776" alt="Screenshot 2025-08-04 at 7 16 09 AM" src="https://github.com/user-attachments/assets/9d570e78-d0d1-4e35-8684-2c53d3314b8d" />

---

## Extending

- Add new limiter strategies by implementing the `RateLimiter` interface in `internal/interfaces`.
- Add new storage backends in `internal/storage`.

---

## Development

### Requirements

- Go 1.18+
- Redis (for distributed mode)
