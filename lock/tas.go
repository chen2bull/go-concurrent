// Package lock provides basic lock lock.
//
// Values containing the types defined in this package should not be copied.
package lock

import "sync/atomic"

const (
	mutexUnlocked = 0
	mutexLocked   = 1
)

// A TASLock is a simple Test and Set lock.
//
// A TASLock must Not be copied.
type TASLock struct {
	state int32
}

func (l *TASLock) Lock() {
	for ; !atomic.CompareAndSwapInt32(&l.state, mutexUnlocked, mutexLocked); {
	}
}

func (l *TASLock) Unlock() {
	if l.state == mutexUnlocked {
		panic("sync: unlock of unlocked mutex")
	}
	l.state = mutexUnlocked
}

// A TTASLock is a simple Test Test and Set lock.
// 每次CompareAndSwapXxx操作都会迫使其他核丢弃它们cache中变量的副本
// TTASLock加入“本地旋转”，因此，性能表现会比TASLock要好
// 本地旋转是指线程反复地重读被缓存的值而不是反复地使用总线。
// 无论TASLock还是TTASLock，在Unlock以后，在某些CPU上引起其他核对总线的争抢
// 而在现代的CPU上,会引起缓存一致性的相关操作。总之，Unlock的时候会有性能损耗。
//
// A TTASLock must Not be copied.
type TTASLock struct {
	state int32
}

func (l *TTASLock) Lock() {
	for {
		for ; l.state == mutexLocked; { // 本地旋转
		}
		if atomic.CompareAndSwapInt32(&l.state, mutexUnlocked, mutexLocked) {
			return
		}
	}
}

func (l *TTASLock) Unlock() {
	if l.state == mutexUnlocked {
		panic("sync: unlock of unlocked mutex")
	}
	l.state = mutexUnlocked
}
