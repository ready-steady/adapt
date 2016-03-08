package equidistant

import (
	"github.com/ready-steady/linear"
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
