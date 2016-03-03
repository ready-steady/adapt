// Package equidistant provides means for working with the Newton–Cotes grid.
//
// Each node of an nd-dimensional grid is given by nd pairs (level, order). Each
// pair is given as a uint64 equal to (level|order<<levelSize) where levelSize
// is set to 6. In this encoding, the maximal level is 2^levelSize, and the
// maximal order is 2^(64-levelSize).
package equidistant

import (
	"github.com/ready-steady/linear"
)

const (
	levelMask = 0x3F
	levelSize = 6
	orderSize = 64 - levelSize
)

func index(lindices []uint64, generate func(uint64) []uint64, nd uint) []uint64 {
	nn := uint(len(lindices)) / nd

	cache := make(map[uint64][]uint64)

	indices1D := make([][]uint64, nd)
	indicesND := make([]uint64, 0)
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nd; j++ {
			level := lindices[i*nd+j]
			indices, ok := cache[level]
			if !ok {
				indices = generate(level)
				cache[level] = indices
			}
			indices1D[j] = indices
		}
		indicesND = append(indicesND, linear.TensorUint64(indices1D...)...)
	}

	return indicesND
}
