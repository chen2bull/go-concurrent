package atomic

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

// unit:Nanosecond
type BackOff struct {
	minDelay int64
	maxDelay int64
	limit    int64
}

func NewBackOff(minDelay int64, maxDelay int64) *BackOff {
	if minDelay > maxDelay {
		panic(fmt.Sprintf("min can not greater than max!minDelay:%v maxDelay:%v", minDelay, maxDelay))
	}
	return &BackOff{minDelay: minDelay, maxDelay: maxDelay, limit:minDelay}
}

// 指数后退
func (b *BackOff) BackOffWait() {
	oldLimit := atomic.LoadInt64(&b.limit)
	var delay = rand.Int63n(oldLimit)
	atomic.AddInt64(&b.limit, Min(b.maxDelay, 2*oldLimit))
	time.Sleep(time.Duration(delay))
}

