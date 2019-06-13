package atomic

import "sync/atomic"

// Atomically sets to the given value and returns the previous value.
func GetAndSetInt64(addr *int64, newValue int64) int64 {
	var oldValue int64
	for oldValue = *addr; !atomic.CompareAndSwapInt64(addr, oldValue, newValue); oldValue = *addr {
	}
	return oldValue
}

// Atomically sets to the given value and returns the previous value.
func GetAndSetInt32(addr *int32, newValue int32) int32 {
	var oldValue int32
	for oldValue = *addr; !atomic.CompareAndSwapInt32(addr, oldValue, newValue); {
	}
	return oldValue
}
