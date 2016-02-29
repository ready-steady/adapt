package global

import (
	"runtime"
)

// Config represents a configuration of the algorithm.
type Config struct {
	MaxLevel   uint // Maximum level of interpolation
	MaxIndices uint // Maximum number of indices

	AbsoluteError float64 // Tolerance on the absolute error
	RelativeError float64 // Tolerance on the relative error

	Workers uint // Number of concurrent workers
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MaxLevel:   10,
		MaxIndices: ^uint(0),

		AbsoluteError: 1e-6,
		RelativeError: 1e-2,

		Workers: uint(runtime.GOMAXPROCS(0)),
	}
}
