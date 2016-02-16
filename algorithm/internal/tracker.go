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

	ni uint
	nn uint

	lmax uint
	imax uint
	rate float64

	norms  []uint64
	scores []float64

	forward  Reference
	backward Reference
}

// NewTracker returns a tracker of active level indices.
func NewTracker(ni, lmax, imax uint, rate float64) *Tracker {
	return &Tracker{
		ni: ni,

		lmax: lmax,
		imax: imax,
		rate: rate,

		forward:  make(Reference),
		backward: make(Reference),
	}
}

// Pull fetches the next level index.
func (self *Tracker) Pull() []uint64 {
	if self.Active == nil {
		return self.pullFirst()
	} else {
		return self.pullSubsequent()
	}
}

// Push updates the tracker with the score of the previously pulled index.
func (self *Tracker) Push(score float64) {
	self.scores = append(self.scores, score)
}

func (self *Tracker) pullFirst() []uint64 {
	self.Indices = make([]uint64, 1*self.ni)
	self.Active = make(Set)
	self.Active[0] = true
	self.nn = 1
	self.norms = make([]uint64, 1)
	return self.Indices
}

func (self *Tracker) pullSubsequent() (indices []uint64) {
	ni, nn := self.ni, self.nn
	active, forward, backward := self.Active, self.forward, self.backward

	min, k := minUint64Set(self.norms, active)
	max := MaxUint64(self.norms)
	if float64(min) > (1.0-self.rate)*float64(max) {
		_, k = maxFloat64Set(self.scores, active)
	}
	delete(active, k)

	index, norm := self.Indices[k*ni:(k+1)*ni], self.norms[k]+1

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
		self.norms = append(self.norms, norm)

		active[nn] = true

		nn++
	}

	indices = self.Indices[self.nn*ni:]
	self.nn = nn

	return
}
