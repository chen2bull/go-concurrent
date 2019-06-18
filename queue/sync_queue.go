package queue

import (
	"sync"
)

type SyncQueue struct {
	sync.Locker
	cond      *sync.Cond
	enqueuing bool
	item      *syncQueueNode
}

type syncQueueNode struct {
	value interface {}
}

func NewSyncQueue() *SyncQueue {
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)
	return &SyncQueue{Locker: &mu, cond: cond, enqueuing: false}
}

func (qu *SyncQueue) Deq() interface{} {
	qu.Lock()
	defer qu.Unlock()
	for ; qu.item == nil; {
		qu.cond.Wait()
	}
	ele := *qu.item
	qu.item = nil
	qu.cond.Broadcast()
	return ele.value
}

func (qu *SyncQueue) Enq(value interface{}) {
	qu.Lock()
	defer qu.Unlock()
	for ; qu.enqueuing; {
		qu.cond.Wait()
	}
	qu.enqueuing = true // my turn starts
	qu.item = &syncQueueNode{value:value}
	qu.cond.Broadcast()
	for ; qu.item != nil; {
		qu.cond.Wait() // 如果元素没有被取出就继续等
	}
	qu.enqueuing = false // my turn ends
	qu.cond.Broadcast()
}
