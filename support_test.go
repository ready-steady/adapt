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
		10.0, 11.0, 12.0,
		13.0, 14.0, 15.0,
		16.0, 17.0, 18.0,
		19.0, 20.0, 21.0,
		22.0, 23.0, 24.0,
	}
	scores := []float64{-1.0, 0.0, 1.0, -1.0, 0.0}

	indices, surpluses, scores = compact(indices, surpluses, scores, ni, no, nn)

	assert.Equal(indices, []uint64{2, 3, 4, 5, 8, 9}, t)
	assert.Equal(surpluses, []float64{13.0, 14.0, 15.0, 16.0, 17.0, 18.0, 22.0, 23.0, 24.0}, t)
	assert.Equal(scores, []float64{0.0, 1.0, 0.0}, t)
}

func TestBalance(t *testing.T) {
	const (
		ni = 2
	)

	history := newHash(ni)
	grid := newcot.NewOpen(ni)

	parents := []uint64{
		0 | 0<<LEVEL_SIZE, 0 | 0<<LEVEL_SIZE,
		1 | 0<<LEVEL_SIZE, 0 | 0<<LEVEL_SIZE,
		1 | 2<<LEVEL_SIZE, 0 | 0<<LEVEL_SIZE,
		0 | 0<<LEVEL_SIZE, 1 | 0<<LEVEL_SIZE,
		0 | 0<<LEVEL_SIZE, 1 | 2<<LEVEL_SIZE,
	}

	children := []uint64{
		1 | 0<<LEVEL_SIZE, 1 | 0<<LEVEL_SIZE,
		1 | 2<<LEVEL_SIZE, 1 | 0<<LEVEL_SIZE,
		0 | 0<<LEVEL_SIZE, 2 | 0<<LEVEL_SIZE,
		0 | 0<<LEVEL_SIZE, 2 | 2<<LEVEL_SIZE,
	}

	siblings := []uint64{
		1 | 0<<LEVEL_SIZE, 1 | 2<<LEVEL_SIZE,
		1 | 2<<LEVEL_SIZE, 1 | 2<<LEVEL_SIZE,
	}

	for i := 0; i < len(parents); i += ni {
		history.push(parents[i : i+ni])
	}

	for i := 0; i < len(children); i += ni {
		history.push(children[i : i+ni])
	}

	assert.Equal(balance(grid, history, children), siblings, t)
}
