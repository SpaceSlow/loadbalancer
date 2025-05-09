package balancer

import (
	"errors"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/config"
)

type Balancer interface {
	Handler(w http.ResponseWriter, r *http.Request)
}

func NewBalancer(cfg *config.BalancerConfig) (Balancer, error) {
	if cfg.Port < 0 || cfg.Port > 65535 {
		return nil, errors.New("[config] incorrect port number (port must be in bounds 0-65535)")
	}

	if !cfg.Strategy.IsValid() {
		return nil, errors.New("[config] unknown load balancer strategy (check README)")
	}

	if len(cfg.Backends) == 0 {
		return nil, errors.New("[config] no specified backend urls")
	}

	switch cfg.Strategy {
	case config.RoundRobinStrategy:
		return nil, errors.New("not implemented")
	}

	return nil, errors.New("unexpected error")
}
