package store

import "example.com/sample-custody/internal/model"

func (s *Store) TransfersForSample(sampleID string) []model.Transfer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := s.transfers[sampleID]
	result := make([]model.Transfer, len(entries))
	for index, entry := range entries {
		result[index] = entry.Clone()
	}
	return result
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

func (s *Store) UpdateSampleAndAppendTransferGuarded(sample model.Sample, transfer model.Transfer, guard model.TransferGuard) error {
	if !guard.Valid() {
		return model.NewError(model.ErrorInvalid, "transfer guard is incomplete")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.samples[guard.SampleID]
	if !exists {
		return model.NewError(model.ErrorNotFound, "sample %q was not found", guard.SampleID)
	}
	if !guard.Matches(current) || current.ID != sample.ID {
		return model.NewError(model.ErrorConflict, "sample %q custody changed before transfer", guard.SampleID)
	}
	s.samples[sample.ID] = sample.Clone()
	s.transfers[sample.ID] = append(s.transfers[sample.ID], transfer.Clone())
	return nil
}
