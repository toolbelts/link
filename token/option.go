package token

import "time"

type option struct {
	expire time.Duration
}

// apply apply options
func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Default default configuration
func (o *option) Default() {
	if o.expire == 0 {
		o.expire = time.Hour * 24
	}
}

type Option func(*option)

// WithExpire set expiration time
func WithExpire(expire time.Duration) Option {
	return func(c *option) {
		c.expire = expire
	}
}
