package queue

import (
	"testing"
)

var SyncDualGoroutineNum = 8
var SyncDualCount = SyncDualGoroutineNum * 64
var SyncDualPerRoutine = SyncDualCount / SyncDualGoroutineNum

func syncDualEnqFunc(qp *SyncDualQueue, value int, doneChan chan bool) {
	for i := 0; i < SyncDualPerRoutine; i++ {
		//fmt.Printf("Enq:start i:%v v:%v\n", value+i, value+i)
		qp.Enq(value + i)
		//fmt.Printf("Enq:end: i:%v v:%v \n", value+i, value+i)
	}
	doneChan <- true
}

func syncDualDeqFunc(qp *SyncDualQueue, intChan chan int, doneChan chan bool) {
	for i := 0; i < SyncDualPerRoutine; i++ {
		//fmt.Printf("Deq:start i:%v\n", i)
		v := qp.Deq()
		//fmt.Printf("Deq:done i:%v v:%v\n", i, v)
		what := v.(int)
		//v := qp.Deq().(int)
		intChan <- what
	}
	doneChan <- true
}

func syncDualCheckIntValid(intChan chan int, t *testing.T) {
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

func TestSyncDualQueue_Both(t *testing.T) {
	qu := NewSyncDualQueue()
	doneChan := make(chan bool)
	var intChan = make(chan int, SyncDualCount)
	go syncDualCheckIntValid(intChan, t)
	for i := 0; i < SyncDualGoroutineNum; i ++ {
		go syncDualEnqFunc(qu, i*SyncDualPerRoutine, doneChan)
	}
	for i := 0; i < SyncDualGoroutineNum; i ++ {
		go syncDualDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < SyncDualGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestSyncDualQueue_Nil(t *testing.T) {
	var qu = NewSyncDualQueue()
	go qu.Enq(nil)
	var v = qu.Deq()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}

func TestSyncDualQueue_BothNil(t *testing.T) {
	qu := NewSyncDualQueue()
	doneChan := make(chan bool)
	var interfaceChan = make(chan interface{}, SyncDualCount)
	go func() {
		for v := range interfaceChan {
			if v != nil {
				t.Errorf("unexpected value|v:%d", v)
			}
		}
	}()
	for i := 0; i < SyncDualGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < SyncDualPerRoutine; i++ {
				qu.Enq(nil)
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < SyncDualGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < SyncDualPerRoutine; i++ {
				v := qu.Deq()
				interfaceChan <- v
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < SyncDualGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(interfaceChan)
}
