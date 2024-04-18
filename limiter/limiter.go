package limiter

import (
	"context"
	"regexp"
	"sync"

	"github.com/bsm/ratelimit"
)

// Limiter is the rate limiter
type Limiter struct {
	mu   sync.RWMutex
	data map[string]*ratelimit.RateLimiter
	opts option
}

// New creates a new rate limiter
func New(opts ...Option) *Limiter {
	o := option{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Limiter{
		opts: o,
		data: make(map[string]*ratelimit.RateLimiter),
	}
}

// Len returns the number of rate limiters
func (l *Limiter) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.data)
}

// SetOptions sets the options for the rate limiter
func (l *Limiter) SetOptions(opts ...Option) {
	l.mu.Lock()
	for _, opt := range opts {
		opt(&l.opts)
	}
	l.mu.Unlock()
}

// getConfigs returns the configurations for the given key
func (l *Limiter) getConfigs(key string) (configs []Config) {
	l.mu.RLock()
	for _, cfg := range l.opts.Configs {
		if key == cfg.Key || regexp.MustCompile(cfg.Key).FindString(key) == key {
			configs = append(configs, cfg)
		}
	}
	l.mu.RUnlock()
	return
}

// getRate returns the rate limiter for the given configuration
func (l *Limiter) getRate(ctx context.Context, cfg *Config) (rl *ratelimit.RateLimiter) {
	key := cfg.Key
	if l.opts.KeyFunc != nil {
		key += l.opts.KeyFunc(ctx, cfg)
	}

	l.mu.Lock()
	rl, ok := l.data[key]
	if !ok {
		rl = ratelimit.New(cfg.Rate, cfg.Duration)
		l.data[key] = rl
	}
	l.mu.Unlock()
	return rl
}

// Limit limits the request
func (l *Limiter) Limit(ctx context.Context, key string) (ok bool) {
	cfgs := l.getConfigs(key)

	for _, cfg := range cfgs {
		rl := l.getRate(ctx, &cfg)
		ok = rl.Limit()
		if ok || cfg.Skip {
			return
		}
	}
	return
}
