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

	n := int(0xFF)
	fmt.Printf("%X\n", n)
	n = int(uint(n) >> 3)
	fmt.Printf("%X\n", n)
}

func TestNewReferenceArray(t *testing.T) {
	arr := NewReferenceArray(10)
	arr.printElements()
	arr.Get(0)
	arr.Set(0, 100)
	arr.Set(1, 200)
	arr.printElements()
	v0 := arr.Get(0)
	v1 := arr.Get(1)
	fmt.Printf("v0:%d v1:%d\n", v0, v1)
	swap := arr.CompareAndSet(0, 90, 100)
	fmt.Printf("swap:%v\n", swap)
	if swap {
		t.Fatalf("unexpected done swaped!")
	}
	arr.printElements()
	swap = arr.CompareAndSet(0, 100, 101)
	if !swap {
		t.Fatalf("unexpected not swaped!")
	}
	arr.printElements()

	swap = arr.CompareAndSet(0, 101, 102)
	if !swap {
		t.Fatalf("unexpected not swaped!")
	}
	arr.printElements()

	swap = arr.CompareAndSet(1, 200, 202)
	if !swap {
		t.Fatalf("unexpected not swaped!")
	}
	arr.printElements()
}

func TestReferenceArrayStruct(t *testing.T) {
	c := cStruct{15, 1}
	arr := NewReferenceArray(10)
	for i := 0; i < 10; i++ {
		v00 := arr.Get(i)
		set := arr.CompareAndSet(i, nil, c)
		if ! set {
			t.Fatalf("set:%v", set)
		}
		v01 := arr.Get(i)
		v02 := v01.(cStruct)
		fmt.Printf("v00:%v v01:%v v02:%v\n", v00, v01, v02)

		set = arr.CompareAndSet(i, cStruct{15, 1}, cStruct{15, 2}) // succ
		if !set {
			t.Fatalf("set:%v", set)
		}
		v := arr.Get(i).(cStruct)
		var expectedValue = cStruct{15, 2}
		if v != expectedValue {
			t.Fatalf("value is error:%v", v)
		}

		set = arr.CompareAndSet(i, cStruct{15, 1}, cStruct{15, 3}) // fail
		if set {
			t.Fatalf("set:%v", set)
		}

		set = arr.CompareAndSet(i, cStruct{15, 1}, cStruct{15, 3}) // fail
		if set {
			t.Fatalf("set:%v", set)
		}
		set = arr.CompareAndSet(i, cStruct{15, 2}, cStruct{15, 3}) // succ
		if !set {
			t.Fatalf("set:%v", set)
		}
		expectedValue = cStruct{15, 3}
		v = arr.Get(i).(cStruct)
		if v != expectedValue {
			t.Fatalf("value is error:%v", v)
		}
		arr.printElements()
	}
}
