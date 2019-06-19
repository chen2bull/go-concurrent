package queue

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/atomic"
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

func (node *syncDualQueueNode) changeNodeType(nodeType syncDualQueueNodeType) {
	if node.nodeType == nodeType {
		panic(fmt.Sprintf("type not change,nodeType:%v", nodeType))
	}
	node.nodeType = nodeType
}

func (qu *SyncDualQueue) Enq(value interface{}) {
	offer := newSyncDualItemNode(value)
	//backoff := lock.NewBackOff(backOffMinDelay, backOffMaxDelay)
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
					qu.tail.CompareAndSet(tail, offer) // 只尝试一次,失败的话下一次Enq会重新调整qu.tail(在LINEA)
					for ; offer.item.Get() != nil; {
						//backoff.BackOffWait() // spin
					}
					headRef := qu.head.Get()
					head := headRef.(*syncDualQueueNode)
					if offer == head.next.Get() {
						qu.head.CompareAndSet(head, offer)
					}
					//fmt.Printf("Enq AAAA:%v\n", value)
					return
				}
			} else { // qu.tail还不是"最后的节点"(所以要先调整)
				next := tNextRef.(*syncDualQueueNode)
				qu.tail.CompareAndSet(tail, next)
			}
		} else {
			hNextRef := head.next.Get()
			if hNextRef == nil {
				continue
			}
			next := hNextRef.(*syncDualQueueNode)
			success := next.item.CompareAndSet(nil, &syncDualValue{v: value})
			qu.head.CompareAndSet(head, next)
			if success {
				//fmt.Printf("Enq BBBB:%v\n", value)
				return
			}
		}
	}
}

func (qu *SyncDualQueue) Deq() interface{} {
	offer := newSyncDualReservationNode()
	//backoff := lock.NewBackOff(backOffMinDelay, backOffMaxDelay)
	for ; true; {
		headRef :=qu.head.Get()
		head := headRef.(*syncDualQueueNode)
		tailRef :=qu.tail.Get()
		tail := tailRef.(*syncDualQueueNode)
		// head == tail 表示队列为空,
		// tail.nodeType == nodeTypeReservation表示tail以及tail的前面有没有nodeTypeItem节点!!!!
		if head == tail || tail.nodeType == nodeTypeReservation {
			tNextRef := tail.next.Get()
			if tNextRef == nil {
				if tail.next.CompareAndSet(nil, offer) {
					qu.tail.CompareAndSet(tail, offer)
					for;offer.item.Get() == nil; {
						//backoff.BackOffWait() // spin
					}
					headRef = qu.head.Get()
					head = headRef.(*syncDualQueueNode)
					if offer == head.next.Get() {
						qu.head.CompareAndSet(head, offer)
					}
					item :=  offer.item.Get().(*syncDualValue)
					//fmt.Printf("Deq AAAAA:%v\n", item.v)
					return item.v
				}
			} else {
				next := tNextRef.(*syncDualQueueNode)
				qu.tail.CompareAndSet(tail, next)
			}
		} else {
			hNextRef := head.next.Get()
			if hNextRef == nil {
				continue
			}
			next := hNextRef.(*syncDualQueueNode)
			itemRef := next.item.Get()
			item := itemRef.(*syncDualValue)
			success := next.item.CompareAndSet(item, nil)
			qu.head.CompareAndSet(head, next)
			if success {
				//fmt.Printf("Deq BBBB:%v\n", item.v)
				return item.v
			}
		}
	}
	panic("never here")
}
