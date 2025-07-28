package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Redis struct {
		Addr     string `envconfig:"REDIS_ADDR"`
		Password string `envconfig:"REDIS_PASSWORD"`
		DB       int    `envconfig:"REDIS_DB"`
	}

	Worker struct {
		Count int `envconfig:"WORKER_COUNT"`
	}

	RateLimits struct {
		DefaultRPS        int `envconfig:"RATE_LIMIT_DEFAULT_RPS"`
		DefaultBurst      int `envconfig:"RATE_LIMIT_DEFAULT_BURST"`
		DefaultConcurrent int `envconfig:"RATE_LIMIT_DEFAULT_CONCURRENT"`
		UserLimits        map[string]UserRateLimit
	}

	FleetUsage struct {
		CriticalReservation float64       `envconfig:"FLEET_USAGE_CRITICAL_RESERVATION"`
		MonitoringWindow    time.Duration `envconfig:"FLEET_USAGE_MONITORING_WINDOW"`
	}

	WorkerThresholds struct {
		Emergency float64 `envconfig:"WORKER_THRESHOLD_EMERGENCY"`
		High      float64 `envconfig:"WORKER_THRESHOLD_HIGH"`
		Medium    float64 `envconfig:"WORKER_THRESHOLD_MEDIUM"`
		Low       float64 `envconfig:"WORKER_THRESHOLD_LOW"`
	}
}

type UserRateLimit struct {
	RequestsPerSecond int
	Burst             int
	Concurrent        int
}

func (c *Config) Validate() error {
	if c.Redis.Addr == "" {
		return fmt.Errorf("redis address is required")
	}

	if c.RateLimits.DefaultRPS <= 0 {
		return fmt.Errorf("default RPS must be greater than 0")
	}

	if c.FleetUsage.CriticalReservation <= 0 || c.FleetUsage.CriticalReservation >= 1 {
		return fmt.Errorf("critical reservation must be between 0 and 1")
	}

	if c.WorkerThresholds.Emergency <= c.WorkerThresholds.High {
		return fmt.Errorf("emergency threshold must be greater than high threshold")
	}

	return nil
}

func GetConfig(env string) (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		if err := loadConfig(instance, env); err != nil {
			panic(err)
		}
	})
	return instance, nil
}

func loadConfig(cfg *Config, env string) error {
	// Load .env file
	if err := godotenv.Load(fmt.Sprintf("internal/config/config.%s.env", env)); err != nil {
		return err
	}

	// Parse environment variables into struct
	if err := envconfig.Process("", cfg); err != nil {
		return err
	}

	return nil
}
