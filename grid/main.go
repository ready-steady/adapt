// Package grid contains code shared by the interpolation grids.
package grid

// Computer returns the nodes corresponding to a set of indices.
type Computer interface {
	Compute([]uint64) []float64
}

// Indexer returns the nodal indices of a set of level indices.
type Indexer interface {
	Index([]uint64) []uint64
}

// Parenter returns the parent index of an index in one dimension.
type Parenter interface {
	Parent(uint64, uint64) (uint64, uint64)
}

// Refiner returns the child indices of a set of indices.
type Refiner interface {
	Refine([]uint64) []uint64
}

// RefinerToward returns the child indices of a set of indices with respect to a
// particular dimension.
type RefinerToward interface {
	RefineToward([]uint64, uint) []uint64
}
