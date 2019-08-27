package atomic

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type ReferenceArray struct {
	ptr uintptr
	len int
}

var noUseInterface interface{}
var sizeOfInterface uintptr

func init() {
	sizeOfInterface = unsafe.Sizeof(noUseInterface)
}

// NOTE:需要和golang中的slice结构保持一致
// TODO:golang切片结构修改时,要同步修改这个结构
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

func NewReferenceArray(len int) *ReferenceArray {
	s := make([]unsafe.Pointer, len)
	s2 := (*slice)(unsafe.Pointer(&s))
	return &ReferenceArray{ptr: uintptr(s2.array), len: len}
}

func (arr *ReferenceArray) calcOffset(i int) uintptr {
	if i < 0 || i >= arr.len {
		panic(fmt.Sprintf("index out of bounds %d, len:%d", i, arr.len))
	}
	return uintptr(i) * sizeOfInterface
}

func (arr *ReferenceArray) Get(i int) interface{} {
	elePtrPtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	elePtr := (*interface{})(atomic.LoadPointer(elePtrPtr))
	return *elePtr
}

func (arr *ReferenceArray) Set(i int, v interface{}) {
	vPtr := unsafe.Pointer(&v)
	elePtrPtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	atomic.StorePointer(elePtrPtr, vPtr)
}

func (arr *ReferenceArray) CompareAndSwap(i int, oldV interface{}, newV interface{}) bool {
	elePtrPtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	var oldEle = atomic.LoadPointer(elePtrPtr)
	var curV = (*interface{})(oldEle)
	oldPtr := unsafe.Pointer(&oldV)
	newPtr := unsafe.Pointer(&newV)
	if *curV == oldV {
		if oldV != newV {
			return atomic.CompareAndSwapPointer(elePtrPtr, oldPtr, newPtr)
		}
	}
	return false
}
