package hash

type Hashable interface {
	hashCode() int32
}

type Map interface {
	size() int

}