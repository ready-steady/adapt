package adapt

// Target is a quantity to be interpolated.
type Target interface {
	// Dimensions returns the number of inputs and the number of outputs.
	Dimensions() (uint, uint)

	// Compute evaluates the quantity at a point in [0, 1]^n.
	Compute(point, value []float64)

	// Monitor keeps track of the interpolation progress. The function is called
	// once for each interpolation step before the evaluation of the quantity at
	// the nodes of that step. The arguments of the function are the step
	// number, number of accepted nodes, number of rejected nodes, and number of
	// current nodes, respectively.
	Monitor(step, accept, reject, current uint)

	// Score guides the local adaptivity. The function takes a node, the
	// hierarchical surplus at the node, and the volume under the basis function
	// corresponding to the node. The function returns the score of the node. A
	// positive score represents the importance of refining the node. A zero
	// score signifies that the node should not be refined. A negative score
	// signifies that the node should be excluded from the interpolant.
	Score(node, surplus []float64, volume float64) float64
}

// GenericTarget is a generic quantity satisfying the Target interface.
type GenericTarget struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(uint, uint, uint, uint)
	ScoreHandler   func([]float64, []float64, float64) float64 // != nil
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

func (t *GenericTarget) Monitor(level, accept, reject, current uint) {
	if t.MonitorHandler != nil {
		t.MonitorHandler(level, accept, reject, current)
	}
}

func (t *GenericTarget) Score(node, surplus []float64, volume float64) float64 {
	return t.ScoreHandler(node, surplus, volume)
}
