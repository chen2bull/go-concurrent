package counting

import "sync/atomic"

type Balancer struct {
	p *int32
}

func NewBlancer() *Balancer {
	a := blancerHigh
	return &Balancer{&a}
}

const (
	blancerHigh = int32(0)
	blancerLow  = int32(1)
)

func (b Balancer) Traverse() int {
	oldValue := atomic.LoadInt32(b.p)
	for ; !atomic.CompareAndSwapInt32(b.p, oldValue, -oldValue); {
		oldValue = atomic.LoadInt32(b.p)
	}
	return int(oldValue)
}
