package local

// Location contains information about a spacial location.
type Location struct {
	Value   []float64 // Target-function value
	Surplus []float64 // Hierarchical surplus
	Volume  float64   // Volume under the basis function
}

// Metric is an accuracy metric.
type Metric interface {
	// Score assigns a score to a location.
	Score(*Location) float64
}

// BasicMetric is a basic accuracy metric.
type BasicMetric struct {
	absolute float64
}

// NewMetric creates a basic accuracy metric.
func NewMetric(_ uint, absolute float64) *BasicMetric {
	return &BasicMetric{
		absolute: absolute,
	}
}

func (self *BasicMetric) Score(location *Location) float64 {
	absolute := self.absolute
	for _, ε := range location.Surplus {
		if ε < 0.0 {
			ε = -ε
		}
		if ε > absolute {
			return 1.0
		}
	}
	return 0.0
}
