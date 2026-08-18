package service_test

import (
	"sync"
	"testing"
	"time"

	"example.com/sample-custody/internal/audit"
	"example.com/sample-custody/internal/clock"
	"example.com/sample-custody/internal/ids"
	"example.com/sample-custody/internal/metrics"
	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/service"
	"example.com/sample-custody/internal/store"
)

type handoffClock struct {
	mu      sync.Mutex
	calls   int
	entered chan struct{}
	release chan struct{}
	now     time.Time
}

func newHandoffClock() *handoffClock {
	return &handoffClock{
		entered: make(chan struct{}, 2),
		release: make(chan struct{}),
		now:     time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
	}
}

func (c *handoffClock) Now() time.Time {
	c.mu.Lock()
	c.calls++
	call := c.calls
	c.mu.Unlock()
	if call == 2 || call == 3 {
		c.entered <- struct{}{}
		<-c.release
	}
	return c.now
}

var _ clock.Clock = (*handoffClock)(nil)

func TestConcurrentHandoffAcceptsOnlyOneTransfer(t *testing.T) {
	clock := newHandoffClock()
	services := service.New(store.New(), audit.NewLog(), clock, ids.New(), metrics.NewCounters())
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "concurrent-sample",
		Code:     "CONCURRENT-001",
		Material: "serum",
		Origin:   "race-lab",
		Quantity: 2,
		Unit:     "ml",
	})
	if err != nil {
		t.Fatalf("register sample: %v", err)
	}

	results := make(chan error, 2)
	for _, nextHolder := range []string{"review-a", "review-b"} {
		nextHolder := nextHolder
		go func() {
			_, _, transferErr := services.Transfers.Transfer(sample.ID, model.TransferInput{
				From:     "intake",
				To:       nextHolder,
				Location: "cold-room",
				Operator: "concurrent-test",
			})
			results <- transferErr
		}()
	}
	<-clock.entered
	<-clock.entered
	close(clock.release)

	successes := 0
	for range 2 {
		if <-results == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("expected exactly one accepted handoff, got %d", successes)
	}
}
