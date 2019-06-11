package util

import (
	"sync"
)

type Semaphore struct {
	capacity int
	mu       *sync.Mutex
	cond     *sync.Cond
	state    int
}

func NewSemaphore(cap int) Semaphore {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	return Semaphore{mu: mu, cond: cond, capacity: cap, state: 0}
}

func (se *Semaphore) Acquire() {
	se.mu.Lock()
	defer se.mu.Unlock()
	for ; se.state == se.capacity; {
		se.cond.Wait()
	}
	se.state++
}

func (se *Semaphore) Release() {
	se.mu.Lock()
	defer se.mu.Unlock()
	se.state--
	se.cond.Broadcast()
}
