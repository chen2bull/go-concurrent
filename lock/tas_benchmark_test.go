package lock

import "testing"

var BenchNum = 8
var BenchPerGo = 1000
var BenchCount = BenchNum*BenchPerGo

func BenchmarkTASLock(b *testing.B) {
	for i :=0; i< b.N; i++ {
		counter := 0
		var lock = &TASLock{}
		var endChan = make(chan bool)
		for i := 0; i < BenchNum; i ++ {
			go func() {
				for i := 0; i < BenchPerGo; i ++ {
					lock.Lock()
					counter = counter + 1
					lock.Unlock()
				}
				endChan <- true
			}()
		}
		for i := 0; i < BenchNum; i ++ {
			<-endChan
		}
		if counter != BenchCount {
			b.Errorf("Not Equal|counter:%d Count:%d", counter, BenchCount)
		}
	}
}

func BenchmarkTTASLock(b *testing.B) {
	for i :=0; i< b.N; i++ {
		counter := 0
		var lock = &TTASLock{}
		var endChan = make(chan bool)
		for i := 0; i < BenchNum; i ++ {
			go func() {
				for i := 0; i < BenchPerGo; i ++ {
					lock.Lock()
					counter = counter + 1
					lock.Unlock()
				}
				endChan <- true
			}()
		}
		for i := 0; i < BenchNum; i ++ {
			<-endChan
		}
		if counter != BenchCount {
			b.Errorf("Not Equal|counter:%d Count:%d", counter, BenchCount)
		}
	}
}
