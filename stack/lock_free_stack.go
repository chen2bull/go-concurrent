package stack

import (
	"github.com/cmingjian/go-concurrent/atomic"
	"time"
)

var stackMinDelay = int64(4 * time.Millisecond)
var stackMaxDelay = int64(1024 * time.Millisecond)

type LockFreeStack struct {
	*atomic.BackOff
	top *atomic.Reference // Reference是nil时表示栈为空
}

func NewLockFreeStack() *LockFreeStack {
	backOff := atomic.NewBackOff(stackMinDelay, stackMaxDelay)
	top := atomic.NewReference(nil)
	return &LockFreeStack{backOff, top}
}

type lockFreeStackNode struct {
	value interface{}
	next  *atomic.Reference
}

func newLockFreeStackNode(value interface{}) *lockFreeStackNode {
	next := atomic.NewReference(nil)
	return &lockFreeStackNode{value: value, next: next}
}

func (stack *LockFreeStack) tryPush(node *lockFreeStackNode) bool {
	oldTopRef := stack.top.Get()
	if oldTopRef == nil {
		node.next.Set(nil)
		return stack.top.CompareAndSet(nil, node)
	}
	oldTop := oldTopRef.(*lockFreeStackNode)
	node.next.Set(oldTop)
	return stack.top.CompareAndSet(oldTopRef, node)
}

var stackMaxSpin = 512

func (stack *LockFreeStack) Push(value interface{}) {
	node := newLockFreeStackNode(value)
	spinTime := 0
	for ; true; {
		if stack.tryPush(node) {
			return
		} else {
			spinTime ++
			if spinTime > stackMaxSpin {
				spinTime = 0
				stack.BackOffWait()
			}
		}
	}
}

func (stack *LockFreeStack) tryPop() (interface{}, bool) {
	oldTopRef := stack.top.Get()
	if oldTopRef == nil {
		return nil, false
	}
	oldTop := oldTopRef.(*lockFreeStackNode)
	newTop := oldTop.next.Get()
	if stack.top.CompareAndSet(oldTop, newTop) {
		return oldTop.value, true
	} else {
		return nil, false
	}
}

func (stack *LockFreeStack) Pop() interface{} {
	spinTime := 0
	for ; true; {
		value, ok := stack.tryPop()
		if ok {
			return value
		} else {
			spinTime ++
			if spinTime > stackMaxSpin {
				spinTime = 0
				stack.BackOffWait()
			}
		}
	}
	panic("never here")
}
