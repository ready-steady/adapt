package local

// Target represents a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point.
	Compute(point, value []float64)

	// Monitor gets called at the beginning of each iteration.
	Monitor(*Progress)

	// Score assigns a score to a spacial local.
	Score(*Location) float64
}

// Location contains information about a spacial location.
type Location struct {
	Surplus []float64 // Hierarchical surplus
	Volume  float64   // Volume under the basis function
}

// GenericTarget is a generic target satisfying the Target interface.
type GenericTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(*Progress)
	ScoreHandler   func(*Location) float64 // != nil
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

func (self *GenericTarget) Score(location *Location) float64 {
	return self.ScoreHandler(location)
}
