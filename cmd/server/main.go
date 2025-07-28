package main

import (
	"log"
	"os"

	"Rate-Limiter/internal/config"
	"Rate-Limiter/internal/constants"
	"Rate-Limiter/internal/limiter"
	"Rate-Limiter/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Default to local config if no environment is specified
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = constants.EnvLocal
	}

	cfg, err := config.GetConfig(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize default user limits if none exist
	if cfg.RateLimits.UserLimits == nil {
		cfg.RateLimits.UserLimits = make(map[string]config.UserRateLimit)
	}

	manager := limiter.NewManager(cfg)

	// Set Gin mode based on environment
	if env != constants.EnvLocal {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.RateLimit(manager))

	// Add your API routes here
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
