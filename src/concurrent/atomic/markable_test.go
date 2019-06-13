package atomic

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewMarkableReference(t *testing.T) {
	sr := NewMarkableReference(15, true)
	v, mark := sr.Get()
	if !mark {
		t.Errorf("mark is error:%v", mark)
	}
	if v != 15 {
		t.Errorf("value is error:%v", v)
	}
	sr2 := NewMarkableReference(bStruct{15, 1}, true)
	v2, mark2 := sr2.Get()
	if !mark2 {
		t.Errorf("mark is error:%v", mark2)
	}
	value := bStruct{15, 1}
	if v2 != value {
		t.Errorf("value is error:%v %v", v2, value)
	}
}

type bStruct struct {
	a int
	b int
}

func TestMarkableReferenceStruct(t *testing.T) {
	sr := NewMarkableReference(bStruct{15, 1}, true)
	set := sr.CompareAndSet(bStruct{15, 1}, bStruct{15, 2}, true, false)
	if ! set {
		t.Fatalf("set:%v", set)
	}
	v, mark := sr.Get()
	if mark {
		t.Fatalf("mark is error:%v", mark)
	}
	expectedValue := bStruct{15, 2}
	if v != expectedValue {
		t.Fatalf("value is error:%v", v)
	}

	set = sr.CompareAndSet(bStruct{15, 1}, bStruct{15, 3}, false, true) // fail
	if set {
		t.Fatalf("set:%v", set)
	}
	set = sr.CompareAndSet(bStruct{15, 2}, bStruct{15, 3}, true, true) // fail
	if set {
		t.Fatalf("set:%v", set)
	}
	set = sr.CompareAndSet(bStruct{15, 2}, bStruct{15, 3}, false, true) // succ
	if !set {
		t.Fatalf("set:%v", set)
	}
}

func TestMarkableReferenceConcurrent(t *testing.T) {
	var succTimes int64 = 0
	var succTimesP *int64
	succTimesP = &succTimes
	sr := NewMarkableReference(bStruct{15, 1}, true)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			if sr.CompareAndSet(bStruct{15, 1}, bStruct{15, 3}, true, false) {
				atomic.AddInt64(succTimesP, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if succTimes != 1 {
		t.Fatalf("wrong succTimes:%v", succTimes)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			if sr.CompareAndSet(bStruct{15, 3}, bStruct{15, 4}, false, true) {
				atomic.AddInt64(succTimesP, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if succTimes != 2 {
		t.Fatalf("wrong succTimes:%v", succTimes)
	}
}
