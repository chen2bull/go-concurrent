package queue

type Pool interface {
	Put(interface{})
	Get() interface{}
}
