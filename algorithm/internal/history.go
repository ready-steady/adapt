package internal

// History is a book-keeper of indices.
type History struct {
	*Hash
	mapping map[string]uint
}

// NewHistory creates a book-keeper.
func NewHistory(ni uint) *History {
	return &History{
		Hash:    NewHash(ni),
		mapping: make(map[string]uint),
	}
}

// GetSet looks up the position of an index and, if not found, assigns one.
func (self *History) GetSet(index []uint64) (position uint, found bool) {
	key := self.Key(index)
	position, found = self.mapping[key]
	if !found {
		position = uint(len(self.mapping))
		self.mapping[key] = position
	}
	return
}
