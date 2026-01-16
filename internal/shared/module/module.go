package module

// Module defines the interface that each feature module must implement.
// It provides a unified way to register handlers across different delivery mechanisms.
type Module interface {
	// Name returns the module name for logging/debugging purposes.
	Name() string
}
