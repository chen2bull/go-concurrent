package queue

import (
	"testing"
)

var GoroutineNum = 8
var Count = GoroutineNum * 64
var PerRoutine = Count / GoroutineNum

func TestLockedQueue_Sequential(t *testing.T) {
	var qu = NewLockedQueue(Count)
	for i := 0; i < Count; i++ {
		qu.Enq(i)
	}
	for i := 0; i < Count; i++ {
		var v = qu.Deq()
		var value = v.(int)
		if value != i {
			t.Errorf("not equal|value:%d i:%d", value, i)
		}
	}

}

func TestLockedQueue_ParallelEnq(t *testing.T) {
	qu := NewLockedQueue(Count)
	doneChan := make(chan bool)
	for i := 0; i < GoroutineNum; i++ {
		go enqFunc(qu, i*PerRoutine, doneChan)
	}
	for i := 0; i < GoroutineNum; i++ {
		<-doneChan
	}
	var intMap = make(map[int]bool)
	for i := 0; i < Count; i++ {
		v := qu.Deq().(int)
		_, ok := intMap[v]
		if ok {
			t.Errorf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
}

func enqFunc(qp *LockedQueue, value int, doneChan chan bool) {
	for i := 0; i < PerRoutine; i++ {
		qp.Enq(value + i)
	}
	doneChan <- true
}

func deqFunc(qp *LockedQueue, intChan chan int, doneChan chan bool) {
	for i := 0; i < PerRoutine; i++ {
		v := qp.Deq().(int)
		intChan <- v
	}
	doneChan <- true
}

func checkIntValid(intChan chan int, t *testing.T) {
	var intMap = make(map[int]bool)
	for v := range intChan {
		_, ok := intMap[v]
		if ok {
			t.Errorf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
}

func TestLockedQueue_ParallelDeq(t *testing.T) {
	qu := NewLockedQueue(Count)
	doneChan := make(chan bool)
	for i := 0; i < Count; i++ {
		qu.Enq(i)
	}
	var intChan = make(chan int, Count)
	go checkIntValid(intChan, t)
	for i := 0; i < GoroutineNum; i++ {
		go deqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < GoroutineNum; i++ {
		<-doneChan
	}
	close(intChan)
}

func TestLockedQueue_Both(t *testing.T) {
	qu := NewLockedQueue(Count)
	doneChan := make(chan bool)
	var intChan = make(chan int, Count)
	go checkIntValid(intChan, t)
	for i := 0; i < GoroutineNum; i ++ {
		go enqFunc(qu, i*PerRoutine, doneChan)
		go deqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < GoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestLockedQueue_Nil(t *testing.T) {
	var qu = NewLockedQueue(Count)
	qu.Enq(nil)
	var v = qu.Deq()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}

func TestLockedQueue_BothNil(t *testing.T) {
	qu := NewLockedQueue(1)
	doneChan := make(chan bool)
	var interfaceChan = make(chan interface{}, LockFreeCount)
	go func() {
		for v := range interfaceChan {
			if v != nil {
				t.Errorf("unexpected value|v:%d", v)
			}
		}
	}()
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < LockFreePerRoutine; i++ {
				qu.Enq(nil)
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < LockFreePerRoutine; i++ {
				v := qu.Deq()
				interfaceChan <- v
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(interfaceChan)
}