package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

const (
	RoundRobinStrategy = "round-robin"
)

type Strategy string

func (s Strategy) IsValid() bool {
	switch string(s) {
	case RoundRobinStrategy:
		return true
	}
	return false
}

type HealthCheckConfig struct {
	Interval time.Duration `yaml:"interval"`
	Path     string        `yaml:"path"`
}

type BackendConfig struct {
	URL         string            `yaml:"url"`
	HealthCheck HealthCheckConfig `yaml:"healthcheck"`
}

type BalancerConfig struct {
	Port     int             `yaml:"port"`
	Strategy Strategy        `yaml:"strategy"`
	Backends []BackendConfig `yaml:"backends"`
}

type BucketConfig struct {
	Capacity  float64 `yaml:"capacity"`
	RefillRPS float64 `yaml:"refill_rps"`
}

type RateLimiterConfig struct {
	DefaultBucket BucketConfig `yaml:"default_bucket"`
}

type Config struct {
	Balancer    BalancerConfig    `yaml:"load_balancer"`
	RateLimiter RateLimiterConfig `yaml:"rate_limiter"`
}

func LoadConfig(filename string) (*Config, error) {
	configData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("[config] read file error: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(configData, &cfg); err != nil {
		return nil, fmt.Errorf("[config] yaml unmarshal error: %w", err)
	}

	return &cfg, nil
}
