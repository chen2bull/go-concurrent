package lock

import (
	atomic2 "github.com/cmingjian/go-concurrent/atomic"
	"sync/atomic"
)

type BackOffLock struct {
	*atomic2.BackOff
	state int32
}

func NewBackOffLock(minDelay, maxDelay int64) BackOffLock {
	bo := atomic2.NewBackOff(minDelay, maxDelay)
	return BackOffLock{bo, mutexUnlocked}
}

func (bl *BackOffLock) Lock() {
	for {
		for ; bl.state == mutexLocked; { // 本地旋转
		}
		if atomic.CompareAndSwapInt32(&bl.state, mutexUnlocked, mutexLocked) {
			return
		} else {
			bl.BackOffWait()
		}
	}
}

func (bl *BackOffLock) Unlock() {
	if bl.state == mutexUnlocked {
		panic("sync: unlock of unlocked mutex")
	}
	bl.state = mutexUnlocked
}

