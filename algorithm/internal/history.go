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

// Get looks up the value of an index.
func (self *History) Get(index []uint64) (uint, bool) {
	current, found := self.mapping[self.Key(index)]
	return current, found
}

// GetSet looks up the value of an index and, if not found, assigns one.
func (self *History) GetSet(index []uint64, value uint) (uint, bool) {
	key := self.Key(index)
	current, found := self.mapping[key]
	if !found {
		self.mapping[key] = value
	}
	return current, found
}

// Set assigns a value to an index.
func (self *History) Set(index []uint64, value uint) {
	self.mapping[self.Key(index)] = value
}
