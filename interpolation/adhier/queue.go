package adhier

import (
	"math"
)

type queue struct {
	ni int
	nn int

	lmin uint
	lmax uint
	rate float64

	root *element
}

type element struct {
	index []uint64
	score float64
	next  *element
}

func newQueue(ni uint, c *Config) *queue {
	return &queue{
		ni: int(ni),
		nn: 1,

		lmin: c.MinLevel,
		lmax: c.MaxLevel,
		rate: c.Rate,

		root: &element{
			index: make([]uint64, ni),
			score: math.Inf(1),
		},
	}
}

func (q *queue) push(indices []uint64, scores []float64) {
	ni := q.ni
	nn, ns := len(indices)/ni, 0

	lmin, lmax := q.lmin, q.lmax

	for i := 0; i < nn; i++ {
		index := indices[i*ni : (i+1)*ni]
		score := scores[i]

		level := uint(0)
		for j := 0; j < ni; j++ {
			level += uint(0xFFFFFFFF & index[j])
		}
		if level < lmin {
			score = math.Inf(1)
		}
		if level >= lmax || score <= 0 {
			continue
		}

		candidate := &element{
			index: index,
			score: score,
		}

		var previous, current *element = nil, q.root
		for {
			if current == nil || current.score < score {
				if previous == nil {
					q.root = candidate
				} else {
					previous.next = candidate
				}
				candidate.next = current
				break
			}
			previous, current = current, current.next
		}

		ns++
	}

	q.nn += ns
}

func (q *queue) pull() []uint64 {
	ni := q.ni
	nn := int(math.Ceil(q.rate * float64(q.nn)))

	indices := make([]uint64, nn*ni)

	current := q.root
	for i := 0; i < nn; i++ {
		copy(indices[i*ni:], current.index)
		current = current.next
	}

	q.root = current
	q.nn -= nn

	return indices
}
