package lock

import (
	"sync"
	"testing"
)

var GoroutineNum = 32
var Count = GoroutineNum * 102400
var PerRoutine = Count / GoroutineNum

func TestTASLock(t *testing.T) {
	counter := 0
	var lock = &TASLock{}
	var endChan = make(chan bool)
	for i := 0; i < GoroutineNum; i ++ {
		go func() {
			for i := 0; i < PerRoutine; i ++ {
				lock.Lock()
				counter = counter + 1
				lock.Unlock()
			}
			endChan <- true
		}()
	}
	for i := 0; i < GoroutineNum; i ++ {
		<-endChan
	}
	if counter != Count {
		t.Errorf("Not Equal|counter:%d Count:%d", counter, Count)
	}
}

func TestTTASLock(t *testing.T) {
	counter := 0
	var lock = &TTASLock{}
	var endChan = make(chan bool)
	for i := 0; i < GoroutineNum; i ++ {
		go func() {
			for i := 0; i < PerRoutine; i ++ {
				lock.Lock()
				counter = counter + 1
				lock.Unlock()
			}
			endChan <- true
		}()
	}
	for i := 0; i < GoroutineNum; i ++ {
		<-endChan
	}
	if counter != Count {
		t.Errorf("Not Equal|counter:%d Count:%d", counter, Count)
	}
}

func TestMutex(t *testing.T) {
	counter := 0
	var lock = &sync.Mutex{}
	var endChan = make(chan bool)
	for i := 0; i < GoroutineNum; i ++ {
		go func() {
			for i := 0; i < PerRoutine; i ++ {
				lock.Lock()
				counter = counter + 1
				lock.Unlock()
			}
			endChan <- true
		}()
	}
	for i := 0; i < GoroutineNum; i ++ {
		<-endChan
	}
	if counter != Count {
		t.Errorf("Not Equal|counter:%d Count:%d", counter, Count)
	}
}

func TestNoMutex(t *testing.T) {
	counter := 0
	var endChan = make(chan bool)
	for i := 0; i < GoroutineNum; i ++ {
		go func() {
			for i := 0; i < PerRoutine; i ++ {
				counter = counter + 1
			}
			endChan <- true
		}()
	}
	for i := 0; i < GoroutineNum; i ++ {
		<-endChan
	}
	if counter == Count {
		t.Errorf("Unexpected equal!Maybe the Count is too small|counter:%d Count:%d", counter, Count)
	}
}
