package external

// Progress is a structure that contains information about the interpolation
// process.
type Progress struct {
	More uint // Number of nodes to be evaluated
	Done uint // Number of nodes evaluated so far
}

// NewProgress returns a empty progress structure.
func NewProgress() *Progress {
	return &Progress{}
}

// Push takes into account new indices.
func (self *Progress) Push(indices []uint64, ni uint) {
	self.Done += self.More
	self.More = uint(len(indices)) / ni
}
