package adhier

// Quantity is a generic quantity satisfying the Target interface.
type Quantity struct {
	Inputs  uint // > 0
	Outputs uint // > 0

	ComputeHandler func([]float64, []float64) // != nil
	MonitorHandler func(uint, uint, uint)
	RefineHandler  func([]float64) bool // != nil
}

// NewQuantity returns a new generic quantity.
func NewQuantity(inputs, outputs uint) *Quantity {
	return &Quantity{
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func (q *Quantity) Dimensions() (uint, uint) {
	return q.Inputs, q.Outputs
}

func (q *Quantity) Compute(node, value []float64) {
	q.ComputeHandler(node, value)
}

func (q *Quantity) Monitor(level, active, passive uint) {
	if q.MonitorHandler != nil {
		q.MonitorHandler(level, active, passive)
	}
}

func (q *Quantity) Refine(surplus []float64) bool {
	return q.RefineHandler(surplus)
}
