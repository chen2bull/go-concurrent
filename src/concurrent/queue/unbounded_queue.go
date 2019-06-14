package queue

import (
	"sync"
	"sync/atomic"
)

type UnBoundedQueue struct {
	enqLock, deqLock *sync.Mutex
	notEmptyCond     *sync.Cond
	head, tail       *unBoundedNode
	size             int64
}

func NewUnBoundedQueue() *UnBoundedQueue {
	enqLock := &sync.Mutex{}
	deqLock := &sync.Mutex{}
	notEmptyCond := sync.NewCond(deqLock)
	head := &unBoundedNode{v: nil}
	tail := head
	return &UnBoundedQueue{enqLock: enqLock, deqLock: deqLock, notEmptyCond: notEmptyCond,
		head: head, tail: tail, size: 0}
}

type unBoundedNode struct {
	v    interface{}
	next *unBoundedNode
}

// enqDeal returns true if the queue must wake up Dequeuers.
func (q *UnBoundedQueue) enqDeal(v interface{}) bool {
	q.enqLock.Lock()
	defer q.enqLock.Unlock()
	e := &unBoundedNode{v: v}
	q.tail.next = e
	q.tail = e
	return atomic.AddInt64(&q.size, 1) == 1
}

func (q *UnBoundedQueue) Enq(v interface{}) {
	if q.enqDeal(v) {
		q.deqLock.Lock()
		defer q.deqLock.Unlock()
		q.notEmptyCond.Broadcast()
	}
}

func (q *UnBoundedQueue) Deq() interface{} {
	q.deqLock.Lock()
	defer q.deqLock.Unlock()
	for ; q.size == 0; {
		q.notEmptyCond.Wait()
	}
	result := q.head.next.v
	q.head = q.head.next
	return result
}
