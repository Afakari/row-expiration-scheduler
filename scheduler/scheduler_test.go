package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	var mu sync.Mutex
	expired := make(map[string]bool)

	cb := func(value string) {
		mu.Lock()
		expired[value] = true
		mu.Unlock()
	}

	s := NewScheduler(cb)
	s.Start()
	defer s.Stop()

	s.Add("task-1", time.Now().Add(100*time.Millisecond))
	s.Add("task-2", time.Now().Add(300*time.Millisecond))

	time.Sleep(150 * time.Millisecond)
	mu.Lock()
	if !expired["task-1"] {
		t.Error("task-1 should have expired")
	}
	mu.Unlock()

	time.Sleep(200 * time.Millisecond)
	mu.Lock()
	if !expired["task-2"] {
		t.Error("task-2 should have expired")
	}
	mu.Unlock()
}
