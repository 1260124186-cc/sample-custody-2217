package service

import (
	"fmt"

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

type TransferService struct {
	store    *store.Store
	log      *audit.Log
	clock    clock.Clock
	ids      *ids.Generator
	counters *metrics.Counters
}

func (s *TransferService) Transfer(sampleID string, input model.TransferInput) (model.Transfer, model.Sample, error) {
	normalized, err := validation.Transfer(input)
	if err != nil {
		return model.Transfer{}, model.Sample{}, err
	}
	sample, err := s.store.GetSample(sampleID)
	if err != nil {
		return model.Transfer{}, model.Sample{}, err
	}
	checklist := review.ForTransfer(sample, normalized.From, normalized.To)
	if !checklist.Allowed {
		return model.Transfer{}, model.Sample{}, model.NewError(model.ErrorConflict, "transfer review rejected: %s", review.Explain(checklist))
	}
	if allowed, reason := policy.CanTransfer(sample, normalized.From, normalized.To); !allowed {
		return model.Transfer{}, model.Sample{}, model.NewError(
			model.ErrorConflict,
			"sample %q cannot be transferred: %s",
			sample.ID,
			reason,
		)
	}
	now := s.clock.Now()
	transfer := model.Transfer{
		ID:        s.ids.Next("transfer"),
		SampleID:  sample.ID,
		From:      normalized.From,
		To:        normalized.To,
		Location:  normalized.Location,
		Operator:  normalized.Operator,
		Note:      normalized.Note,
		CreatedAt: now,
	}
	sample.CurrentHolder = normalized.To
	sample.UpdatedAt = now
	if err := s.store.UpdateSampleAndAppendTransfer(sample, transfer); err != nil {
		// 使用 %w 包装以保留 store 层 DomainError 的 Kind
		return model.Transfer{}, model.Sample{}, fmt.Errorf("save transfer failed: %w", err)
	}
	s.log.Append(model.AuditEvent{
		ID:        s.ids.Next("event"),
		Entity:    "sample",
		EntityID:  sample.ID,
		Action:    "transferred",
		Actor:     normalized.Operator,
		Detail:    "custody moved to " + normalized.To + " at " + normalized.Location,
		CreatedAt: now,
	})
	s.counters.Increment("custody_transfers")
	return transfer, sample, nil
}
