package util

import (
	"sync/atomic"
	"unsafe"
)

type pair struct {
	v       interface{}
	stamped int64
}

type StampedReference struct {
	p unsafe.Pointer
}

func NewStampedReference(v interface{}, stamped int64) StampedReference {
	var p = unsafe.Pointer(&pair{v: v, stamped: stamped})
	return StampedReference{p: p}
}

func (sr *StampedReference) GetReference() interface{} {
	var a = (*pair)(atomic.LoadPointer(&sr.p))
	return a.v
}

func (sr *StampedReference) GetStamp() interface{} {
	var a = (*pair)(atomic.LoadPointer(&sr.p))
	return a.stamped
}

func (sr *StampedReference) Get() (interface{}, int64) {
	var a = (*pair)(atomic.LoadPointer(&sr.p))
	return a.v, a.stamped
}

func (sr *StampedReference) CompareAndSet(expectedV interface{}, newV interface{}, expectedStamp int64, newStamp int64) bool {
	var cur = (*pair)(atomic.LoadPointer(&sr.p))
	if cur.v == expectedV && cur.stamped == expectedStamp {
		if cur.v != newV || cur.stamped != newStamp {
			return atomic.CompareAndSwapPointer(&sr.p, unsafe.Pointer(&cur), unsafe.Pointer(&pair{v: newV, stamped: newStamp}))
		}
		return true
	}
	return false
}
