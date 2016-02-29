package local

import (
	"runtime"
)

// Config represents a configuration of the algorithm.
type Config struct {
	MinLevel uint // Minimum level of interpolation
	MaxLevel uint // Maximum level of interpolation

	Workers uint // Number of concurrent workers
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MinLevel: 1,
		MaxLevel: 10,

		Workers: uint(runtime.GOMAXPROCS(0)),
	}
}
