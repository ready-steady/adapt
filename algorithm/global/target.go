package global

import (
	"math"
)

var (
	infinity = math.Inf(1.0)
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Before gets called once per iteration before involving Compute. If the
	// function returns false, the interpolation process is terminated.
	Before(*Progress) bool

	// Compute evaluates the target function at a point. The function is called
	// multiple times per iteration, depending on the number of active nodes.
	Compute(point, value []float64)

	// Score assigns a score to a location. The function is called after
	// Compute, and it is called as many times as Compute.
	Score(*Location) float64

	// After gets called once per iteration after involving Compute and Score.
	// The argument of the function is the set of currently active indices. If
	// the function returns false, the interpolation process is terminated.
	After(Set) bool
}

// Location contains information about a dimensional location.
type Location struct {
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level     uint // Reached level
	Active    uint // Number of active level indices
	Passive   uint // Number of passive level indices
	Requested uint // Number of requested function evaluations
	Performed uint // Number of performed function evaluations
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	Absolute float64 // ≥ 0
	Relative float64 // ≥ 0

	BeforeHandler  func(*Progress) bool
	ComputeHandler func([]float64, []float64) // != nil
	ScoreHandler   func(*Location) float64
	AfterHandler   func(Set) bool

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

func (self *BasicTarget) Before(progress *Progress) bool {
	if self.BeforeHandler != nil {
		return self.BeforeHandler(progress)
	} else {
		return true
	}
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

func (self *BasicTarget) After(active Set) bool {
	if self.AfterHandler != nil {
		return self.AfterHandler(active)
	} else {
		return self.defaultAfter(active)
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

func (self *BasicTarget) defaultAfter(active Set) bool {
	no, errors := self.Outputs, self.errors
	δ := threshold(self.lower, self.upper, self.Absolute, self.Relative)
	for i := range active {
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > δ[j] {
				return true
			}
		}
	}
	return false
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
