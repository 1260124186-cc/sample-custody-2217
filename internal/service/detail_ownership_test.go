package service_test

import (
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestDetailDoesNotExposeMutableHistory(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "own-a",
		Code:     "SC-OWN-001",
		Material: "serum",
		Origin:   "north-lab",
		Quantity: 1,
		Unit:     "ml",
	})
	if err != nil {
		t.Fatalf("register sample: %v", err)
	}
	if _, _, err := services.Transfers.Transfer(sample.ID, model.TransferInput{
		From:     "intake",
		To:       "review",
		Location: "cold-room",
		Operator: "op-1",
		Note:     "handoff",
	}); err != nil {
		t.Fatalf("transfer sample: %v", err)
	}
	detail, err := services.Queries.Detail(sample.ID)
	if err != nil {
		t.Fatalf("get detail: %v", err)
	}
	if len(detail.Transfers) != 1 || len(detail.Events) < 1 {
		t.Fatalf("detail history is incomplete: %+v", detail)
	}
	detail.Sample.Code = "MUTATED"
	detail.Transfers[0].To = "MUTATED"
	detail.Events[0].Detail = "MUTATED"

	again, err := services.Queries.Detail(sample.ID)
	if err != nil {
		t.Fatalf("get detail again: %v", err)
	}
	if again.Sample.Code == "MUTATED" || again.Transfers[0].To == "MUTATED" || again.Events[0].Detail == "MUTATED" {
		t.Fatalf("detail exposed mutable state: %+v", again)
	}
}
