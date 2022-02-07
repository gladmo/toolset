package grate

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

// RateLimiter rate limiter with ttl
type RateLimiter struct {
	ttl    time.Duration
	cache  *cache.Cache
	bursts int
	limit  rate.Limit
}

// NewRateLimiter create rate limiter with ttl,
// every *rate.Limiter allows events up to
// rate limit and permits bursts of at most b tokens.
func NewRateLimiter(bursts int, limit rate.Limit, ttl time.Duration) *RateLimiter {
	return &RateLimiter{
		ttl:    ttl,
		cache:  cache.New(ttl, time.Minute*5),
		bursts: bursts,
		limit:  limit,
	}
}

// OnEvicted a callback function called when lim is purged from the cache.
func (th RateLimiter) OnEvicted(fn func(item string, lim *rate.Limiter)) {
	th.cache.OnEvicted(func(key string, value interface{}) {
		fn(key, value.(*rate.Limiter))
	})
}

// getLimiter return rate.Limiter
func (th RateLimiter) getLimiter(item string) (limiter *rate.Limiter) {
	defer func() {
		th.cache.Set(item, limiter, th.ttl)
	}()

	lim, ok := th.cache.Get(item)
	if !ok {
		return rate.NewLimiter(th.limit, th.bursts)
	}

	return lim.(*rate.Limiter)
}

// Remove rate.Limiter
func (th RateLimiter) Remove(item string) {
	th.cache.Delete(item)
}

// Allow is shorthand for getLimiter(item).Allow().
func (th RateLimiter) Allow(item string) bool {
	return th.getLimiter(item).Allow()
}

// Reserve is shorthand for getLimiter(item).Reserve().
func (th RateLimiter) Reserve(item string) *rate.Reservation {
	return th.getLimiter(item).Reserve()
}

// Wait is shorthand for getLimiter(item).Wait(ctx).
func (th RateLimiter) Wait(item string, ctx context.Context) error {
	return th.getLimiter(item).Wait(ctx)
}
