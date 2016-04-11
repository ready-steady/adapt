package external

// Element contains information about an interpolation element.
type Element struct {
	Index []uint64 // Nodal index

	Volume  float64   // Basis-function volume
	Value   []float64 // Target-function value
	Surplus []float64 // Hierarchical surplus
}
