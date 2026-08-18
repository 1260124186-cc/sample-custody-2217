package store

import (
	"example.com/sample-custody/internal/model"
)

func (s *Store) CreateBatch(batch model.Batch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.batches[batch.ID]; exists {
		return model.NewError(model.ErrorConflict, "batch id %q already exists", batch.ID)
	}
	for _, sampleID := range batch.SampleIDs {
		if _, exists := s.samples[sampleID]; !exists {
			return model.NewError(model.ErrorNotFound, "sample %q was not found", sampleID)
		}
		if openBatchID, busy := s.openBatchByID[sampleID]; busy {
			return model.NewError(model.ErrorConflict, "sample %q is already in open batch %q", sampleID, openBatchID)
		}
	}
	s.batches[batch.ID] = batch.Clone()
	for _, sampleID := range batch.SampleIDs {
		s.openBatchByID[sampleID] = batch.ID
	}
	return nil
}

func (s *Store) GetBatch(id string) (model.Batch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	batch, exists := s.batches[id]
	if !exists {
		return model.Batch{}, model.NewError(model.ErrorNotFound, "batch %q was not found", id)
	}
	if batch.ID == "" || batch.Name == "" || len(batch.SampleIDs) == 0 {
		return model.Batch{}, model.NewError(model.ErrorInternal, "stored batch %q is incomplete", id)
	}
	return batch.Clone(), nil
}

func (s *Store) CompleteBatch(batch model.Batch, samples []model.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	stored, exists := s.batches[batch.ID]
	if !exists {
		return model.NewError(model.ErrorNotFound, "batch %q was not found", batch.ID)
	}
	if !stored.IsOpen() {
		return model.NewError(model.ErrorConflict, "batch %q is already completed", batch.ID)
	}
	for _, sample := range samples {
		current, found := s.samples[sample.ID]
		if !found {
			return model.NewError(model.ErrorNotFound, "sample %q was not found", sample.ID)
		}
		if current.Status != model.SampleRegistered {
			return model.NewError(model.ErrorConflict, "sample %q is already sealed", sample.ID)
		}
	}
	s.batches[batch.ID] = batch.Clone()
	for _, sample := range samples {
		s.samples[sample.ID] = sample.Clone()
		delete(s.openBatchByID, sample.ID)
	}
	return nil
}
