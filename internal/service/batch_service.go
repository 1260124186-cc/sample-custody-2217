package service

import (
	"example.com/sample-custody/internal/audit"
	"example.com/sample-custody/internal/clock"
	"example.com/sample-custody/internal/ids"
	"example.com/sample-custody/internal/metrics"
	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/policy"
	"example.com/sample-custody/internal/review"
	"example.com/sample-custody/internal/store"
	"example.com/sample-custody/internal/validation"
)

type BatchService struct {
	store    *store.Store
	log      *audit.Log
	clock    clock.Clock
	ids      *ids.Generator
	counters *metrics.Counters
}

func (s *BatchService) Create(input model.CreateBatchInput) (model.Batch, error) {
	normalized, err := validation.CreateBatch(input)
	if err != nil {
		return model.Batch{}, err
	}
	if normalized.ID == "" {
		normalized.ID = s.ids.Next("batch")
	}
	samples, err := s.store.SamplesByIDs(normalized.SampleIDs)
	if err != nil {
		return model.Batch{}, err
	}
	checklist := review.ForBatch(samples)
	if !checklist.Allowed {
		return model.Batch{}, model.NewError(model.ErrorConflict, "batch review rejected: %s", review.Explain(checklist))
	}
	for _, sample := range samples {
		if !policy.CanEnterBatch(sample) {
			return model.Batch{}, model.NewError(model.ErrorConflict, "sample %q cannot be added to a batch", sample.ID)
		}
	}
	now := s.clock.Now()
	batch := model.Batch{
		ID:        normalized.ID,
		Name:      normalized.Name,
		Purpose:   normalized.Purpose,
		SampleIDs: normalized.SampleIDs,
		Status:    model.BatchOpen,
		CreatedAt: now,
	}
	if err := s.store.CreateBatch(batch); err != nil {
		return model.Batch{}, err
	}
	for _, sample := range samples {
		nextStatus := policy.TransitionForBatch(sample)
		if nextStatus == sample.Status {
			continue
		}
		sample.Status = nextStatus
		sample.UpdatedAt = now
		if err := s.store.UpdateSample(sample); err != nil {
			return model.Batch{}, err
		}
		s.log.Append(model.AuditEvent{
			ID:        s.ids.Next("event"),
			Entity:    "sample",
			EntityID:  sample.ID,
			Action:    "entered_review",
			Actor:     "review",
			Detail:    "specimen was assigned to inspection batch",
			CreatedAt: now,
		})
	}
	s.log.Append(model.AuditEvent{
		ID:        s.ids.Next("event"),
		Entity:    "batch",
		EntityID:  batch.ID,
		Action:    "opened",
		Actor:     "review",
		Detail:    "inspection batch was created",
		CreatedAt: now,
	})
	s.counters.Increment("batches_opened")
	return batch, nil
}

func (s *BatchService) Get(id string) (model.Batch, error) {
	return s.store.GetBatch(id)
}

func (s *BatchService) Complete(id string) (model.BatchCompletion, error) {
	batch, err := s.store.GetBatch(id)
	if err != nil {
		return model.BatchCompletion{}, err
	}
	if !batch.IsOpen() {
		return model.BatchCompletion{}, model.NewError(model.ErrorConflict, "batch %q is already completed", batch.ID)
	}
	samples, err := s.store.SamplesByIDs(batch.SampleIDs)
	if err != nil {
		return model.BatchCompletion{}, err
	}
	checklist := review.ForCompletion(samples)
	if !checklist.Allowed {
		return model.BatchCompletion{}, model.NewError(model.ErrorConflict, "batch completion rejected: %s", review.Explain(checklist))
	}
	now := s.clock.Now()
	for index := range samples {
		if !policy.CanSeal(samples[index]) {
			return model.BatchCompletion{}, model.NewError(model.ErrorConflict, "sample %q cannot be sealed", samples[index].ID)
		}
		samples[index].Status = model.SampleSealed
		samples[index].UpdatedAt = now
	}
	batch.Status = model.BatchCompleted
	batch.CompletedAt = &now
	if err := s.store.CompleteBatch(batch, samples); err != nil {
		return model.BatchCompletion{}, err
	}
	for _, sample := range samples {
		s.log.Append(model.AuditEvent{
			ID:        s.ids.Next("event"),
			Entity:    "sample",
			EntityID:  sample.ID,
			Action:    "sealed",
			Actor:     "review",
			Detail:    "inspection batch " + batch.ID + " completed",
			CreatedAt: now,
		})
	}
	s.log.Append(model.AuditEvent{
		ID:        s.ids.Next("event"),
		Entity:    "batch",
		EntityID:  batch.ID,
		Action:    "completed",
		Actor:     "review",
		Detail:    "inspection batch completed and specimens sealed",
		CreatedAt: now,
	})
	s.counters.Increment("batches_completed")
	return model.BatchCompletion{Batch: batch, SealedCount: len(samples)}, nil
}
