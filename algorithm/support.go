package algorithm

import (
	"github.com/ready-steady/adapt/grid"

	ainternal "github.com/ready-steady/adapt/algorithm/internal"
	rinternal "github.com/ready-steady/adapt/internal"
)

// Validate checks if an index set is admissible and contains no repetitions.
func Validate(indices []uint64, ni uint, parent grid.Parenter) bool {
	nn := uint(len(indices)) / ni

	hash := ainternal.NewHash(ni)
	mapping := make(map[string]bool)
	for i := uint(0); i < nn; i++ {
		key := hash.Key(indices[i*ni : (i+1)*ni])
		if _, ok := mapping[key]; ok {
			return false
		}
		mapping[key] = true
	}

	for i := uint(0); i < nn; i++ {
		root, found := true, false
		index := indices[i*ni : (i+1)*ni]
		for j := uint(0); !found && j < ni; j++ {
			level := rinternal.LEVEL_MASK & index[j]
			if level == 0 {
				continue
			} else {
				root = false
			}

			order := index[j] >> rinternal.LEVEL_SIZE
			plevel, porder := parent.Parent(level, order)

			index[j] = porder<<rinternal.LEVEL_SIZE | plevel
			_, found = mapping[hash.Key(index)]
			index[j] = order<<rinternal.LEVEL_SIZE | level
		}
		if !found && !root {
			return false
		}
	}

	return true
}
