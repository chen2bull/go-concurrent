package hash

import "testing"

func TestReverse(t *testing.T) {
	max := 1<<24 - 1
	for i := 0; i < max; i++ {
		a := reverse(i)
		b := lookupReverse(i)
		if a != b {
			t.Fatalf("error a:%v b:%v", a, b)
		}
	}
}

func BenchmarkLookupReverse(b *testing.B) {
	max := 1<<24 - 1
	for i := 0; i < b.N; i++ {
		for i := 0; i < max; i++ {
			lookupReverse(i)
		}
	}
}

func BenchmarkReverse(b *testing.B) {
	max := 1<<24 - 1
	for i := 0; i < b.N; i++ {
		for i := 0; i < max; i++ {
			reverse(i)
		}
	}
}
