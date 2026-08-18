package metrics

import "sync"

// Counters records request-independent workflow totals for the health response.
type Counters struct {
	mu     sync.RWMutex
	values map[string]int64
}

func NewCounters() *Counters {
	return &Counters{values: make(map[string]int64)}
}

func (c *Counters) Increment(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[name]++
}

func (c *Counters) Snapshot() map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	copyValues := make(map[string]int64, len(c.values))
	for name, value := range c.values {
		copyValues[name] = value
	}
	return copyValues
}
