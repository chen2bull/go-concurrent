package queue

import (
	"sync"
	"sync/atomic"
)

type BoundedQueue struct {
	enqLock, deqLock          *sync.Mutex
	notEmptyCond, notFullCond *sync.Cond
	head, tail                *boundedNode
	size, capacity            int64
}

func NewBoundedQueue(capacity int64) *BoundedQueue {
	enqLock := &sync.Mutex{}
	deqLock := &sync.Mutex{}
	notFullCond := sync.NewCond(enqLock)
	notEmptyCond := sync.NewCond(deqLock)
	// 最开始是这个节点作为哨兵，每次Deq操作以后，被Deq的节点作为哨兵了
	head := &boundedNode{v: nil}
	tail := head
	return &BoundedQueue{enqLock: enqLock, deqLock: deqLock, notEmptyCond: notEmptyCond, notFullCond: notFullCond,
		head: head, tail: tail, size: 0, capacity: capacity}
}

type boundedNode struct {
	v    interface{}
	next *boundedNode
}

// enqDeal returns true if the queue must wake up Dequeuers.
func (qu *BoundedQueue) enqDeal(v interface{}) bool {
	qu.enqLock.Lock()
	defer qu.enqLock.Unlock()
	for ; atomic.LoadInt64(&qu.size) == atomic.LoadInt64(&qu.capacity); {
		qu.notFullCond.Wait()
	}
	e := &boundedNode{v: v}
	qu.tail.next = e
	qu.tail = e
	//return atomic.GetAndAddInt64(&qu.size, 1) == 0
	// 加完以后等于1,等价于加之前为0(go 没有GetAndAddInt64的原生支持,runtime/internal/atomic)
	return atomic.AddInt64(&qu.size, 1) == 1
}

func (qu *BoundedQueue) Enq(v interface{}) {
	if qu.enqDeal(v) {
		qu.deqLock.Lock()
		defer qu.deqLock.Unlock()
		qu.notEmptyCond.Broadcast()
	}
}

func (qu *BoundedQueue) deqDeal() (interface{}, bool) {
	qu.deqLock.Lock()
	defer qu.deqLock.Unlock()
	for ; atomic.LoadInt64(&qu.size) == 0; {
		qu.notEmptyCond.Wait()
	}
	result := qu.head.next.v
	qu.head = qu.head.next
	return result, atomic.AddInt64(&qu.size, -1) == atomic.LoadInt64(&qu.capacity) - 1
}

func (qu *BoundedQueue) Deq() interface{} {
	result, mustWakeEnqueuers := qu.deqDeal()
	if mustWakeEnqueuers {
		// 惯用的稳妥作法,notFullCond.Wait的调用发生在获得锁enqLock以后
		// 因此，notFullCond.Broadcast也在获得同一把锁的情况下调用
		qu.enqLock.Lock()
		defer qu.enqLock.Unlock()
		qu.notFullCond.Broadcast()
	}
	return result
}
