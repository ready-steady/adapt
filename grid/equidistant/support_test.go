package equidistant

import (
	"github.com/ready-steady/adapt/internal"
)

func compose(levels []uint64, orders []uint64) []uint64 {
	indices := make([]uint64, len(levels))

	for i := range levels {
		indices[i] = uint64(levels[i]) | uint64(orders[i])<<internal.LEVEL_SIZE
	}

	return indices
}
