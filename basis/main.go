package basis

type Interface interface {
	ComputeOrders(uint8) []uint16
	ComputeNodes([]uint8, []uint16) []float64
	ComputeChildren([]uint8, []uint16) ([]uint8, []uint16)
}
