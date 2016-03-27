package internal

import (
	"testing"

	"github.com/ready-steady/adapt/internal"
	"github.com/ready-steady/assert"
)

func TestIsAdmissible(t *testing.T) {
	const (
		ni = 2
	)

	var indices []uint64

	indices = []uint64{
		0, 0,
		0, 1,
		1, 0,
		1, 1,
	}
	assert.Equal(IsAdmissible(indices, ni), true, t)

	indices = []uint64{
		0, 0,
		0, 1,
		1, 0,
		1, 1,
		1, 2,
	}
	assert.Equal(IsAdmissible(indices, ni), false, t)
}

func TestIsUnique(t *testing.T) {
	const (
		ni = 2
	)

	var levels, orders []uint64

	levels = []uint64{
		1, 2,
		3, 4,
		5, 6,
	}
	orders = []uint64{
		6, 5,
		4, 3,
		2, 1,
	}
	assert.Equal(IsUnique(internal.Compose(levels, orders), ni), true, t)

	levels = []uint64{
		1, 2,
		3, 4,
		5, 6,
		1, 2,
	}
	orders = []uint64{
		6, 5,
		4, 3,
		2, 1,
		6, 5,
	}
	assert.Equal(IsUnique(internal.Compose(levels, orders), ni), false, t)
}

func TestLevelize(t *testing.T) {
	const (
		ni = 3
	)

	indices := []uint64{
		1 | 1<<internal.LEVEL_SIZE, 4 | 1<<internal.LEVEL_SIZE, 7 | 1<<internal.LEVEL_SIZE,
		2 | 2<<internal.LEVEL_SIZE, 5 | 2<<internal.LEVEL_SIZE, 8 | 2<<internal.LEVEL_SIZE,
		3 | 3<<internal.LEVEL_SIZE, 6 | 3<<internal.LEVEL_SIZE, 9 | 3<<internal.LEVEL_SIZE,
	}

	assert.Equal(Levelize(indices, ni), []uint64{12, 15, 18}, t)
}
