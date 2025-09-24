package scheduler

import (
	"container/heap"
	"sync"
	"time"
)

type Task struct {
	Value    string
	Deadline time.Time
}

type Scheduler struct {
	mu       sync.Mutex
	cond     *sync.Cond
	pq       PriorityQueue
	stopCh   chan struct{}
	callback func(string)
	wg       sync.WaitGroup
}

func NewScheduler(cb func(string)) *Scheduler {
	s := &Scheduler{
		pq:       PriorityQueue{},
		stopCh:   make(chan struct{}),
		callback: cb,
	}
	s.cond = sync.NewCond(&s.mu)
	heap.Init(&s.pq)
	return s
}

func (s *Scheduler) Add(value string, deadline time.Time) {
	item := &Item{
		Value:    value,
		Priority: deadline.UnixMilli(),
	}
	s.mu.Lock()
	heap.Push(&s.pq, item)
	s.cond.Signal()
	s.mu.Unlock()
}

func (s *Scheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			s.mu.Lock()
			for s.pq.Len() == 0 {
				s.cond.Wait()
			}

			item := s.pq[0]
			now := time.Now().UnixMilli()

			if item.Priority <= now {
				expired := heap.Pop(&s.pq).(*Item)
				s.mu.Unlock()
				s.callback(expired.Value)
			} else {
				sleepDuration := time.Duration(item.Priority-now) * time.Millisecond
				timer := time.NewTimer(sleepDuration)
				s.mu.Unlock()

				select {
				case <-timer.C:
				case <-s.stopCh:
					timer.Stop()
					return
				}
			}
		}
	}()
}

func (s *Scheduler) Load(tasks []Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range tasks {
		item := &Item{
			Value:    task.Value,
			Priority: task.Deadline.UnixMilli(),
		}
		heap.Push(&s.pq, item)
	}

	s.cond.Signal()
}
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}
