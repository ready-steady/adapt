package internal

import (
	"reflect"
	"unsafe"
)

const (
	sizeOfUint64 = 8
)

// Hash is a structure for keeping track of unique indices.
type Hash struct {
	ni      uint
	mapping map[string]bool
}

func NewHash(ni uint) *Hash {
	return &Hash{
		ni:      ni,
		mapping: make(map[string]bool),
	}
}

func (self *Hash) Filter(indices []uint64) []uint64 {
	ni := self.ni
	nn := uint(len(indices)) / ni
	nb := ni * sizeOfUint64

	var bytes []byte

	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Cap = int(nb)
	header.Len = int(nb)

	offset := ((*reflect.SliceHeader)(unsafe.Pointer(&indices))).Data

	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		header.Data = offset + uintptr(j*nb)
		key := string(bytes)
		if _, ok := self.mapping[key]; ok {
			j++
			continue
		}
		self.mapping[key] = true
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
