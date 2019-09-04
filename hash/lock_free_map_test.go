package hash

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
)

func TestLockFreeMap_Put_Concurrent(t *testing.T) {
	freeMap := NewLockFreeMap(reflect.Int)
	wg := sync.WaitGroup{}
	// PushBack
	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(m *LockFreeMap, i2 int) {
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
}

func TestLockFreeMap_CalcHash(t *testing.T) {
	aTestValue3 := Abc{a: int(1000), b: 1}
	aTestValue4 := int(1000)
	m := NewLockFreeMap(reflect.Int)
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

func TestLockFreeMap_CalcHash2(t *testing.T) {
	m := NewLockFreeMap(reflect.Int)
	for i := 0; i < 0xFFF; i++ {
		n := int(rand.Int63())
		hashV := m.calcHash(n)
		if hashV&0x1 != 1 {
			t.Fatalf("illeagal hash value. hashV:%v n:%v\n", hashV, n)
		}
	}
}

func TestLockFreeMap_Put(t *testing.T) {
	m := NewLockFreeMap(reflect.Int)
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

func BenchmarkLockFreeMap_Put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < goroutineCount / 25; j ++ {
			wg.Add(1)
			go func() {
				for k:=0; k < timesPerGoroutine / 10; k ++ {
					v := valueArrayForTest[i%(timesPerGoroutine*goroutineCount)]
					freeMap.Put(i, v)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkFixedBucketLockFreeMap_Put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < goroutineCount / 25; j ++ {
			wg.Add(1)
			go func() {
				for k:=0; k < timesPerGoroutine / 10; k ++ {
					v := valueArrayForTest[i%(timesPerGoroutine*goroutineCount)]
					fixedFreeMap.Put(i, v)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSyncMap_Put(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < goroutineCount / 25; j ++ {
			wg.Add(1)
			go func() {
				for k:=0; k < timesPerGoroutine / 10; k ++ {
					v := valueArrayForTest[i%(timesPerGoroutine*goroutineCount)]
					syncMap.Store(i, v)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

var freeMap *LockFreeMap

func init() {
	freeMap = NewLockFreeMap(reflect.Int)
}
