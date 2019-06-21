package combine

import "fmt"

//var threads = 8
//var tries = 1024 * 1024
//var testCount = threads * tries

type Tree struct {
	testArray []bool
	leaf      []*combineNode
}

func NewTree(size int) *Tree {
	nodes := make([]*combineNode, size-1)
	nodesLen := len(nodes)
	nodes[0] = newCombineRootNode()
	for i := 1; i < nodesLen; i++ {
		nodes[i] = newCombineNode(nodes[(i-1)/2])
	}
	leaf := make([]*combineNode, (size+1)/2)
	for i := 0; i < len(leaf); i++ {
		leaf[i] = nodes[nodesLen-i-1] // 协程i被分配到
	}
	return &Tree{leaf: leaf}
}

func (t *Tree) GetAndIncrement(goID int) int {
	stack := make([]*combineNode, 0, 0)
	myLeaf := t.leaf[goID/2]
	node := myLeaf
	for ; node.preCombine(); {
		node = node.parent
	}
	stop := node
	node = myLeaf
	combined := 1
	for ; node != stop; {
		combined = node.combine(combined)
		stack = append(stack, node)
		node = node.parent
	}
	prior := stop.op(combined)
	for i := len(stack) - 1; i >= 0; i-- {
		stack[i].distribute(prior)
	}
	return prior
}

func (t *Tree) printAll() {
	for i := 0; i < len(t.leaf); i++ {
		fmt.Printf("i:%v v:%v\n", i, *t.leaf[i])
	}
}
