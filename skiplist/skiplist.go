package skiplist

import "github.com/cmingjian/go-concurrent/atomic"

const MaxLevel = 32
const Probability = 0.25

type LockFreeSkipList struct {
}

type Comparable interface {
	Compare(value interface{}) int
}

type skiplistNode struct {
	key      Comparable
	value    interface{}
	next     []*atomic.MarkableReference
	topLevel int32
}

func NewSkiplistSentinel(key Comparable) *skiplistNode {
	next := make([]*atomic.MarkableReference, MaxLevel+1, MaxLevel+1)
	for i := 0; i < len(next); i++ {
		next[i] = atomic.NewMarkableReference(nil, false)
	}
	return &skiplistNode{key: key, value: nil, next: next, topLevel: MaxLevel}
}

