package queue

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/atomic"
	"github.com/cmingjian/go-concurrent/lock"
)

type RecycleLockFreeQueue struct {
	head, tail *atomic.StampedReference
}

type recycleLockFreeQueueNode struct {
	v    interface{}
	next *atomic.StampedReference
}

func newRecycleLockFreeQueueNode(v interface{}) *recycleLockFreeQueueNode {
	// next不会为nil，而是value为nil且stamped为0的StampedReference
	defaultNext := atomic.NewStampedReference(nil, 0)
	return &recycleLockFreeQueueNode{v: v, next: defaultNext}
}

func NewRecycleLockFreeQueue() *RecycleLockFreeQueue {
	sentinelNode := newRecycleLockFreeQueueNode(nil)
	head := atomic.NewStampedReference(sentinelNode, 0)
	tail := atomic.NewStampedReference(sentinelNode, 0)
	return &RecycleLockFreeQueue{head: head, tail: tail}
}

func (qu *RecycleLockFreeQueue) Enq(v interface{}) {
	node := newRecycleLockFreeQueueNode(v)
	for ; true; {
		lastRef, lastStamp := qu.tail.Get()
		last := lastRef.(*recycleLockFreeQueueNode)
		nextRef, nextStamp := last.next.Get()
		if nextRef == nil {
			if last.next.CompareAndSet(nil, node, nextStamp, nextStamp+1) {
				// 入队已完成, 尝试将tail向前移动,只尝试一次
				qu.tail.CompareAndSet(last, node, lastStamp, lastStamp+1)
				return
			}
		} else {
			// try to swing tail to next node
			next := nextRef.(*recycleLockFreeQueueNode)
			qu.tail.CompareAndSet(last, next, lastStamp, lastStamp+1)
		}
	}
}

func (qu *RecycleLockFreeQueue) Deq() interface{} {
	backoff := lock.NewBackOff(backOffMinDelay, backOffMaxDelay)
	for ; true; {
		firstRef, firstStamp := qu.head.Get()
		first := firstRef.(*recycleLockFreeQueueNode)
		lastRef, lastStamp := qu.tail.Get()
		last := lastRef.(*recycleLockFreeQueueNode)
		nextRef, _ := first.next.Get()
		if first == last {
			if nextRef == nil {
				backoff.BackOffWait()
				continue
			}
			next := nextRef.(*recycleLockFreeQueueNode)
			qu.tail.CompareAndSet(last, next, lastStamp, lastStamp+1)
		} else {
			next := nextRef.(*recycleLockFreeQueueNode)
			value := next.v
			if qu.head.CompareAndSet(first, next, firstStamp, firstStamp+1) {
				// TODO: 添加Recycle 功能 
				return value
			}
		}
	}
	panic("never here")
}

func (qu * RecycleLockFreeQueue) PrintAllElement() {
	head, headStamp := qu.head.Get()
	tail, tailStamp := qu.tail.Get()
	fmt.Printf("head:%v %v\n", head, headStamp)
	fmt.Printf("tail:%v %v\n", tail, tailStamp)
	head, _ = qu.head.Get()
	var cur = head.(*recycleLockFreeQueueNode)
	for ;cur.next.GetReference() != nil ; {
		next, nextStamp := cur.next.Get()
		cur = next.(*recycleLockFreeQueueNode)
		fmt.Printf("cur.v,%v, stamp, %v\n", cur.v, nextStamp)
	}
}