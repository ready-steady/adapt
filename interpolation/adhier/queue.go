package adhier

import (
	"math"
)

type queue interface {
	process([]uint64, []float64) ([]uint64, []bool)
}

type fakeQueue struct {
	ni  int
	min uint
	max uint
}

type realQueue struct {
	ni   int
	nn   int
	min  uint
	max  uint
	rate float64
	root *element
}

type element struct {
	total    float64
	index    []uint64
	selector []bool
	next     *element
}

func newQueue(ni, min, max uint, rate float64) queue {
	if rate <= 0 || rate >= 1 {
		return &fakeQueue{
			ni:  int(ni),
			min: min,
			max: max,
		}
	} else {
		return &realQueue{
			ni:   int(ni),
			min:  min,
			max:  max,
			rate: rate,
		}
	}
}

func (q *fakeQueue) process(indices []uint64, scores []float64) ([]uint64, []bool) {
	ni := q.ni
	nn := len(indices) / ni

	min, max := q.min, q.max

	selectors := make([]bool, nn*ni)

	ns := 0

	for i, k := 0, 0; i < nn; i++ {
		index := indices[k*ni : (k+1)*ni]
		score := scores[i*ni : (i+1)*ni]
		selector := selectors[ns*ni : (ns+1)*ni]

		skip := true
		level := uint(0)
		for j := 0; j < ni; j++ {
			level += uint(0xFFFFFFFF & index[j])
			selector[j] = score[j] > 0
			if selector[j] {
				skip = false
			}
		}
		if level < min {
			for j := 0; j < ni; j++ {
				selector[j] = true
			}
			skip = false
		}
		if level >= max || skip {
			k++
			continue
		}

		if k > ns {
			copy(indices[ns*ni:], indices[k*ni:])
			k = ns
		}

		k++
		ns++
	}

	return indices[:ns*ni], selectors[:ns*ni]
}

func (q *realQueue) process(indices []uint64, scores []float64) ([]uint64, []bool) {
	q.push(indices, scores)
	return q.pull()
}

func (q *realQueue) push(indices []uint64, scores []float64) {
	ni := q.ni
	nn, ns := len(indices)/ni, 0

	min, max := q.min, q.max

	selectors := make([]bool, nn*ni)

	for i := 0; i < nn; i++ {
		index := indices[i*ni : (i+1)*ni]
		score := scores[i*ni : (i+1)*ni]
		selector := selectors[i*ni : (i+1)*ni]

		total := 0.0
		level := uint(0)
		for j := 0; j < ni; j++ {
			level += uint(0xFFFFFFFF & index[j])
			selector[j] = score[j] > 0
			if selector[j] {
				selector[j] = true
				total += score[j]
			}
		}
		if level < min {
			for j := 0; j < ni; j++ {
				selector[j] = true
			}
			total = math.Inf(1)
		}
		if level >= max || total <= 0 {
			continue
		}

		candidate := &element{
			total:    total,
			index:    index,
			selector: selector,
		}

		var previous, current *element = nil, q.root
		for {
			if current == nil || current.total < total {
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

func (q *realQueue) pull() ([]uint64, []bool) {
	ni := q.ni
	nn := int(q.rate * float64(q.nn))
	if nn == 0 && q.nn > 0 {
		nn++
	}

	indices := make([]uint64, nn*ni)
	selectors := make([]bool, nn*ni)

	current := q.root
	for i := 0; i < nn; i++ {
		copy(indices[i*ni:], current.index)
		copy(selectors[i*ni:], current.selector)
		current = current.next
	}

	q.root = current
	q.nn -= nn

	return indices, selectors
}
