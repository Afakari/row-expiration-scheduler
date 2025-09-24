package main

import (
	"container/heap"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var mu sync.Mutex
var cond = sync.NewCond(&mu)

func main() {
	pq := &PriorityQueue{}
	expChan := make(chan string, 100)
	heap.Init(pq)

	go watchdog(pq, expChan)

	go func() {
		i := 0
		for {
			addItem(pq, &Item{
				value:    strconv.Itoa(i),
				priority: time.Now().Add(time.Duration(2+i%5) * time.Second).UnixMilli(),
			})
			i++
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for id := range expChan {
			fmt.Println("Expired:", id, "at", time.Now())
		}
	}()

	select {}
}

func addItem(pq *PriorityQueue, item *Item) {
	mu.Lock()
	heap.Push(pq, item)
	cond.Signal()
	mu.Unlock()
}

func watchdog(pq *PriorityQueue, expChan chan<- string) {
	for {
		mu.Lock()
		for pq.Len() == 0 {
			cond.Wait()
		}

		item := (*pq)[0]
		now := time.Now().UnixMilli()

		if item.priority <= now {
			expired := heap.Pop(pq).(*Item)
			mu.Unlock()

			select {
			case expChan <- expired.value:
			default:
				fmt.Println("channel is full")
			}
		} else {
			sleepDuration := time.Duration(item.priority-now) * time.Millisecond
			timer := time.NewTimer(sleepDuration)
			mu.Unlock()

			<-timer.C
		}
	}
}
