package adapt

// Target represents a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point in [0, 1]^n.
	Compute(point, value []float64)

	// Monitor keeps track of the interpolation progress. The function is called
	// once for each iteration before evaluating the target function at the
	// nodes of that iteration.
	Monitor(*Progress)

	// Score guides the local adaptivity. The function assigns a score to the
	// behavior of the target function at a particular node of the underlying
	// grid. A positive score signifies that the node should be refined, and the
	// score is the importance of this refinement. A zero score signifies that
	// the node should not be refined. A negative score signifies that the node
	// should be excluded from the interpolant.
	Score(*Location, *Progress) float64
}

// Location contains information about a spacial location.
type Location struct {
	Node    []float64 // Collocation node
	Surplus []float64 // Hierarchical surplus
	Volume  float64   // Volume under the basis function
}

// GenericTarget is a generic target satisfying the Target interface.
type GenericTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(*Progress)
	ScoreHandler   func(*Location, *Progress) float64 // != nil
}

// NewTarget returns a new generic target.
func NewTarget(inputs, outputs uint) *GenericTarget {
	return &GenericTarget{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (self *GenericTarget) Dimensions() (uint, uint) {
	return self.Inputs, self.Outputs
}

func (self *GenericTarget) Compute(node, value []float64) {
	self.ComputeHandler(node, value)
}

func (self *GenericTarget) Monitor(progress *Progress) {
	if self.MonitorHandler != nil {
		self.MonitorHandler(progress)
	}
}

func (self *GenericTarget) Score(location *Location, progress *Progress) float64 {
	return self.ScoreHandler(location, progress)
}
