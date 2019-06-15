package queue

import "github.com/cmingjian/go-concurrent/atomic"

type LockFreeQueue struct {
	head,tail *atomic.StampedReference
}

type lockFreeQueueNode struct {
	v interface{}
	next * lockFreeQueueNode
}

func NewLockFreeQueue() *LockFreeQueue {
	return &LockFreeQueue{head: head}
}
