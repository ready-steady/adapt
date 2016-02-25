package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/internal"
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Continue is called at the end of each iteration. If the function returns
	// false, the interpolation process is terminated. The first argument is the
	// set of currently active indices.
	Continue(*Active, *Progress) bool

	// Compute evaluates the target function at a point. The function is called
	// once for each node of the admissible neighbors.
	Compute([]float64, []float64)

	// Score assigns a score to a location. The function is called once for each
	// admissible neighbor.
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
	ContinueHandler func(*Active, *Progress) bool
	ComputeHandler  func([]float64, []float64) // != nil
	ScoreHandler    func(*Location) float64

	ni uint
	no uint

	absolute float64
	relative float64

	errors []float64
	lower  []float64
	upper  []float64
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, absolute, relative float64,
	compute func([]float64, []float64)) *BasicTarget {

	return &BasicTarget{
		ni: inputs,
		no: outputs,

		absolute: absolute,
		relative: relative,

		ComputeHandler: compute,

		lower: internal.RepeatFloat64(infinity, outputs),
		upper: internal.RepeatFloat64(-infinity, outputs),
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.ni, self.no
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
	no, errors := self.no, self.errors
	ne := uint(len(errors)) / no
	if ne == 0 {
		return true
	}
	δ := threshold(self.lower, self.upper, self.absolute, self.relative)
	for i := range active.Positions {
		if i >= ne {
			continue
		}
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > δ[j] {
				return true
			}
		}
	}
	return false
}

func (self *BasicTarget) defaultScore(location *Location) float64 {
	no := self.no

	self.updateBounds(location.Values)
	self.errors = append(self.errors, error(location.Surpluses, no)...)

	score := 0.0
	for _, value := range location.Surpluses {
		score += math.Abs(value)
	}
	score /= float64(uint(len(location.Values)) / no)

	return score
}

func (self *BasicTarget) updateBounds(values []float64) {
	no := self.no
	for i, point := range values {
		j := uint(i) % no
		self.lower[j] = math.Min(self.lower[j], point)
		self.upper[j] = math.Max(self.upper[j], point)
	}
}

func error(surpluses []float64, no uint) []float64 {
	error := internal.RepeatFloat64(-infinity, no)
	for i, value := range surpluses {
		j := uint(i) % no
		error[j] = math.Max(error[j], math.Abs(value))
	}
	return error
}

func threshold(lower, upper []float64, absolute, relative float64) []float64 {
	threshold := make([]float64, len(lower))
	for i := range threshold {
		threshold[i] = math.Max(relative*(upper[i]-lower[i]), absolute)
	}
	return threshold
}
