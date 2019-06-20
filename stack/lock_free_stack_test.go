package stack

import "testing"

var LockFreeGoroutineNum = 8
var LockFreeCount = LockFreeGoroutineNum * 64
var LockFreePerRoutine = LockFreeCount / LockFreeGoroutineNum

func TestLockFreeStack_Sequential(t *testing.T) {
	var stack = NewLockFreeStack()
	for i := 0; i < LockFreeCount; i++ {
		stack.Push(i)
	}
	for i := 0; i < LockFreeCount; i++ {
		var v = stack.Pop()
		var value = v.(int)
		if value + i != LockFreeCount - 1 {
			t.Fatalf("not equal|value:%d i:%d LockFreeCount:%d", value, i, LockFreeCount)
		}
	}
}

func TestLockFreeStack_ParallelEnq(t *testing.T) {
	stack := NewLockFreeStack()
	doneChan := make(chan bool)
	for i := 0; i < LockFreeGoroutineNum; i++ {
		go lockFreeEnqFunc(stack, i*LockFreePerRoutine, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
	}
	var intMap = make(map[int]bool)
	for i := 0; i < LockFreeCount; i++ {
		v := stack.Pop().(int)
		_, ok := intMap[v]
		if ok {
			t.Fatalf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
}

func lockFreeEnqFunc(stack *LockFreeStack, value int, doneChan chan bool) {
	for i := 0; i < LockFreePerRoutine; i++ {
		stack.Push(value + i)
	}
	doneChan <- true
}

func lockFreeDeqFunc(stack *LockFreeStack, intChan chan int, doneChan chan bool) {
	for i := 0; i < LockFreePerRoutine; i++ {
		v := stack.Pop().(int)
		intChan <- v
	}
	doneChan <- true
}

func lockFreeCheckIntValid(intChan chan int, t *testing.T) {
	var intMap = make(map[int]bool)
	for v := range intChan {
		_, ok := intMap[v]
		if ok {
			t.Fatalf("duplicate pop|v:%d", v)
		} else {
			intMap[v] = true
		}
	}
	if len(intMap) != LockFreeCount {
		t.Fatalf("leak element |len:%v count:%v", len(intMap), LockFreeCount)
	}
}

func TestLockFreeStack_ParallelDeq(t *testing.T) {
	stack := NewLockFreeStack()
	doneChan := make(chan bool)
	for i := 0; i < LockFreeCount; i++ {
		stack.Push(i)
	}
	var intChan = make(chan int, LockFreeCount)
	go lockFreeCheckIntValid(intChan, t)
	for i := 0; i < LockFreeGoroutineNum; i++ {
		go lockFreeDeqFunc(stack, intChan, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
	}
	close(intChan)
}

func TestLockFreeStack_Both(t *testing.T) {
	stack := NewLockFreeStack()
	doneChan := make(chan bool)
	var intChan = make(chan int, LockFreeCount)
	go lockFreeCheckIntValid(intChan, t)
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go lockFreeEnqFunc(stack, i*LockFreePerRoutine, doneChan)
		go lockFreeDeqFunc(stack, intChan, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestLockFreeStack_Both2(t *testing.T) {
	stack := NewLockFreeStack()
	doneChan := make(chan bool)
	var intChan = make(chan int, LockFreeCount)
	go lockFreeCheckIntValid(intChan, t)
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go lockFreeEnqFunc(stack, i*LockFreePerRoutine, doneChan)
		go lockFreeDeqFunc(stack, intChan, doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(intChan)
}

func TestLockFreeStack_Nil(t *testing.T) {
	var stack = NewLockFreeStack()
	stack.Push(nil)
	var v = stack.Pop()
	if v != nil {
		t.Fatalf("v is not nil, v:%v", v)
	}
}

func TestLockFreeStack_BothNil(t *testing.T) {
	stack := NewLockFreeStack()
	doneChan := make(chan bool)
	var interfaceChan = make(chan interface{}, LockFreeCount)
	go func() {
		for v := range interfaceChan {
			if v != nil {
				t.Fatalf("unexpected value|v:%d", v)
			}
		}
	}()
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < LockFreePerRoutine; i++ {
				stack.Push(nil)
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i ++ {
		go func(dChan chan bool) {
			for i := 0; i < LockFreePerRoutine; i++ {
				v := stack.Pop()
				interfaceChan <- v
			}
			dChan <- true
		}(doneChan)
	}
	for i := 0; i < LockFreeGoroutineNum; i++ {
		<-doneChan
		<-doneChan
	}
	close(interfaceChan)
}
