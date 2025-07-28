package limiter

import (
	"time"
)

// Priority levels for requests
type Priority int

const (
	Critical Priority = iota
	PostRequest
	GetRequest
	TestMode
)

// Request represents an incoming API request for rate limiting
type Request struct {
	UserID    string
	Endpoint  string
	Method    string
	Priority  Priority
	RequestID string
	Timestamp time.Time
}

// Decision represents the result of a rate limit check
type Decision struct {
	Allowed    bool
	StatusCode int
	Reason     string
	RetryAfter time.Duration
	Headers    map[string]string
}

// Metrics for tracking limiter statistics
type Metrics struct {
	RequestsTotal  int64
	RequestsDenied int64
	AverageLatency time.Duration
}

// RateLimitConfig holds per-user or per-endpoint rate limit settings
type RateLimitConfig struct {
	RequestsPerSecond int
	// Add other rate limit parameters as needed
}
