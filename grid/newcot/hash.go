package newcot

import (
	"reflect"
	"unsafe"
)

type hash struct {
	depth   int
	mapping map[string]bool
}

func newHash(depth uint, capacity uint) *hash {
	return &hash{
		depth:   int(depth),
		mapping: make(map[string]bool, capacity),
	}
}

func (h *hash) tap(trace []uint64) bool {
	const (
		sizeOfUint64 = 8
	)

	sliceHeader := *(*reflect.SliceHeader)(unsafe.Pointer(&trace))

	stringHeader := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sizeOfUint64 * h.depth,
	}

	key := *(*string)(unsafe.Pointer(&stringHeader))

	if _, ok := h.mapping[key]; ok {
		return true
	} else {
		h.mapping[key] = true
		return false
	}
}
