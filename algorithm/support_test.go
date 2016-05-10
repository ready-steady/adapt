package algorithm

import (
	"testing"

	"github.com/ready-steady/adapt/grid/equidistant"
	"github.com/ready-steady/adapt/internal"
	"github.com/ready-steady/assert"
)

func TestValidate(t *testing.T) {
	const (
		ni = 2
	)

	cases := []struct {
		levels []uint64
		orders []uint64
		result bool
	}{
		{
			[]uint64{
				0, 0,
				0, 1,
				1, 0,
				1, 1,
				1, 2,
			},
			[]uint64{
				0, 0,
				0, 2,
				2, 0,
				2, 2,
				2, 3,
			},
			true,
		},
		{
			[]uint64{
				0, 0,
				0, 1,
				1, 0,
				1, 1,
				1, 1,
				1, 2,
			},
			[]uint64{
				0, 0,
				0, 2,
				2, 0,
				2, 2,
				2, 2,
				2, 3,
			},
			false,
		},
		{
			[]uint64{
				0, 0,
				0, 1,
				1, 0,
				1, 1,
				1, 2,
			},
			[]uint64{
				0, 0,
				0, 2,
				2, 0,
				2, 2,
				2, 1,
			},
			false,
		},
		{
			[]uint64{
				0, 0,
				0, 1,
				1, 0,
				1, 1,
				2, 2,
			},
			[]uint64{
				0, 0,
				0, 2,
				2, 0,
				2, 2,
				3, 3,
			},
			false,
		},
	}

	for _, c := range cases {
		indices := internal.Compose(c.levels, c.orders)
		assert.Equal(Validate(indices, ni, equidistant.ClosedParent), c.result, t)
	}
}
