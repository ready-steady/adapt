package local

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point.
	Compute(point, value []float64)

	// Monitor gets called at the beginning of each iteration.
	Monitor(*Progress)

	// Score assigns a score to a location.
	Score(*Location) float64
}

// Location contains information about a spacial location.
type Location struct {
	Surplus []float64 // Hierarchical surplus
	Volume  float64   // Volume under the basis function
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(*Progress)
	ScoreHandler   func(*Location) float64 // != nil
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint) *BasicTarget {
	return &BasicTarget{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.Inputs, self.Outputs
}

func (self *BasicTarget) Compute(node, value []float64) {
	self.ComputeHandler(node, value)
}

func (self *BasicTarget) Monitor(progress *Progress) {
	if self.MonitorHandler != nil {
		self.MonitorHandler(progress)
	}
}

func (self *BasicTarget) Score(location *Location) float64 {
	return self.ScoreHandler(location)
}
