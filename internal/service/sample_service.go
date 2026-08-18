package service

import (
	"fmt"

	"example.com/sample-custody/internal/audit"
	"example.com/sample-custody/internal/clock"
	"example.com/sample-custody/internal/ids"
	"example.com/sample-custody/internal/metrics"
	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/review"
	"example.com/sample-custody/internal/store"
	"example.com/sample-custody/internal/validation"
)

type SampleService struct {
	store    *store.Store
	log      *audit.Log
	clock    clock.Clock
	ids      *ids.Generator
	counters *metrics.Counters
}

func (s *SampleService) Register(input model.RegisterSampleInput) (model.Sample, error) {
	checklist := review.ForRegistration(input)
	if !checklist.Allowed {
		return model.Sample{}, model.NewError(model.ErrorInvalid, "registration review rejected: %s", review.Explain(checklist))
	}
	normalized, err := validation.RegisterSample(input)
	if err != nil {
		return model.Sample{}, err
	}
	if normalized.ID == "" {
		normalized.ID = s.ids.Next("sample")
	}
	now := s.clock.Now()
	sample := model.Sample{
		ID:            normalized.ID,
		Code:          normalized.Code,
		Material:      normalized.Material,
		Origin:        normalized.Origin,
		Quantity:      normalized.Quantity,
		Unit:          normalized.Unit,
		Status:        model.SampleRegistered,
		CurrentHolder: "intake",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateSample(sample); err != nil {
		// 使用 %w 包装以保留 store 层 DomainError 的 Kind，让 HTTP 层能映射为正确的状态码
		return model.Sample{}, fmt.Errorf("create sample failed: %w", err)
	}
	s.log.Append(model.AuditEvent{
		ID:        s.ids.Next("event"),
		Entity:    "sample",
		EntityID:  sample.ID,
		Action:    "registered",
		Actor:     "intake",
		Detail:    "specimen was accepted into custody",
		CreatedAt: now,
	})
	s.counters.Increment("samples_registered")
	return sample, nil
}

func (s *SampleService) Get(id string) (model.Sample, error) {
	return s.store.GetSample(id)
}

func (s *SampleService) List(filter model.SampleFilter) []model.Sample {
	return s.store.ListSamples(filter)
}

func (s *SampleService) Seal(id, actor, detail string) (model.Sample, error) {
	sample, err := s.store.GetSample(id)
	if err != nil {
		return model.Sample{}, err
	}
	if sample.Status == model.SampleSealed {
		return model.Sample{}, model.NewError(model.ErrorConflict, "sample %q is already sealed", id)
	}
	now := s.clock.Now()
	sample.Status = model.SampleSealed
	sample.UpdatedAt = now
	if err := s.store.UpdateSample(sample); err != nil {
		return model.Sample{}, err
	}
	s.log.Append(model.AuditEvent{
		ID:        s.ids.Next("event"),
		Entity:    "sample",
		EntityID:  sample.ID,
		Action:    "sealed",
		Actor:     actor,
		Detail:    detail,
		CreatedAt: now,
	})
	s.counters.Increment("samples_sealed")
	return sample, nil
}
