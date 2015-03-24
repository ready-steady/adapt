package adhier

// Config represents a configuration of the algorithm.
type Config struct {
	// The refinement rate of the algorithm. The parameter specifies the
	// fraction of the nodes queued for refinement to be taken from the queue at
	// each iteration.
	Rate float64 // ⊆ (0, 1]
	// The minimal level of interpolation. The nodes that belong to lower levels
	// are unconditionally included in the surrogate.
	MinLevel uint
	// The maximal level of interpolation. The nodes that belong to this level
	// are not refined, and, thus, the algorithm stops.
	MaxLevel uint
	// The flag controlling grid balancing. If the flag is true, additional
	// nodes are added at each iteration to balance the underlying grid.
	Balance bool
	// The number of concurrent workers. The evaluation of the target function
	// and the surrogate itself is distributed among this many goroutines.
	Workers uint
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		Rate:     1,
		MinLevel: 1,
		MaxLevel: 9,
	}
}
