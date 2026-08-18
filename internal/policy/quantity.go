package policy

import (
	"math"

	"example.com/sample-custody/internal/model"
)

const (
	MinQuantity = 0.0001
	MaxQuantity = 1000000
)

type QuantityAssessment struct {
	Acceptable bool
	Reason     string
}

func AssessQuantity(quantity float64, unit string) QuantityAssessment {
	if math.IsNaN(quantity) || math.IsInf(quantity, 0) {
		return QuantityAssessment{Reason: "quantity must be finite"}
	}
	if quantity < MinQuantity {
		return QuantityAssessment{Reason: "quantity is below the measurable threshold"}
	}
	if quantity > MaxQuantity {
		return QuantityAssessment{Reason: "quantity exceeds the acceptance range"}
	}
	if unit == "" {
		return QuantityAssessment{Reason: "unit is required"}
	}
	return QuantityAssessment{Acceptable: true}
}

func CanSeal(sample model.Sample) bool {
	return sample.Status.AllowsCompletion() &&
		(sample.Status == model.SampleRegistered || sample.Status == model.SampleInReview)
}

func CanEnterBatch(sample model.Sample) bool {
	return sample.Status.AllowsBatchAssignment() && sample.IsReviewable() && sample.Quantity >= MinQuantity
}
