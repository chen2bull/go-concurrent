package hash

import (
	"fmt"
	atomic2 "github.com/cmingjian/go-concurrent/atomic"
	"math/rand"
	"reflect"
	"sync/atomic"
)

type FixedBucketLockFreeMap struct {
	bucket     *atomic2.ReferenceArray
	t          *typeFuncs
	hashSeed   uintptr
	bucketSize int64
	bucketCap  int64
	tabSize    int64
}

func NewFixedBucketLockFreeMap(bucketCap int, keyType reflect.Kind) *FixedBucketLockFreeMap {
	t := keyTypeMap[keyType]
	if !isKeyTypeInit {
		panic("can not create map while not init\n")
	}
	if t.hash == nil {
		panic(fmt.Sprintf("unsupported type:%v\n", keyType.String()))
	}
	bucketCap = tableSizeFor(bucketCap)
	bucket := atomic2.NewReferenceArray(bucketCap)
	bucket.Set(0, NewBucketList())
	var bucketSize int64 = 128
	var tabSize int64 = 0
	hashSeed := uintptr(uint32(rand.Int31())) // do not change this line.
	return &FixedBucketLockFreeMap{
		bucket:     bucket,
		hashSeed:   hashSeed,
		bucketSize: bucketSize,
		bucketCap:  int64(bucketCap),
		tabSize:    tabSize,
		t:          &t,
	}
}

func (lfMap *FixedBucketLockFreeMap) calcHash(key interface{}) int64 {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	return makeRegularKey(hashCode)
}

func (lfMap *FixedBucketLockFreeMap) isEqual(key interface{}, key2 interface{}) bool {
	addr := lfMap.t.getInterfaceValueAddr(key)
	addr2 := lfMap.t.getInterfaceValueAddr(key2)
	return lfMap.t.equal(addr, addr2)
}

func (lfMap *FixedBucketLockFreeMap) Put(key interface{}, value interface{}) bool {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	bl := lfMap.getBucketList(myBucket)
	bl.Put(makeRegularKey(hashCode), key, value) // 必然成功

	setSizeNow := atomic.AddInt64(&lfMap.tabSize, 1)
	bucketSizeNow := atomic.LoadInt64(&lfMap.bucketSize)
	bucketCap := atomic.LoadInt64(&lfMap.bucketCap)

	if float64(setSizeNow)/float64(bucketSizeNow) > DefaultLoadFactor && 2*bucketSizeNow <= bucketCap {
		atomic.CompareAndSwapInt64(&lfMap.bucketSize, bucketSizeNow, 2*bucketSizeNow) // 如果失败,表示已经在别处添加
	}
	return true
}

func (lfMap *FixedBucketLockFreeMap) Remove(key interface{}) bool {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	if !b.Remove(makeRegularKey(hashCode), key) {
		return false
	}
	return true
}

func (lfMap *FixedBucketLockFreeMap) Contains(key interface{}) bool {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	return b.Contains(makeRegularKey(hashCode), key)
}

func (lfMap *FixedBucketLockFreeMap) Get(key interface{}) interface{} {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	//fmt.Printf("keyHash:%56b key:%v Get value throw bucket %v bucketSize:%v \n",
	//	b.head.keyHash, b.head.key, myBucket, bucketSize)
	return b.Get(makeRegularKey(hashCode), key)
}

func (lfMap *FixedBucketLockFreeMap) Find(key interface{}) (interface{}, bool){
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	return b.Find(makeRegularKey(hashCode), key)
}

func (lfMap *FixedBucketLockFreeMap) printAllElements() {
	lfMap.bucket.Get(0).(*BucketList).printAllElements()
}

func (lfMap *FixedBucketLockFreeMap) getBucketList(myBucket int64) *BucketList {
	index := int(myBucket)
	if lfMap.bucket.Get(index) == nil {
		lfMap.initializeBucket(myBucket)
	}
	bl := lfMap.bucket.Get(index).(*BucketList)
	return bl
}

func (lfMap *FixedBucketLockFreeMap) initializeBucket(myBucket int64) {
	parent := lfMap.getParent(myBucket)
	parentIndex := int(parent)
	if lfMap.bucket.Get(parentIndex) == nil {
		lfMap.initializeBucket(parent)
	}
	parentBl := lfMap.bucket.Get(parentIndex).(*BucketList)
	bl := parentBl.getSentinelByBucket(myBucket)
	if bl != nil {
		lfMap.bucket.Set(int(myBucket), bl)
	}
}

func (lfMap *FixedBucketLockFreeMap) getParent(myBucket int64) int64 {
	bitVal := atomic.LoadInt64(&lfMap.bucketSize) // bucketSize must pow of 2
	for ; bitVal > myBucket; { // 循环过后 bitVal的值为myBucket二进制的最高位
		bitVal = bitVal >> 1
	}
	return myBucket - bitVal
}
