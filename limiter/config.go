package limiter

import (
	"context"
	"time"
)

const (
	TypeUserAll  = "all"
	TypeUserIp   = "user_ip"
	TypeUserId   = "user_id"
	TypeDeviceId = "device_id"
)

// Config is the configuration for the rate limiter
type Config struct {
	Key      string
	Type     string
	Rate     int
	Skip     bool
	Duration time.Duration
}

type KeyFunc func(context.Context, *Config) string

// option is the option for the rate limiter
type option struct {
	KeyFunc KeyFunc
	Configs []Config
}

// apply applies the options
func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(*option)

// WithKeyFunc sets the key function for the rate limiter
func WithKeyFunc(kf KeyFunc) Option {
	return func(o *option) {
		o.KeyFunc = kf
	}
}

// WithConfigs sets the configurations for the rate limiter
func WithConfigs(cfgs ...Config) Option {
	return func(o *option) {
		o.Configs = cfgs
	}
}
