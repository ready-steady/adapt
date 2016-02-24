package hybrid

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1)
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Continue is called at the end of each iteration. If the function returns
	// false, the interpolation process is terminated. The first argument is the
	// set of currently active indices.
	Continue(*external.Active, *Progress) bool

	// Compute evaluates the target function at a point. The function is called
	// once for each admissible node of the admissible neighbors.
	Compute([]float64, []float64)

	// Score assigns a global score and a set of local scores to a location. The
	// function is called once for each admissible neighbor.
	Score(*Location) (float64, []float64)
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

	Global float64 // ≥ 0
	Local  float64 // ≥ 0

	ContinueHandler func(*external.Active, *Progress) bool
	ComputeHandler  func([]float64, []float64) // != nil
	ScoreHandler    func(*Location) (float64, []float64)

	scores []float64
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, global, local float64,
	compute func([]float64, []float64)) *BasicTarget {

	return &BasicTarget{
		Inputs:  inputs,
		Outputs: outputs,

		Global: global,
		Local:  local,

		ComputeHandler: compute,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.Inputs, self.Outputs
}

func (self *BasicTarget) Continue(active *external.Active, progress *Progress) bool {
	if self.ContinueHandler != nil {
		return self.ContinueHandler(active, progress)
	} else {
		return self.defaultContinue(active, progress)
	}
}

func (self *BasicTarget) Compute(node, value []float64) {
	self.ComputeHandler(node, value)
}

func (self *BasicTarget) Score(location *Location) (float64, []float64) {
	if self.ScoreHandler != nil {
		return self.ScoreHandler(location)
	} else {
		return self.defaultScore(location)
	}
}

func (self *BasicTarget) defaultContinue(active *external.Active, progress *Progress) bool {
	return true
}

func (self *BasicTarget) defaultScore(location *Location) (float64, []float64) {
	no := self.Outputs
	nn := uint(len(location.Volumes))

	global := -infinity
	local := internal.RepeatFloat64(-infinity, nn)
	for i := uint(0); i < no; i++ {
		Γ := 0.0
		for j := uint(0); j < nn; j++ {
			γ := location.Volumes[j] * location.Surpluses[j*no+i]
			local[j] = math.Max(local[j], math.Abs(γ))
			Γ += γ
		}
		global = math.Max(global, math.Abs(Γ))
	}

	return global, local
}
