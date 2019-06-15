package queue

import (
	"sync"
)

type LockedQueue struct {
	sync.Locker
	notFull           *sync.Cond
	notEmpty          *sync.Cond
	items             []interface{}
	tail, head, count int
}

func NewLockedQueue(cap int) *LockedQueue {
	var l sync.Locker = &sync.Mutex{}
	notFull := sync.NewCond(l)
	notEmpty := sync.NewCond(l)
	items := make([]interface{}, cap, cap)
	return &LockedQueue{l, notFull, notEmpty, items, 0, 0, 0}
}

func (queue *LockedQueue) Enq(v interface{}) {
	queue.Lock()
	defer queue.Unlock()
	for ; cap(queue.items) == queue.count; {
		queue.notFull.Wait()
	}
	queue.items[queue.tail] = v
	queue.tail ++
	if queue.tail == cap(queue.items) {
		queue.tail = 0
	}
	queue.count ++
	queue.notEmpty.Broadcast()	// 注意不是Signal,否则有可能发生唤醒丢失
}

func (queue *LockedQueue) Deq() interface{} {
	queue.Lock()
	defer queue.Unlock()
	for ; queue.count == 0; {
		queue.notEmpty.Wait()
	}
	v := queue.items[queue.head]
	queue.head ++
	if queue.head == cap(queue.items) {
		queue.head = 0
	}
	queue.count --
	queue.notFull.Broadcast()	// 注意不是Signal,否则有可能发生唤醒丢失
	return v
}
