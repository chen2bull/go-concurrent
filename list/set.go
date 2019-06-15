package list

type Set interface {
	Add(v interface{}) bool
	Remove(v interface{}) bool
	contains(v interface{}) bool
}
