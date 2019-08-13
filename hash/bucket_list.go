package hash

import (
	"github.com/cmingjian/go-concurrent/atomic"
	"math"
)

type BucketList struct {
	head *bucketListNode
}

func NewBucketList() *BucketList {
	head := newBucketListSentinelNode(0)
	head.next = atomic.NewMarkableReference(newBucketListSentinelNode(math.MaxInt32), false)
	return &BucketList{head: head}
}

func (bl *BucketList) Add(value Hashable) bool {
	key := makeRegularKey(value)
	splice := false
	for ; true; {
		pred, curr := bl.head.find(key)
		if curr.key == key { // is the key present?
			return false
		} else {
			entry := newBucketListNode(key, value)
			entry.next.Set(curr, false)
			splice = pred.next.CompareAndSet(curr, entry, false, false)
			if splice {
				return true
			} /*else {
				continue
			}*/
		}
	}
	panic("never here")
}

func (bl *BucketList) Contains(value Hashable) bool {
	key := makeRegularKey(value)
	_, curr := bl.head.find(key)
	return curr.key == key
}

func (bl *BucketList) Remove(value Hashable) bool {
	key := makeRegularKey(value)
	for ; true; {
		pred, curr := bl.head.find(key)
		if curr.key != key {
			return false
		} else {
			snip := pred.next.AttemptMark(curr, true)
			if snip {
				return true
			} /*else {
				continue
			}*/
		}
	}
	panic("never here")
}

func (bl *BucketList) getSentinel(index int32) *BucketList {
	key := makeSentinelKey(index)
	for ;true; {
		pred, curr := bl.head.find(key)
		if curr.key == key {
			return curr.warpAsBucketList()
		} else {
			entry := newBucketListSentinelNode(key)
			entry.next.Set(pred.next.GetReference(), false)
			splice := pred.next.CompareAndSet(curr, entry, false, false)
			if splice {
				return entry.warpAsBucketList()
			} /*else {

			}*/
		}
	}
	panic("never here")
}

type bucketListNode struct {
	key   int32
	value interface{}
	next  *atomic.MarkableReference
}

func newBucketListNode(key int32, value interface{}) *bucketListNode { // usual constructor
	next := atomic.NewMarkableReference(nil, false) // next永远不为nil，next结构中的value有可能为nil
	return &bucketListNode{key: key, value: value, next: next}
}

func newBucketListSentinelNode(key int32) *bucketListNode { // sentinel constructor
	next := atomic.NewMarkableReference(nil, false) // next永远不为nil，next结构中的value有可能为nil
	return &bucketListNode{key: key, next: next}
}

func (node *bucketListNode) getNext() *bucketListNode {
	var currMarked bool
	entryRef, currMarked := node.next.Get()
	for ; currMarked; {
		// 如果marked表示该元素已经删除了,既然已标记删除,那么value肯定不会为nil
		// 列表肯定会有值,因为列表中,有数值为0的Sentinel节点,也有数值为math.MaxInt32的Sentinel节点
		entry := entryRef.(bucketListNode)
		succRef, succMarked := entry.next.Get()
		node.next.CompareAndSet(entryRef, succRef, true, succMarked)
		entryRef, currMarked = node.next.Get()
	}
	entry := entryRef.(bucketListNode)
	return &entry
}

// 只能用与head
func (node *bucketListNode) find(key int32) (*bucketListNode, *bucketListNode) {
	pred := node
	curr := pred.getNext()
	for ; curr.key < key; { // 列表有值为0和math.MaxInt32的哨兵
		pred = curr
		curr = pred.getNext()
	}
	return pred, curr
}

func (node *bucketListNode) warpAsBucketList() *BucketList {
	return &BucketList{head: node}
}
