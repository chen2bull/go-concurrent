package lock

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

var N1 = 2341234
var N2 = 34234
var CriticalValue = N1

var ReadTimesPerRoutine = 100000
var WriteTimePerRoutine = 100000
var ReadRoutineCount = 3

func TestUselessRWLock(t *testing.T) {
	lo := NewUselessRWLock()
	wg := &sync.WaitGroup{}
	outChan := make(chan int, ReadRoutineCount*ReadTimesPerRoutine)
	wg.Add(ReadRoutineCount + 2)
	go validGoroutine(outChan, wg)
	go writeGoroutine(lo, wg)
	for i := 0; i < ReadRoutineCount; i++ {
		go readGoroutine(lo, outChan, wg)
	}
	wg.Wait()
}

func TestNewFifoRWLock(t *testing.T) {
	lo := NewFifoRWLock()
	wg := &sync.WaitGroup{}
	outChan := make(chan int, 100)
	wg.Add(ReadRoutineCount + 2)
	for i := 0; i < ReadRoutineCount; i++ {
		go readGoroutine(lo, outChan, wg)
	}
	go writeGoroutine(lo, wg)
	go validGoroutine(outChan, wg)
	wg.Wait()
}

func validGoroutine(outChan chan int, wg *sync.WaitGroup) {
	totalReadTimes := ReadRoutineCount * ReadTimesPerRoutine
	for i := 0; i < totalReadTimes; i++ {
		val := <-outChan
		if val != N1 && val != N2 {
			panic(fmt.Sprintf("unexpected value:%d", val))
		}
	}
	wg.Done()
}

func readGoroutine(lo RWLocker, outChan chan int, wg *sync.WaitGroup) {
	for i := 0; i < ReadTimesPerRoutine; i ++ {
		lo.RLock()
		outChan <- CriticalValue
		lo.RUnlock()
	}
	wg.Done()
}

func writeGoroutine(lo RWLocker, wg *sync.WaitGroup) {
	for i := 0; i < WriteTimePerRoutine; i ++ {
		lo.WLock()
		randNumber := rand.Intn(100)
		if randNumber%2 == 1 {
			CriticalValue = N1
		} else {
			CriticalValue = N2
		}
		lo.WUnlock()
	}
	wg.Done()
}

