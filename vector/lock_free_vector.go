package vector

import (
	"github.com/cmingjian/go-concurrent/atomic"
	"math/bits"
)

const FBS = 64            // First bucket size; can be any power of 2.
const highestBitOfFBS = 6 // highestBit(FBS)

type LockFreeVector struct {
	desc *atomic.StampedReference
	vals *atomic.ReferenceArray // ReferenceArray of ReferenceArray
}

func NewEmptyLockFreeVector() *LockFreeVector {
	desc := atomic.NewStampedReference(newDescriptor(0, nil), 0)
	vals := atomic.NewReferenceArray(bits.UintSize)
	vals.Set(0, atomic.NewReferenceArray(FBS))
	return &LockFreeVector{desc: desc, vals: vals}
}

func NewLockFreeVector(size int) *LockFreeVector {
	vecPtr := NewEmptyLockFreeVector()
	vecPtr.Reserve(size)
	desc := atomic.NewStampedReference(newDescriptor(0, nil), 0)
	vecPtr.desc.Set(desc, 0)
	return vecPtr
}

func (vec *LockFreeVector) Reserve(newSize int) {
	descRef, _ := vec.desc.Get()
	desc := descRef.(*descriptor)
	bucketIdx := getBucket(desc.size - 1)
	if bucketIdx < 0 {
		bucketIdx = 0
	}
	for ; bucketIdx < getBucket(newSize-1); {
		bucketIdx ++
		vec.allocateBucket(bucketIdx)
	}
}

func (vec *LockFreeVector) allocateBucket(bucketIdx int) {
	bucketSize := 1 << (uint)(bucketIdx+highestBitOfFBS)
	newBucket := atomic.NewReferenceArray(bucketSize)
	if !vec.vals.CompareAndSet(bucketIdx, nil, newBucket) {
	}
}

func (vec *LockFreeVector) WriteAt(idx int, v interface{}) {
	bucketIdx, withinIdx := getBucketAndIndex(idx)
	bucket := vec.vals.Get(int(bucketIdx)).(*atomic.ReferenceArray)
	bucket.Set(withinIdx, v)
}

func (vec *LockFreeVector) ReadAt(idx int) interface{} {
	bucketIdx, withinIdx := getBucketAndIndex(idx)
	bucket := vec.vals.Get(int(bucketIdx)).(*atomic.ReferenceArray)
	return bucket.Get(withinIdx)
}

func (vec *LockFreeVector) Size(idx int) int {
	currDesc := vec.desc.GetReference().(*descriptor)
	size := currDesc.size
	if currDesc.writeOp != nil && currDesc.writeOp.pending {
		size --
	}
	return size
}

func (vec *LockFreeVector) PushBack(newElement interface{}) {
	var writeDesc *writeDescriptor
	for ; ; {
		currDescRef, stamped := vec.desc.Get()
		currDesc := currDescRef.(*descriptor)
		vec.completeWrite(currDesc.writeOp)

		bucketIdx := highestBit(currDesc.size+FBS) - highestBitOfFBS
		if vec.vals.Get(bucketIdx) == nil {
			vec.allocateBucket(bucketIdx)
		}
		writeDesc = newWriteDescriptor(vec.ReadAt(currDesc.size), newElement, currDesc.size)
		newDesc := newDescriptor(currDesc.size+1, writeDesc)
		if vec.desc.CompareAndSet(currDesc, newDesc, stamped, stamped+1) {
			break
		}
	}
	vec.completeWrite(writeDesc)
}

func (vec *LockFreeVector) PopBack() interface{} {
	var element interface{}
	for ; ; {
		currDescRef, stamped := vec.desc.Get()
		currDesc := currDescRef.(*descriptor)
		vec.completeWrite(currDesc.writeOp)

		if currDesc.size == 0 {
			return nil
		}
		element = vec.ReadAt(currDesc.size - 1)
		newDesc := newDescriptor(currDesc.size, nil)
		if vec.desc.CompareAndSet(currDesc, newDesc, stamped, stamped+1) {
			break
		}
	}
	return element
}

func (vec *LockFreeVector) completeWrite(writeDesc *writeDescriptor) {
	if writeDesc != nil && writeDesc.pending {
		bucketIdx, withinIdx := getBucketAndIndex(writeDesc.idx)
		array := vec.vals.Get(bucketIdx).(*atomic.ReferenceArray)
		array.CompareAndSet(withinIdx, writeDesc.oldV, writeDesc.newV)
		writeDesc.pending = false
	}
}

func highestBit(n int) int {
	return bits.TrailingZeros(highestOneBit(n))
}

func highestOneBit64(i int64) uint {
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	return uint(uint64(i) - (uint64(i) >> 1))
}

func highestOneBit32(i int32) uint {
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	return uint(uint32(i) - (uint32(i) >> 1))
}

func highestOneBit(x int) uint {
	if bits.UintSize == 32 {
		return highestOneBit32(int32(x))
	}
	return highestOneBit64(int64(x))
}

// Returns the index of the bucket
func getBucket(i int) int {
	pos := i + FBS
	hiBit := highestBit(pos)
	return int(hiBit - highestBitOfFBS)
}

// Returns the index within the bucket
func getIdxWithinBucket(i int) int {
	pos := i + FBS
	hiBit := highestBit(pos)
	return int(pos ^ (1 << (uint)(hiBit)))
}

// Returns the index of and within the bucket
func getBucketAndIndex(i int) (int, int) {
	pos := i + FBS
	hiBit := highestBit(pos)
	return int(hiBit - highestBitOfFBS), int(pos ^ (1 << (uint)(hiBit)))
}

type descriptor struct {
	size    int
	writeOp *writeDescriptor
}

func newDescriptor(size int, writeOp *writeDescriptor) *descriptor {
	return &descriptor{size: size, writeOp: writeOp}
}

type writeDescriptor struct {
	oldV    interface{}
	newV    interface{}
	idx     int
	pending bool
}

func newWriteDescriptor(oldV interface{}, newV interface{}, idx int) *writeDescriptor {
	return &writeDescriptor{oldV: oldV, newV: newV, idx: idx, pending: true}
}
