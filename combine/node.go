package combine

import (
	"fmt"
	"sync"
)

const (
	combineFirst = iota // 一个主动线程已访问该节点
	combineSecond
	combineResult
	combineIdle
	combineRoot
)

type combineNode struct {
	mu                      *sync.Mutex
	cond                    *sync.Cond
	locked                  bool
	cStatus                 int
	firstValue, secondValue int
	result                  int
	parent                  *combineNode
}

func newCombineRootNode() *combineNode {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	cStatus := combineRoot
	return &combineNode{mu: mu, cond: cond, cStatus: cStatus, locked: false}
}

func newCombineNode(myParent *combineNode) *combineNode {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	return &combineNode{mu: mu, cond: cond, cStatus: combineIdle, locked: false, parent: myParent}
}

func (n *combineNode) preCombine() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	for ; n.locked; {
		n.cond.Wait()
	}
	switch n.cStatus {
	case combineIdle:
		n.cStatus = combineFirst
		return true
	case combineFirst:
		n.locked = true
		n.cStatus = combineSecond
		return false
	case combineRoot:
		return false
	default:
		panic(fmt.Sprintf("unexpected Node status:%v", n.cStatus))
	}
}

func (n *combineNode) combine(combined int) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	for ; n.locked; {
		n.cond.Wait()
	}
	n.locked = true
	n.firstValue = combined
	switch n.cStatus {
	case combineFirst:
		return n.firstValue
	case combineSecond:
		return n.firstValue + n.secondValue
	default:
		panic(fmt.Sprintf("unexpected Node status:%v", n.cStatus))
	}
}

func (n *combineNode) op(combined int) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	switch n.cStatus {
	case combineRoot:
		oldValue := n.result
		n.result += combined
		return oldValue
	case combineSecond:
		n.secondValue = combined
		n.locked = false
		n.cond.Broadcast()
		for ; n.cStatus != combineResult; {
			n.cond.Wait()
		}
		n.locked = false
		n.cond.Broadcast()
		n.cStatus = combineIdle
		return n.result
	default:
		panic(fmt.Sprintf("unexpected Node status:%v", n.cStatus))
	}
}

func (n *combineNode) distribute(prior int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	switch n.cStatus {
	case combineFirst:
		n.cStatus = combineIdle
		n.locked = false
	case combineSecond:
		n.result = prior + n.firstValue
		n.cStatus = combineResult
	default:
		panic(fmt.Sprintf("unexpected Node status:%v", n.cStatus))
	}
	n.cond.Broadcast()
}
