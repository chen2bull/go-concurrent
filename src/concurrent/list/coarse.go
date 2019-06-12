package list

type Element struct {
	v interface{}
	key int
	next *Element
}
// A CoarseList is a list with single mutex.
//
// The head element's key is always math.MinInt64 and the last element's key is alwasys math.MaxInt64
func NewElement(v interface{}) *Element {

	return &Element{v: v}
}


