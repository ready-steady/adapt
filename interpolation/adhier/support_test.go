package adhier

import (
	"testing"

	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestBalance(t *testing.T) {
	const (
		ni = 2
	)

	grid := newcot.NewOpen(ni)

	parents := []uint64{
		0 | 0<<32, 0 | 0<<32,
		1 | 0<<32, 0 | 0<<32,
		1 | 2<<32, 0 | 0<<32,
		0 | 0<<32, 1 | 0<<32,
		0 | 0<<32, 1 | 2<<32,
	}

	children := []uint64{
		1 | 0<<32, 1 | 0<<32,
		1 | 2<<32, 1 | 0<<32,
		0 | 0<<32, 2 | 0<<32,
		0 | 0<<32, 2 | 2<<32,
	}

	indices := append(parents, children...)

	find := func(index []uint64) bool {
		nn := len(indices) / ni

		for i := 0; i < nn; i++ {
			match := true
			for j := 0; j < ni; j++ {
				if indices[i*ni+j] != index[j] {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}

		return false
	}

	push := func(index []uint64) {
		indices = append(indices, index...)
	}

	balance(grid, children, ni, find, push)

	siblings := []uint64{
		1 | 0<<32, 1 | 2<<32,
		1 | 2<<32, 1 | 2<<32,
	}

	assert.Equal(indices, append(append(parents, children...), siblings...), t)
}
