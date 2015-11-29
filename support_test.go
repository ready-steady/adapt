package adapt

import (
	"testing"

	"github.com/ready-steady/adapt/grid/newcot"
	"github.com/ready-steady/assert"
)

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
