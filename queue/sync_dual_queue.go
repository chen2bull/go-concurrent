package queue

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/atomic"
)

type SyncDualQueue struct {
	head, tail *atomic.Reference
}

func NewSyncDualQueue() *SyncDualQueue {
	sentinel := newSyncDualNode(nodeTypeItem, nil)
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

func newSyncDualNode(nodeType syncDualQueueNodeType, value interface{}) *syncDualQueueNode {
	item := atomic.NewReference(value)
	next := atomic.NewReference(nil) // next永远不为nil,但next中的value可能为nil(也可能是*syncDualQueueNode)
	return &syncDualQueueNode{nodeType: nodeType, item: item, next: next}
}

func (node *syncDualQueueNode) changeNodeType(nodeType syncDualQueueNodeType) {
	if node.nodeType == nodeType {
		panic(fmt.Sprintf("type not change,nodeType:%v", nodeType))
	}
	node.nodeType = nodeType
}

