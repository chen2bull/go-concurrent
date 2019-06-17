package atomic

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewStampedReference(t *testing.T) {
	sr := NewStampedReference(15, 1)
	v, stamp := sr.Get()
	if stamp != 1 {
		t.Errorf("stamp is error:%v", stamp)
	}
	if v != 15 {
		t.Errorf("value is error:%v", v)
	}
	sr2 := NewStampedReference(aStruct{15, 1}, 100)
	v2, stamp2 := sr2.Get()
	if stamp2 != 100 {
		t.Errorf("stamp is error:%v", stamp2)
	}
	value := aStruct{15, 1}
	if v2 != value {
		t.Errorf("value is error:%v", v2)
	}
}

type aStruct struct {
	a int
	b int
}

func TestStampedReferenceStruct(t *testing.T) {
	sr := NewStampedReference(aStruct{15, 1}, 100)
	set := sr.CompareAndSet(aStruct{15, 1}, aStruct{15, 2}, 100, 101)
	if ! set {
		t.Fatalf("set:%v", set)
	}
	v, stamp := sr.Get()
	if stamp != 101 {
		t.Fatalf("stamp is error:%v", stamp)
	}
	expectedValue := aStruct{15, 2}
	if v != expectedValue {
		t.Fatalf("value is error:%v", v)
	}

	set = sr.CompareAndSet(aStruct{15, 1}, aStruct{15, 3}, 101, 102) // fail
	if set {
		t.Fatalf("set:%v", set)
	}
	set = sr.CompareAndSet(aStruct{15, 2}, aStruct{15, 3}, 100, 102) // fail
	if set {
		t.Fatalf("set:%v", set)
	}
	set = sr.CompareAndSet(aStruct{15, 2}, aStruct{15, 3}, 101, 102) // succ
	if !set {
		t.Fatalf("set:%v", set)
	}
}

func TestStampedReferenceConcurrent(t *testing.T) {
	var succTimes int64 = 0
	var succTimesP *int64
	succTimesP = &succTimes
	sr := NewStampedReference(aStruct{15, 1}, 100)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			if sr.CompareAndSet(aStruct{15, 1}, aStruct{15, 3}, 100, 101) {
				atomic.AddInt64(succTimesP, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if succTimes != 1 {
		t.Fatalf("wrong succTimes:%v", succTimes)
	}
}
