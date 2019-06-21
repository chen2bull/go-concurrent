package stack

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var ExchangerGoroutineNum = 2
var ExchangerCount = ExchangerGoroutineNum * 64
var ExchangerPerRoutine = ExchangerCount / ExchangerGoroutineNum

func runExchange(exchanger *LockFreeExchanger, goID int, wg *sync.WaitGroup, timeOut time.Duration,
	iChan chan interface{}) {
	wg2 := sync.WaitGroup{}
	for j := 0; j < ExchangerPerRoutine; j++ {
		v := goID*ExchangerPerRoutine + j
		//fmt.Printf("start exchange|v:%v i:%v j:%v\n", v, goID, j)
		wg2.Add(1)
		go func() {
			exValue, ok := exchanger.Exchange(v, timeOut)
			//fmt.Printf("end exchange|v:%v i:%v j:%v\n", v, goID, j)
			if ok {
				iChan <- exValue
			} else {
				panic(fmt.Sprintf("unexpected timeout|v:%v i:%v j:%v\n", v, goID, j))
			}
			wg2.Done()
		}()
		wg2.Wait()
	}
	wg.Done()
}

func TestLockFreeExchanger_Exchange(t *testing.T) {
	exchanger := NewLockFreeExchanger()
	var interfaceChan = make(chan interface{}, ExchangerCount)
	wg := sync.WaitGroup{}
	timeOut := time.Second * 5
	for i := 0; i < ExchangerGoroutineNum; i ++ {
		wg.Add(1)
		go runExchange(exchanger, i, &wg, timeOut, interfaceChan)
	}
	doneChan := make(chan bool)
	go func(iChan chan interface{}) {
		var intMap = make(map[int]bool)
		for v := range iChan {
			fmt.Printf("value,%v\n", v)
			value := v.(int)
			if intMap[value] {
				t.Errorf("duplicate pop|v:%d", v)
			} else {
				intMap[value] = true
			}
		}
		doneChan <- true
	}(interfaceChan)
	wg.Wait()
	close(interfaceChan)
	<-doneChan
}

func runExchangeNil(exchanger *LockFreeExchanger, wg *sync.WaitGroup, timeOut time.Duration, iChan chan interface{}) {
	wg2 := sync.WaitGroup{}
	for j := 0; j < ExchangerPerRoutine; j++ {
		wg2.Add(1)
		go func() {
			//fmt.Printf("start exchange|v:%v i:%v j:%v\n", v, goID, j)
			exValue, ok := exchanger.Exchange(nil, timeOut)
			if ok {
				//fmt.Printf("end exchange|v:%v i:%v j:%v\n", v, goID, j)
				iChan <- exValue
			} else {
				panic(fmt.Sprintf("unexpected timeout|j:%v\n", j))
			}
			wg2.Done()
		}()
		wg2.Wait()
	}
	wg.Done()
}

func TestLockFreeExchanger_ExchangeNil(t *testing.T) {
	exchanger := NewLockFreeExchanger()
	var interfaceChan = make(chan interface{}, ExchangerCount)
	wg := sync.WaitGroup{}
	timeOut := time.Second * 5
	for i := 0; i < ExchangerGoroutineNum; i ++ {
		wg.Add(1)
		go runExchangeNil(exchanger, &wg, timeOut, interfaceChan)
	}
	doneChan := make(chan bool)
	go func(iChan chan interface{}) {
		for v := range iChan {
			if v != nil {
				t.Fatalf("unexpected value|v:%d", v)
			}
		}
		doneChan <- true
	}(interfaceChan)
	wg.Wait()
	close(interfaceChan)
	<-doneChan
}
