package queue

import "github.com/cmingjian/go-concurrent/atomic"

type LockFreeQueue struct {
	head, tail *atomic.StampedReference
}

type lockFreeQueueNode struct {
	v    interface{}
	next *atomic.StampedReference
}

func newLockFreeQueueNode(v interface{}) *lockFreeQueueNode {
	// next不会为nil，而是value为nil且stamped为0的StampedReference
	defaultNext := atomic.NewStampedReference(nil, 0)
	return &lockFreeQueueNode{v: v, next: defaultNext}
}

func NewLockFreeQueue() *LockFreeQueue {
	sentinelNode := newLockFreeQueueNode(nil)
	head := atomic.NewStampedReference(sentinelNode, 0)
	tail := atomic.NewStampedReference(sentinelNode, 0)
	return &LockFreeQueue{head: head, tail: tail}
}

func (q *LockFreeQueue) Enq(v interface{}) {
	node := newLockFreeQueueNode(v)
	for ; true; {
		lastRef, lastStamp := q.tail.Get()
		last := lastRef.(*lockFreeQueueNode)
		nextRef, nextStamp := last.next.Get()
		//next := nextRef.(*lockFreeQueueNode)
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

func (q *LockFreeQueue) Deq() interface{} {
	for ; true; {
		firstRef, firstStamp := q.head.Get()
		first := firstRef.(*lockFreeQueueNode)
		lastRef, lastStamp := q.tail.Get()
		last := lastRef.(*lockFreeQueueNode)
		nextRef, _ := last.next.Get()
		if firstRef == lastRef {
			if nextRef == nil {
				panic("todo:add condition")
			}
			next := nextRef.(*lockFreeQueueNode)
			q.tail.CompareAndSet(last, next, lastStamp, lastStamp+1)
		} else {
			next := nextRef.(*lockFreeQueueNode)
			value := next.v
			if q.head.CompareAndSet(first, next, firstStamp, firstStamp+1) {
				return value
			}
		}
	}
	panic("never here")
}
