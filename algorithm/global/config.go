package global

// Config represents a configuration of the algorithm.
type Config struct {
	MinLevel uint // Minimum level of interpolation
	MaxLevel uint // Maximum level of interpolation

	AbsoluteError float64 // Tolerance on the absolute error
	RelativeError float64 // Tolerance on the relative error
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MinLevel: 1,
		MaxLevel: 10,

		AbsoluteError: 1e-6,
		RelativeError: 1e-3,
	}
}
