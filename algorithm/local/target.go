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

// Element contains information about an interpolation element.
type Element struct {
	Index   []uint64  // Nodal index
	Volume  float64   // Basis-function volume
	Value   []float64 // Target-function value
	Surplus []float64 // Hierarchical surplus
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

func (self *BasicTarget) Score(element *Element) (score float64) {
	if self.ScoreHandler != nil {
		score = self.ScoreHandler(element)
	} else {
		for _, ε := range element.Surplus {
			if math.Abs(ε) > self.ε {
				score = 1.0
				break
			}
		}
	}
	return
}
