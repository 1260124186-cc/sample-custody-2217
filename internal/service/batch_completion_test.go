package service_test

import (
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestCompletingReviewBatchSealsEverySample(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID: "completion-sample", Code: "COMPLETE-001", Material: "plasma",
		Origin: "central-lab", Quantity: 1, Unit: "ml",
	})
	if err != nil {
		t.Fatalf("register sample: %v", err)
	}
	batch, err := services.Batches.Create(model.CreateBatchInput{
		ID: "completion-batch", Name: "release inspection", Purpose: "release",
		SampleIDs: []string{sample.ID},
	})
	if err != nil {
		t.Fatalf("create batch: %v", err)
	}
	completion, err := services.Batches.Complete(batch.ID)
	if err != nil {
		t.Fatalf("complete review batch: %v", err)
	}
	if completion.SealedCount != 1 || completion.Batch.Status != model.BatchCompleted {
		t.Fatalf("unexpected completion: %+v", completion)
	}
}
