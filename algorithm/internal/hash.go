package internal

import (
	"reflect"
	"unsafe"
)

const (
	sizeOfUint64 = 8
)

// Hash is a means of creating hash keys from indices.
type Hash struct {
	ni     uint
	bytes  []byte
	header *reflect.SliceHeader
}

// NewHash creates a Hash.
func NewHash(ni uint) *Hash {
	hash := &Hash{ni: ni, bytes: make([]byte, 0)}
	hash.header = (*reflect.SliceHeader)(unsafe.Pointer(&hash.bytes))
	hash.header.Cap = int(ni * sizeOfUint64)
	hash.header.Len = hash.header.Cap
	return hash
}

// Key creates a hash key from an index.
func (self *Hash) Key(index []uint64) string {
	self.header.Data = uintptr(((*reflect.SliceHeader)(unsafe.Pointer(&index))).Data)
	key := string(self.bytes)
	self.header.Data = 0
	return key
}
