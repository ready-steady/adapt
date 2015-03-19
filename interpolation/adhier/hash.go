package adhier

import (
	"reflect"
	"unsafe"
)

type hash struct {
	ni      int
	mapping map[string]bool
}

func newHash(ni uint, capacity uint) *hash {
	return &hash{
		ni:      int(ni),
		mapping: make(map[string]bool, capacity),
	}
}

func (h *hash) unique(indices []uint64) []uint64 {
	const (
		sizeOfUint64 = 8
	)

	ni := h.ni
	nn := len(indices) / ni
	nb := ni * sizeOfUint64

	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Cap = nb
	header.Len = nb

	offset := ((*reflect.SliceHeader)(unsafe.Pointer(&indices))).Data

	ns := 0

	for i, j := 0, 0; i < nn; i++ {
		header.Data = offset + uintptr(j*nb)
		key := string(bytes)

		if _, ok := h.mapping[key]; !ok {
			h.mapping[key] = true
			if j > ns {
				copy(indices[ns*ni:], indices[j*ni:])
				j = ns
			}
			ns++
		}

		j++
	}

	return indices[:ns*ni]
}
