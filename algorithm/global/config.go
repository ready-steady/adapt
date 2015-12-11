package global

import (
	"runtime"
)

// Config represents a configuration of the algorithm.
type Config struct {
	// The maximum level of interpolation.
	MaxLevel uint8

	// The maximum number of indices.
	MaxIndices uint

	// The maximum number of evaluations.
	MaxEvaluations uint

	// The degree of adaptivity.
	Adaptivity float64

	// The number of concurrent workers. The evaluation of the target function
	// and the surrogate itself is distributed among this many goroutines.
	Workers uint
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MaxLevel:       9,
		MaxIndices:     ^uint(0),
		MaxEvaluations: ^uint(0),

		Adaptivity: 1.0,

		Workers: uint(runtime.GOMAXPROCS(0)),
	}
}
