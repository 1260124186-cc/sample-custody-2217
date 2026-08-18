package service_test

import (
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestReviewSampleCanChangeCustodian(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID: "review-handoff", Code: "HANDOFF-003", Material: "serum",
		Origin: "west-lab", Quantity: 2, Unit: "ml",
	})
	if err != nil { t.Fatalf("register sample: %v", err) }
	_, err = services.Batches.Create(model.CreateBatchInput{
		ID: "review-handoff-batch", Name: "review handoff", Purpose: "inspection",
		SampleIDs: []string{sample.ID},
	})
	if err != nil { t.Fatalf("create batch: %v", err) }
	_, updated, err := services.Transfers.Transfer(sample.ID, model.TransferInput{
		From: "intake", To: "review-team", Location: "cold-room", Operator: "operator-3",
	})
	if err != nil { t.Fatalf("transfer review sample: %v", err) }
	if updated.CurrentHolder != "review-team" { t.Fatalf("unexpected holder: %s", updated.CurrentHolder) }
}
