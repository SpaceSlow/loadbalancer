package config

import (
	"fmt"
	"os"

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

type BackendConfig struct {
	URL string `yaml:"url"`
}

type BalancerConfig struct {
	Port     int             `yaml:"port"`
	Strategy Strategy        `yaml:"strategy"`
	Backends []BackendConfig `yaml:"backends"`
}

type Config struct {
	Balancer BalancerConfig `yaml:"load_balancer"`
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
