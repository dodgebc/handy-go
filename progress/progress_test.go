package progress

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestIteration(t *testing.T) {

	fmt.Println("Expecting:\t10 it\t4 it/s")
	mon := NewMonitor("Progress")
	for i := 0; i < 10; i++ {
		mon.Increment(1)
		time.Sleep(250 * time.Millisecond)
	}
	mon.Close()
	fmt.Println()
}

func TestCounter(t *testing.T) {

	fmt.Println("Expecting:\t10 it\t4 it/s\terror: 4")
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

	fmt.Println("Expecting:\t10 it\t4 it/s\ttoo big: 1\ttoo small: 2")
	mon := NewMonitor("Progress")
	mon.RegisterCounters([]string{"too big", "too small"})
	for i := 0; i < 10; i++ {
		mon.Increment(1)
		if i < 2 {
			mon.IncrementCounter("too small", 1)
		}
		if i > 8 {
			mon.IncrementCounter("too big", 1)
		}
		time.Sleep(250 * time.Millisecond)
	}
	mon.Close()
	fmt.Println()
}

func TestConcurrentCounter(t *testing.T) {

	fmt.Println("Expecting:\t20 it\t8 it/s\terror: 8")
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
