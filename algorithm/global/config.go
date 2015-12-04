package global

// Config represents a configuration of the algorithm.
type Config struct {
	// The maximum level of interpolation.
	MaxLevel uint

	// The maximum number of indices.
	MaxIndices uint

	// The maximum number of evaluations.
	MaxEvaluations uint

	// The degree of adaptivity.
	Adaptivity float64

	// The absolute-error tolerance.
	AbsTolerance float64

	// The relative-error tolerance.
	RelTolerance float64
}

// NewConfig returns a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		MaxLevel:       9,
		MaxIndices:     ^uint(0),
		MaxEvaluations: ^uint(0),

		Adaptivity:   1.0,
		AbsTolerance: 1e-6,
		RelTolerance: 1e-2,
	}
}
