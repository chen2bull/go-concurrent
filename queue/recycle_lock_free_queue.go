package queue

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/atomic"
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

func (q *RecycleLockFreeQueue) Enq(v interface{}) {
	node := newRecycleLockFreeQueueNode(v)
	for ; true; {
		lastRef, lastStamp := q.tail.Get()
		last := lastRef.(*recycleLockFreeQueueNode)
		nextRef, nextStamp := last.next.Get()
		//next := nextRef.(*recycleLockFreeQueueNode)
		if nextRef == nil {
			if last.next.CompareAndSet(nil, node, nextStamp, nextStamp+1) {
				// 入队已完成, 尝试将tail向前移动,只尝试一次
				q.tail.CompareAndSet(last, node, lastStamp, lastStamp+1)
				return
			}
		} else {
			// try to swing tail to next node
			q.tail.CompareAndSet(last, node, lastStamp, lastStamp+1)
		}
	}
}

func (q *RecycleLockFreeQueue) Deq() interface{} {
	for ; true; {
		firstRef, firstStamp := q.head.Get()
		first := firstRef.(*recycleLockFreeQueueNode)
		lastRef, lastStamp := q.tail.Get()
		last := lastRef.(*recycleLockFreeQueueNode)
		nextRef, _ := last.next.Get()
		if first == last {
			if nextRef == nil {
				panic("todo:add condition")
			}
			next := nextRef.(*recycleLockFreeQueueNode)
			q.tail.CompareAndSet(last, next, lastStamp, lastStamp+1)
		} else {
			next := nextRef.(*recycleLockFreeQueueNode)
			value := next.v
			if q.head.CompareAndSet(first, next, firstStamp, firstStamp+1) {
				return value
			}
		}
	}
	panic("never here")
}

func (q * RecycleLockFreeQueue) PrintAllElement() {
	head, headStamp := q.head.Get()
	tail, tailStamp := q.tail.Get()
	fmt.Printf("head:%v %v\n", head, headStamp)
	fmt.Printf("tail:%v %v\n", tail, tailStamp)
	head, _ = q.head.Get()
	var cur = head.(*recycleLockFreeQueueNode)
	for ;true ; {
		nextRef, stamp := cur.next.Get()
		cur := nextRef.(*recycleLockFreeQueueNode)
		fmt.Printf("cur.v,%v, stamp, %v\n", cur.v, stamp)
	}
}