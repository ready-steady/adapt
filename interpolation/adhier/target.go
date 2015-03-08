package adhier

// Target is a quantity to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute calculates the value of the quantity at a point.
	Compute(point, value []float64)

	// Monitor is called once on each level before evaluating the quantity at
	// the nodes of that level. The arguments are the current level, number of
	// active nodes, and number of passive nodes, respectively.
	Monitor(level, active, passive uint)

	// Refine checks if a node of the underlying sparse grid should be refined
	// based on its hierarchical surplus, which is the difference between the
	// true value of the quantity at the node and its current approximation.
	Refine(surplus []float64) bool
}

// GenericTarget is a generic quantity to be interpolated.
type GenericTarget struct {
	Inputs  uint
	Outputs uint

	ComputeFunc func([]float64, []float64)
	MonitorFunc func(uint, uint, uint)
	RefineFunc  func([]float64) bool
}

// NewTarget returns a generic target.
func NewTarget(inputs, outputs uint) *GenericTarget {
	return &GenericTarget{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (t *GenericTarget) Dimensions() (uint, uint) {
	return t.Inputs, t.Outputs
}

func (t *GenericTarget) Compute(node, value []float64) {
	t.ComputeFunc(node, value)
}

func (t *GenericTarget) Monitor(level, active, passive uint) {
	if t.MonitorFunc != nil {
		t.MonitorFunc(level, active, passive)
	}
}

func (t *GenericTarget) Refine(surplus []float64) bool {
	return t.RefineFunc(surplus)
}
