package queue

import (
	"testing"
)

var BoundedGoroutineNum = 8
var BoundedCount = BoundedGoroutineNum * 64
var BoundedPerRoutine = BoundedCount / BoundedGoroutineNum

func TestBoundedQueue_Sequential(t *testing.T) {
	var qu = NewBoundedQueue(int64(BoundedCount))
	for i := 0; i < BoundedCount; i++ {
		qu.Enq(i)
	}
	for i := 0; i < BoundedCount; i++ {
		var v = qu.Deq()
		var value = v.(int)
		if value != i {
			t.Errorf("not equal|value:%d i:%d", value, i)
		}
	}
}

func TestBoundedQueue_ParallelEnq(t *testing.T) {
	qu := NewBoundedQueue(int64(BoundedCount))
	doneChan := make(chan bool)
	for i := 0; i < BoundedGoroutineNum; i++ {
		go boundedEnqFunc(qu, i*BoundedPerRoutine, doneChan)
	}
	for i := 0; i < BoundedGoroutineNum; i++ {
		<-doneChan
	}
	var intMap = make(map[int]bool)
	for i := 0; i < BoundedCount; i++ {
		v := qu.Deq().(int)
		_, ok := intMap[v]
		if ok {
			t.Errorf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
}

func boundedEnqFunc(qp *BoundedQueue, value int, doneChan chan bool) {
	for i := 0; i < BoundedPerRoutine; i++ {
		qp.Enq(value + i)
	}
	doneChan <- true
}

func boundedDeqFunc(qp *BoundedQueue, intChan chan int, doneChan chan bool) {
	for i := 0; i < BoundedPerRoutine; i++ {
		v := qp.Deq().(int)
		intChan <- v
	}
	doneChan <- true
}

func boundedCheckIntValid(intChan chan int, t *testing.T) {
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

func TestBoundedQueue_ParallelDeq(t *testing.T) {
	qu := NewBoundedQueue(int64(BoundedCount))
	doneChan := make(chan bool)
	for i := 0; i < BoundedCount; i++ {
		qu.Enq(i)
	}
	var intChan = make(chan int, BoundedCount)
	go boundedCheckIntValid(intChan, t)
	for i := 0; i < BoundedGoroutineNum; i++ {
		go boundedDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < BoundedGoroutineNum; i++ {
		<-doneChan
	}
	close(intChan)
}

func TestBoundedQueue_Both(t *testing.T) {
	qu := NewBoundedQueue(int64(BoundedCount))
	doneChan := make(chan bool)
	var intChan = make(chan int, BoundedCount)
	go boundedCheckIntValid(intChan, t)
	for i := 0; i < BoundedGoroutineNum; i ++ {
		go boundedEnqFunc(qu, i*BoundedPerRoutine, doneChan)
		go boundedDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < BoundedGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestBoundedQueue_Both2(t *testing.T) {
	qu := NewBoundedQueue(1) // 大大提高Enq和Deq阻塞的概率
	doneChan := make(chan bool)
	var intChan = make(chan int, BoundedCount)
	go boundedCheckIntValid(intChan, t)
	for i := 0; i < BoundedGoroutineNum; i ++ {
		go boundedEnqFunc(qu, i*BoundedPerRoutine, doneChan)
		go boundedDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < BoundedGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestBoundedQueue_Nil(t *testing.T) {
	var qu = NewBoundedQueue(int64(BoundedCount))
	qu.Enq(nil)
	var v = qu.Deq()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}
