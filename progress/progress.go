// Package progress provides simple concurrent progress monitoring
package progress

import (
	"fmt"
	"sync"
	"time"
)

// Monitor tracks iterations and other statistics in real time
type Monitor struct {
	description     string
	startTime       time.Time
	lastTime        time.Time
	iteration       int
	lastIteration   int
	countKeys       []string
	countIterations []int
	mux             sync.Mutex
}

// NewMonitor initializes a Monitor at the current time
func NewMonitor(description string) *Monitor {
	mon := Monitor{
		description: description,
		startTime:   time.Now(),
	}
	return &mon
}

// Increment adds to the main iteration counter
func (mon *Monitor) Increment(n int) {
	mon.mux.Lock()
	defer mon.mux.Unlock()

	mon.iteration += n

	if time.Since(mon.lastTime).Seconds() > 0.5 {
		fmt.Printf("\r%s:\t%d it", mon.description, mon.iteration)
		fmt.Printf("\t%.f it/s", float64(mon.iteration-mon.lastIteration)/time.Since(mon.lastTime).Seconds())
		for i := range mon.countKeys {
			fmt.Printf("\t%s: %d", mon.countKeys[i], mon.countIterations[i])
		}
		fmt.Printf("\t")
		mon.lastTime = time.Now()
		mon.lastIteration = mon.iteration
	}
}

// IncrementCounter adds to (or creates) a secondary counter
func (mon *Monitor) IncrementCounter(name string, n int) {
	mon.mux.Lock()
	defer mon.mux.Unlock()

	for i := range mon.countKeys {
		if name == mon.countKeys[i] {
			mon.countIterations[i] += n
			return
		}
	}
	mon.countKeys = append(mon.countKeys, name)
	mon.countIterations = append(mon.countIterations, 1)
}

// RegisterCounters creates multiple secondary counters (not required unless a specific order is desired)
func (mon *Monitor) RegisterCounters(names []string) {
	mon.mux.Lock()
	defer mon.mux.Unlock()

	for _, name := range names {
		mon.countKeys = append(mon.countKeys, name)
		mon.countIterations = append(mon.countIterations, 0)
	}
}

// Iteration returns the current iteration
func (mon *Monitor) Iteration() int {
	mon.mux.Lock()
	defer mon.mux.Unlock()
	return mon.iteration
}

// Count returns the current value of a counter
func (mon *Monitor) Count(name string) int {
	mon.mux.Lock()
	defer mon.mux.Unlock()
	for i := range mon.countKeys {
		if name == mon.countKeys[i] {
			return mon.countIterations[i]
		}
	}
	return 0
}

// Close prints final statistics
func (mon *Monitor) Close() {
	mon.mux.Lock()
	defer mon.mux.Unlock()

	fmt.Printf("\r%s:\t%d it", mon.description, mon.iteration)
	fmt.Printf("\t%.f it/s", float64(mon.iteration)/time.Since(mon.startTime).Seconds())
	for i := range mon.countKeys {
		fmt.Printf("\t%s: %d", mon.countKeys[i], mon.countIterations[i])
	}
	fmt.Printf("\t\n")
}
