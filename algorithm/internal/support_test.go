package internal

import (
	"testing"

	"github.com/ready-steady/adapt/internal"
	"github.com/ready-steady/assert"
)

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
