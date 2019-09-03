package atomic

import (
	"sync/atomic"
	"unsafe"
)

type stampedPair struct {
	value interface{}
	stamp int64
}

type StampedReference struct {
	p *unsafe.Pointer
}

func NewStampedReference(value interface{}, stamped int64) *StampedReference {
	var p = unsafe.Pointer(&stampedPair{value: value, stamp: stamped})
	return &StampedReference{p: &p}
}

func (sr *StampedReference) GetReference() interface{} {
	var a = (*stampedPair)(atomic.LoadPointer(sr.p))
	return a.value
}

func (sr *StampedReference) GetStamp() interface{} {
	var a = (*stampedPair)(atomic.LoadPointer(sr.p))
	return a.stamp
}

func (sr *StampedReference) Get() (interface{}, int64) {
	var a = (*stampedPair)(atomic.LoadPointer(sr.p))
	return a.value, a.stamp
}

func (sr *StampedReference) CompareAndSet(expectedV interface{}, newV interface{}, expectedStamp int64, newStamp int64) bool {
	// 惯用方式:先把old的指针值读出来,做一系列判断以后,改变的时候,用old值作为基准
	var old = atomic.LoadPointer(sr.p)
	var cur = (*stampedPair)(old)

	if cur.value == expectedV && cur.stamp == expectedStamp {
		//fmt.Printf("*sr.p:%value old:%value\n", *sr.p, old)
		return atomic.CompareAndSwapPointer(sr.p, old, unsafe.Pointer(&stampedPair{value: newV, stamp: newStamp}))
	}
	return false
}

func (sr *StampedReference) AttemptStamp(expectedV interface{}, newStamp int64) bool {
	var old = atomic.LoadPointer(sr.p)
	var cur = (*stampedPair)(old)
	if cur.value == expectedV && cur.stamp != newStamp {
		return atomic.CompareAndSwapPointer(sr.p, old, unsafe.Pointer(&stampedPair{value: expectedV, stamp: newStamp}))
	}
	return false
}

// Unconditionally sets both the value and stamp.
func (sr *StampedReference) Set(newV interface{}, newStamp int64) {
	var newP = unsafe.Pointer(&stampedPair{value: newV, stamp: newStamp})
	var old = atomic.LoadPointer(sr.p)
	var cur = (*stampedPair)(old)
	if newV != cur.value || newStamp != cur.stamp {
		atomic.StorePointer(sr.p, newP)
	}
}
