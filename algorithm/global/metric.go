package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Location contains information about a dimensional location.
type Location struct {
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
}

// Metric is an accuracy metric.
type Metric interface {
	// Done checks if the accuracy requirements have been satiated.
	Done(active internal.Set) bool

	// Score assigns a score to a location.
	Score(*Location) float64
}

// BasicMetric is a basic accuracy metric.
type BasicMetric struct {
	no       uint
	absolute float64
	relative float64

	errors []float64
	lower  []float64
	upper  []float64
}

// NewMetric creates a basic accuracy metric.
func NewMetric(no uint, absolute, relative float64) *BasicMetric {
	return &BasicMetric{
		no:       no,
		absolute: absolute,
		relative: relative,

		lower: repeatFloat64(infinity, no),
		upper: repeatFloat64(-infinity, no),
	}
}

func (self *BasicMetric) Done(active internal.Set) bool {
	no, errors := self.no, self.errors
	δ := threshold(self.lower, self.upper, self.absolute, self.relative)
	for i := range active {
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > δ[j] {
				return false
			}
		}
	}
	return true
}

func (self *BasicMetric) Score(location *Location) float64 {
	no := self.no
	nn := uint(len(location.Values)) / no

	for i, point := range location.Values {
		j := uint(i) % no
		if self.lower[j] > point {
			self.lower[j] = point
		}
		if self.upper[j] < point {
			self.upper[j] = point
		}
	}
	self.errors = append(self.errors, error(location.Surpluses, no)...)

	score := 0.0
	for _, value := range location.Surpluses {
		if value < 0.0 {
			value = -value
		}
		score += value
	}

	return score / float64(nn)
}

func error(surpluses []float64, no uint) []float64 {
	error := repeatFloat64(-infinity, no)
	for i, value := range surpluses {
		j := uint(i) % no
		if value < 0.0 {
			value = -value
		}
		if value > error[j] {
			error[j] = value
		}
	}
	return error
}

func threshold(lower, upper []float64, absolute, relative float64) []float64 {
	no := uint(len(lower))
	threshold := make([]float64, no)
	for i := uint(0); i < no; i++ {
		threshold[i] = relative * (upper[i] - lower[i])
		if threshold[i] < absolute {
			threshold[i] = absolute
		}
	}
	return threshold
}
