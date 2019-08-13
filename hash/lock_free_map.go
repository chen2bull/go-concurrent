package hash

import (
	"sync/atomic"
)

const THRESHOLD = 4.0

type LockFreeMap struct {
	bucket     []*BucketList
	bucketSize int32
	setSize    int32
}

func NewLockFreeMap(capacity int) *LockFreeMap {
	bucket := make([]*BucketList, capacity, capacity)
	bucket[0] = NewBucketList()
	var bucketSize int32 = 2
	var setSize int32 = 0
	return &LockFreeMap{bucket: bucket, bucketSize: bucketSize, setSize: setSize}
}

func (lfMap *LockFreeMap) Add(value Hashable) bool {
	bucketSize := atomic.LoadInt32(&lfMap.bucketSize)
	myBucket := Abs(value.hashCode() % bucketSize)
	bl := lfMap.getBucketList(myBucket)
	if !bl.Add(value) {
		return false
	}
	setSizeNow := atomic.AddInt32(&lfMap.setSize, 1)
	bucketSizeNow := atomic.LoadInt32(&lfMap.bucketSize)

	if float64(setSizeNow)/float64(bucketSizeNow) > THRESHOLD {
		atomic.CompareAndSwapInt32(&lfMap.bucketSize, bucketSizeNow, 2*bucketSizeNow) // 如果失败,表示已经在别处添加
	}
	return true
}

func (lfMap *LockFreeMap) Remove(value Hashable) bool {
	bucketSize := atomic.LoadInt32(&lfMap.bucketSize)
	myBucket := Abs(value.hashCode() % bucketSize)
	b := lfMap.getBucketList(myBucket)
	if !b.Remove(value) {
		return false
	}
	return true
}

func (lfMap *LockFreeMap) Contains(value Hashable) bool {
	bucketSize := atomic.LoadInt32(&lfMap.bucketSize)
	myBucket := Abs(value.hashCode() % bucketSize)
	b := lfMap.getBucketList(myBucket)
	return b.Contains(value)
}

func (lfMap *LockFreeMap) getBucketList(myBucket int32) *BucketList {
	if lfMap.bucket[myBucket] == nil {
		lfMap.initializeBucket(myBucket)
	}
	return lfMap.bucket[myBucket]
}

func (lfMap *LockFreeMap) initializeBucket(myBucket int32) {
	parent := lfMap.getParent(myBucket)
	if lfMap.bucket[parent] == nil {
		lfMap.initializeBucket(parent)
	}
	bl := lfMap.bucket[parent].getSentinel(myBucket)
	if bl != nil {
		lfMap.bucket[myBucket] = bl
	}
}

func (lfMap *LockFreeMap) getParent(myBucket int32) int32 {
	parent := atomic.LoadInt32(&lfMap.bucketSize)
	for ; parent > myBucket; {
		parent = parent >> 1
	}
	parent = myBucket - parent
	return parent
}
