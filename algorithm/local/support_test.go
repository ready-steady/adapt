package local

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestLevelize(t *testing.T) {
	const (
		ni = 3
	)

	indices := []uint64{
		1 | 1<<LEVEL_SIZE, 4 | 1<<LEVEL_SIZE, 7 | 1<<LEVEL_SIZE,
		2 | 2<<LEVEL_SIZE, 5 | 2<<LEVEL_SIZE, 8 | 2<<LEVEL_SIZE,
		3 | 3<<LEVEL_SIZE, 6 | 3<<LEVEL_SIZE, 9 | 3<<LEVEL_SIZE,
	}

	assert.Equal(levelize(indices, ni), []uint{12, 15, 18}, t)
}
