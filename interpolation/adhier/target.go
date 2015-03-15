package adhier

// Target is a quantity to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the quantity at a point in [0, 1]^n.
	Compute(point, value []float64)

	// Monitor is called once on each level before evaluating the quantity at
	// the nodes of that level. The arguments are the current level, number of
	// passive nodes, and number of active nodes, respectively.
	Monitor(level, passive, active uint)

	// Refine identifies the dimensions of the underlying sparse grid that
	// should be refined based on a hierarchical surplus, which is the
	// difference between the true value of the quantity at a node and its
	// current approximation.
	Refine(surplus []float64, dimensions []bool)
}

// GenericTarget is a generic quantity satisfying the Target interface.
type GenericTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(uint, uint, uint)
	RefineHandler  func([]float64, []bool) // != nil
}

// NewTarget returns a new generic quantity.
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
	t.ComputeHandler(node, value)
}

func (t *GenericTarget) Monitor(level, passive, active uint) {
	if t.MonitorHandler != nil {
		t.MonitorHandler(level, passive, active)
	}
}

func (t *GenericTarget) Refine(surplus []float64, dimensions []bool) {
	t.RefineHandler(surplus, dimensions)
}
