package ratelimiter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/SpaceSlow/loadbalancer/config"
	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
)

type RateLimiter struct {
	cfg *config.RateLimiterConfig

	buckets sync.Map
}

func NewRateLimiter(ctx context.Context, cfg *config.RateLimiterConfig) *RateLimiter {
	limiter := &RateLimiter{
		cfg: cfg,
	}
	go limiter.refillBucketTokensJob(ctx)
	return limiter
}

func (rl *RateLimiter) AddBucket(clientID string, capacity, rps float64) {
	rl.buckets.Store(clientID, &Bucket{
		remainTokens: capacity,
		capacity:     capacity,
		refillRPS:    rps,
	})
}

func (rl *RateLimiter) UpdateBucket(clientID string, capacity, rps float64) {
	rl.DeleteBucket(clientID)
	rl.AddBucket(clientID, capacity, rps)
}

func (rl *RateLimiter) DeleteBucket(clientID string) {
	rl.buckets.Delete(clientID)
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		apiKey := request.URL.Query().Get("api_key")
		if apiKey == "" {
			dto.WriteErrorResponse(writer, http.StatusForbidden, "not rights")
			return
		}
		clientID, err := clients.ParseClientIDFromAPIKey(apiKey)
		if err != nil {
			dto.WriteErrorResponse(writer, http.StatusForbidden, err.Error())
			return
		}
		b, ok := rl.buckets.Load(clientID)
		var bucket *Bucket
		if !ok {
			dto.WriteErrorResponse(writer, http.StatusForbidden, err.Error())
			return
		}
		bucket, ok = b.(*Bucket)
		if !ok {
			slog.Error(
				"Failed type assertion to *Bucket",
				slog.String("actual_type", fmt.Sprintf("%T", b)),
			)
			dto.WriteErrorResponse(writer, http.StatusInternalServerError, "Internal error")
			return
		}

		if !bucket.TakeToken() {
			dto.WriteErrorResponse(writer, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (rl *RateLimiter) refillBucketTokensJob(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	slog.Info("Refill tokens in buckets job started")
	for {
		select {
		case <-ctx.Done():
			slog.Info("Refill tokens in buckets job stopped")
		case <-ticker.C:
			slog.Info("Refill tokens in buckets")
			rl.buckets.Range(func(ip, b any) bool {
				bucket, ok := b.(*Bucket)
				if !ok {
					ip, ok := ip.(string)
					var ipAttr slog.Attr
					if ok {
						ipAttr = slog.String("ip", ip)
					}
					slog.Error("Corrupted bucket", ipAttr)
					return true
				}
				bucket.RefillTokens(bucket.refillRPS)
				return true
			})
		}
	}
}
