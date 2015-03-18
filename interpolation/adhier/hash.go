package adhier

import (
	"reflect"
	"unsafe"
)

type hash struct {
	ni      uint
	mapping map[string]bool
}

func newHash(ni uint, capacity uint) *hash {
	return &hash{
		ni:      ni,
		mapping: make(map[string]bool, capacity),
	}
}

func (h *hash) unique(indices []uint64) []uint64 {
	const (
		sizeOfUint64 = 8
	)

	ni := h.ni
	nn := uint(len(indices)) / ni
	nb := uintptr(ni) * sizeOfUint64

	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Data = ((*reflect.SliceHeader)(unsafe.Pointer(&indices))).Data
	header.Cap = int(nb)
	header.Len = int(nb)

	newIndices := make([]uint64, nn*ni)

	ns := uint(0)

	for i := uint(0); i < nn; i++ {
		key := string(bytes)

		if _, ok := h.mapping[key]; !ok {
			copy(newIndices[ns*ni:], indices[i*ni:(i+1)*ni])
			h.mapping[key] = true
			ns++
		}

		header.Data += nb
	}

	return newIndices[:ns*ni]
}
