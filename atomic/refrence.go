package atomic

import (
	"sync/atomic"
	"unsafe"
)

type valueRef struct {
	value interface{}
}

type Reference struct {
	p *unsafe.Pointer
}

func NewReference(value interface{}) *Reference {
	var p = unsafe.Pointer(&valueRef{value: value})
	return &Reference{p: &p}
}

func (sr *Reference) Get() interface{} {
	var a = (*valueRef)(atomic.LoadPointer(sr.p))
	return a.value
}

func (sr *Reference) CompareAndSet(expectedV interface{}, newV interface{}) bool {
	var old = atomic.LoadPointer(sr.p)
	var cur = (*valueRef)(old)

	if cur.value == expectedV {
		if cur.value != newV {
			return atomic.CompareAndSwapPointer(sr.p, old, unsafe.Pointer(&valueRef{value: newV}))
		}
	}
	return false
}

// Unconditionally sets the value
func (sr *Reference) Set(newV interface{}) {
	var old = atomic.LoadPointer(sr.p)
	var cur = (*valueRef)(old)
	if newV != cur.value {
		p := unsafe.Pointer(&valueRef{value: newV}) // 注意，需要保证每次对p的修改都会把地址值改掉
		sr.p = &p
	}
}
