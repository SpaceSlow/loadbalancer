package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Balancer                  BalancerConfig    `yaml:"load_balancer"`
	RateLimiter               RateLimiterConfig `yaml:"rate_limiter"`
	DB                        DBConfig          `yaml:"db"`
	MaxTimeoutShutdown        time.Duration
	HTTPServerTimeoutShutdown time.Duration
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
	cfg.MaxTimeoutShutdown = 15 * time.Second
	cfg.HTTPServerTimeoutShutdown = 5 * time.Second

	return &cfg, nil
}
