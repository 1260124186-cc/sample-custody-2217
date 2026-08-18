package store

import (
	"sync"

	"example.com/sample-custody/internal/model"
)

// Store owns all mutable workflow state. Methods acquire the same lock so the
// service layer can compose sample and batch operations without partial writes.
type Store struct {
	mu            sync.RWMutex
	samples       map[string]model.Sample
	sampleByCode  map[string]string
	transfers     map[string][]model.Transfer
	batches       map[string]model.Batch
	openBatchByID map[string]string
}

func New() *Store {
	return &Store{
		samples:       make(map[string]model.Sample),
		sampleByCode:  make(map[string]string),
		transfers:     make(map[string][]model.Transfer),
		batches:       make(map[string]model.Batch),
		openBatchByID: make(map[string]string),
	}
}

func (s *Store) Counts() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]int{
		"samples":   len(s.samples),
		"batches":   len(s.batches),
		"transfers": s.transferCountLocked(),
	}
}

func (s *Store) transferCountLocked() int {
	count := 0
	for _, entries := range s.transfers {
		count += len(entries)
	}
	return count
}
