package container

// Provider is the interface that wraps the basic methods of a provider.
type Provider interface {
	Name() string
	Load() error
	Exit()
}
