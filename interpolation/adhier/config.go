package adhier

// Config represents a configuration of the algorithm.
type Config struct {
	MinLevel uint8   // The minimal level of interpolation
	MaxLevel uint8   // The maximal level of interpolation
	MaxNodes uint32  // The maximal number of nodes
	AbsError float64 // The absolute error
	RelError float64 // The relative error
}

// DefaultConfig returns the default configuration of the algorithm.
func DefaultConfig() Config {
	return Config{
		MinLevel: 1,
		MaxLevel: 9,
		MaxNodes: 10000,
		AbsError: 1e-4,
		RelError: 1e-2,
	}
}
