package hash

import (
	"math/rand"
	"testing"
)

func TestReverse(t *testing.T) {
	for i := int32(0); i < 0xFFFF; i++ {
		v := rand.Int31n(Mask)
		a := reverse(v)
		b := lookupReverse(v)
		if a != b {
			t.Fatalf("error i:%v a:%v b:%v", i, a, b)
		}
	}
}

const aLen = 0xFFFF

var a [aLen]int32

func init() {
	for i := int32(0); i < int32(aLen); i ++ {
		a[i] = rand.Int31()
	}
}

func BenchmarkLookupReverse(b *testing.B) {
	for i := int32(0); int(i) < b.N; i++ {
		lookupReverse(a[i % aLen])
	}
}

func BenchmarkReverse(b *testing.B) {
	for i := int32(0); int(i) < b.N; i++ {
		reverse(a[i % aLen])
	}
}

func TestTableSizeFor(t *testing.T) {
	for i := int32(0); i < Mask; i++ {
		size := tableSizeFor(i)
		if !isPowerOfTwo(size) {
			t.Fatalf("i:%d size:%d\n", i, size)
		}
	}
}