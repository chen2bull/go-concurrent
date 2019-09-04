package hash

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

type Abc struct {
	a int
	b int
}
//type Comparer interface {
//	Compare(j interface{}) int
//}

func (abc *Abc) Compare(abc2 *Abc) int {
	if abc.a > abc2.a {
		return 1
	}
	if abc.a < abc2.a {
		return -1
	}
	if abc.b > abc2.b {
		return 1
	}
	if abc.b < abc2.b {
		return -1
	}
	return 0
}

func TestHello(t *testing.T) {
	v11 := Abc{a: 1, b: 1}
	v12 := Abc{a: 1, b: 2}
	v21 := Abc{a: 2, b: 1}
	v22 := Abc{a: 2, b: 2}
	if v11.Compare(&v12) == 0 {
		t.Fatalf("fail")
	}
	if v21.Compare(&v22) == 0 {
		t.Fatalf("fail")
	}
}

func TestBucketList_NewBucketList(t *testing.T) {
	bucketLs := NewBucketList()
	if bucketLs.Contains(1, 1) {
		t.Fatalf("contains fail!")
	}
	if bucketLs.Contains(math.MaxInt64-1, math.MaxInt64-1) {
		t.Fatalf("contains fail!")
	}
	sentinel := bucketLs.getSentinelByBucket(0x10)
	fmt.Printf("sentinel:%v", sentinel)
}

func  TestBucketList_Put_Get(t *testing.T) {
	bucketLs := NewBucketList()
	bucketLs.Put(1000, "abc", 10)
	//bucketLs.printAllElements()
	v0 := bucketLs.Get(1000, "abc")
	if v0 != 10 {
		t.Fatalf("error v0:%v expected:%v\n", v0, 10)
	}
	bucketLs.Put(1000, "abc", Abc{a: 1, b: 1})
	v1 := bucketLs.Get(1000, "abc")
	s1 := Abc{a: 1, b: 1}
	if v1 != s1 {
		t.Fatalf("error value:%v expected:%v\n", v1, s1)
	}
	bucketLs.Put(1000, "bcd", 12300)
	v2 := bucketLs.Get(1000, "abc")
	v3 := bucketLs.Get(1000, "bcd")
	v4 := bucketLs.Get(1000, "bcdef")
	v5 := bucketLs.Get(999, "abc")
	if v1 != v2 {
		t.Fatalf("error v1:%v v2:%v\n", v1, s1)
	}
	if v3 != 12300 {
		t.Fatalf("error v3:%v\n", v3)
	}
	bucketLs.Put(1000, "bcd", 32100)
	if v4 != nil {
		t.Fatalf("error v4:%v\n", v4)
	}
	if v5 != nil {
		t.Fatalf("error v5:%v\n", v5)
	}
	bucketLs.printAllElements()
}

func TestBucketList_Put(t *testing.T) {
	bl := NewBucketList()
	m := NewFixedBucketLockFreeMap(128, reflect.Int)
	offset := 1000
	for i := 0; i < 0xFF; i++ {
		keyHash := m.calcHash(i)
		bl.Put(keyHash, i, i + offset)
	}

	//bl.printAllElements()
}

func  TestBucketList_Remove(t *testing.T) {
	bucketLs := NewBucketList()
	if bucketLs.Remove(100, 100) {
		t.Fatalf("unexpected succ")
	}
	bucketLs.Put(100, 100, 1654564)

	if bucketLs.Remove(100, 99) {
		t.Fatalf("unexpected succ")
	}

	if !bucketLs.Remove(100, 100) {
		t.Fatalf("unexpected fail")
	}
}

func TestBucketList_getSentinel(t *testing.T) {
	bucketLs := NewBucketList()
	bucketLs.Put(1, 1, 2342)
	bucketLs.Put(1000, 101, 1654564)
	bucketLs.Put(2, 2, 213123)
	bucketLs.Put(3, 3, 12123)
	bucketLs.Put(101, 101, 1654564)
	bucketLs.Put(101, 102, 1654564)
	bucketLs.Put(101, 103, 1654564)
	bucketLs.Put(101, 101, 417417)
	for i:= uint(0); i < uint(55); i++ {
		bl1 := bucketLs.getSentinelByHash(100)
		if !bl1.Contains(101, 101) {
			fmt.Printf("=======================================================================\n")
			bucketLs.printAllElements()
			fmt.Printf("=======================================================================\n")
			bl1.printAllElements()
			t.Fatalf("unexpected fail")
		}
	}

	bucketLs.printAllElements()
}
