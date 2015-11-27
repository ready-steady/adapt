package adapt

import (
	"testing"

	"github.com/ready-steady/adapt/grid/newcot"
	"github.com/ready-steady/assert"
)

func TestCompact(t *testing.T) {
	const (
		ni = 2
		no = 3
		nn = 5
	)

	indices := []uint64{
		0, 1,
		2, 3,
		4, 5,
		6, 7,
		8, 9,
	}
	surpluses := []float64{
		10, 11, 12,
		13, 14, 15,
		16, 17, 18,
		19, 20, 21,
		22, 23, 24,
	}
	scores := []float64{-1, 0, 1, -1, 0}

	indices, surpluses, scores = compact(indices, surpluses, scores, ni, no, nn)

	assert.Equal(indices, []uint64{2, 3, 4, 5, 8, 9}, t)
	assert.Equal(surpluses, []float64{13, 14, 15, 16, 17, 18, 22, 23, 24}, t)
	assert.Equal(scores, []float64{0, 1, 0}, t)
}

func TestBalance(t *testing.T) {
	const (
		ni = 2
	)

	history := newHash(ni)
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

	siblings := []uint64{
		1 | 0<<32, 1 | 2<<32,
		1 | 2<<32, 1 | 2<<32,
	}

	for i := 0; i < len(parents); i += ni {
		history.push(parents[i : i+ni])
	}

	for i := 0; i < len(children); i += ni {
		history.push(children[i : i+ni])
	}

	assert.Equal(balance(grid, history, children), siblings, t)
}
