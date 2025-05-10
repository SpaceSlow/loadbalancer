package ratelimiter

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/SpaceSlow/loadbalancer/config"
	"github.com/SpaceSlow/loadbalancer/pkg/networks"
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

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ip := networks.ParseIP(request.RemoteAddr)
		b, ok := rl.buckets.Load(ip)
		var bucket *Bucket
		if !ok {
			// create default bucket for anonymous ip
			bucket = NewBucket(&rl.cfg.DefaultBucket)
			rl.buckets.Store(ip, bucket)
		} else {
			bucket, ok = b.(*Bucket)
			if !ok {
				http.Error(writer, "Internal Error", http.StatusInternalServerError)
				return
			}
		}

		if !bucket.TakeToken() {
			http.Error(writer, "Rate limit exceeded", http.StatusTooManyRequests)
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
				slog.Info("Refill tokens in buckets")
				return true
			})
		}
	}
}
