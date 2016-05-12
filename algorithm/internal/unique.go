package internal

// Unique is a structure for keeping track of unique indices.
type Unique struct {
	*History
}

// NewUnique creates a Unique.
func NewUnique(ni uint) *Unique {
	return &Unique{NewHistory(ni)}
}

// Distil eliminates the indices that have already been seen.
func (self *Unique) Distil(indices []uint64) []uint64 {
	ni := self.ni
	nn := uint(len(indices)) / ni
	k, unique := uint(0), []uint64{}
	for i := uint(0); i < nn; i++ {
		if _, found := self.GetSet(indices[i*ni:(i+1)*ni], 0); found {
			unique = append(unique, indices[k*ni:i*ni]...)
			k = i + 1
		}
	}
	return append(unique, indices[k*ni:]...)
}
