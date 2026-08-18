package ids

import (
	"fmt"
	"sync"
)

// Generator creates process-local identifiers with a readable entity prefix.
type Generator struct {
	mu      sync.Mutex
	counter uint64
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Next(prefix string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counter++
	return fmt.Sprintf("%s-%06d", prefix, g.counter)
}
