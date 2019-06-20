package queue

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/atomic"
	"time"
)

type LockFreeQueue struct {
	head, tail *atomic.Reference
}

type lockFreeQueueNode struct {
	v    interface{}
	next *atomic.Reference
}

func newLockFreeQueueNode(v interface{}) *lockFreeQueueNode {
	next := atomic.NewReference(nil) // next永远不为nil，next结构中的value有可能为nil
	return &lockFreeQueueNode{v: v, next: next}
}

func NewLockFreeQueue() *LockFreeQueue {
	sentinel := newLockFreeQueueNode(nil)
	head := atomic.NewReference(sentinel)
	tail := atomic.NewReference(sentinel)
	return &LockFreeQueue{head: head, tail: tail}
}

func (qu *LockFreeQueue) Enq(v interface{}) {
	node := newLockFreeQueueNode(v)
	for ; true; {
		lastRef := qu.tail.Get()
		last := lastRef.(*lockFreeQueueNode)
		nextRef := last.next.Get()
		if nextRef == nil { // "最后一个节点"的下一节点为nil
			if last.next.CompareAndSet(nil, node) {
				// 只尝试一次,失败的话表示另一个协程已经调整过qu.tail,下一次Enq会重新调整qu.tail(在LINEA)
				qu.tail.CompareAndSet(last, node)
				return
			}
		} else { // qu.tail还不是"最后的节点"(所以要先调整对)
			next := nextRef.(*lockFreeQueueNode)
			qu.tail.CompareAndSet(last, next) // LINEA:上面调整qu.tail失败了,要重新调整
		}
	}
}

var backOffMinDelay = int64(4 * time.Millisecond)
var backOffMaxDelay = int64(1024 * time.Millisecond)

func (qu *LockFreeQueue) Deq() interface{} {
	backoff := atomic.NewBackOff(backOffMinDelay, backOffMaxDelay)
	for ; true; {
		firstRef := qu.head.Get()
		first := firstRef.(*lockFreeQueueNode)
		lastRef := qu.tail.Get()
		last := lastRef.(*lockFreeQueueNode)
		nextRef := first.next.Get()
		if first == last {
			if nextRef == nil {
				backoff.BackOffWait()
				continue
			}
			next := nextRef.(*lockFreeQueueNode)
			qu.tail.CompareAndSet(last, next)
		} else {
			next := nextRef.(*lockFreeQueueNode)
			value := next.v
			if qu.head.CompareAndSet(first, next) {
				return value
			}
		}
	}
	panic("never here")
}

func (qu *LockFreeQueue) PrintAllElement() {
	fmt.Printf("head:%v\n", qu.head.Get())
	fmt.Printf("tail:%v\n", qu.tail.Get())
	first := qu.head.Get()
	var cur = first.(*lockFreeQueueNode)
	for ; cur.next.Get() != nil; {
		next := cur.next.Get()
		cur = next.(*lockFreeQueueNode)
		fmt.Printf("cur.v,%v\n", cur.v)
	}
}
