package counting

import "sync/atomic"

type blancer struct {
	p *int32
}

func newBlancer() *blancer {
	a := blancerHigh
	return &blancer{&a}
}

const (
	blancerHigh = int32(1)
	blancerLow  = int32(-1)
)

func (b blancer) traverse() int32 {
	oldValue := atomic.LoadInt32(b.p)
	for ; !atomic.CompareAndSwapInt32(b.p, oldValue, -oldValue); {
		oldValue = atomic.LoadInt32(b.p)
	}
	return oldValue
}
