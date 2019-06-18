package queue

import "testing"

var SyncQueueGoroutineNum = 8
var SyncQueueCount = SyncQueueGoroutineNum * 64
var SyncQueuePerRoutine = SyncQueueCount / SyncQueueGoroutineNum

func syncQueueEnqFunc(qp *SyncQueue, value int, doneChan chan bool) {
	for i := 0; i < SyncQueuePerRoutine; i++ {
		qp.Enq(value + i)
	}
	doneChan <- true
}

func syncQueueDeqFunc(qp *SyncQueue, intChan chan int, doneChan chan bool) {
	for i := 0; i < SyncQueuePerRoutine; i++ {
		v := qp.Deq().(int)
		intChan <- v
	}
	doneChan <- true
}

func syncQueueCheckIntValid(intChan chan int, t *testing.T) {
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

func TestSyncQueue_Both(t *testing.T) {
	qu := NewSyncQueue()
	doneChan := make(chan bool)
	var intChan = make(chan int, SyncQueueCount)
	go syncQueueCheckIntValid(intChan, t)
	for i := 0; i < SyncQueueGoroutineNum; i ++ {
		go syncQueueEnqFunc(qu, i*SyncQueuePerRoutine, doneChan)
	}
	for i := 0; i < SyncQueueGoroutineNum; i ++ {
		go syncQueueDeqFunc(qu, intChan, doneChan)
	}
	for i := 0; i < SyncQueueGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestSyncQueue_Nil(t *testing.T) {
	var qu = NewSyncQueue()
	go qu.Enq(nil)
	var v = qu.Deq()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}

func TestSyncQueue_BothNil(t *testing.T) {
	qu := NewSyncQueue()
	doneChan := make(chan bool)
	var interfaceChan = make(chan interface{}, SyncQueueCount)
	go func() {
		for v := range interfaceChan {
			if v != nil {
				t.Errorf("unexpected value|v:%d", v)
			}
		}
	}()
	for i := 0; i < SyncQueueGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < SyncQueuePerRoutine; i++ {
				qu.Enq(nil)
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < SyncQueueGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < SyncQueuePerRoutine; i++ {
				v := qu.Deq()
				interfaceChan <- v
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < SyncQueueGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(interfaceChan)
}
