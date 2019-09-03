package vector

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func TestNewEmptyLockFreeVector(t *testing.T) {
	vec := NewEmptyLockFreeVector()
	size := 10000
	vec.Reserve(size)
	constOffset := 10000000
	for i := 0; i < size; i ++ {
		vec.WriteAt(i, i+constOffset)
	}

	for i := 0; i < size; i ++ {
		value := vec.ReadAt(i)
		if value != i+constOffset {
			t.Fatalf("Error index:%v value:%v", i, value)
		}
	}
}

func TestNewLockFreeVector(t *testing.T) {

}

func TestLockFreeVector_PushBack(t *testing.T) {
	vec := NewEmptyLockFreeVector()
	wg := sync.WaitGroup{}
	timesPerGoroutine := 10000
	goroutineCount := 100
	// PushBack
	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(v *LockFreeVector) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v.PushBack(rand.Int31())
			}
			wg.Done()
		}(vec)
	}
	vec.tryCompleteWrite()
	wg.Wait()
	if vec.Size() != goroutineCount*timesPerGoroutine {
		t.Fatalf("size is wrong:%v", vec.Size())
	}
}

func TestList_PushBack(t *testing.T) {
	ls := list.New()
	wg := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	timesPerGoroutine := 10000
	goroutineCount := 100
	ch := make(chan int, 10)
	wg2.Add(1)
	go func(l *list.List, ch chan int) {
		for value := range ch {
			l.PushBack(value)
		}
		wg2.Done()
	}(ls, ch)

	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(v *list.List) {
			for j := 0; j < timesPerGoroutine; j ++ {
				ch <- rand.Int()
			}
			wg.Done()
		}(ls)
	}
	wg.Wait()
	close(ch)
	wg2.Wait()
	if ls.Len() != goroutineCount*timesPerGoroutine {
		t.Fatalf("size is wrong:%v", ls.Len())
	}
}

func TestLockFreeVector_PushBack_PopBack(t *testing.T) {
	vec := NewEmptyLockFreeVector()
	wg := sync.WaitGroup{}
	timesPerGoroutine := 10000
	goroutineCount := 100
	// PushBack
	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(v *LockFreeVector) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v.PushBack(rand.Int31())
			}
			wg.Done()
		}(vec)
	}
	vec.tryCompleteWrite()
	wg.Wait()
	if vec.Size() != goroutineCount*timesPerGoroutine {
		t.Fatalf("size is wrong:%v", vec.Size())
	}

	// read
	wg2 := sync.WaitGroup{}
	for i := 0; i < goroutineCount; i ++ {
		wg2.Add(1)
		go func(v *LockFreeVector, i int) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v.ReadAt(i*timesPerGoroutine + j)
			}
			wg2.Done()
		}(vec, i)
	}

	// PopBack
	wg3 := sync.WaitGroup{}
	for i := 0; i < goroutineCount; i ++ {
		wg3.Add(1)
		go func(v *LockFreeVector) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v.PopBack()
			}
			wg3.Done()
		}(vec)
	}
	wg3.Wait()
	wg2.Wait()
}

func TestLockFreeVector_PushBack_PopBack_Simultaneously(t *testing.T) {
	wg := sync.WaitGroup{}
	timesPerGoroutine := 10000
	goroutineCount := 100
	vec := NewLockFreeVector(timesPerGoroutine * goroutineCount)
	// PushBack
	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(v *LockFreeVector) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v.PushBack(rand.Int31())
			}
			wg.Done()
		}(vec)
	}

	// read
	wg2 := sync.WaitGroup{}
	for i := 0; i < goroutineCount; i ++ {
		wg2.Add(1)
		go func(v *LockFreeVector, i int) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v.ReadAt(i*timesPerGoroutine + j)
			}
			for j := 0; j < timesPerGoroutine; j ++ {
				v.ReadAt(i*timesPerGoroutine + j)
			}
			for j := 0; j < timesPerGoroutine; j ++ {
				v.ReadAt(i*timesPerGoroutine + j)
			}
			for j := 0; j < timesPerGoroutine; j ++ {
				v.ReadAt(i*timesPerGoroutine + j)
			}
			for j := 0; j < timesPerGoroutine; j ++ {
				v.ReadAt(i*timesPerGoroutine + j)
			}
			for j := 0; j < timesPerGoroutine; j ++ {
				v.ReadAt(i*timesPerGoroutine + j)
			}
			wg2.Done()
		}(vec, i)
	}

	// PopBack
	wg3 := sync.WaitGroup{}
	for i := 0; i < goroutineCount; i ++ {
		wg3.Add(1)
		go func(v *LockFreeVector) {
			for j := 0; j < timesPerGoroutine; j ++ {
				v.PopBack()
			}
			wg3.Done()
		}(vec)
	}
	wg.Wait()
	wg3.Wait()
	wg2.Wait()
}

func TestLockFreeVector_ReadAt_WriteAt(t *testing.T) {
	timesPerGoroutine := 10000
	goroutineCount := 100
	wg := sync.WaitGroup{}
	vec := NewLockFreeVector(timesPerGoroutine * goroutineCount)
	// PushBack
	for i := 0; i < goroutineCount; i ++ {
		wg.Add(1)
		go func(v *LockFreeVector, i2 int) {
			for j := 0; j < timesPerGoroutine; j ++ {
				n := rand.Int()
				idx := i2*timesPerGoroutine + j
				v.WriteAt(idx, n)
				r := v.ReadAt(idx)
				if n != r {
					panic(fmt.Sprintf("err value idx:%v n:%v r:%v\n", idx, n, r))
				}

			}
			wg.Done()
		}(vec, i)
	}
}
