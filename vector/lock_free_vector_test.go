package vector

import (
	"testing"
)

func TestNewEmptyLockFreeVector(t *testing.T) {
	vec := NewEmptyLockFreeVector()
	size := 10000
	vec.Reserve(size)
	constOffset := 10000000
	for i := 0; i < size; i ++ {
		vec.WriteAt(i, i + constOffset)
	}

	for i := 0; i < size; i ++ {
		value := vec.ReadAt(i)
		if value != i + constOffset {
			t.Fatalf("Error index:%v value:%v", i, value)
		}
	}
}

func TestNewLockFreeVector(t *testing.T) {

}