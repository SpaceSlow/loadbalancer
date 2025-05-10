package ratelimiter

import (
	"sync"

	"github.com/SpaceSlow/loadbalancer/config"
)

type Bucket struct {
	mu sync.Mutex

	counter  float64
	capacity float64

	refillRPS float64
}

func NewBucket(cfg *config.BucketConfig) *Bucket {
	return &Bucket{
		counter:   cfg.Capacity,
		capacity:  cfg.Capacity,
		refillRPS: cfg.RefillRPS,
	}
}

func (b *Bucket) TakeToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.counter < 1 {
		return false
	}

	b.counter--
	return true
}

func (b *Bucket) RefillTokens(tokensNum float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.counter = min(b.counter+tokensNum, b.capacity)
}
