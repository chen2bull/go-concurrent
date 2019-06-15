package atomic

import (
	"sync/atomic"
	"unsafe"
)

type markablePair struct {
	v    interface{}
	mark bool
}

type MarkableReference struct {
	p *unsafe.Pointer
}

func NewMarkableReference(v interface{}, mark bool) *MarkableReference {
	var p = unsafe.Pointer(&markablePair{v: v, mark: mark})
	return &MarkableReference{p: &p}
}

func (sr *MarkableReference) GetReference() interface{} {
	var a = (*markablePair)(atomic.LoadPointer(sr.p))
	return a.v
}

func (sr *MarkableReference) IsMarked() interface{} {
	var a = (*markablePair)(atomic.LoadPointer(sr.p))
	return a.mark
}

func (sr *MarkableReference) Get() (interface{}, bool) {
	var a = (*markablePair)(atomic.LoadPointer(sr.p))
	return a.v, a.mark
}

func (sr *MarkableReference) CompareAndSet(expectedV interface{}, newV interface{}, expectedMark bool, newMark bool) bool {
	// 惯用方式:先把old的指针值读出来,做一系列判断以后,改变的时候,用old值作为基准
	var old = atomic.LoadPointer(sr.p)
	var cur = (*markablePair)(old)

	if cur.v == expectedV && cur.mark == expectedMark {
		if cur.v != newV || cur.mark != newMark {
			return atomic.CompareAndSwapPointer(sr.p, old, unsafe.Pointer(&markablePair{v: newV, mark: newMark}))
		}
	}
	return false
}

func (sr *MarkableReference) AttemptMark(expectedV interface{}, newMark bool) bool {
	var old = atomic.LoadPointer(sr.p)
	var cur = (*markablePair)(old)
	if cur.v == expectedV && cur.mark != newMark {
		return atomic.CompareAndSwapPointer(sr.p, old, unsafe.Pointer(&markablePair{v: expectedV, mark: newMark}))
	}
	return false
}

// Unconditionally sets both the value and mark.
func (sr *MarkableReference) Set(newV interface{}, newMark bool) {
	var old = atomic.LoadPointer(sr.p)
	var cur = (*markablePair)(old)
	if newV != cur.v || newMark != cur.mark {
		p := unsafe.Pointer(&markablePair{v: newV, mark: newMark})
		sr.p = &p
	}
}
