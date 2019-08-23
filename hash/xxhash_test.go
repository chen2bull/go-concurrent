package hash

import (
	"fmt"
	"math/rand"
	"testing"
	"unsafe"
)


func TestPtrSize(t *testing.T) {
	fmt.Printf("Expr:uintptr(0):\t%T:%v\n", uintptr(0), uintptr(0))
	fmt.Printf("Expr:^uintptr(0):\t%T:%v\n", ^uintptr(0), ^uintptr(0))
	fmt.Printf("Expr:^uintptr(0)>>63:\t%T:%v\n", ^uintptr(0)>>63, ^uintptr(0)>>63)
	fmt.Printf("Expr:4 << (^uintptr(0) >> 63):\t%T:%v\n", 4 << (^uintptr(0) >> 63), 4 << (^uintptr(0) >> 63))
}

func TestInitFunc(t *testing.T) {
	fmt.Printf("hashkey:\t%T:%v\n", hashkey, hashkey)
}

func Test_memhash(t *testing.T) {
	seed := rand.Int31()
	for i:= int64(0); i < int64(3000); i++ {
		val1 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 1)) & Mask
		val2 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 2)) & Mask
		val3 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 4)) & Mask
		val4 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 8)) & Mask

		val12 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 1)) & Mask
		val22 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 2)) & Mask
		val32 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 4)) & Mask
		val42 := int32(memhash(unsafe.Pointer(&i), uintptr(seed), 8)) & Mask

		if val1 != val12 {
			t.Fatalf("error  i:%v val1:%v val12:%v\n", i, val1, val12)
		}
		if val2 != val22 {
			t.Fatalf("error  i:%v val2:%v val22:%v\n", i, val2, val22)
		}
		if val3 != val32 {
			t.Fatalf("error  i:%v val3:%v val32:%v\n", i, val3, val32)
		}
		if val4 != val42 {
			t.Fatalf("error  i:%v val4:%v val42:%v\n", i, val4, val42)
		}
	}
	type stringInnerStruct struct {
		str unsafe.Pointer
		len int
	}
	// random string
	for i:= 0; i < 3000; i++ {
		randomString := RandStringRunes()
		strPtr := unsafe.Pointer(&randomString)
		x := (*stringInnerStruct)(strPtr)
		val1 := memhash(x.str, uintptr(seed), uintptr(x.len))
		val2 := memhash(x.str, uintptr(seed), uintptr(x.len))
		if val1 != val2 {
			t.Fatalf("error randomString:%v val1:%v val2:%v\n", randomString, val1, val2)
		}
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ一二三四五六七八九十百千万亿")

func RandStringRunes() string {
	n := rand.Int31n(1024)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Test_memhash32(t *testing.T) {
	seed := rand.Int31()
	addNum := rand.Int31()
	for i:= int32(0); i < int32(3000); i++ {
		c := i + addNum
		val1 := int32(memhash32(unsafe.Pointer(&c), uintptr(seed))) & Mask
		val2 := int32(memhash32(unsafe.Pointer(&c), uintptr(seed))) & Mask
		if val1 != val2 {
			t.Fatalf("error  val1:%v val2:%v\n", val1, val2)
		}
	}
}

func Test_memhash32_conflict(t *testing.T) {
	seed := rand.Int31()
	conflictCount := 0
	existMap := make(map[int32]bool)
	testCount := int32(0xFFFFFF)
	for i:= int32(0); i < testCount; i++ {
		val1 := int32(memhash32(unsafe.Pointer(&i), uintptr(seed))) & Mask
		exist, _ := existMap[val1]
		if exist {
			conflictCount = conflictCount + 1
		}
		existMap[val1] = true
	}
	fmt.Printf("conflictCount:%v conflictRate:%v\n", conflictCount, float64(conflictCount)/float64(testCount))
}


func Test_memhash64(t *testing.T) {
	seed := rand.Int63()
	addNum := rand.Int63()
	for i:= int64(0); i < int64(3000); i++ {
		c := i + addNum
		val1 := int32(memhash64(unsafe.Pointer(&c), uintptr(seed))) & Mask
		val2 := int32(memhash64(unsafe.Pointer(&c), uintptr(seed))) & Mask
		if val1 != val2 {
			t.Fatalf("error  val1:%v val2:%v\n", val1, val2)
		}
	}
}
