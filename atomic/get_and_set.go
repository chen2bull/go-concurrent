package atomic

import "sync/atomic"

// Atomically sets to the given value and returns the previous value.
func GetAndSetInt64(addr *int64, newValue int64) int64 {
	var oldValue int64
	for oldValue = atomic.LoadInt64(addr); !atomic.CompareAndSwapInt64(addr, oldValue, newValue);
	oldValue = atomic.LoadInt64(addr) {
	}
	return oldValue
}

// Atomically sets to the given value and returns the previous value.
func GetAndSetInt32(addr *int32, newValue int32) int32 {
	var oldValue int32
	for oldValue = atomic.LoadInt32(addr); !atomic.CompareAndSwapInt32(addr, oldValue, newValue);
	oldValue = atomic.LoadInt32(addr) {
	}
	return oldValue
}

func GetAndAddInt64(addr *int64, delta int64) int64 {
	var oldValue int64
	for oldValue = atomic.LoadInt64(addr); !atomic.CompareAndSwapInt64(addr, oldValue, oldValue+delta);
	oldValue = atomic.LoadInt64(addr) {
	}
	return oldValue
}

func GetAndAddInt32(addr *int32, delta int32) int32 {
	var oldValue int32
	for oldValue = atomic.LoadInt32(addr); !atomic.CompareAndSwapInt32(addr, oldValue, oldValue+delta);
	oldValue = atomic.LoadInt32(addr) {
	}
	return oldValue
}
