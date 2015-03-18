package adhier

// Target is a quantity to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the quantity at a point in [0, 1]^n.
	Compute(point, value []float64)

	// Monitor is called once for each iteration before evaluating the quantity
	// at the active nodes of that iteration. The arguments are the iteration
	// number, number of passive nodes, and number of active nodes,
	// respectively.
	Monitor(iteration, passive, active uint)

	// Refine takes a node and its hierarchical surplus and identifies the
	// dimensions that should be refined by assigning scores to them.
	Refine(node, surplus, scores []float64)
}

// GenericTarget is a generic quantity satisfying the Target interface.
type GenericTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(uint, uint, uint)
	RefineHandler  func([]float64, []float64, []float64) // != nil
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

func (t *GenericTarget) Refine(node, surplus, score []float64) {
	t.RefineHandler(node, surplus, score)
}
