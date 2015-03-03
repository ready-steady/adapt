package newcot

import (
	"reflect"
	"unsafe"
)

type hash struct {
	depth   int
	mapping map[string]bool
}

func newHash(depth int, capacity int) *hash {
	return &hash{
		depth:   depth,
		mapping: make(map[string]bool, capacity),
	}
}

func (h *hash) tap(trace []uint64) bool {
	const (
		sizeOfUint64 = 8
	)

	var bytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Data = ((*reflect.SliceHeader)(unsafe.Pointer(&trace))).Data
	header.Cap = sizeOfUint64 * h.depth
	header.Len = header.Cap

	key := string(bytes)

	if _, ok := h.mapping[key]; ok {
		return true
	} else {
		h.mapping[key] = true
		return false
	}
}
