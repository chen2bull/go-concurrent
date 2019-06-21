package stack

import (
	"math/rand"
	"time"
)

type EliminationArray struct {
	exchangers []*LockFreeExchanger
	duration   time.Duration
}

func NewEliminationArray(capacity int) *EliminationArray {
	exchangers := make([]*LockFreeExchanger, capacity, capacity)
	duration := time.Second
	for i := 0; i < capacity; i++ {
		exchangers = append(exchangers, NewLockFreeExchanger())
	}
	return &EliminationArray{exchangers: exchangers, duration: duration}
}

func (eArray *EliminationArray) visit(value interface{}, randRange int) (interface{}, bool) {
	var slot = rand.Intn(randRange)
	capacity := cap(eArray.exchangers)
	slot = slot % capacity
	return eArray.exchangers[slot].Exchange(value, eArray.duration)
}
