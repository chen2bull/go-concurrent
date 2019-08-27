package atomic

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestReflectInfo(t *testing.T) {
	i := 5565
	s := "any string that is long long long long long "
	var v interface{}
	v = i
	fmt.Printf("SizeOf(v)=%v\n", unsafe.Sizeof(v))
	v = s
	fmt.Printf("SizeOf(v)=%v\n", unsafe.Sizeof(v))
	v = t
	fmt.Printf("SizeOf(v)=%v\n", unsafe.Sizeof(v))
	a := make([]interface{}, 10)
	fmt.Printf("Type of a=%T\n", a)


	a2 := (*slice)(unsafe.Pointer(&a))
	fmt.Printf("Type of a2=%T\n", a2)
	fmt.Printf("Type of a2.array=%T\n", a2.array)
}