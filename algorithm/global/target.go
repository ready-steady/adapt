package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point.
	Compute(point, value []float64)

	// Score assigns a score to a location.
	Score(*Location) float64

	// Done checks if the accuracy requirements have been satiated.
	Done(internal.Set) bool

	// Monitor gets called at the beginning of each iteration.
	Monitor(*Progress)
}

// Location contains information about a dimensional location.
type Location struct {
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	Absolute float64 // ≥ 0
	Relative float64 // ≥ 0

	ComputeHandler func([]float64, []float64) // != nil
	ScoreHandler   func(*Location) float64
	DoneHandler    func(internal.Set) bool
	MonitorHandler func(*Progress)

	errors []float64
	lower  []float64
	upper  []float64
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, absolute, relative float64,
	compute func([]float64, []float64)) *BasicTarget {

	return &BasicTarget{
		Inputs:  inputs,
		Outputs: outputs,

		Absolute: absolute,
		Relative: relative,

		ComputeHandler: compute,

		lower: repeatFloat64(infinity, outputs),
		upper: repeatFloat64(-infinity, outputs),
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.Inputs, self.Outputs
}

func (self *BasicTarget) Compute(node, value []float64) {
	self.ComputeHandler(node, value)
}

func (self *BasicTarget) Score(location *Location) float64 {
	if self.ScoreHandler != nil {
		return self.ScoreHandler(location)
	} else {
		return self.defaultScore(location)
	}
}

func (self *BasicTarget) Done(active internal.Set) bool {
	if self.DoneHandler != nil {
		return self.DoneHandler(active)
	} else {
		return self.defaultDone(active)
	}
}

func (self *BasicTarget) Monitor(progress *Progress) {
	if self.MonitorHandler != nil {
		self.MonitorHandler(progress)
	}
}

func (self *BasicTarget) defaultScore(location *Location) float64 {
	no := self.Outputs
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

func (self *BasicTarget) defaultDone(active internal.Set) bool {
	no, errors := self.Outputs, self.errors
	δ := threshold(self.lower, self.upper, self.Absolute, self.Relative)
	for i := range active {
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > δ[j] {
				return false
			}
		}
	}
	return true
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
