package adapt

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestQueueCompress(t *testing.T) {
	fixture := func() *queue {
		return &queue{
			ni:   2,
			no:   3,
			nn:   5,
			lmax: 20,

			Indices: []uint64{
				1, 2,
				3, 4,
				5, 6,
				7, 8,
				9, 10,
			},
			Nodes: []float64{
				1.1, 1.2,
				2.1, 2.2,
				3.1, 3.2,
				4.1, 4.2,
				5.1, 5.2,
			},
			Values: []float64{
				1.1, 1.2, 1.3,
				2.1, 2.2, 2.3,
				3.1, 3.2, 3.3,
				4.1, 4.2, 4.3,
				5.1, 5.2, 5.3,
			},
			Scores: []float64{0.0, 1.0, 2.0, -1.0, 3.0},
		}
	}

	q := fixture()
	q.compress(0)

	assert.Equal(q, &queue{
		ni:   2,
		no:   3,
		nn:   3,
		lmax: 20,

		Indices: []uint64{
			3, 4,
			5, 6,
			9, 10,
		},
		Nodes: []float64{
			2.1, 2.2,
			3.1, 3.2,
			5.1, 5.2,
		},
		Values: []float64{
			2.1, 2.2, 2.3,
			3.1, 3.2, 3.3,
			5.1, 5.2, 5.3,
		},
		Scores: []float64{1.0, 2.0, 3.0},
	}, t)

	q = fixture()
	q.compress(2)

	assert.Equal(q, &queue{
		ni:   2,
		no:   3,
		nn:   4,
		lmax: 20,

		Indices: []uint64{
			1, 2,
			3, 4,
			5, 6,
			9, 10,
		},
		Nodes: []float64{
			1.1, 1.2,
			2.1, 2.2,
			3.1, 3.2,
			5.1, 5.2,
		},
		Values: []float64{
			1.1, 1.2, 1.3,
			2.1, 2.2, 2.3,
			3.1, 3.2, 3.3,
			5.1, 5.2, 5.3,
		},
		Scores: []float64{0.0, 1.0, 2.0, 3.0},
	}, t)
}

func TestQueuePush(t *testing.T) {
	fixture := func() *queue {
		return &queue{
			ni:   2,
			no:   3,
			nn:   3,
			lmax: 20,

			Indices: []uint64{
				3, 4,
				5, 6,
				9, 10,
			},
			Nodes: []float64{
				2.1, 2.2,
				3.1, 3.2,
				5.1, 5.2,
			},
			Values: []float64{
				2.1, 2.2, 2.3,
				3.1, 3.2, 3.3,
				5.1, 5.2, 5.3,
			},
			Scores: []float64{1.0, 2.0, 3.0},
		}
	}

	q := fixture()
	q.push([]uint64{
		1, 2,
		7, 8,
		11, 12,
	}, []float64{
		1.1, 1.2,
		4.1, 4.2,
		6.1, 6.2,
	}, []float64{
		1.1, 1.2, 1.3,
		4.1, 4.2, 4.3,
		6.1, 6.2, 6.3,
	}, []float64{
		1.0,
		0.0,
		2.0,
	})

	assert.Equal(q, &queue{
		ni:   2,
		no:   3,
		nn:   4,
		lmax: 20,

		Indices: []uint64{
			3, 4,
			5, 6,
			9, 10,
			1, 2,
		},
		Nodes: []float64{
			2.1, 2.2,
			3.1, 3.2,
			5.1, 5.2,
			1.1, 1.2,
		},
		Values: []float64{
			2.1, 2.2, 2.3,
			3.1, 3.2, 3.3,
			5.1, 5.2, 5.3,
			1.1, 1.2, 1.3,
		},
		Scores: []float64{1.0, 2.0, 3.0, 1.0},
	}, t)
}

func TestQueuePull(t *testing.T) {
	fixture := func(rate float64) *queue {
		return &queue{
			ni:   2,
			no:   3,
			nn:   3,
			lmax: 20,
			rate: rate,

			Indices: []uint64{
				1, 2,
				3, 4,
				5, 6,
			},
			Nodes: []float64{
				1.1, 1.2,
				2.1, 2.2,
				3.1, 3.2,
			},
			Values: []float64{
				1.1, 1.2, 1.3,
				2.1, 2.2, 2.3,
				3.1, 3.2, 3.3,
			},
			Scores: []float64{3.0, 1.0, 2.0},
		}
	}

	q := fixture(0.5)

	assert.Equal(q.pull(), []uint64{5, 6, 1, 2}, t)
	assert.Equal(q, &queue{
		ni:   2,
		no:   3,
		nn:   1,
		lmax: 20,
		rate: 0.5,

		Indices: []uint64{3, 4},
		Nodes:   []float64{2.1, 2.2},
		Values:  []float64{2.1, 2.2, 2.3},
		Scores:  []float64{1.0},
	}, t)

	assert.Equal(q.pull(), []uint64{3, 4}, t)
	assert.Equal(q, &queue{
		ni:   2,
		no:   3,
		lmax: 20,
		rate: 0.5,

		Indices: []uint64{},
		Nodes:   []float64{},
		Values:  []float64{},
		Scores:  []float64{},
	}, t)
}

func TestQueueUpdate(t *testing.T) {
	fixture := func() *queue {
		return &queue{
			ni:   2,
			no:   3,
			nn:   3,
			lmax: 20,

			Indices: []uint64{
				3, 4,
				5, 6,
				9, 10,
			},
			Nodes: []float64{
				2.1, 2.2,
				3.1, 3.2,
				5.1, 5.2,
			},
			Values: []float64{
				2.1, 2.2, 2.3,
				3.1, 3.2, 3.3,
				5.1, 5.2, 5.3,
			},
			Scores: []float64{1.0, 2.0, 3.0},
		}
	}

	q := fixture()
	q.update([]float64{4.0, 0.0, 5.0})

	assert.Equal(q, &queue{
		ni:   2,
		no:   3,
		nn:   2,
		lmax: 20,

		Indices: []uint64{
			3, 4,
			9, 10,
		},
		Nodes: []float64{
			2.1, 2.2,
			5.1, 5.2,
		},
		Values: []float64{
			2.1, 2.2, 2.3,
			5.1, 5.2, 5.3,
		},
		Scores: []float64{4.0, 5.0},
	}, t)

	q = fixture()
	q.update([]float64{0.0, 0.0, 0.0})

	assert.Equal(q, &queue{
		ni:   2,
		no:   3,
		nn:   0,
		lmax: 20,

		Indices: []uint64{},
		Nodes:   []float64{},
		Values:  []float64{},
		Scores:  []float64{},
	}, t)
}
