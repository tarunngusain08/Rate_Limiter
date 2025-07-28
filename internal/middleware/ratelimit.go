package middleware

import (
	"time"

	"Rate-Limiter/internal/limiter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RateLimit(manager *limiter.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &limiter.Request{
			UserID:    getUserID(c),
			Endpoint:  c.Request.URL.Path,
			Method:    c.Request.Method,
			Priority:  getPriority(c),
			RequestID: uuid.New().String(),
			Timestamp: time.Now(),
		}

		decision, err := manager.CheckAllLimiters(c.Request.Context(), req)
		if err != nil {
			c.Next() // Fail open
			return
		}

		for k, v := range decision.Headers {
			c.Header(k, v)
		}

		if !decision.Allowed {
			c.JSON(decision.StatusCode, gin.H{
				"error":       decision.Reason,
				"retry_after": decision.RetryAfter.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func getUserID(c *gin.Context) string {
	// Try JWT token first
	if token := c.GetHeader("Authorization"); token != "" {
		if id := extractUserIDFromToken(token); id != "" {
			return id
		}
	}

	// Try API key
	if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
		return apiKey
	}

	return c.ClientIP()
}

func getPriority(c *gin.Context) limiter.Priority {
	// Check request headers for priority
	if c.GetHeader("X-Critical") == "true" {
		return limiter.Critical
	}

	switch c.Request.Method {
	case "POST", "PUT", "DELETE":
		return limiter.PostRequest
	case "GET":
		return limiter.GetRequest
	default:
		return limiter.TestMode
	}
}

func extractUserIDFromToken(token string) string {
	// TODO: Implement JWT parsing
	return ""
}
