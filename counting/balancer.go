package counting

import "sync/atomic"

type Balancer struct {
	p *int32
}

func NewBlancer() *Balancer {
	a := balancerZero
	return &Balancer{&a}
}

const (
	balancerZero  = int32(0)
	balancerOne= int32(1)
)

func (b Balancer) Traverse() int {
	for ;; {
		oldValue := atomic.LoadInt32(b.p)
		if oldValue == balancerZero {
			if atomic.CompareAndSwapInt32(b.p, oldValue, balancerOne) {
				return int(balancerZero)
			}
		} else {
			if atomic.CompareAndSwapInt32(b.p, oldValue, balancerZero) {
				return int(balancerOne)
			}
		}
	}
}
