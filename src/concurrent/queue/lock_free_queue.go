package queue


type lockFreeQueueNode struct {
	v interface{}
	next * lockFreeQueueNode
}
