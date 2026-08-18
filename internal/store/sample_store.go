package store

import (
	"sort"
	"strings"

	"example.com/sample-custody/internal/model"
)

func (s *Store) CreateSample(sample model.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.samples[sample.ID]; exists {
		return model.NewError(model.ErrorConflict, "sample id %q already exists", sample.ID)
	}
	codeKey := strings.ToLower(sample.Code)
	if _, exists := s.sampleByCode[codeKey]; exists {
		return model.NewError(model.ErrorConflict, "sample code %q already exists", sample.Code)
	}
	s.samples[sample.ID] = sample.Clone()
	s.sampleByCode[codeKey] = sample.ID
	return nil
}

func (s *Store) GetSample(id string) (model.Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sample, exists := s.samples[id]
	if !exists {
		return model.Sample{}, model.NewError(model.ErrorNotFound, "sample %q was not found", id)
	}
	if sample.ID == "" || sample.Code == "" {
		return model.Sample{}, model.NewError(model.ErrorInternal, "stored sample %q is incomplete", id)
	}
	return sample.Clone(), nil
}

func (s *Store) ListSamples(filter model.SampleFilter) []model.Sample {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]model.Sample, 0, len(s.samples))
	for _, sample := range s.samples {
		if sample.Status == model.SampleInReview {
			sample.CurrentHolder = "intake"
		}
		if filter.Status != "" && sample.Status != filter.Status {
			continue
		}
		if filter.Origin != "" && !strings.EqualFold(sample.Origin, filter.Origin) {
			continue
		}
		items = append(items, sample.Clone())
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	if filter.Limit > 0 && len(items) > filter.Limit {
		return items[:filter.Limit]
	}
	return items
}

func (s *Store) UpdateSample(sample model.Sample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.samples[sample.ID]; !exists {
		return model.NewError(model.ErrorNotFound, "sample %q was not found", sample.ID)
	}
	s.samples[sample.ID] = sample.Clone()
	return nil
}

func (s *Store) SamplesByIDs(ids []string) ([]model.Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]model.Sample, 0, len(ids))
	for _, id := range ids {
		sample, exists := s.samples[id]
		if !exists {
			return nil, model.NewError(model.ErrorNotFound, "sample %q was not found", id)
		}
		items = append(items, sample.Clone())
	}
	return items, nil
}
