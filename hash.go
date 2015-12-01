package adapt

import (
	"reflect"
	"unsafe"
)

const (
	sizeOfUint64 = 8
)

type hash struct {
	ni      uint
	mapping map[string]bool
}

func newHash(ni uint) *hash {
	return &hash{
		ni:      ni,
		mapping: make(map[string]bool),
	}
}

func (self *hash) find(index []uint64) bool {
	header := reflect.StringHeader{
		Data: (*reflect.SliceHeader)(unsafe.Pointer(&index)).Data,
		Len:  int(self.ni) * sizeOfUint64,
	}
	_, ok := self.mapping[*(*string)(unsafe.Pointer(&header))]
	return ok
}

func (self *hash) push(index []uint64) {
	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Data = ((*reflect.SliceHeader)(unsafe.Pointer(&index))).Data
	header.Cap = int(self.ni) * sizeOfUint64
	header.Len = header.Cap
	self.mapping[string(bytes)] = true
}

func (self *hash) unseen(indices []uint64) []uint64 {
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
