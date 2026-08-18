package service

import (
	"example.com/sample-custody/internal/audit"
	"example.com/sample-custody/internal/clock"
	"example.com/sample-custody/internal/ids"
	"example.com/sample-custody/internal/metrics"
	"example.com/sample-custody/internal/store"
)

// Services exposes workflow-focused service objects backed by one shared store.
type Services struct {
	Samples   *SampleService
	Transfers *TransferService
	Batches   *BatchService
	Queries   *QueryService
}

func New(store *store.Store, log *audit.Log, clock clock.Clock, ids *ids.Generator, counters *metrics.Counters) Services {
	return Services{
		Samples:   &SampleService{store: store, log: log, clock: clock, ids: ids, counters: counters},
		Transfers: &TransferService{store: store, log: log, clock: clock, ids: ids, counters: counters},
		Batches:   &BatchService{store: store, log: log, clock: clock, ids: ids, counters: counters},
		Queries:   &QueryService{store: store, log: log},
	}
}
