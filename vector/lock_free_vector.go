package vector

import (
	atomic2 "github.com/cmingjian/go-concurrent/atomic"
	"math/bits"
	"sync"
	"sync/atomic"
)

const (
	FBS             = 64 // First bucket size; can be any power of 2.
	highestBitOfFBS = 6  // highestBit(FBS)
	isNotPending    = 0
	isPending       = 1
)

type LockFreeVector struct {
	desc    *atomic2.StampedReference
	vals    *atomic2.ReferenceArray // ReferenceArray of ReferenceArray
	resLock *sync.Mutex             // see method allocateBucket
}

func NewEmptyLockFreeVector() *LockFreeVector {
	desc := atomic2.NewStampedReference(newDescriptor(0, nil), 0)
	vals := atomic2.NewReferenceArray(bits.UintSize)
	vals.Set(0, atomic2.NewReferenceArray(FBS))
	return &LockFreeVector{desc: desc, vals: vals, resLock: &sync.Mutex{}}
}

func NewLockFreeVector(size int) *LockFreeVector {
	vecPtr := NewEmptyLockFreeVector()
	vecPtr.Reserve(size)
	descVal := newDescriptor(0, nil)
	vecPtr.desc.Set(descVal, 0)
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
	vec.resLock.Lock() // It's better not to allocate the same bucket in more than one thread/routine
	if vec.vals.Get(bucketIdx) == nil {
		bucketSize := 1 << (uint)(bucketIdx+highestBitOfFBS)
		newBucket := atomic2.NewReferenceArray(bucketSize)
		if !vec.vals.CompareAndSet(bucketIdx, nil, newBucket) {
		}
	}
	vec.resLock.Unlock()
}

func (vec *LockFreeVector) WriteAt(idx int, v interface{}) {
	bucketIdx, withinIdx := getBucketAndIndex(idx)
	bucket := vec.vals.Get(int(bucketIdx)).(*atomic2.ReferenceArray)
	bucket.Set(withinIdx, v)
}

func (vec *LockFreeVector) ReadAt(idx int) interface{} {
	bucketIdx, withinIdx := getBucketAndIndex(idx)
	bucket := vec.vals.Get(int(bucketIdx)).(*atomic2.ReferenceArray)
	return bucket.Get(withinIdx)
}

func (vec *LockFreeVector) Size() int {
	currDesc := vec.desc.GetReference().(*descriptor)
	size := currDesc.size
	if currDesc.writeOp != nil && atomic.LoadInt32(&currDesc.writeOp.pending) == isPending {
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
	if writeDesc != nil && atomic.LoadInt32(&writeDesc.pending) == isPending {
		bucketIdx, withinIdx := getBucketAndIndex(writeDesc.idx)
		array := vec.vals.Get(bucketIdx).(*atomic2.ReferenceArray)
		// 情况1.如果多个线程同时获得修改前的地址,那么只有一个线程的"CAS A"那一行成功执行
		// 情况2.如果线程P1在另一个线程P2成功执行"CAS A"那一行后,读取地址，那么在线程P1中下面的if语句一定会失败
		addr := array.GetAddress(withinIdx)
		if atomic.LoadInt32(&writeDesc.pending) == isPending { // can not omit
			atomic.StoreInt32(&writeDesc.pending, isNotPending)                           // 不能和下面一行换位置
			array.CompareAddrValueAndSet(withinIdx, addr, writeDesc.oldV, writeDesc.newV) // CAS A
		}
	}
}

func (vec *LockFreeVector) tryCompleteWrite() {
	currDescRef, _ := vec.desc.Get()
	currDesc := currDescRef.(*descriptor)
	vec.completeWrite(currDesc.writeOp)
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
	pending int32
}

func newWriteDescriptor(oldV interface{}, newV interface{}, idx int) *writeDescriptor {
	return &writeDescriptor{oldV: oldV, newV: newV, idx: idx, pending: isPending}
}
