package service_test

import (
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestBatchOwnsItsSampleIDSlice(t *testing.T) {
	services := newServices()
	sampleA, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "slice-a",
		Code:     "SC-SLICE-A",
		Material: "serum",
		Origin:   "north-lab",
		Quantity: 1,
		Unit:     "ml",
	})
	if err != nil {
		t.Fatalf("register sample a: %v", err)
	}
	sampleB, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "slice-b",
		Code:     "SC-SLICE-B",
		Material: "plasma",
		Origin:   "east-lab",
		Quantity: 2,
		Unit:     "ml",
	})
	if err != nil {
		t.Fatalf("register sample b: %v", err)
	}
	batch, err := services.Batches.Create(model.CreateBatchInput{
		ID:        "batch-slice",
		Name:      "slice batch",
		Purpose:   "ownership test",
		SampleIDs: []string{sampleB.ID, sampleA.ID},
	})
	if err != nil {
		t.Fatalf("create batch: %v", err)
	}
	if len(batch.SampleIDs) != 2 {
		t.Fatalf("unexpected sample ids: %v", batch.SampleIDs)
	}
	batch.SampleIDs[0] = "mutated-id"
	fetched, err := services.Batches.Get(batch.ID)
	if err != nil {
		t.Fatalf("get batch: %v", err)
	}
	if fetched.SampleIDs[0] == "mutated-id" {
		t.Fatalf("mutating returned batch changed stored sample ids: %v", fetched.SampleIDs)
	}
}
