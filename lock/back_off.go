package lock

import (
	"math/rand"
	"sync/atomic"
	"time"
)

// unit:Nanosecond
type BackOff struct {
	minDelay int64
	maxDelay int64
	limit    int64
}

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

// 指数后退
func (b *BackOff) backoff() {
	var delay = rand.Int63n(b.limit)
	b.limit = Min(b.maxDelay, 2*b.limit)
	time.Sleep(time.Duration(delay))
}

type BackOffLock struct {
	BackOff
	state int32
}

func NewBackOffLock(minDelay, maxDelay int64) BackOffLock {
	return BackOffLock{BackOff{minDelay, maxDelay, minDelay}, mutexUnlocked}
}

func (bl *BackOffLock) Lock() {
	for {
		for ; bl.state == mutexLocked; { // 本地旋转
		}
		if atomic.CompareAndSwapInt32(&bl.state, mutexUnlocked, mutexLocked) {
			return
		} else {
			bl.backoff()
		}
	}
}

func (bl *BackOffLock) Unlock() {
	if bl.state == mutexUnlocked {
		panic("sync: unlock of unlocked mutex")
	}
	bl.state = mutexUnlocked
}
