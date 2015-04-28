package adapt

import (
	"math"
)

type queue struct {
	ni int
	nn int

	lnow uint
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

		lmin: c.MinLevel,
		lmax: c.MaxLevel,

		rate: c.Rate,
	}
}

func (q *queue) push(indices []uint64, scores []float64) []bool {
	ni := q.ni
	nn, nq := len(indices)/ni, 0

	accept := make([]bool, nn)

	lnow, lmin, lmax := q.lnow, q.lmin, q.lmax

	for i := 0; i < nn; i++ {
		score := scores[i]
		accept[i] = score >= 0

		if !accept[i] {
			continue // should be discarded
		}

		index := indices[i*ni : (i+1)*ni]

		l := uint(0)
		for j := 0; j < ni; j++ {
			l += uint(0xFFFFFFFF & index[j])
		}
		if l > lnow {
			lnow = l
		}
		if l >= lmin && (score == 0 || l == lmax) {
			continue // should not be refined
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

		nq++
	}

	q.nn += nq
	q.lnow = lnow

	return accept
}

func (q *queue) pull() []uint64 {
	ni, nn := q.ni, q.nn
	if q.lnow >= q.lmin {
		nn = int(math.Ceil(q.rate * float64(nn)))
	}

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
