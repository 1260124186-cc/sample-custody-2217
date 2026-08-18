package service_test

import (
	"errors"
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

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

var _ clock.Clock = fixedClock{}

func newServices() service.Services {
	repository := store.New()
	log := audit.NewLog()
	counter := metrics.NewCounters()
	return service.New(repository, log, fixedClock{now: time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)}, ids.New(), counter)
}

func TestCustodyWorkflow(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "sample-test-001",
		Code:     "SC-001",
		Material: "serum",
		Origin:   "north-lab",
		Quantity: 2,
		Unit:     "ml",
	})
	if err != nil {
		t.Fatalf("register sample: %v", err)
	}
	if sample.CurrentHolder != "intake" {
		t.Fatalf("unexpected initial holder: %s", sample.CurrentHolder)
	}

	transfer, updated, err := services.Transfers.Transfer(sample.ID, model.TransferInput{
		From:     "intake",
		To:       "review",
		Location: "cold-room",
		Operator: "operator-1",
		Note:     "handoff",
	})
	if err != nil {
		t.Fatalf("transfer sample: %v", err)
	}
	if transfer.To != "review" || updated.CurrentHolder != "review" {
		t.Fatalf("transfer did not update holder: %+v", updated)
	}

	batch, err := services.Batches.Create(model.CreateBatchInput{
		ID:        "batch-test-001",
		Name:      "first inspection",
		Purpose:   "release review",
		SampleIDs: []string{sample.ID},
	})
	if err != nil {
		t.Fatalf("create batch: %v", err)
	}
	completion, err := services.Batches.Complete(batch.ID)
	if err != nil {
		t.Fatalf("complete batch: %v", err)
	}
	if completion.SealedCount != 1 || completion.Batch.Status != model.BatchCompleted {
		t.Fatalf("unexpected completion: %+v", completion)
	}
	final, err := services.Samples.Get(sample.ID)
	if err != nil {
		t.Fatalf("get final sample: %v", err)
	}
	if final.Status != model.SampleSealed {
		t.Fatalf("sample was not sealed: %s", final.Status)
	}
}

func TestCustodyRejectsStaleHolder(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "sample-test-002",
		Code:     "SC-002",
		Material: "plasma",
		Origin:   "east-lab",
		Quantity: 1,
		Unit:     "ml",
	})
	if err != nil {
		t.Fatalf("register sample: %v", err)
	}
	_, _, err = services.Transfers.Transfer(sample.ID, model.TransferInput{
		From:     "wrong-holder",
		To:       "review",
		Location: "cold-room",
		Operator: "operator-2",
	})
	var domainError *model.DomainError
	if !errors.As(err, &domainError) || domainError.Kind != model.ErrorConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
}
