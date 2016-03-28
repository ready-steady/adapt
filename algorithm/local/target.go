package local

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and outputs.
	Dimensions() (uint, uint)

	// Done checks if the stopping criteria have been satisfied.
	Done(*external.Progress) bool

	// Compute evaluates the target function at a point.
	Compute([]float64, []float64)

	// Score assigns a score to an interpolation element.
	Score(*Element) float64
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	DoneHandler    func(*external.Progress) bool
	ComputeHandler func([]float64, []float64) // != nil
	ScoreHandler   func(*Element) float64

	ni uint
	no uint

	ε float64
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, absolute float64,
	compute func([]float64, []float64)) *BasicTarget {

	return &BasicTarget{
		ComputeHandler: compute,

		ni: inputs,
		no: outputs,

		ε: absolute,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.ni, self.no
}

func (self *BasicTarget) Done(progress *external.Progress) bool {
	if self.DoneHandler != nil {
		return self.DoneHandler(progress)
	} else {
		return false
	}
}

func (self *BasicTarget) Compute(node, value []float64) {
	self.ComputeHandler(node, value)
}

func (self *BasicTarget) Score(element *Element) float64 {
	if self.ScoreHandler != nil {
		return self.ScoreHandler(element)
	}

	for _, ε := range element.Surplus {
		if math.Abs(ε) > self.ε {
			return 1.0
		}
	}

	return 0.0
}
