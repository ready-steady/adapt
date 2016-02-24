package external

// Active is a book-keeper of active level indices.
type Active struct {
	// All level indices considered so far.
	Indices []uint64
	// The positions of active level indices.
	Positions Set

	ni   uint
	nn   uint
	lmax uint
	imax uint

	forward  reference
	backward reference
}

type reference map[uint]uint

// NewActive creates a book-keeper.
func NewActive(ni, lmax, imax uint) *Active {
	return &Active{
		Indices:   make([]uint64, 1*ni),
		Positions: Set{0: true},

		ni:   ni,
		nn:   1,
		lmax: lmax,
		imax: imax,

		forward:  make(reference),
		backward: make(reference),
	}
}

// Advance identifies, activates, and returns admissible indices from the
// forward neighborhood of a level index.
func (self *Active) Advance(k uint) (indices []uint64) {
	ni, nn := self.ni, self.nn
	positions, forward, backward := self.Positions, self.forward, self.backward

	index := self.Indices[k*ni : (k+1)*ni]

outer:
	for i := uint(0); i < ni && nn < self.imax; i++ {
		if index[i] >= uint64(self.lmax) {
			continue
		}

		newBackward := make(reference)
		for j := uint(0); j < ni; j++ {
			if i == j || index[j] == 0 {
				continue
			}
			if l, ok := forward[backward[k*ni+j]*ni+i]; !ok || positions[l] {
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

		positions[nn] = true

		nn++
	}

	indices = self.Indices[self.nn*ni:]
	self.nn = nn

	return
}

// Forget deactivates a level index.
func (self *Active) Forget(k uint) {
	delete(self.Positions, k)
}
