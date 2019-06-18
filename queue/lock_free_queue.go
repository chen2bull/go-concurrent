package queue

import (
	"fmt"
	atomic2 "github.com/cmingjian/go-concurrent/atomic"
	"github.com/cmingjian/go-concurrent/lock"
	"time"
)

type LockFreeQueue struct {
	head, tail *atomic2.Reference
}

type lockFreeQueueNode struct {
	v    interface{}
	next *atomic2.Reference
}

func newLockFreeQueueNode(v interface{}) *lockFreeQueueNode {
	next := atomic2.NewReference(nil) // next永远不为nil，next结构中的value有可能为nil
	return &lockFreeQueueNode{v: v, next: next}
}

func NewLockFreeQueue() *LockFreeQueue {
	sentinel := newLockFreeQueueNode(nil)
	head := atomic2.NewReference(sentinel)
	tail := atomic2.NewReference(sentinel)
	return &LockFreeQueue{head: head, tail: tail}
}

func (queue *LockFreeQueue) Enq(v interface{}) {
	node := newLockFreeQueueNode(v)
	for ; true; {
		lastRef := queue.tail.Get()
		last := lastRef.(*lockFreeQueueNode)
		nextRef := last.next.Get()
		if nextRef == nil {
			if last.next.CompareAndSet(nil, node) {
				queue.tail.CompareAndSet(last, node)
				return
			}
		} else {
			next := nextRef.(*lockFreeQueueNode)
			queue.tail.CompareAndSet(last, next)
		}
	}
}

var backOffMinDelay  = int64(8 * time.Millisecond)
var backOffMaxDelay  = int64(1024 * time.Millisecond)

func (queue *LockFreeQueue) Deq() interface{} {
	backoff := lock.NewBackOff(backOffMinDelay, backOffMaxDelay)
	for ; true; {
		firstRef := queue.head.Get()
		first := firstRef.(*lockFreeQueueNode)
		lastRef := queue.tail.Get()
		last := lastRef.(*lockFreeQueueNode)
		nextRef := first.next.Get()
		if first == last {
			if nextRef == nil {
				backoff.BackOffWait()
				continue
			}
			next := nextRef.(*lockFreeQueueNode)
			queue.tail.CompareAndSet(last, next)
		} else {
			next := nextRef.(*lockFreeQueueNode)
			value := next.v
			if queue.head.CompareAndSet(first, next) {
				return value
			}
		}
	}
	panic("never here")
}

func (queue * LockFreeQueue) PrintAllElement() {
	fmt.Printf("head:%v\n", queue.head.Get())
	fmt.Printf("tail:%v\n", queue.tail.Get())
	first := queue.head.Get()
	var cur = first.(*lockFreeQueueNode)
	for ;cur.next.Get() != nil ; {
		next := cur.next.Get()
		cur = next.(*lockFreeQueueNode)
		fmt.Printf("cur.v,%v\n", cur.v)
	}
}
