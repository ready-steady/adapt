package internal

import (
	"math"
)

var (
	infinity = math.Inf(1.0)
)

// Reference is a relation on ordered elements.
type Reference map[uint]uint

// Set is a subset of ordered elements.
type Set map[uint]bool

// Tracker is a structure for keeping track of active level indices.
type Tracker struct {
	// All the level indices considered so far.
	Indices []uint64
	// The positions of the active level indices.
	Active Set

	ni   uint
	nn   uint
	lmax uint
	imax uint

	forward  Reference
	backward Reference

	initialized bool
}

// NewTracker returns a tracker of active level indices.
func NewTracker(ni, lmax, imax uint) *Tracker {
	return &Tracker{
		Indices: make([]uint64, 1*ni),
		Active:  Set{0: true},

		ni:   ni,
		nn:   1,
		lmax: lmax,
		imax: imax,

		forward:  make(Reference),
		backward: make(Reference),
	}
}

// Pull fetches the next level index.
func (self *Tracker) Pull(k uint) (indices []uint64) {
	if !self.initialized {
		self.initialized = true
		indices = self.Indices
		return
	}

	ni, nn := self.ni, self.nn
	active, forward, backward := self.Active, self.forward, self.backward
	index := self.Indices[k*ni : (k+1)*ni]

	delete(active, k)

outer:
	for i := uint(0); i < ni && nn < self.imax; i++ {
		if index[i] >= uint64(self.lmax) {
			continue
		}

		newBackward := make(Reference)
		for j := uint(0); j < ni; j++ {
			if i == j || index[j] == 0 {
				continue
			}
			if l, ok := forward[backward[k*ni+j]*ni+i]; !ok || active[l] {
				continue outer
			} else {
				newBackward[j] = l
			}
		}
		newBackward[i] = k
		for j, l := range newBackward {
			forward[l*ni+j] = nn
			backward[nn*ni+j] = l
		}

		self.Indices = append(self.Indices, index...)
		self.Indices[nn*ni+i]++

		active[nn] = true

		nn++
	}

	indices = self.Indices[self.nn*ni:]
	self.nn = nn

	return
}
