package progress

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestIteration(t *testing.T) {

	fmt.Println("Expecting:   10 it    4 it/s")
	mon := NewMonitor("Progress")
	for i := 0; i < 10; i++ {
		mon.Increment(1)
		time.Sleep(250 * time.Millisecond)
	}
	mon.Close()
	fmt.Println()
}

func TestCounter(t *testing.T) {

	fmt.Println("Expecting:   10 it    4 it/s    error: 4")
	mon := NewMonitor("Progress")
	for i := 0; i < 10; i++ {
		mon.Increment(1)
		if i%3 == 0 {
			mon.IncrementCounter("error", 1)
		}
		time.Sleep(250 * time.Millisecond)
	}
	mon.Close()
	fmt.Println()
}

func TestTwoCounters(t *testing.T) {

	fmt.Println("Expecting:   10 it    4 it/s    too big: 3    too small: 4")
	mon := NewMonitor("Progress")
	mon.StartCounter("too big")
	mon.StartCounter("too small")
	for i := 0; i < 10; i++ {
		mon.Increment(1)
		if i < 4 {
			mon.IncrementCounter("too small", 1)
		}
		if i > 6 {
			mon.IncrementCounter("too big", 1)
		}
		time.Sleep(250 * time.Millisecond)
	}
	mon.Close()
	fmt.Println()
}

func TestConcurrentCounter(t *testing.T) {

	fmt.Println("Expecting:   20 it    8 it/s    error: 8")
	var wg sync.WaitGroup
	mon := NewMonitor("Progress")
	for j := 0; j < 2; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				mon.Increment(1)
				if i%3 == 0 {
					mon.IncrementCounter("error", 1)
				}
				time.Sleep(250 * time.Millisecond)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	mon.Close()
	fmt.Println()

}

func BenchmarkIteration(b *testing.B) {

	fmt.Println("Benchmarking...")
	mon := NewMonitor("Progress")
	for i := 0; i < b.N; i++ {
		mon.Increment(1)
	}
	mon.Close()
	fmt.Println()
}
