package lock

import (
	"sync"
)

type RWLocker interface {
	RLock()
	RUnlock()
	WLock()
	WUnlock()
}

type UselessRWLock struct {
	mu *sync.Mutex // 不要使用非指针类型的Mutex
	cond *sync.Cond
	readers int
	writing bool
}

func NewUselessRWLock() *UselessRWLock {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	return &UselessRWLock{mu: mu, cond: cond}
}

func (rwl *UselessRWLock) RLock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	for ; rwl.writing;  { // 常见用法,等待某个条件成立,否则等待,后面就能认为某个条件已经成立了
		rwl.cond.Wait()
	}
	rwl.readers++
}

func (rwl *UselessRWLock) RUnlock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	rwl.readers--
	if rwl.readers == 0 {
		rwl.cond.Broadcast()
	}
}

func (rwl *UselessRWLock) WLock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	for ; rwl.readers > 0;  {
		rwl.cond.Wait()
	}
	rwl.writing = true
	rwl.cond.Broadcast()
}

func (rwl *UselessRWLock) WUnlock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	rwl.writing = false
	rwl.cond.Broadcast()
}

type FifoRWLock struct {
	mu *sync.Mutex
	cond *sync.Cond
	acquireReaders int64
	releaseReaders int64
	writing bool
}

func NewFifoRWLock() * FifoRWLock {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	return &FifoRWLock{mu: mu, cond: cond}
}

func (rwl *FifoRWLock) RLock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	for ; rwl.writing;  { // 常见用法,等待某个条件成立,否则等待,后面就能认为某个条件已经成立了
		rwl.cond.Wait()
	}
	rwl.acquireReaders++
}

func (rwl *FifoRWLock) RUnlock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	rwl.releaseReaders++
	if rwl.releaseReaders == rwl.acquireReaders { // 只用==和!=,即使发生溢出,也能正常工作
		rwl.cond.Broadcast()
	}
}

func (rwl *FifoRWLock) WLock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	// 惯用代码块 start
	for ; rwl.writing;  {
		rwl.cond.Wait()
	}
	rwl.writing = true // 必需在下一次Wait前面修改rwl.writing的值
	// 惯用代码块 end

	for ; rwl.acquireReaders != rwl.releaseReaders;  {
		rwl.cond.Wait()
	}
}

func (rwl *FifoRWLock) WUnlock()  {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()
	rwl.writing  = false
	rwl.cond.Broadcast()
}
