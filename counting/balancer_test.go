package counting

import (
	"fmt"
	"testing"
)

const TestSize = 256

func TestBalancer_Traverse(t *testing.T) {
	b := NewBlancer()
	m := []int{0, 0}
	for i := 0; i < TestSize; i++ {
		m[b.Traverse()]++
	}
	checkStep(m)
}

func checkStep(m []int) {
	step := false
	for i := 1; i < len(m); i++ {
		if m[i] != m[i-1] {
			if !step && m[i] == m[i-1]+1 {
				step = true
			} else {
				fmt.Printf("m:%v", m)
				panic("Step property failed")
			}
		}
	}
}
