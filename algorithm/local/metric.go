package local

import (
	"math"
)

// Location contains information about a spacial location.
type Location struct {
	Value   []float64 // Target-function value
	Surplus []float64 // Hierarchical surplus
	Volume  float64   // Volume under the basis function
}

// Metric is an accuracy metric.
type Metric interface {
	// Score assigns a score to a location. If the score is positive, the
	// corresponding node is refined; otherwise, no refinement is performed.
	Score(*Location) float64
}

// BasicMetric is a basic accuracy metric.
type BasicMetric struct {
	tolerance float64
}

// NewMetric creates a basic accuracy metric.
func NewMetric(_ uint, tolerance float64) *BasicMetric {
	return &BasicMetric{
		tolerance: tolerance,
	}
}

func (self *BasicMetric) Score(location *Location) float64 {
	for _, ε := range location.Surplus {
		if math.Abs(ε*location.Volume) > self.tolerance {
			return 1.0
		}
	}
	return 0.0
}
