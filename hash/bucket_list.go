package hash

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/atomic"
	"math"
	"reflect"
)

type BucketList struct {
	head *bucketListNode
}

func NewBucketList() *BucketList {
	head := newBucketListSentinelNode(0)
	// math.MaxInt64 是独一无二的哨兵值,比其他所有‘哨兵节点’和‘元素节点’的哈希值都大，永远都不会被删除
	head.next = atomic.NewMarkableReference(newBucketListSentinelNode(math.MaxInt64), false)
	// 注意：
	// 1、列表中math.MaxInt64这个哨兵节点其实是‘哨兵的哨兵’，运行过程中，最大的输入keyHash绝对会小于math.MaxInt64，
	// 有math.MaxInt64这个哨兵的存在，可以保证对于所有合法的输入keyHash，bucketListNode的find方法都不会报错
	// (math.MaxInt64不是合法输入)
	// 2、0节点和math.MaxInt64一样，是唯二不能用getSentinel获取的哨兵
	// 实现LockFreeHashMap的时候，会确保不会用getSentinel获取0节点
	return &BucketList{head: head}
}

// notes: Put nil value is not equivalent as Remove
func (bl *BucketList) Put(keyHash int64, key interface{}, value interface{}) {
	splice := false
	for ; true; {
		pred, curr := bl.head.find(keyHash, key)
		if curr.keyHash == keyHash && reflect.DeepEqual(curr.key, key) { // is the key present?
			entry := newBucketListNode(keyHash, key, value)
			entry.next.Set(curr, true) // !!!entry.next的marked是true,因此插入完成的时候,curr也被删除了
			splice = pred.next.CompareAndSet(curr, entry, false, false)
			if splice {
				entry.getNext() // 触发真的删除
				return
			} /*else {
				continue
			}*/
		} else {
			entry := newBucketListNode(keyHash, key, value)
			entry.next.Set(curr, false)
			splice = pred.next.CompareAndSet(curr, entry, false, false)
			if splice {
				return
			} /*else {
				continue
			}*/
		}
	}
	panic("never here")
}

func (bl *BucketList) Contains(keyHash int64, key interface{}) bool {
	_, curr := bl.head.find(keyHash, key)
	return curr.keyHash == keyHash && reflect.DeepEqual(curr.key, key)
}

func (bl *BucketList) Remove(keyHash int64, key interface{}) bool {
	for ; true; {
		pred, curr := bl.head.find(keyHash, key)
		if curr.keyHash == keyHash && reflect.DeepEqual(curr.key, key) {
			snip := pred.next.AttemptMark(curr, true)
			if snip {
				pred.getNext() // 触发真的删除
				return true
			} /*else {
				continue
			}*/
		} else {
			return false
		}
	}
	panic("never here")
}

func (bl *BucketList) Get(keyHash int64, key interface{}) interface{} {
	_, curr := bl.head.find(keyHash, key)
	if curr.keyHash == keyHash && reflect.DeepEqual(curr.key, key) { // is the key present?
		return curr.value
	} else {
		return nil
	}
}
func (bl *BucketList) getSentinelByBucket(bucketIdx int64) *BucketList {
	keyHash := makeSentinelKey(bucketIdx)
	return bl.getSentinelByHash(keyHash)
}

func (bl *BucketList) getSentinelByHash(keyHash int64) *BucketList {
	for ; true; {
		pred, curr := bl.head.find(keyHash, keyHash)
		if curr.keyHash == keyHash {
			return curr.warpAsBucketList()
		} else {
			entry := newBucketListSentinelNode(keyHash)
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

func (bl *BucketList) printAllElements() {
	curr := bl.head
	for ; true; {
		currRef, marked := curr.next.Get()
		if currRef == nil {
			fmt.Printf("iterate done!\n")
			return
		}
		if marked {
			fmt.Printf("iterate deleted node!\n")
			return
		}
		curr = currRef.(*bucketListNode)
		fmt.Printf("keyHash:%56b key:%v value:%v\n", curr.keyHash, curr.key, curr.value)
	}
}

type bucketListNode struct {
	keyHash int64
	key     interface{}
	value   interface{}
	next    *atomic.MarkableReference
}

func newBucketListNode(keyHash int64, key interface{}, value interface{}) *bucketListNode { // usual constructor
	next := atomic.NewMarkableReference(nil, false) // next永远不为nil，next结构中的value有可能为nil
	return &bucketListNode{keyHash: keyHash, key: key, value: value, next: next}
}

// 哨兵节点有个性质: keyHash值等于key值
func newBucketListSentinelNode(keyHash int64) *bucketListNode { // sentinel constructor
	next := atomic.NewMarkableReference(nil, false) // next永远不为nil，next结构中的value有可能为nil
	return &bucketListNode{keyHash: keyHash, key: keyHash, next: next}
}

func (node *bucketListNode) getNext() *bucketListNode {
	var currMarked bool
	entryRef, currMarked := node.next.Get()
	for ; currMarked; {
		// 如果marked表示该元素已经删除了,既然已标记删除,那么value肯定不会为nil
		// 列表肯定会有值,因为列表中,有数值为0的Sentinel节点,也有数值为math.MaxInt64的Sentinel节点
		entry := entryRef.(*bucketListNode)
		succRef, succMarked := entry.next.Get()
		node.next.CompareAndSet(entryRef, succRef, true, succMarked)
		entryRef, currMarked = node.next.Get()
	}
	entry := entryRef.(*bucketListNode)
	return entry
}

func (node *bucketListNode) find(keyHash int64, key interface{}) (*bucketListNode, *bucketListNode) {
	pred := node
	curr := pred.getNext()
	for ; curr.keyHash < keyHash || (curr.keyHash == keyHash && !reflect.DeepEqual(curr.key, key)); {
		pred = curr
		curr = pred.getNext()
	}
	return pred, curr
}

func (node *bucketListNode) warpAsBucketList() *BucketList {
	return &BucketList{head: node}
}
