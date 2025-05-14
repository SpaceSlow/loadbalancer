package ratelimiter

import (
	"sync"
)

type Bucket struct {
	mu sync.Mutex

	remainTokens float64
	capacity     float64

	refillRPS float64
}

func (b *Bucket) TakeToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.remainTokens < 1 {
		return false
	}

	b.remainTokens--
	return true
}

func (b *Bucket) RefillTokens(tokensNum float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.remainTokens = min(b.remainTokens+tokensNum, b.capacity)
}
