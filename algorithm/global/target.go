package global

// Target represents a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point.
	Compute(point, value []float64)
}

// GenericTarget is a generic target satisfying the Target interface.
type GenericTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
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
