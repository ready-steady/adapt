package hybrid

import (
	"math"
)

var (
	infinity = math.Inf(1)
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Continue gets called at the end of each iteration. If the function
	// returns false, the interpolation process is terminated. The first
	// argument is the set of currently active indices.
	Continue(*Active, *Progress) bool

	// Compute evaluates the target function at a point. The function is called
	// multiple times per iteration, depending on the number of active nodes.
	Compute(point, value []float64)

	// Score assigns a score to a location. The function is called after
	// Compute, and it is called as many times as Compute.
	Score(*Location) float64
}

// Location contains information about a dimensional location.
type Location struct {
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
	Volumes   []float64 // Volumes under the basis functions
}

// Progress contains information about the interpolation process.
type Progress struct {
	More uint // Number of nodes to be evaluated
	Done uint // Number of nodes evaluated so far
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ContinueHandler func(*Active, *Progress) bool
	ComputeHandler  func([]float64, []float64) // != nil
	ScoreHandler    func(*Location) float64

	scores []float64
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, compute func([]float64, []float64)) *BasicTarget {
	return &BasicTarget{
		Inputs:  inputs,
		Outputs: outputs,

		ComputeHandler: compute,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.Inputs, self.Outputs
}

func (self *BasicTarget) Continue(active *Active, progress *Progress) bool {
	if self.ContinueHandler != nil {
		return self.ContinueHandler(active, progress)
	} else {
		return self.defaultContinue(active, progress)
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

func (self *BasicTarget) defaultContinue(active *Active, progress *Progress) bool {
	return true
}

func (self *BasicTarget) defaultScore(location *Location) float64 {
	no := self.Outputs
	nn := uint(len(location.Volumes))

	score := -infinity
	for i := uint(0); i < no; i++ {
		value := 0.0
		for j := uint(0); j < nn; j++ {
			value += location.Volumes[j] * location.Surpluses[j*no+i]
		}
		score = math.Max(score, math.Abs(value))
	}
	self.scores = append(self.scores, score)

	return score
}
