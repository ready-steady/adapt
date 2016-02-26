package hybrid

import (
	"runtime"
)

// Config represents a configuration of the algorithm.
type Config struct {
	// The maximum level of interpolation.
	MaxLevel uint

	// The maximum number of indices.
	MaxIndices uint

	// The number of concurrent workers. The evaluation of the target function
	// and the surrogate itself is distributed among this many goroutines.
	Workers uint
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MaxLevel:   9,
		MaxIndices: ^uint(0),
		Workers:    uint(runtime.GOMAXPROCS(0)),
	}
}
