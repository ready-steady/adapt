package adapt

// Config represents a configuration of the algorithm.
type Config struct {
	// The refinement rate of the algorithm. The parameter specifies the
	// fraction of the nodes queued for refinement to be taken from the queue at
	// each iteration.
	Rate float64 // ⊆ (0, 1]

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

	// The maximum number of iterations. Depending on Rate, an iteration may or
	// may not correspond to a level.
	MaxIterations uint

	// The number of concurrent workers. The evaluation of the target function
	// and the surrogate itself is distributed among this many goroutines.
	Workers uint
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		Rate:           1.0,
		MinLevel:       1,
		MaxLevel:       9,
		MaxEvaluations: ^uint(0),
		MaxIterations:  ^uint(0),
	}
}
