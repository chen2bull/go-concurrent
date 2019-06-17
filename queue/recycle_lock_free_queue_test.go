package queue

import (
	"fmt"
	"testing"
)

var RecycleLockFreeGoroutineNum = 8
var RecycleLockFreeCount = RecycleLockFreeGoroutineNum * 64
var RecycleLockFreePerRoutine = RecycleLockFreeCount / RecycleLockFreeGoroutineNum

func TestRecycleLockFreeQueue_Sequential(t *testing.T) {
	var qu = NewRecycleLockFreeQueue()
	for i := 0; i < RecycleLockFreeCount; i++ {
		qu.Enq(i)
	}
	for i := 0; i < RecycleLockFreeCount; i++ {
		var v = qu.Deq()
		var value = v.(int)
		if value != i {
			t.Errorf("not equal|value:%d i:%d", value, i)
		}
	}
}

func TestRecycleLockFreeQueue_ParallelEnq(t *testing.T) {
	qu := NewRecycleLockFreeQueue()
	doneChan := make(chan bool)
	for i := 0; i < RecycleLockFreeGoroutineNum; i++ {
		go recycleLockFreeEnqFunc(qu, i*RecycleLockFreePerRoutine, doneChan)
	}
	for i := 0; i < RecycleLockFreeGoroutineNum; i++ {
		<-doneChan
	}
	fmt.Printf("elements start\n")
	qu.PrintAllElement()
	fmt.Printf("elements end\n")
	var intMap = make(map[int]bool)
	for i := 0; i < RecycleLockFreeCount; i++ {
		v := qu.Deq().(int)
		_, ok := intMap[v]
		if ok {
			t.Errorf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
}

func recycleLockFreeEnqFunc(qp *RecycleLockFreeQueue, value int, doneChan chan bool) {
	for i := 0; i < RecycleLockFreePerRoutine; i++ {
		qp.Enq(value + i)
	}
	doneChan <- true
}

func recycleLockFreeDeqFunc(qp *RecycleLockFreeQueue, intChan chan int, doneChan chan bool) {
	for i := 0; i < RecycleLockFreePerRoutine; i++ {
		v := qp.Deq().(int)
		intChan <- v
	}
	doneChan <- true
}

func recycleLockFreeCheckIntValid(intChan chan int, t *testing.T) {
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

func TestRecycleLockFreeQueue_ParallelDeq(t *testing.T) {
	qu := NewRecycleLockFreeQueue()
	doneChan := make(chan bool)
	for i := 0; i < RecycleLockFreeCount; i++ {
		qu.Enq(i)
	}
	var intChan = make(chan int, RecycleLockFreeCount)
	go recycleLockFreeCheckIntValid(intChan, t)
	for i := 0; i < RecycleLockFreeGoroutineNum; i++ {
		go recycleLockFreeDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < RecycleLockFreeGoroutineNum; i++ {
		<-doneChan
	}
	close(intChan)
}

func TestRecycleLockFreeQueue_Both(t *testing.T) {
	qu := NewRecycleLockFreeQueue()
	doneChan := make(chan bool)
	var intChan = make(chan int, RecycleLockFreeCount)
	go recycleLockFreeCheckIntValid(intChan, t)
	for i := 0; i < RecycleLockFreeGoroutineNum; i ++ {
		go recycleLockFreeEnqFunc(qu, i*RecycleLockFreePerRoutine, doneChan)
	}
	for i := 0; i < RecycleLockFreeGoroutineNum; i ++ {
		go recycleLockFreeDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < RecycleLockFreeGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestRecycleLockFreeQueue_Nil(t *testing.T) {
	var qu = NewRecycleLockFreeQueue()
	qu.Enq(nil)
	var v = qu.Deq()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}
