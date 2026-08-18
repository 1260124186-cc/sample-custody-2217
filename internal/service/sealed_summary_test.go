package service_test

import (
	"testing"

	"example.com/sample-custody/internal/model"
)

func TestSummaryCountsSealedSampleForItsOrigin(t *testing.T) {
	services := newServices()
	sample, err := services.Samples.Register(model.RegisterSampleInput{
		ID: "sealed-summary", Code: "SEALED-004", Material: "plasma",
		Origin: "archive-lab", Quantity: 1, Unit: "ml",
	})
	if err != nil { t.Fatalf("register sample: %v", err) }
	if _, err := services.Samples.Seal(sample.ID, "review", "finalized"); err != nil { t.Fatalf("seal sample: %v", err) }
	summary := services.Queries.Summary()
	if summary.SealedCount != 1 || len(summary.Origins) != 1 || summary.Origins[0].SealedCount != 1 || summary.Origins[0].OpenCount != 0 {
		t.Fatalf("unexpected sealed summary: %+v", summary)
	}
}
