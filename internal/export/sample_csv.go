package export

import (
	"encoding/csv"
	"strconv"
	"strings"

	"example.com/sample-custody/internal/model"
)

// SampleCSV emits a deterministic portable catalog extract for review clients.
func SampleCSV(samples []model.Sample) string {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	_ = writer.Write([]string{
		"id",
		"code",
		"material",
		"origin",
		"quantity",
		"unit",
		"status",
		"current_holder",
		"created_at",
		"updated_at",
	})
	for _, sample := range samples {
		_ = writer.Write([]string{
			sample.ID,
			sample.Code,
			sample.Material,
			sample.Origin,
			strconv.FormatFloat(sample.Quantity, 'f', -1, 64),
			sample.Unit,
			string(sample.Status),
			sample.CurrentHolder,
			sample.CreatedAt.Format("2006-01-02T15:04:05Z"),
			sample.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	writer.Flush()
	return builder.String()
}
