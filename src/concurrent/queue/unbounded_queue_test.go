package queue

import (
	"testing"
)

var UnBoundedGoroutineNum = 8
var UnBoundedCount = UnBoundedGoroutineNum * 64
var UnBoundedPerRoutine = UnBoundedCount / UnBoundedGoroutineNum

func TestUnBoundedQueue_Sequential(t *testing.T) {
	var qu = NewUnBoundedQueue()
	for i := 0; i < UnBoundedCount; i++ {
		qu.Enq(i)
	}
	for i := 0; i < UnBoundedCount; i++ {
		var v = qu.Deq()
		var value = v.(int)
		if value != i {
			t.Errorf("not equal|value:%d i:%d", value, i)
		}
	}
}

func TestUnBoundedQueue_ParallelEnq(t *testing.T) {
	qu := NewUnBoundedQueue()
	doneChan := make(chan bool)
	for i := 0; i < UnBoundedGoroutineNum; i++ {
		go unBoundedEnqFunc(qu, i*UnBoundedPerRoutine, doneChan)
	}
	for i := 0; i < UnBoundedGoroutineNum; i++ {
		<-doneChan
	}
	var intMap = make(map[int]bool)
	for i := 0; i < UnBoundedCount; i++ {
		v := qu.Deq().(int)
		_, ok := intMap[v]
		if ok {
			t.Errorf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
}

func unBoundedEnqFunc(qp *UnBoundedQueue, value int, doneChan chan bool) {
	for i := 0; i < UnBoundedPerRoutine; i++ {
		qp.Enq(value + i)
	}
	doneChan <- true
}

func unBoundedDeqFunc(qp *UnBoundedQueue, intChan chan int, doneChan chan bool) {
	for i := 0; i < UnBoundedPerRoutine; i++ {
		v := qp.Deq().(int)
		intChan <- v
	}
	doneChan <- true
}

func unBoundedCheckIntValid(intChan chan int, t *testing.T) {
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

func TestUnBoundedQueue_ParallelDeq(t *testing.T) {
	qu := NewUnBoundedQueue()
	doneChan := make(chan bool)
	for i := 0; i < UnBoundedCount; i++ {
		qu.Enq(i)
	}
	var intChan = make(chan int, UnBoundedCount)
	go unBoundedCheckIntValid(intChan, t)
	for i := 0; i < UnBoundedGoroutineNum; i++ {
		go unBoundedDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < UnBoundedGoroutineNum; i++ {
		<-doneChan
	}
	close(intChan)
}

func TestUnBoundedQueue_Both(t *testing.T) {
	qu := NewUnBoundedQueue()
	doneChan := make(chan bool)
	var intChan = make(chan int, UnBoundedCount)
	go unBoundedCheckIntValid(intChan, t)
	for i := 0; i < UnBoundedGoroutineNum; i ++ {
		go unBoundedEnqFunc(qu, i*UnBoundedPerRoutine, doneChan)
		go unBoundedDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < UnBoundedGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestUnBoundedQueue_Nil(t *testing.T) {
	var qu = NewUnBoundedQueue()
	qu.Enq(nil)
	var v = qu.Deq()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}
