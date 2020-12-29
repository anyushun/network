package counter

import (
	"sync"
)

// Counter struct
type Counter struct {
	lock    sync.Mutex
	counter float64
}

// Add value
func (c *Counter) Add(value float64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.counter += value
}

// Sub value
func (c *Counter) Sub(value float64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.counter -= value
}

// Set counter
func (c *Counter) Set(value float64) {
	c.counter = value
}

// Get counter
func (c *Counter) Get() float64 {
	return c.counter
}
