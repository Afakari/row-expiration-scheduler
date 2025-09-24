package main

import (
	"fmt"
	"row-expiration-scheduler/scheduler"
	"time"
)

func main() {
	cb := func(value string) {
		fmt.Println("Expired:", value, "at", time.Now())
	}

	s := scheduler.NewScheduler(cb)
	s.Start()

	for i := range 5 {
		s.Add(fmt.Sprintf("task-%d", i), time.Now().Add(time.Duration(2+i)*time.Second))
	}

	time.Sleep(10 * time.Second)
	s.Stop()
}
