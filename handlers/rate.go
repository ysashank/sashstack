package handlers

import (
	"sync"
	"time"
)

var rateLimiter = &ipRateLimiter{visits: make(map[string]time.Time)}

type ipRateLimiter struct {
	sync.Mutex
	visits map[string]time.Time
}

func (r *ipRateLimiter) Allow(ip string) bool {
	r.Lock()
	defer r.Unlock()

	now := time.Now()
	last, exists := r.visits[ip]
	if exists && now.Sub(last) < 10*time.Second {
		return false
	}
	r.visits[ip] = now
	return true
}
