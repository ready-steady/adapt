package internal

// Active is a book-keeper of active level indices.
type Active struct {
	Indices   []uint64      // Level indices considered so far
	Norms     []uint64      // Manhattan norms of level indices
	Positions map[uint]bool // Positions of active level indices

	ni uint

	history  *History
	forward  reference
	backward reference
}

type reference map[uint]uint

// NewActive creates a book-keeper.
func NewActive(ni uint) *Active {
	return &Active{
		ni: ni,
	}
}

// Next identifies, activates, and returns admissible level indices from the
// forward neighborhood of a level index.
func (self *Active) Next(k uint) (indices []uint64) {
	if self.Indices == nil {
		self.Indices, self.Norms = make([]uint64, 1*self.ni), []uint64{0}
		self.Positions = map[uint]bool{0: true}
		self.history = NewHistory(self.ni)
		self.forward, self.backward = make(reference), make(reference)
		return self.Indices
	}

	ni, no := self.ni, uint(len(self.Norms))
	forward, backward := self.forward, self.backward

	index, norm := self.Indices[k*ni:(k+1)*ni], self.Norms[k]

outer:
	for i, nn := uint(0), no; i < ni; i++ {
		index[i]++
		_, found := self.history.Get(index)
		index[i]--

		if found {
			// The forward neighbor in dimension i has already been considered.
			continue
		}

		newBackward := make(reference)
		for j := uint(0); j < ni; j++ {
			if index[j] == 0 {
				// The level of dimension j is the lowest possible, so there is
				// no backward neighbor.
				continue
			}
			if i == j {
				// The dimension is the one that we would like to bump up, so
				// the backward neighbor obviously exists.
				continue
			}
			l, found := forward[backward[k*ni+j]*ni+i]
			if !found {
				// The backward neighbor in dimension j has not been bumped up
				// in dimension i. So the candidate index has no backward
				// neighbor in dimension j and, hence, is not admissible.
				continue outer
			}
			newBackward[j] = l
		}
		newBackward[i] = k

		index[i]++
		self.Indices = append(self.Indices, index...)
		self.history.Set(index, 0)
		index[i]--

		self.Norms = append(self.Norms, norm+1)
		self.Positions[nn] = true

		for j, l := range newBackward {
			forward[l*ni+j] = nn
			backward[nn*ni+j] = l
		}

		nn++
	}

	indices = self.Indices[no*ni:]
	return
}

// Drop deactivates a level index.
func (self *Active) Drop(k uint) {
	delete(self.Positions, k)
}
