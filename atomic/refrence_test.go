package atomic

import (
	"testing"
)

func TestNewReference(t *testing.T) {
	ref := NewReference(100)
	value := ref.Get().(int)
	if value != 100 {
		t.Fatalf("wrong value:%v", value)
	}

	ok := ref.CompareAndSet(100, 101)
	if !ok {
		t.Fatalf("unexpected fail!")
	}
	value = ref.Get().(int)
	if value != 101 {
		t.Fatalf("wrong value:%v", value)
	}

	ok2 := ref.CompareAndSet(100, 102)
	if ok2 {
		t.Fatalf("unexpected succ!")
	}
	value = ref.Get().(int)
	if value != 101 {
		t.Fatalf("wrong value:%v", value)
	}

	ok3 := ref.CompareAndSet(101, 102)
	if !ok3 {
		t.Fatalf("unexpected fail!")
	}
	value = ref.Get().(int)
	if value != 102 {
		t.Fatalf("wrong value:%v", value)
	}
}