package internal

// Active is a book-keeper of active level indices.
type Active struct {
	Indices   []uint64      // Level indices considered so far
	Positions map[uint]bool // Positions of active level indices

	ni   uint
	nn   uint
	lmax uint

	forward  reference
	backward reference
}

type reference map[uint]uint

// NewActive creates a book-keeper.
func NewActive(ni, lmax uint) *Active {
	return &Active{
		ni:   ni,
		lmax: lmax,
	}
}

// Next identifies, activates, and returns admissible level indices from the
// forward neighborhood of a level index.
func (self *Active) Next(k uint) (indices []uint64) {
	if self.Indices == nil {
		self.Indices, self.Positions = make([]uint64, 1*self.ni), map[uint]bool{0: true}
		self.nn, self.forward, self.backward = 1, make(reference), make(reference)
		return self.Indices
	}

	ni, nn := self.ni, self.nn
	positions, forward, backward := self.Positions, self.forward, self.backward

	index := self.Indices[k*ni : (k+1)*ni]

outer:
	for i := uint(0); i < ni; i++ {
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

// Drop deactivates a level index.
func (self *Active) Drop(k uint) {
	delete(self.Positions, k)
}
