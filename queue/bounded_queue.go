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
func (q *BoundedQueue) enqDeal(v interface{}) bool {
	q.enqLock.Lock()
	defer q.enqLock.Unlock()
	for ; atomic.LoadInt64(&q.size) == atomic.LoadInt64(&q.capacity); {
		q.notFullCond.Wait()
	}
	e := &boundedNode{v: v}
	q.tail.next = e
	q.tail = e
	//return atomic.GetAndAddInt64(&q.size, 1) == 0
	// 加完以后等于1,等价于加之前为0(go 没有GetAndAddInt64的原生支持,runtime/internal/atomic)
	return atomic.AddInt64(&q.size, 1) == 1
}

func (q *BoundedQueue) Enq(v interface{}) {
	if q.enqDeal(v) {
		q.deqLock.Lock()
		defer q.deqLock.Unlock()
		q.notEmptyCond.Broadcast()
	}
}

func (q *BoundedQueue) deqDeal() (interface{}, bool) {
	q.deqLock.Lock()
	defer q.deqLock.Unlock()
	for ; atomic.LoadInt64(&q.size) == 0; {
		q.notEmptyCond.Wait()
	}
	result := q.head.next.v
	q.head = q.head.next
	return result, atomic.AddInt64(&q.size, -1) == atomic.LoadInt64(&q.capacity) - 1
}

func (q *BoundedQueue) Deq() interface{} {
	result, mustWakeEnqueuers := q.deqDeal()
	if mustWakeEnqueuers {
		// 惯用的稳妥作法,notFullCond.Wait的调用发生在获得锁enqLock以后
		// 因此，notFullCond.Broadcast也在获得同一把锁的情况下调用
		q.enqLock.Lock()
		defer q.enqLock.Unlock()
		q.notFullCond.Broadcast()
	}
	return result
}
