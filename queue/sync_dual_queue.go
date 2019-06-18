package queue

type SyncDualQueue struct {

}

const (
	nodeTypeItem = iota + 1
	nodeTypeReservation
)

type syncDualQueueNodeType int

type syncDualQueueNode struct {
	nodeType syncDualQueueNodeType

}