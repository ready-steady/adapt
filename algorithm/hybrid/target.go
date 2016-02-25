package hybrid

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Continue decides if the interpolation process should go on. The function
	// is called at the beginning of each iteration.
	Continue(*external.Active, *Progress) bool

	// Compute evaluates the target function at a point. The function is called
	// for each node of the admissible neighbors of the currently refined index.
	Compute([]float64, []float64)

	// Score assesses a location. The function is called for each admissible
	// neighbor of the currently refined index.
	Score(*Location)

	// Select decides which of the active indices should be refined next. The
	// function is called at the end of each iteration.
	Select(*external.Active) uint
}

// Location contains information about a dimensional location.
type Location struct {
	Indices   []uint64  // Indices of the grid nodes
	Volumes   []float64 // Volumes of the basis functions
	Values    []float64 // Target-function values
	Surpluses []float64 // Hierarchical surpluses
}

// Progress contains information about the interpolation process.
type Progress struct {
	More uint // Number of nodes to be evaluated
	Done uint // Number of nodes evaluated so far
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	ContinueHandler func(*external.Active, *Progress) bool
	ComputeHandler  func([]float64, []float64) // != nil
	ScoreHandler    func(*Location)
	SelectHandler   func(*external.Active) uint

	ni uint
	no uint

	εg float64
	εl float64

	global []float64
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, global, local float64,
	compute func([]float64, []float64)) *BasicTarget {

	return &BasicTarget{
		ComputeHandler: compute,

		ni: inputs,
		no: outputs,

		εg: global,
		εl: local,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.ni, self.no
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

func (self *BasicTarget) Score(location *Location) {
	if self.ScoreHandler != nil {
		self.ScoreHandler(location)
	} else {
		self.defaultScore(location)
	}
}

func (self *BasicTarget) Select(active *external.Active) uint {
	if self.SelectHandler != nil {
		return self.SelectHandler(active)
	} else {
		return self.defaultSelect(active)
	}
}

func (self *BasicTarget) defaultContinue(active *external.Active, progress *Progress) bool {
	Σ := 0.0
	for i := range active.Positions {
		Σ += self.global[i]
	}
	return Σ > self.εg
}

func (self *BasicTarget) defaultScore(location *Location) {
	no := self.no
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
	self.global = append(self.global, global)
}

func (self *BasicTarget) defaultSelect(active *external.Active) uint {
	return internal.LocateMaxFloat64s(self.global, active.Positions)
}
