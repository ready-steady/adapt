package hybrid

import (
	"runtime"
)

// Config represents a configuration of the algorithm.
type Config struct {
	MaxLevel   uint // Maximum level of interpolation
	MaxIndices uint // Maximum number of indices

	GlobalError float64 // Tolerance on the global error
	LocalError  float64 // Tolerance on the local error

	Workers uint // Number of concurrent workers
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MaxLevel:   10,
		MaxIndices: ^uint(0),

		GlobalError: 1e-2,
		LocalError:  1e-2,

		Workers: uint(runtime.GOMAXPROCS(0)),
	}
}
