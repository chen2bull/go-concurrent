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

func (qu *LockedQueue) Enq(v interface{}) {
	qu.Lock()
	defer qu.Unlock()
	for ; cap(qu.items) == qu.count; {
		qu.notFull.Wait()
	}
	qu.items[qu.tail] = v
	qu.tail ++
	if qu.tail == cap(qu.items) {
		qu.tail = 0
	}
	qu.count ++
	qu.notEmpty.Broadcast() // 注意不是Signal,否则有可能发生唤醒丢失
}

func (qu *LockedQueue) Deq() interface{} {
	qu.Lock()
	defer qu.Unlock()
	for ; qu.count == 0; {
		qu.notEmpty.Wait()
	}
	v := qu.items[qu.head]
	qu.head ++
	if qu.head == cap(qu.items) {
		qu.head = 0
	}
	qu.count --
	qu.notFull.Broadcast() // 注意不是Signal,否则有可能发生唤醒丢失
	return v
}
