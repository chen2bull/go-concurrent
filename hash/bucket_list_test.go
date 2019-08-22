package hash

import (
	"testing"
)

type Abc struct {
	a int
	b int
}
//type Comparer interface {
//	Compare(j interface{}) int
//}

func (abc *Abc) Compare(abc2 *Abc) int {
	if abc.a > abc2.a {
		return 1
	}
	if abc.a < abc2.a {
		return -1
	}
	if abc.b > abc2.b {
		return 1
	}
	if abc.b < abc2.b {
		return -1
	}
	return 0
}

func TestHello(t *testing.T) {
	v11 := Abc{a: 1, b: 1}
	v12 := Abc{a: 1, b: 2}
	v21 := Abc{a: 2, b: 1}
	v22 := Abc{a: 2, b: 2}
	if v11.Compare(&v12) == 0 {
		t.Fatalf("fail")
	}
	if v21.Compare(&v22) == 0 {
		t.Fatalf("fail")
	}
}

