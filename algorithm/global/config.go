package global

// Config represents a configuration of the algorithm.
type Config struct {
	// The minimum level of interpolation. The nodes that belong to lower levels
	// are unconditionally included in the surrogate.
	MinLevel uint

	// The maximum level of interpolation. The nodes that belong to this level
	// are never refined.
	MaxLevel uint

	// The maximum number of target-function evaluations. The limit is not
	// enforced precisely. The process stops when undertaking the next iteration
	// would violate the limit.
	MaxEvaluations uint
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MinLevel:       1,
		MaxLevel:       9,
		MaxEvaluations: ^uint(0),
	}
}
