package search

import (
	"sort"
	"strings"

	"example.com/sample-custody/internal/model"
)

func Samples(samples []model.Sample, query string, limit int) []model.SearchResult {
	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return []model.SearchResult{}
	}
	results := make([]model.SearchResult, 0)
	for _, sample := range samples {
		fields := matchingFields(sample, needle)
		if len(fields) == 0 {
			continue
		}
		results = append(results, model.SearchResult{
			ID:            sample.ID,
			Code:          sample.Code,
			Material:      sample.Material,
			Origin:        sample.Origin,
			Status:        sample.Status,
			CurrentHolder: sample.CurrentHolder,
			MatchFields:   fields,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Code == results[j].Code {
			return results[i].ID < results[j].ID
		}
		return results[i].Code < results[j].Code
	})
	if limit > 0 && len(results) > limit {
		return results[:limit]
	}
	return results
}

func matchingFields(sample model.Sample, needle string) []string {
	fields := make([]string, 0, 4)
	if strings.Contains(sample.Code, needle) {
		fields = append(fields, "code")
	}
	if strings.Contains(strings.ToLower(sample.Material), needle) {
		fields = append(fields, "material")
	}
	if strings.Contains(strings.ToLower(sample.Origin), needle) {
		fields = append(fields, "origin")
	}
	if strings.Contains(strings.ToLower(sample.CurrentHolder), needle) {
		fields = append(fields, "current_holder")
	}
	return fields
}
