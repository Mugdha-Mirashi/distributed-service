package node

import (
	"sync"
)

// Counter represents the in-memory counter with deduplication support
type Counter struct {
	Value    int
	Mutex    sync.Mutex
	Applied  map[string]bool // Track applied increment IDs
}

// NewCounter initializes a new counter
func NewCounter() *Counter {
	return &Counter{
		Value:   0,
		Applied: make(map[string]bool),
	}
}

// Increment safely increments if ID is new (returns true if applied)
func (c *Counter) Increment(id string) bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.Applied[id] {
		// Duplicate increment
		return false
	}

	c.Value++
	c.Applied[id] = true
	return true
}

// Get returns the current counter value
func (c *Counter) Get() int {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Value
}


