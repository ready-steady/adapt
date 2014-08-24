package basis

type Interface interface {
	ComputeOrders(uint8) []uint32
	ComputeNodes([]uint8, []uint32) []float64
	ComputeChildren([]uint8, []uint32) ([]uint8, []uint32)
}
