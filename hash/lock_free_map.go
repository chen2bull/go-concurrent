package hash

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/vector"
	"math/rand"
	"reflect"
	"sync/atomic"
)

type LockFreeMap struct {
	bucket     *vector.LockFreeVector
	t          *typeFuncs
	hashSeed   uintptr
	bucketSize int64
	tabSize    int64
}

const lockFreeMapMaxCap = maxBucketSize64 >> 16 // seems we don't need such big bucket capacity

func NewLockFreeMap(keyType reflect.Kind) *LockFreeMap {
	t := keyTypeMap[keyType]
	if !isKeyTypeInit {
		panic("can not create map while not init\n")
	}
	if t.hash == nil {
		panic(fmt.Sprintf("unsupported type:%v\n", keyType.String()))
	}
	var bucketSize int64 = 128
	bucket := vector.NewEmptyLockFreeVector()
	bucket.Reserve(int(bucketSize))
	bucket.WriteAt(0, NewBucketList())
	return &LockFreeMap{
		bucket:     bucket,
		hashSeed:   uintptr(uint32(rand.Int31())),
		bucketSize: bucketSize,
		tabSize:    0,
		t:          &t,
	}
}

func (lfMap *LockFreeMap) calcHash(key interface{}) int64 {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	return makeRegularKey(hashCode)
}

func (lfMap *LockFreeMap) isEqual(key interface{}, key2 interface{}) bool {
	addr := lfMap.t.getInterfaceValueAddr(key)
	addr2 := lfMap.t.getInterfaceValueAddr(key2)
	return lfMap.t.equal(addr, addr2)
}

func (lfMap *LockFreeMap) Put(key interface{}, value interface{}) bool {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	bl := lfMap.getBucketList(myBucket)
	bl.Put(makeRegularKey(hashCode), key, value) // 必然成功

	setSizeNow := atomic.AddInt64(&lfMap.tabSize, 1)
	bucketSizeNow := atomic.LoadInt64(&lfMap.bucketSize)

	if (setSizeNow/bucketSizeNow > DefaultLoadFactor) && (2*bucketSizeNow <= lockFreeMapMaxCap) {
		lfMap.bucket.Reserve(int(2 * bucketSizeNow))
		atomic.CompareAndSwapInt64(&lfMap.bucketSize, bucketSizeNow, 2*bucketSizeNow) // 如果失败,表示已经在别处添加
	}
	return true
}

func (lfMap *LockFreeMap) Remove(key interface{}) bool {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	if !b.Remove(makeRegularKey(hashCode), key) {
		return false
	}
	return true
}

func (lfMap *LockFreeMap) Contains(key interface{}) bool {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	return b.Contains(makeRegularKey(hashCode), key)
}

func (lfMap *LockFreeMap) Get(key interface{}) interface{} {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	return b.Get(makeRegularKey(hashCode), key)
}

func (lfMap *LockFreeMap) Find(key interface{}) (interface{}, bool) {
	hashCode := abs64(int64(lfMap.t.hash(lfMap.t.getInterfaceValueAddr(key), lfMap.hashSeed)))
	bucketSize := atomic.LoadInt64(&lfMap.bucketSize)
	myBucket := hashCode % bucketSize
	b := lfMap.getBucketList(myBucket)
	return b.Find(makeRegularKey(hashCode), key)
}

func (lfMap *LockFreeMap) printAllElements() {
	iter := lfMap.Iter()
	for ; true; {
		ok, key, value := iter.Next()
		if !ok {
			break
		}
		fmt.Printf("key:%v value:%v\n", key, value)
	}
}

func (lfMap *LockFreeMap) printAllNodes() {
	lfMap.bucket.ReadAt(0).(*BucketList).printAllElements()
}

func (lfMap *LockFreeMap) getBucketList(myBucket int64) *BucketList {
	index := int(myBucket)
	if lfMap.bucket.ReadAt(index) == nil {
		lfMap.initializeBucket(myBucket)
	}
	bl := lfMap.bucket.ReadAt(index).(*BucketList)
	return bl
}

func (lfMap *LockFreeMap) initializeBucket(myBucket int64) {
	parent := lfMap.getParent(myBucket)
	parentIndex := int(parent)
	if lfMap.bucket.ReadAt(parentIndex) == nil {
		lfMap.initializeBucket(parent)
	}
	parentBl := lfMap.bucket.ReadAt(parentIndex).(*BucketList)
	bl := parentBl.getSentinelByBucket(myBucket)
	if bl != nil {
		lfMap.bucket.WriteAt(int(myBucket), bl)
	}
}

func (lfMap *LockFreeMap) getParent(myBucket int64) int64 {
	bitVal := atomic.LoadInt64(&lfMap.bucketSize) // bucketSize must pow of 2
	for ; bitVal > myBucket; { // 循环过后 bitVal的值为myBucket二进制的最高位
		bitVal = bitVal >> 1
	}
	return myBucket - bitVal
}

func (lfMap *LockFreeMap) Iter() MapIterator {
	bucket := lfMap.getBucketList(0)
	return MapIterator{curr: bucket.head}
}

type MapIterator struct {
	curr *bucketListNode
}

func (iter *MapIterator) Next() (bool, interface{}, interface{}) {
	for ; true; {
		if iter.curr.keyHash == SentinelOfSentinelHash {
			return false, nil, nil
		}
		curr := iter.curr.getNext()
		iter.curr = curr
		if curr.keyHash == SentinelOfSentinelHash {
			return false, nil, nil
		}
		if isRegularKey(curr.keyHash) {
			iter.curr = curr
			return true, curr.key, curr.value
		}
	}
	panic("never here")
}
