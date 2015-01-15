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

// DefaultConfig is the default configuration of the algorithm.
var DefaultConfig = Config{
	MinLevel: 1,
	MaxLevel: 9,
	AbsError: 1e-4,
	RelError: 1e-2,
}
