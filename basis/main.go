package basis

type Interface interface {
	ComputeOrders(level uint8) []uint32
	ComputeNodes(levels []uint8, orders []uint32) []float64
	ComputeChildren(levels []uint8, orders []uint32) ([]uint8, []uint32)
	Evaluate(points []float64, levels []uint8, orders []uint32,
		surpluses []float64) []float64
}
