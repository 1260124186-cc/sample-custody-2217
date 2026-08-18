package policy

import (
	"strings"

	"example.com/sample-custody/internal/model"
)

func HolderMatches(current, provided string) bool {
	return strings.EqualFold(strings.TrimSpace(current), strings.TrimSpace(provided))
}

func CanTransfer(sample model.Sample, from, to string) (bool, string) {
	if !sample.Status.AllowsTransfer() || !sample.IsTransferable() || sample.Status != model.SampleRegistered {
		return false, "sealed specimens cannot be transferred"
	}
	if !HolderMatches(sample.CurrentHolder, from) {
		return false, "the stated holder does not match the recorded holder"
	}
	if strings.EqualFold(strings.TrimSpace(from), strings.TrimSpace(to)) {
		return false, "the next holder must differ from the current holder"
	}
	return true, ""
}

func TransitionForBatch(sample model.Sample) model.SampleStatus {
	if sample.Status == model.SampleRegistered {
		if sample.Status.CanTransitionTo(model.SampleInReview) {
			return model.SampleInReview
		}
	}
	return sample.Status
}
