package hash

import (
	"math/rand"
	"testing"
)

func TestReverse32(t *testing.T) {
	for i := int32(mask32 | hiMask32); i > int32(0); i = i - rand.Int31n(0xFFF) {
		a := reverse32(i)
		b := lookupReverse32(i)
		if a != b {
			t.Fatalf("error i:%v a32:%v b:%v", i, a, b)
		}
	}
}

func TestReverse64(t *testing.T) {
	for i := int64(mask64 | hiMask64); i > int64(0); i = i - rand.Int63n(0xFFFFFFFFF) {
		a := reverse64(i)
		b := lookupReverse64(i)
		if a != b {
			t.Fatalf("error i:%v a32:%v b:%v", i, a, b)
		}
		//fmt.Printf("a32:%v b:%v\n", a32, b)
	}
}

func TestReverse(t *testing.T) {
	for i := mask | hiMask; i > 0; i = i - int(rand.Int63n(0xFFFFFFFFF)) {
		a := reverse(i)
		b := lookupReverse(i)
		if a != b {
			t.Fatalf("error i:%v a32:%v b:%v", i, a, b)
		}
	}
}

const aLen = 0xFFFF

var a32 [aLen]int32
var a64 [aLen]int64
var a [aLen]int

func init() {
	for i := 0; i < aLen; i ++ {
		a32[i] = rand.Int31()
		a64[i] = rand.Int63()
		a[i] = rand.Int()
	}
}

func BenchmarkLookupReverse32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupReverse32(a32[i%aLen])
	}
}

func BenchmarkReverse32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reverse32(a32[i%aLen])
	}
}

func BenchmarkLookupReverse64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupReverse64(a64[i%aLen])
	}
}

func BenchmarkReverse64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reverse64(a64[i%aLen])
	}
}

func BenchmarkLookupReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupReverse(a[i%aLen])
	}
}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reverse(a[i%aLen])
	}
}

func TestTableSizeFor32(t *testing.T) {
	for i := int32(mask32 | hiMask32); i > int32(0); i = i - rand.Int31n(0xFFF) {
		size := tableSizeFor32(i)
		if !isPowerOfTwo64(int64(size)) {
			t.Fatalf("error i:%v size:%v", i, size)
		}
	}
}

func TestTableSizeFor64(t *testing.T) {
	for i := int64(mask64 | hiMask64); i > int64(0); i = i - rand.Int63n(0xFFFFFFFFF) {
		size := tableSizeFor64(i)
		if !isPowerOfTwo64(size) {
			t.Fatalf("error i:%v size:%v", i, size)
		}
	}
}

func TestTableSizeFor(t *testing.T) {
	for i := mask | hiMask; i > 0; i = i - int(rand.Int63n(0xFFFFFFFFF)) {
		size := tableSizeFor(i)
		if !isPowerOfTwo64(int64(size)) {
			t.Fatalf("error i:%v size:%v", i, size)
		}
	}
}
