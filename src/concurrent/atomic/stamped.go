package atomic

import (
	"sync/atomic"
	"unsafe"
)

type pair struct {
	v       interface{}
	stamped int64
}

type StampedReference struct {
	p *unsafe.Pointer
}

func NewStampedReference(v interface{}, stamped int64) *StampedReference {
	var p = unsafe.Pointer(&pair{v: v, stamped: stamped})
	return &StampedReference{p: &p}
}

func (sr *StampedReference) GetReference() interface{} {
	var a = (*pair)(atomic.LoadPointer(sr.p))
	return a.v
}

func (sr *StampedReference) GetStamp() interface{} {
	var a = (*pair)(atomic.LoadPointer(sr.p))
	return a.stamped
}

func (sr *StampedReference) Get() (interface{}, int64) {
	var a = (*pair)(atomic.LoadPointer(sr.p))
	return a.v, a.stamped
}

func (sr *StampedReference) CompareAndSet(expectedV interface{}, newV interface{}, expectedStamp int64, newStamp int64) bool {
	// 惯用方式:先把old的指针值读出来,做一系列判断以后,改变的时候,用old值作为基准
	var old = atomic.LoadPointer(sr.p)
	var cur = (*pair)(old)

	if cur.v == expectedV && cur.stamped == expectedStamp {
		if cur.v != newV || cur.stamped != newStamp {
			//fmt.Printf("*sr.p:%v old:%v\n", *sr.p, old)
			return atomic.CompareAndSwapPointer(sr.p, old, unsafe.Pointer(&pair{v: newV, stamped: newStamp}))
		}
	}
	return false
}
