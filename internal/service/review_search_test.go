package service_test

import (
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestSearchFindsReviewSampleByCurrentHolder(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID: "review-search", Code: "SEARCH-005", Material: "serum",
		Origin: "city-lab", Quantity: 1, Unit: "ml",
	})
	if err != nil { t.Fatalf("register sample: %v", err) }
	_, err = services.Batches.Create(model.CreateBatchInput{
		ID: "review-search-batch", Name: "search review", Purpose: "inspection",
		SampleIDs: []string{sample.ID},
	})
	if err != nil { t.Fatalf("create batch: %v", err) }
	if _, _, err := services.Transfers.Transfer(sample.ID, model.TransferInput{From: "intake", To: "review-team", Location: "cold", Operator: "operator-5"}); err != nil { t.Fatalf("transfer: %v", err) }
	results := services.Queries.Search("review-team", 10)
	if len(results) != 1 || results[0].CurrentHolder != "review-team" { t.Fatalf("unexpected search results: %+v", results) }
}
