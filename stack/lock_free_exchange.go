package stack

import (
	"github.com/cmingjian/go-concurrent/atomic"
	"time"
)

type LockFreeExchanger struct {
	slot *atomic.StampedReference
}

func NewLockFreeExchange() *LockFreeExchanger {
	slot := atomic.NewStampedReference(nil, exchangeWAITING)
	return &LockFreeExchanger{slot: slot}
}

const (
	exchangeEmpty = iota
	exchangeWAITING
	exchangeBUSY
)

func (e *LockFreeExchanger) Exchange(myItem interface{}, timeout time.Duration) (interface{}, bool) {
	timeBound := time.Now().Add(timeout)
	for ; true; {
		if timeBound.Before(time.Now()) {
			return nil, false
		}
		yrItem, stamp := e.slot.Get()
		switch stamp {
		case exchangeEmpty:
			if e.slot.CompareAndSet(yrItem, myItem, exchangeEmpty, exchangeWAITING) {
				
			}
		case exchangeWAITING:
			if e.slot.CompareAndSet(yrItem, myItem, exchangeWAITING, exchangeBUSY) {
				return yrItem, true
			}
		case exchangeBUSY:

		}

	}
	panic("never here")
}
