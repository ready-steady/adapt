package adhier

// Config represents a configuration of the algorithm.
type Config struct {
	// The minimal level of interpolation. The nodes that belong to lower levels
	// are unconditionally included in the surrogate.
	MinLevel uint
	// The maximal level of interpolation. The nodes that belong to this level
	// are not refined, and, thus, the algorithm stops.
	MaxLevel uint
	// The maximal number of nodes. The algorithm stops after reaching this many
	// nodes.
	MaxNodes uint
	// The number of concurrent workers. The evaluation of the target function
	// and the surrogate itself is distributed among this many goroutines.
	Workers uint
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MinLevel: 1,
		MaxLevel: 9,
		MaxNodes: 10000,
	}
}
