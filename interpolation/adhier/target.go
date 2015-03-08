package adhier

// Target is a quantity to be interpolated.
type Target interface {
	// Inputs returns the number of inputs.
	Inputs() uint

	// Outputs returns the number of outputs.
	Outputs() uint

	// Compute calculates the quantity of interest at a node.
	Compute(node, value []float64)

	// Monitor keeps track of the interpolation progress. The function is called
	// once on each level before evaluating the target function at the nodes of
	// that level. The arguments are the current level, number of active nodes,
	// and number of passive nodes, respectively.
	Monitor(level, active, passive uint)

	// Refine guides the local adaptivity of interpolation. The function is
	// called for each hierarchical surplus in order to check if the
	// corresponding node should be refined.
	Refine(surplus []float64) bool
}

// AbsErrorTarget is a target whose local adaptivity guide is the absolute error
// at the nodes of the underlying sparse grid.
type AbsErrorTarget struct {
	inputs    uint
	outputs   uint
	tolerance float64

	ComputeFunc func([]float64, []float64)
	MonitorFunc func(uint, uint, uint)
}

// NewAbsErrorTarget returns an absolute-error-driven target.
func NewAbsErrorTarget(inputs, outputs uint, tolerance float64,
	compute func([]float64, []float64)) *AbsErrorTarget {

	return &AbsErrorTarget{
		inputs:    inputs,
		outputs:   outputs,
		tolerance: tolerance,

		ComputeFunc: compute,
	}
}

func (t *AbsErrorTarget) Inputs() uint {
	return t.inputs
}

func (t *AbsErrorTarget) Outputs() uint {
	return t.outputs
}

func (t *AbsErrorTarget) Compute(node, value []float64) {
	t.ComputeFunc(node, value)
}

func (t *AbsErrorTarget) Monitor(level, active, passive uint) {
	if t.MonitorFunc != nil {
		t.MonitorFunc(level, active, passive)
	}
}

func (t *AbsErrorTarget) Refine(surplus []float64) bool {
	for _, ε := range surplus {
		if ε < 0 {
			ε = -ε
		}
		if ε > t.tolerance {
			return true
		}
	}
	return false
}
