package export

import (
	"fmt"
	"sort"
	"strings"

	"example.com/sample-custody/internal/model"
)

// Manifest returns a deterministic plain-text review record for one inspection batch.
func Manifest(batch model.Batch, samples []model.Sample) string {
	items := append([]model.Sample(nil), samples...)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})
	lines := []string{
		"inspection batch: " + batch.Name,
		"batch id: " + batch.ID,
		"purpose: " + batch.Purpose,
		"status: " + string(batch.Status),
		"specimens:",
	}
	for _, sample := range items {
		lines = append(lines, fmt.Sprintf(
			"- %s | %s | %.3f %s | holder=%s | status=%s (%s; %s)",
			sample.Code,
			sample.Material,
			sample.Quantity,
			sample.Unit,
			sample.CurrentHolder,
			sample.Status,
			sample.Status.Label(),
			sample.Status.TransitionHint(),
		))
	}
	return strings.Join(lines, "\n") + "\n"
}
