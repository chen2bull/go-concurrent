package stack

import (
	"github.com/cmingjian/go-concurrent/atomic"
	"time"
)

type LockFreeExchanger struct {
	slot *atomic.StampedReference
}

func NewLockFreeExchanger() *LockFreeExchanger {
	slot := atomic.NewStampedReference(nil, exchangeEmpty)
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
				for ; timeBound.After(time.Now()); {
					yrItem, stamp = e.slot.Get()
					if stamp == exchangeBUSY {
						e.slot.Set(nil, exchangeEmpty)
						return yrItem, true
					}
				}
				if e.slot.CompareAndSet(myItem, nil, exchangeWAITING, exchangeEmpty) {
					return nil, false
				} else { // 本线程A超时以后,另一线程B却把WAITING状态改成BUSY状态了
					yrItem, _ = e.slot.Get()
					e.slot.Set(nil, exchangeEmpty)
					return yrItem, true
				}
			}
		case exchangeWAITING:
			if e.slot.CompareAndSet(yrItem, myItem, exchangeWAITING, exchangeBUSY) {
				return yrItem, true
			}
		case exchangeBUSY:
			// spin
		}

	}
	panic("never here")
}
