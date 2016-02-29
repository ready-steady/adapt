package hybrid

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and outputs.
	Dimensions() (uint, uint)

	// Continue decides if the interpolation process should go on.
	Continue(*external.Progress) bool

	// Compute evaluates the target function at a point.
	Compute([]float64, []float64)

	// Score assigns a score to a location.
	Score(*Location) float64
}

// Location contains information about a dimensional location.
type Location struct {
	Indices   []uint64  // Indices of the grid nodes
	Volumes   []float64 // Volumes of the basis functions
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	ContinueHandler func(*external.Progress) bool
	ComputeHandler  func([]float64, []float64) // != nil
	ScoreHandler    func(*Location) float64

	ni uint
	no uint
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, compute func([]float64, []float64)) *BasicTarget {
	return &BasicTarget{
		ComputeHandler: compute,

		ni: inputs,
		no: outputs,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.ni, self.no
}

func (self *BasicTarget) Continue(progress *external.Progress) bool {
	if self.ContinueHandler != nil {
		return self.ContinueHandler(progress)
	} else {
		return true
	}
}

func (self *BasicTarget) Compute(node, value []float64) {
	self.ComputeHandler(node, value)
}

func (self *BasicTarget) Score(location *Location) (score float64) {
	if self.ScoreHandler != nil {
		score = self.ScoreHandler(location)
	} else {
		for _, value := range location.Surpluses {
			score += math.Abs(value)
		}
		score /= float64(uint(len(location.Values)) / self.no)
	}
	return
}
