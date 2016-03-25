package internal

// Unique is a book-keeper of indices.
type Unique struct {
	*History
}

// NewUnique creates a book-keeper.
func NewUnique(ni uint) *Unique {
	return &Unique{NewHistory(ni)}
}

// Distil eliminates in place the indices that have already been seen.
func (self *Unique) Distil(indices []uint64) []uint64 {
	ni := self.ni
	nn := uint(len(indices)) / ni
	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if _, found := self.GetSet(indices[j*ni:(j+1)*ni], 0); found {
			j++
			continue
		}
		if j > na {
			copy(indices[na*ni:], indices[j*ni:ne*ni])
			ne -= j - na
			j = na
		}
		na++
		j++
	}
	return indices[:na*ni]
}
