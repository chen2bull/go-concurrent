package atomic

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewReference(t *testing.T) {
	ref := NewReference(100)
	value := ref.Get().(int)
	if value != 100 {
		t.Fatalf("wrong value:%v", value)
	}

	ok := ref.CompareAndSet(100, 101)
	if !ok {
		t.Fatalf("unexpected fail!")
	}
	value = ref.Get().(int)
	if value != 101 {
		t.Fatalf("wrong value:%v", value)
	}

	ok2 := ref.CompareAndSet(100, 102)
	if ok2 {
		t.Fatalf("unexpected succ!")
	}
	value = ref.Get().(int)
	if value != 101 {
		t.Fatalf("wrong value:%v", value)
	}

	ok3 := ref.CompareAndSet(101, 102)
	if !ok3 {
		t.Fatalf("unexpected fail!")
	}
	value = ref.Get().(int)
	if value != 102 {
		t.Fatalf("wrong value:%v", value)
	}
}


type cStruct struct {
	a int
	b int
}

func TestReferenceStruct(t *testing.T) {
	sr := NewReference(cStruct{15, 1})
	set := sr.CompareAndSet(cStruct{15, 1}, cStruct{15, 2})
	if ! set {
		t.Fatalf("set:%v", set)
	}
	v := sr.Get().(cStruct)
	var expectedValue = cStruct{15, 2}
	if v != expectedValue {
		t.Fatalf("value is error:%v", v)
	}

	set = sr.CompareAndSet(cStruct{15, 1}, cStruct{15, 3}) // fail
	if set {
		t.Fatalf("set:%v", set)
	}
	set = sr.CompareAndSet(cStruct{15, 2}, cStruct{15, 3}) // succ
	if !set {
		t.Fatalf("set:%v", set)
	}
	expectedValue = cStruct{15, 3}
	v = sr.Get().(cStruct)
	if v != expectedValue {
		t.Fatalf("value is error:%v", v)
	}
}

func TestReferenceConcurrent(t *testing.T) {
	var succTimes int64 = 0
	var succTimesP *int64
	succTimesP = &succTimes
	sr := NewReference(cStruct{15, 1})
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			if sr.CompareAndSet(cStruct{15, 1}, cStruct{15, 3}) {
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