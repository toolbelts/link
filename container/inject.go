package container

import (
	"sync"
)

var mapper sync.Map

// Get a value from the container
func Get[T any](key string) T {
	if v, ok := mapper.Load(key); ok {
		return v.(T)
	}
	return *new(T)
}

// Set a value to the container
func Set[T any](key string, value T) {
	mapper.Store(key, value)
}
