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
func (qu *UnBoundedQueue) enqDeal(v interface{}) bool {
	qu.enqLock.Lock()
	defer qu.enqLock.Unlock()
	e := &unBoundedNode{v: v}
	qu.tail.next = e
	qu.tail = e
	return atomic.AddInt64(&qu.size, 1) == 1
}

func (qu *UnBoundedQueue) Enq(v interface{}) {
	if qu.enqDeal(v) {
		qu.deqLock.Lock()
		defer qu.deqLock.Unlock()
		qu.notEmptyCond.Broadcast()
	}
}

func (qu *UnBoundedQueue) Deq() interface{} {
	qu.deqLock.Lock()
	defer qu.deqLock.Unlock()
	for ; atomic.LoadInt64(&qu.size) == 0; {
		qu.notEmptyCond.Wait()
	}
	result := qu.head.next.v
	qu.head = qu.head.next
	atomic.AddInt64(&qu.size, -1)
	return result
}
