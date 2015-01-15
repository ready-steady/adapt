package adhier

// Config represents a configuration of the algorithm.
type Config struct {
	// The minimal level of interpolation.
	MinLevel uint8
	// The maximal level of interpolation.
	MaxLevel uint8
	// The tolerance of the absolute error.
	AbsError float64
	// The tolerance of the relative error.
	RelError float64
}

// DefaultConfig returns the default configuration of the algorithm.
func DefaultConfig() Config {
	return Config{
		MinLevel: 1,
		MaxLevel: 9,
		AbsError: 1e-4,
		RelError: 1e-2,
	}
}
