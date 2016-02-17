package local

import (
	"math"
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point.
	Compute(point, value []float64)

	// Score assigns a score to a location. If the score is positive, the
	// corresponding node is refined; otherwise, no refinement is performed.
	Score(*Location) float64

	// Monitor gets called at the beginning of each iteration.
	Monitor(*Progress)
}

// Location contains information about a spacial location.
type Location struct {
	Value   []float64 // Target-function value
	Surplus []float64 // Hierarchical surplus
	Volume  float64   // Volume under the basis function
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level    uint      // Reached level
	Active   uint      // Number of active nodes
	Passive  uint      // Number of passive nodes
	Refined  uint      // Number of refined nodes
	Integral []float64 // Integral over the whole domain
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	Tolerance float64 // ≥ 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(*Progress)
	ScoreHandler   func(*Location) float64
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, tolerance float64,
	compute func([]float64, []float64)) *BasicTarget {

	return &BasicTarget{
		Inputs:  inputs,
		Outputs: outputs,

		Tolerance: tolerance,

		ComputeHandler: compute,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.Inputs, self.Outputs
}

func (self *BasicTarget) Compute(node, value []float64) {
	self.ComputeHandler(node, value)
}

func (self *BasicTarget) Score(location *Location) float64 {
	if self.MonitorHandler != nil {
		return self.ScoreHandler(location)
	} else {
		return self.defaultScore(location)
	}
}

func (self *BasicTarget) Monitor(progress *Progress) {
	if self.MonitorHandler != nil {
		self.MonitorHandler(progress)
	}
}

func (self *BasicTarget) defaultScore(location *Location) float64 {
	for _, ε := range location.Surplus {
		if math.Abs(ε) > self.Tolerance {
			return 1.0
		}
	}
	return 0.0
}
