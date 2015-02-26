package adhier

// Config represents a configuration of the algorithm.
type Config struct {
	// The number of inputs.
	Inputs uint
	// The number of outputs.
	Outputs uint
	// The minimal level of interpolation. The nodes that belong to lower levels
	// are unconditionally included in the surrogate.
	MinLevel uint32
	// The maximal level of interpolation. The nodes that belong to this level
	// are not refined, and, thus, the algorithm stops.
	MaxLevel uint32
	// The maximal number of nodes. The algorithm stops after reaching this many
	// nodes.
	MaxNodes uint
	// The absolute error tolerance. The parameter is used for local refinement
	// and is given in absolute units.
	AbsError float64
	// The relative error tolerance. The parameter is used for local refinement
	// and is given in relative units.
	RelError float64
	// The number of concurrent workers. The evaluation of the target function
	// and the surrogate itself is distributed among this many goroutines.
	Workers uint
}

// DefaultConfig returns the default configuration of the algorithm.
func DefaultConfig(inputs, outputs uint) *Config {
	return &Config{
		Inputs:   inputs,
		Outputs:  outputs,
		MinLevel: 1,
		MaxLevel: 9,
		MaxNodes: 10000,
		AbsError: 1e-4,
		RelError: 1e-2,
		Workers:  0,
	}
}
