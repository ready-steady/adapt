// Package equidistant provides means for working with the Newton–Cotes grid.
//
// Each node in the grid is identified by a sequence of levels and orders. Such
// a sequence is encoded as a sequence of uint64s where each uint64 is
// (level|order<<6). Consequently, the maximal level is 2^6, and the maximal
// order is 2^58.
package equidistant

import (
	"github.com/ready-steady/linear"
)

const (
	levelMask = 0x3F
	levelSize = 6
	orderSize = 64 - levelSize
)

func index(levels []uint8, generate func(uint8) []uint64, nd uint) []uint64 {
	nn := uint(len(levels)) / nd

	cache := make(map[uint8][]uint64)

	indices1D := make([][]uint64, nd)
	indicesND := make([]uint64, 0)
	for i := uint(0); i < nn; i++ {
		for j := uint(0); j < nd; j++ {
			level := levels[i*nd+j]
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
