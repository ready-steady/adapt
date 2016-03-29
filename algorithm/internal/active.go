package internal

// Active is a structure for keeping track of active level indices.
type Active struct {
	Indices   []uint64      // Level indices considered so far
	Positions map[uint]bool // Positions of active level indices

	ni uint

	history  *History
	forward  reference
	backward reference
}

type reference map[uint]uint

// NewActive creates an Active.
func NewActive(ni uint) *Active {
	return &Active{
		ni: ni,
	}
}

// First returns the initial level indices.
func (self *Active) First() []uint64 {
	self.Indices = make([]uint64, 1*self.ni)
	self.Positions = map[uint]bool{0: true}
	self.history = NewHistory(self.ni)
	self.forward, self.backward = make(reference), make(reference)

	return self.Indices
}

// Next returns admissible forward neighbors of a level index.
func (self *Active) Next(k uint) (indices []uint64) {
	ni := self.ni
	no := uint(len(self.Indices)) / ni

	forward, backward := self.forward, self.backward
	index := self.Indices[k*ni : (k+1)*ni]

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
