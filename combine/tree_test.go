package combine

import (
	"fmt"
	"github.com/cmingjian/go-concurrent/atomic"
	"sync"
	"testing"
)

const THREADS = 128
const TRIES = 1024 * 16

// 测试发现: 在我的4核x86机器上,即使开到128个goroutine,CombineTree性能依然不如简单的GetAndAddInt64
func TestTree_GetAndIncrement(t *testing.T) {
	tree := NewTree(THREADS)
	intChan := make(chan int, THREADS * TRIES)
	wg := sync.WaitGroup{}
	go func() {
		boolArray := [THREADS * TRIES]bool{}
		for v := range intChan {
			if boolArray[v] {
				panic(fmt.Sprintf("ERROR duplicate value %d\n", v))
			} else {
				boolArray[v] = true
			}
		}
		//tree.printAll()
		for v := 0; v < len(boolArray); v++ {
			if !boolArray[v] {
				panic(fmt.Sprintf("missing value at %d\n", v))
			}
		}
	}()
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go func(goID int) {
			for j := 0; j < TRIES; j++ {
				v := tree.GetAndIncrement(goID)
				intChan <- v
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(intChan)
}

var incrValue int64

func TestTree_Incr(t *testing.T) {
	incrValue = 0
	intChan := make(chan int64, THREADS * TRIES)
	wg := sync.WaitGroup{}
	go func() {
		boolArray := [THREADS * TRIES]bool{}
		for v := range intChan {
			if boolArray[v] {
				panic(fmt.Sprintf("ERROR duplicate value %d\n", v))
			} else {
				boolArray[v] = true
			}
		}
		for v := 0; v < len(boolArray); v++ {
			if !boolArray[v] {
				panic(fmt.Sprintf("missing value at %d\n", v))
			}
		}
	}()

	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go func(goID int) {
			for j := 0; j < TRIES; j++ {
				v := atomic.GetAndAddInt64(&incrValue, 1)
				intChan <- v
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(intChan)
}
