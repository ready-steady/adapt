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

func (self *queue) push(indices []uint64, scores []float64) {
	ni := self.ni
	nn, nq := len(indices)/ni, 0

	lnow, lmin, lmax := self.lnow, self.lmin, self.lmax

	for i := 0; i < nn; i++ {
		index := indices[i*ni : (i+1)*ni]
		score := scores[i]

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

		var previous, current *element = nil, self.root
		for {
			if current == nil || current.score < score {
				if previous == nil {
					self.root = candidate
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

	self.nn += nq
	self.lnow = lnow
}

func (self *queue) pull() []uint64 {
	ni, nn := self.ni, self.nn
	if self.lnow >= self.lmin {
		nn = int(math.Ceil(self.rate * float64(nn)))
	}

	indices := make([]uint64, nn*ni)

	current := self.root
	for i := 0; i < nn; i++ {
		copy(indices[i*ni:], current.index)
		current = current.next
	}

	self.root = current
	self.nn -= nn

	return indices
}
