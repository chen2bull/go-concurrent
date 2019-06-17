package queue

import (
	"fmt"
	"testing"
)

func TestNewRecycleLockFreeQueue(t *testing.T) {
	qu := NewLockFreeQueue()
	tailRef := qu.tail.GetReference()
	tail := tailRef.(*recycleLockFreeQueueNode)
	nextRef := tail.next.GetReference()
	if nextRef != nil {
		next := nextRef.(*recycleLockFreeQueueNode)
		fmt.Printf("next:%v\n", next)
	}
	fmt.Printf("nextRef:%v\n", nextRef)
}
