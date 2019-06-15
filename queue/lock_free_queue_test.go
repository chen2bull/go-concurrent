package queue

import (
	"fmt"
	"testing"
)

func TestNewLockFreeQueue(t *testing.T) {
	qu := NewLockFreeQueue()
	tailRef := qu.tail.GetReference()
	tail := tailRef.(*lockFreeQueueNode)
	nextRef := tail.next.GetReference()
	if nextRef != nil {
		next := nextRef.(*lockFreeQueueNode)
		fmt.Printf("next:%v\n", next)
	}
	fmt.Printf("nextRef:%v\n", nextRef)
}
