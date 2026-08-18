package service_test

import (
	"strings"
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestOpenBatchManifestDoesNotPanic(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID:       "manifest-a",
		Code:     "SC-MAN-001",
		Material: "serum",
		Origin:   "north-lab",
		Quantity: 1,
		Unit:     "ml",
	})
	if err != nil {
		t.Fatalf("register sample: %v", err)
	}
	batch, err := services.Batches.Create(model.CreateBatchInput{
		ID:        "batch-manifest",
		Name:      "manifest batch",
		Purpose:   "manifest test",
		SampleIDs: []string{sample.ID},
	})
	if err != nil {
		t.Fatalf("create batch: %v", err)
	}
	manifest, err := services.Queries.Manifest(batch.ID)
	if err != nil {
		t.Fatalf("get manifest: %v", err)
	}
	if !strings.Contains(manifest, "status: open") {
		t.Fatalf("open batch manifest is missing open status: %s", manifest)
	}
}
