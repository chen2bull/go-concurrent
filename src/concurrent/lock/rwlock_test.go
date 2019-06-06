package lock

import (
	"testing"
)

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

	chanA := make(chan int)
	chanB := make(chan int)
	chanC := make(chan int)
	chanD := make(chan int)
	tryRead(chanA)
	tryRead(chanB)
	tryWrite(chanC, N2)
	tryRead(chanD)
	releaseRead(chanA, N1)
	releaseRead(chanA, N1)
	releaseRead(chanA, N1)

	// 协程A'试图'获取读锁
	// 协程B'试图'获取读锁
	// 协程C'试图'获取写锁,试图将值改成N2
	// 协程D'试图'获取读锁
	// 协程A'试图'释放读锁(必成功),数值为N1
	// 协程C'试图'释放写锁
	// 协程B'试图'释放读锁(必成功),数值为N1
}

//
//func opAndCheck(ch chan int, validFunc func(expected int) bool) {
//	opCode := <-ch
//	switch opCode {
//	case rLockOp:
//
//	case rUnlockOp:
//
//	case wLockOp:
//
//	case wUnlockOp:
//
//	default:
//		panic("unexpected error!")
//	}
//}
