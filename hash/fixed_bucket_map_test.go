package hash

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestFixedBucketLockFreeMap_make_key(t *testing.T) {
	for i := 0; i < 0xFFFFF; i++ {
		n := rand.Int63()
		sentinelKey := makeSentinelKey(n)
		regularKey := makeRegularKey(n)
		if sentinelKey&0x1 != 0 {
			t.Fatalf("invalid sentinelKey:%v", sentinelKey)
		}
		if regularKey&0x1 != 1 {
			t.Fatalf("invalid sentinelKey:%v", regularKey)
		}
	}

	for i := 0; i < 0xFFFF; i++ {
		sentinelKey := makeSentinelKey(int64(i))
		if sentinelKey&0x1 != 0 {
			t.Fatalf("invalid sentinelKey:%v", sentinelKey)
		}
		//fmt.Printf("%56b\n", sentinelKey)
	}
}

func assertStringPanic(t *testing.T, f func(), s string) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("The code did not panic")
		}
		rs := r.(string)
		if !strings.HasPrefix(rs, s) {
			t.Fatalf("recover return:%v | expected prefix:%v", rs, s)
		}
		fmt.Printf("XFAIL:%v\n", rs)
	}()
	f()
}

func TestFixedBucketLockFreeMap_New(t *testing.T) {
	Func := func() {
		NewFixedBucketLockFreeMap(123, reflect.Array)
	}
	assertStringPanic(t, Func, "unsupported type:array")
}

var aTestValue1 = int(1000)

const aTestValue2 = int(1000)

func TestFixedBucketLockFreeMap_CalcHash(t *testing.T) {
	aTestValue3 := Abc{a: int(1000), b: 1}
	aTestValue4 := int(1000)
	m := NewFixedBucketLockFreeMap(128, reflect.Int)
	v1Hash := m.calcHash(aTestValue1)
	v2Hash := m.calcHash(aTestValue2)
	v3Hash := m.calcHash(aTestValue3.a)
	v4Hash := m.calcHash(aTestValue4)
	if v1Hash != v2Hash || v2Hash != v3Hash || v3Hash != v4Hash {
		t.Fatalf("v1Hash:%v v2Hash:%v v3Hash:%v v4Hash:%v\n", v1Hash, v2Hash, v3Hash, v4Hash)
	}
	fmt.Printf("v1Hash:%v v2Hash:%v v3Hash:%v v4Hash:%v\n", v1Hash, v2Hash, v3Hash, v4Hash)

	if !m.isEqual(aTestValue1, aTestValue2) {
		t.Fatalf("unexpected equal1")
	}
	if !m.isEqual(aTestValue2, aTestValue3.a) {
		t.Fatalf("unexpected equal2")
	}
	if !m.isEqual(aTestValue3.a, aTestValue4) {
		t.Fatalf("unexpected equal3")
	}
}

func TestFixedBucketLockFreeMap_CalcHash2(t *testing.T) {
	m := NewFixedBucketLockFreeMap(128, reflect.Int)
	for i := 0; i < 0xFFF; i++ {
		n := int(rand.Int63())
		hashV := m.calcHash(n)
		if hashV&0x1 != 1 {
			t.Fatalf("illeagal hash value. hashV:%v n:%v\n", hashV, n)
		}
	}
}

func TestFixedBucketLockFreeMap_Put(t *testing.T) {
	m := NewFixedBucketLockFreeMap(128, reflect.Int)
	offset := 1000
	for i := 0; i < 0xFFF; i++ {
		m.Put(i, i+offset)
		v := m.Get(i)
		if v != offset+i {
			m.printAllElements()
			m.Get(i)
			t.Fatalf("unexpected value. v:%v i:%v\n", v, i)
		}
	}
	//m.printAllElements()
}

const (
	goroutineCount    = 1000
	timesPerGoroutine = 5000
)

var valueArrayForTest [timesPerGoroutine * goroutineCount]int
var fixedFreeMap *FixedBucketLockFreeMap
var syncMap = &sync.Map{}

func TestFixedBucketLockFreeMap_Put_Concurrent(t *testing.T) {
	freeMap := NewFixedBucketLockFreeMap(timesPerGoroutine*goroutineCount/3, reflect.Int)
	wg := sync.WaitGroup{}

	// PushBack
	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(m *FixedBucketLockFreeMap, i2 int) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v := valueArrayForTest[i2*timesPerGoroutine+j]
				m.Put(i2*timesPerGoroutine+j, v)
				v2 := m.Get(i2*timesPerGoroutine + j)
				if v != v2 {
					panic(fmt.Sprintf("not equal|v:%v v2:%v\n", v, v2))
				}
			}
			wg.Done()
		}(freeMap, i)
	}
	wg.Wait()
	if freeMap.tabSize != int64(timesPerGoroutine*goroutineCount) {
		t.Fatalf("tabsize:%v timesPerGoroutine:%v goroutineCount:%v\n",
			freeMap.tabSize, timesPerGoroutine, goroutineCount)
	}
}

func TestSyncMap_Put_Concurrent(t *testing.T) {
	conMap := sync.Map{}
	wg := sync.WaitGroup{}

	// PushBack
	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(m *sync.Map, i2 int) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v := valueArrayForTest[i2*timesPerGoroutine+j]
				m.Store(i2*timesPerGoroutine+j, v)
				v2, _ := m.Load(i2*timesPerGoroutine + j)
				if v != v2 {
					panic(fmt.Sprintf("not equal|v:%v v2:%v\n", v, v2))
				}
			}
			wg.Done()
		}(&conMap, i)
	}
	wg.Wait()
}


func init() {
	valueArrayForTest = [timesPerGoroutine * goroutineCount] int{}
	for i := 0; i < goroutineCount; i ++ {
		for j := 0; j < timesPerGoroutine; j++ {
			valueArrayForTest[i*timesPerGoroutine+j] = rand.Int()
		}
	}
	fixedFreeMap = NewFixedBucketLockFreeMap(timesPerGoroutine*goroutineCount/3, reflect.Int)
}
