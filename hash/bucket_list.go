package hash

import (
	"github.com/cmingjian/go-concurrent/atomic"
)

type BucketList struct {
}

type bucketListNode struct {
	key   int
	value interface{}
	next  *atomic.MarkableReference
}

func newBucketListNode(key int, value interface{}) *bucketListNode { // usual constructor
	next := atomic.NewMarkableReference(nil, false) // next永远不为nil，next结构中的value有可能为nil
	return &bucketListNode{key: key, value: value, next: next}
}

func newBucketListSentinelNode(key int) *bucketListNode { // sentinel constructor
	next := atomic.NewMarkableReference(nil, false) // next永远不为nil，next结构中的value有可能为nil
	return &bucketListNode{key: key, next: next}
}

func (node *bucketListNode) getNext() interface{} {
	var currMarked bool
	entryRef, currMarked := node.next.Get()
	for ; currMarked; { // 如果marked表示该元素已经删除了,既然已标记删除,那么value肯定不会为nil
		entry := entryRef.(bucketListNode)
		succRef, succMarked := entry.next.Get()
		node.next.CompareAndSet(entryRef, succRef, true, succMarked)
		entryRef, currMarked = node.next.Get()
	}
	return entryRef
}
