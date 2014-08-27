package basis

type Interface interface {
	ComputeNodes(levels []uint8, orders []uint32) []float64
	ComputeChildren(levels []uint8, orders []uint32) ([]uint8, []uint32)
	Evaluate(point float64, levels []uint8, orders []uint32, surpluses []float64) float64
}
