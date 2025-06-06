package config

import "time"

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
