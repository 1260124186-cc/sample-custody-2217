package service_test

import (
	"errors"
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestDomainErrorsStayClassified(t *testing.T) {
	services := newServices()
	base := model.RegisterSampleInput{
		ID:       "domain-a",
		Code:     "SC-DOM-001",
		Material: "serum",
		Origin:   "north-lab",
		Quantity: 1,
		Unit:     "ml",
	}
	if _, err := services.Samples.Register(base); err != nil {
		t.Fatalf("register first sample: %v", err)
	}
	duplicate := base
	duplicate.ID = "domain-b"
	_, err := services.Samples.Register(duplicate)
	var domainError *model.DomainError
	if !errors.As(err, &domainError) || domainError.Kind != model.ErrorConflict {
		t.Fatalf("expected conflict domain error, got %v", err)
	}
}
