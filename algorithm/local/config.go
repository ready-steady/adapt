package local

import (
	"runtime"
)

// Config represents a configuration of the algorithm.
type Config struct {
	// The minimum level of interpolation. The nodes that belong to lower levels
	// are unconditionally included in the surrogate.
	MinLevel uint

	// The maximum level of interpolation. The nodes that belong to this level
	// are never refined.
	MaxLevel uint

	// The number of concurrent workers. The evaluation of the target function
	// and the surrogate itself is distributed among this many goroutines.
	Workers uint
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MinLevel: 1,
		MaxLevel: 9,
		Workers:  uint(runtime.GOMAXPROCS(0)),
	}
}
