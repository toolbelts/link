package container

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Container is a container for providers.
type Container struct {
	mu         sync.RWMutex
	exists     map[string]struct{}
	extensions []Provider
}

func New() *Container {
	return &Container{
		exists: make(map[string]struct{}),
	}
}

// Add adds a provider to the container.
func (c *Container) Add(p Provider) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.exists[p.Name()]; !ok {
		c.exists[p.Name()] = struct{}{}
		c.extensions = append(c.extensions, p)
	}
}

// Load loads all providers.
func (c *Container) Load() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, p := range c.extensions {
		start := time.Now()
		if err := p.Load(); err != nil {
			return err
		}

		log.Info().Str("extension", p.Name()).Dur("spent", time.Since(start)).Msg("extension load")
	}
	return nil
}

// Exit exits all providers.
func (c *Container) Exit() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, p := range c.extensions {
		start := time.Now()
		p.Exit()

		log.Info().Str("extension", p.Name()).Dur("spent", time.Since(start)).Msg("extension exit")
	}
}

var defaultContainer = New()

// Add adds a provider to the default container.
func Add(p Provider) {
	defaultContainer.Add(p)
}

// Load loads all providers in the default container.
func Load() error {
	return defaultContainer.Load()
}

// Exit exits all providers in the default container.
func Exit() {
	defaultContainer.Exit()
}
