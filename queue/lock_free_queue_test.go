package queue

import (
	"testing"
)

var LockFreeGoroutineNum = 8
var LockFreeCount = LockFreeGoroutineNum * 64
var LockFreePerRoutine = LockFreeCount / LockFreeGoroutineNum

func TestLockFreeQueue_Sequential(t *testing.T) {
	var qu = NewLockFreeQueue()
	for i := 0; i < LockFreeCount; i++ {
		qu.Enq(i)
	}
	for i := 0; i < LockFreeCount; i++ {
		var v = qu.Deq()
		var value = v.(int)
		if value != i {
			t.Errorf("not equal|value:%d i:%d", value, i)
		}
	}
}

func TestLockFreeQueue_ParallelEnq(t *testing.T) {
	qu := NewLockFreeQueue()
	doneChan := make(chan bool)
	for i := 0; i < LockFreeGoroutineNum; i++ {
		go lockFreeEnqFunc(qu, i*LockFreePerRoutine, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
	}
	//fmt.Printf("elements start=====\n")
	//qu.PrintAllElement()
	//fmt.Printf("elements end=====\n")
	var intMap = make(map[int]bool)
	for i := 0; i < LockFreeCount; i++ {
		v := qu.Deq().(int)
		_, ok := intMap[v]
		if ok {
			t.Errorf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
}

func lockFreeEnqFunc(qp *LockFreeQueue, value int, doneChan chan bool) {
	for i := 0; i < LockFreePerRoutine; i++ {
		qp.Enq(value + i)
	}
	doneChan <- true
}

func lockFreeDeqFunc(qp *LockFreeQueue, intChan chan int, doneChan chan bool) {
	for i := 0; i < LockFreePerRoutine; i++ {
		v := qp.Deq().(int)
		intChan <- v
	}
	doneChan <- true
}

func lockFreeCheckIntValid(intChan chan int, t *testing.T) {
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

func TestLockFreeQueue_ParallelDeq(t *testing.T) {
	qu := NewLockFreeQueue()
	doneChan := make(chan bool)
	for i := 0; i < LockFreeCount; i++ {
		qu.Enq(i)
	}
	var intChan = make(chan int, LockFreeCount)
	go lockFreeCheckIntValid(intChan, t)
	for i := 0; i < LockFreeGoroutineNum; i++ {
		go lockFreeDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
	}
	close(intChan)
}

func TestLockFreeQueue_Both(t *testing.T) {
	qu := NewLockFreeQueue()
	doneChan := make(chan bool)
	var intChan = make(chan int, LockFreeCount)
	go lockFreeCheckIntValid(intChan, t)
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go lockFreeEnqFunc(qu, i*LockFreePerRoutine, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go lockFreeDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestLockFreeQueue_Nil(t *testing.T) {
	var qu = NewLockFreeQueue()
	qu.Enq(nil)
	var v = qu.Deq()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}

func TestLockFreeQueue_BothNil(t *testing.T) {
	qu := NewLockFreeQueue()
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
