package internal

import (
	"reflect"
	"unsafe"
)

const (
	sizeOfUint64 = 8
)

// Hash is a means of converting indices into strings.
type Hash struct {
	bytes  []byte
	header *reflect.SliceHeader
}

// NewHash creates a hash.
func NewHash(ni uint) *Hash {
	hash := &Hash{bytes: make([]byte, 0)}
	hash.header = (*reflect.SliceHeader)(unsafe.Pointer(&hash.bytes))
	hash.header.Cap = int(ni * sizeOfUint64)
	hash.header.Len = hash.header.Cap
	return hash
}

// Key converts an index into a string.
func (self *Hash) Key(index []uint64) string {
	self.header.Data = uintptr(((*reflect.SliceHeader)(unsafe.Pointer(&index))).Data)
	key := string(self.bytes)
	self.header.Data = 0
	return key
}
