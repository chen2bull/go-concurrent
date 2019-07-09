package counting

import (
	"fmt"
	"math/rand"
	"testing"
)

const bitonicSize = 1024

//func bitonicN(size int) {
//	b := NewBitonic(size)
//	m := make([]int, size, size)
//	old := size - 1
//	for i := 0; i < bitonicSize; i++ {
//		ret := b.Traverse(0)
//		if (old+1)%size != ret {
//			panic(fmt.Errorf("wrong value|ret:%v old:%v", ret, old))
//		}
//		old = ret
//		m[ret]++
//	}
//	//fmt.Print(m)
//	checkStep(m[:])
//}
//
//func TestBitonicTraverse(t *testing.T) {
//	for size := 2; size <= bitonicSize; size = size << 1 {
//		bitonicN(size)
//	}
//}

func bitonicNRandom(size int) {
	b := NewBitonic(size)
	m := make([]int, size, size)
	for i := 0; i < bitonicSize; i++ {
		ret := b.Traverse(rand.Int() % size)
		fmt.Printf("%v\n",ret)
		m[ret]++
	}
	checkStep(m)
}

func TestBitonicTraverseRandom(t *testing.T) {
	for size := 2; size <= 8; size = size << 1 {
		bitonicNRandom(size)
	}
}

