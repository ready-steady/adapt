package local

// Config represents a configuration of the algorithm.
type Config struct {
	MinLevel uint // Minimum level of interpolation
	MaxLevel uint // Maximum level of interpolation
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MinLevel: 1,
		MaxLevel: 10,
	}
}
