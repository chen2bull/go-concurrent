package queue

import (
	"github.com/cmingjian/go-concurrent/atomic"
	"github.com/cmingjian/go-concurrent/lock"
)

type SyncDualQueue struct {
	head, tail *atomic.Reference
}

func NewSyncDualQueue() *SyncDualQueue {
	sentinel := newSyncDualItemNode(nil)
	head := atomic.NewReference(sentinel)
	tail := atomic.NewReference(sentinel)
	return &SyncDualQueue{head: head, tail: tail}
}

const (
	nodeTypeItem = iota + 1
	nodeTypeReservation
)

type syncDualQueueNodeType int

type syncDualQueueNode struct {
	nodeType syncDualQueueNodeType
	item     *atomic.Reference
	next     *atomic.Reference
}

type syncDualValue struct {
	v interface{}
}

func newSyncDualItemNode(value interface{}) *syncDualQueueNode {
	// item永远不为nil,但item中的value可能是nil,可能是syncDualValue
	// item的value为nil时,表示元素已经被取出
	// 注意输入参数为nil时,item的value值为syncDualValue{v:nil}(并不是nil)
	item := atomic.NewReference(&syncDualValue{v: value})
	next := atomic.NewReference(nil) // next永远不为nil,但next中的value可能为nil(也可能是*syncDualQueueNode)
	return &syncDualQueueNode{nodeType: nodeTypeItem, item: item, next: next}
}

func newSyncDualReservationNode() *syncDualQueueNode {
	next := atomic.NewReference(nil)
	item := atomic.NewReference(nil)
	return &syncDualQueueNode{nodeType: nodeTypeReservation, item: item, next: next}
}

func (qu *SyncDualQueue) Enq(value interface{}) {
	offer := newSyncDualItemNode(value)
	backoff := lock.NewBackOff(backOffMinDelay, backOffMaxDelay)
	var counter = 0
	for ; true; {
		tailRef := qu.tail.Get()
		tail := tailRef.(*syncDualQueueNode)
		headRef := qu.head.Get()
		head := headRef.(*syncDualQueueNode)
		// head == tail 表示队列为空,
		// tail.nodeType == nodeTypeItem表示tail或者tail的前面有nodeTypeItem的节点
		if head == tail || tail.nodeType == nodeTypeItem {
			tNextRef := tail.next.Get()
			if tNextRef == nil {
				if tail.next.CompareAndSet(tail, offer) {
					qu.tail.CompareAndSet(tail, offer)
					for ; offer.item.Get() != nil; {
						counter ++ // spin
						if counter > syncDualSpinCount {
							counter = 0
							backoff.BackOffWait()
						}
					}
					headRef := qu.head.Get()
					head := headRef.(*syncDualQueueNode)
					if offer == head.next.Get() {
						qu.head.CompareAndSet(head, offer)
					}
					return
				}
			} else { // qu.tail还不是"最后的节点"(所以要先调整)
				next := tNextRef.(*syncDualQueueNode)
				qu.tail.CompareAndSet(tail, next)
			}
		} else {
			hNextRef := head.next.Get()
			if hNextRef == nil {
				counter = counter + 5 // spins here is costly than continuously call offer.item.Get()
				if counter > syncDualSpinCount {
					counter = 0
					backoff.BackOffWait()
				}
				continue
			}
			next := hNextRef.(*syncDualQueueNode)
			success := next.item.CompareAndSet(nil, &syncDualValue{v: value})
			qu.head.CompareAndSet(head, next)
			if success {
				return
			}
		}
	}
}

var syncDualSpinCount = 500

func (qu *SyncDualQueue) Deq() interface{} {
	offer := newSyncDualReservationNode()
	backoff := lock.NewBackOff(backOffMinDelay, backOffMaxDelay)
	var counter = 0
	for ; true; {
		headRef := qu.head.Get()
		head := headRef.(*syncDualQueueNode)
		tailRef := qu.tail.Get()
		tail := tailRef.(*syncDualQueueNode)
		// head == tail 表示队列为空,
		// 元素都是从头往后取的,且Enq时只要tail是nodeTypeItem都会在tail后面加入节点,
		if head == tail || tail.nodeType == nodeTypeReservation {
			tNextRef := tail.next.Get()
			if tNextRef == nil {
				if tail.next.CompareAndSet(nil, offer) {
					qu.tail.CompareAndSet(tail, offer)
					for ; offer.item.Get() == nil; {
						counter ++ // spin
						if counter > syncDualSpinCount {
							counter = 0
							backoff.BackOffWait()
						}
					}
					headRef = qu.head.Get()
					head = headRef.(*syncDualQueueNode)
					if offer == head.next.Get() {
						qu.head.CompareAndSet(head, offer)
					}
					item := offer.item.Get().(*syncDualValue)
					return item.v
				}
			} else {
				next := tNextRef.(*syncDualQueueNode)
				qu.tail.CompareAndSet(tail, next)
			}
		} else {
			hNextRef := head.next.Get()
			if hNextRef == nil {
				counter = counter + 5 // spins here is costly than continuously call offer.item.Get()
				if counter > syncDualSpinCount {
					counter = 0
					backoff.BackOffWait()
				}
				continue
			}
			next := hNextRef.(*syncDualQueueNode)
			itemRef := next.item.Get()
			item := itemRef.(*syncDualValue)
			success := next.item.CompareAndSet(item, nil) // 如果失败了,一定是另一个协程Deq成功了
			qu.head.CompareAndSet(head, next)
			if success {
				return item.v
			}
		}
	}
	panic("never here")
}
