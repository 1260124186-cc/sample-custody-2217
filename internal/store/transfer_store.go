package store

import "example.com/sample-custody/internal/model"

func (s *Store) TransfersForSample(sampleID string) []model.Transfer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transfers[sampleID]
}

func (s *Store) UpdateSampleAndAppendTransfer(sample model.Sample, transfer model.Transfer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.samples[sample.ID]; !exists {
		return model.NewError(model.ErrorNotFound, "sample %q was not found", sample.ID)
	}
	s.samples[sample.ID] = sample.Clone()
	s.transfers[sample.ID] = append(s.transfers[sample.ID], transfer.Clone())
	return nil
}
