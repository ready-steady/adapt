package internal

// Active is a book-keeper of active level indices.
type Active struct {
	Indices   []uint64      // Level indices considered so far
	Norms     []uint64      // Manhattan norms of level indices
	Positions map[uint]bool // Positions of active level indices

	ni uint

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
		self.forward, self.backward = make(reference), make(reference)
		return self.Indices
	}

	ni, no := self.ni, uint(len(self.Norms))
	positions, forward, backward := self.Positions, self.forward, self.backward

	index, norm := self.Indices[k*ni:(k+1)*ni], self.Norms[k]

outer:
	for i, nn := uint(0), no; i < ni; i++ {
		newBackward := make(reference)
		for j := uint(0); j < ni; j++ {
			if index[j] == 0 {
				// It is the lowest possible level, so there is no ancestor.
				continue
			}
			if i == j {
				// It is the dimension along which we would like move forward,
				// so `index` is the ancestor, and it obviously exists.
				continue
			}
			l, ok := forward[backward[k*ni+j]*ni+i]
			if !ok {
				// The ancestor in the jth dimension has not been bumped up in
				// the ith dimension. So the candidate index has no ancestor in
				// the jth dimension and, hence, is not admissible.
				continue outer
			}
			if positions[l] {
				// The candidate index has an ancestor in the jth dimension;
				// however, this ancestor is still active, and, hence, the
				// candidate index is not admissible.
				continue outer
			}
			newBackward[j] = l
		}
		newBackward[i] = k

		self.Indices = append(self.Indices, index...)
		self.Indices[nn*ni+i]++
		self.Norms = append(self.Norms, norm+1)

		for j, l := range newBackward {
			forward[l*ni+j] = nn
			backward[nn*ni+j] = l
		}
		positions[nn] = true

		nn++
	}

	indices = self.Indices[no*ni:]
	return
}

// Drop deactivates a level index.
func (self *Active) Drop(k uint) {
	delete(self.Positions, k)
}
