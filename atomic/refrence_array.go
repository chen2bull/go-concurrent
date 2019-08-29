package atomic

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type ReferenceArray struct {
	ptr uintptr
	len int
	_s  []unsafe.Pointer
}

var sizeOfInterface uintptr

func init() {
	var v unsafe.Pointer
	sizeOfInterface = unsafe.Sizeof(v)
}

// NOTE:需要和golang中的slice结构保持一致
// TODO:golang切片结构修改时,要同步修改这个结构
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

var _defaultInterface interface{}

func NewReferenceArray(len int) *ReferenceArray {
	s := make([]unsafe.Pointer, len, len)
	for i := 0; i < len; i ++ {
		s[i] = (unsafe.Pointer)(&_defaultInterface)
	}
	s2 := (*slice)(unsafe.Pointer(&s))
	return &ReferenceArray{ptr: uintptr(s2.array), len: len, _s: s}
}

func (arr *ReferenceArray) calcOffset(i int) uintptr {
	if i < 0 || i >= arr.len {
		panic(fmt.Sprintf("index out of bounds %d, len:%d", i, arr.len))
	}
	return uintptr(i) * sizeOfInterface
}

func (arr *ReferenceArray) Get(i int) interface{} {
	elePtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	facePtr := (*interface{})(atomic.LoadPointer(elePtr))
	return *facePtr
}

func (arr *ReferenceArray) GetAddress(i int) unsafe.Pointer {
	elePtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	return atomic.LoadPointer(elePtr)
}

func (arr *ReferenceArray) Set(i int, v interface{}) {
	vPtr := unsafe.Pointer(&v)
	elePtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	atomic.StorePointer(elePtr, vPtr)
}

func (arr *ReferenceArray) CompareAndSet(i int, oldV interface{}, newV interface{}) bool {
	elePtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	var oldEle = atomic.LoadPointer(elePtr)
	var curVPtr = (*interface{})(oldEle)
	newPtr := unsafe.Pointer(&newV)
	if *curVPtr == oldV {
		return atomic.CompareAndSwapPointer(elePtr, oldEle, newPtr)
	}
	return false
}

func (arr *ReferenceArray) CompareAddrValueAndSet(i int, expectOldAddr unsafe.Pointer, oldV interface{},
	newV interface{}) bool {
	elePtr := (*unsafe.Pointer)(unsafe.Pointer(arr.ptr + arr.calcOffset(i)))
	var oldEle = atomic.LoadPointer(elePtr)
	var curVPtr = (*interface{})(oldEle)
	newPtr := unsafe.Pointer(&newV)
	if *curVPtr == oldV {
		return atomic.CompareAndSwapPointer(elePtr, expectOldAddr, newPtr)
	}
	return false
}

func (arr *ReferenceArray) printElements() {
	fmt.Printf("at %p slice:%v\n", unsafe.Pointer(arr.ptr), arr._s)
}
