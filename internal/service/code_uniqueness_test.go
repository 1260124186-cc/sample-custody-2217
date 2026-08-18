package service_test

import (
	"errors"
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestSampleCodeUniquenessIsCaseInsensitive(t *testing.T) {
	services := newServices()
	if _, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "case-a",
		Code:     "SC-CASE-001",
		Material: "serum",
		Origin:   "north-lab",
		Quantity: 1,
		Unit:     "ml",
	}); err != nil {
		t.Fatalf("register first sample: %v", err)
	}
	_, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "case-b",
		Code:     "sc-case-001",
		Material: "plasma",
		Origin:   "east-lab",
		Quantity: 2,
		Unit:     "ml",
	})
	var domainError *model.DomainError
	if !errors.As(err, &domainError) || domainError.Kind != model.ErrorConflict {
		t.Fatalf("expected conflict for case-insensitive duplicate code, got %v", err)
	}
}
