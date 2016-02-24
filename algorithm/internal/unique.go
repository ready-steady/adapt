package internal

import (
	"reflect"
	"unsafe"
)

const (
	sizeOfUint64 = 8
)

// Unique is a book-keeper of unique indices.
type Unique struct {
	ni      uint
	mapping map[string]bool
}

// NewUnique creates a book-keeper.
func NewUnique(ni uint) *Unique {
	return &Unique{
		ni:      ni,
		mapping: make(map[string]bool),
	}
}

// Distil eliminates in place the indices that have already been seen.
func (self *Unique) Distil(indices []uint64) []uint64 {
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