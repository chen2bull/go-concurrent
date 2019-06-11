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

/*
const (
	rLockOp = iota
	rUnlockOp
	wLockOp
	wUnlockOp
)

// 初始值为N1
var N1 = 2341234
var N2 = 34234
var CriticalValue = N1

func TestFifoRWLock_RLock(t *testing.T) {
	mu := NewUselessRWLock()
	chanA := make(chan int)
	chanB := make(chan int)
	chanC := make(chan int)
	chanD := make(chan int)
	inputChan := make(chan int)
	outputChan := make(chan int)
	go runChan(chanA, &mu, inputChan, outputChan)
	go runChan(chanB, &mu, inputChan, outputChan)
	go runChan(chanC, &mu, inputChan, outputChan)
	go runChan(chanD, &mu, inputChan, outputChan)
	tryRead(chanA)
	tryRead(chanB)
	tryWrite(chanC, N2, inputChan)
	tryRead(chanD)
	releaseRead(chanA, N1, outputChan)
	releaseRead(chanB, N1, outputChan)
	releaseRead(chanD, N2, outputChan)

	// 协程A'试图'获取读锁
	// 协程B'试图'获取读锁
	// 协程C'试图'获取写锁,试图将值改成N2
	// 协程D'试图'获取读锁
	// 协程A'试图'释放读锁(必成功),数值为N1
	// 协程C'试图'释放写锁
	// 协程B'试图'释放读锁(必成功),数值为N1
	close(chanA)
	close(chanB)
	close(chanC)
	close(chanD)
}

func runChan(ch chan int, mu RWLocker, inputChan chan int, outputChan chan int) {
	for v := range ch {
		switch v {
		case rLockOp:
			mu.RLock()
		case rUnlockOp:
			val := CriticalValue
			mu.RUnlock()
			outputChan <- val
		case wLockOp:
			val := <-inputChan
			mu.WLock()
			CriticalValue = val
		case wUnlockOp:
			mu.WUnlock()
		}
	}
}

func tryRead(ch chan int) {
	ch <- rLockOp
}

func releaseRead(ch chan int, N1 int, outputChan chan int) {
	ch <- rUnlockOp
	val := <-outputChan
	if val != N1 {
		panic(fmt.Sprintf("val:%d is not equal to N1:%d", val, N1))
	}
}

func tryWrite(ch chan int, N2 int, inputChan chan int) {
	ch <- wLockOp
	inputChan <- N2
	ch <- wUnlockOp
}
*/
