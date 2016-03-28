package global

// Target is a function to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the target function at a point.
	Compute([]float64, []float64)
}

// BasicTarget is a basic target satisfying the Target interface.
type BasicTarget struct {
	ni uint
	no uint

	handler func([]float64, []float64)
}

// NewTarget creates a basic target.
func NewTarget(inputs, outputs uint, compute func([]float64, []float64)) *BasicTarget {
	return &BasicTarget{
		ni: inputs,
		no: outputs,

		handler: compute,
	}
}

func (self *BasicTarget) Dimensions() (uint, uint) {
	return self.ni, self.no
}

func (self *BasicTarget) Compute(node, value []float64) {
	self.handler(node, value)
}
