package token

import "time"

type config struct {
	expire time.Duration
}

// apply
func (c *config) apply(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

// Default default configuration
func (c *config) Default() {
	if c.expire == 0 {
		c.expire = time.Hour * 24
	}
}

type Option func(*config)

// WithExpire set expiration time
func WithExpire(expire time.Duration) Option {
	return func(c *config) {
		c.expire = expire
	}
}
