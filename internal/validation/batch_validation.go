package validation

import (
	"sort"

	"example.com/sample-custody/internal/model"
)

func CreateBatch(input model.CreateBatchInput) (model.CreateBatchInput, error) {
	var err error
	if input.ID != "" {
		input.ID, err = Identifier(input.ID, "id")
		if err != nil {
			return model.CreateBatchInput{}, err
		}
	}
	input.Name, err = Required(Compact(input.Name), "name")
	if err != nil {
		return model.CreateBatchInput{}, err
	}
	input.Purpose, err = Required(Compact(input.Purpose), "purpose")
	if err != nil {
		return model.CreateBatchInput{}, err
	}
	if len(input.SampleIDs) == 0 {
		return model.CreateBatchInput{}, invalid("sample_ids must contain at least one specimen")
	}
	if len(input.SampleIDs) > 100 {
		return model.CreateBatchInput{}, invalid("sample_ids may contain at most 100 specimens")
	}
	seen := make(map[string]struct{}, len(input.SampleIDs))
	normalized := make([]string, 0, len(input.SampleIDs))
	for _, sampleID := range input.SampleIDs {
		value, validateErr := Identifier(sampleID, "sample_id")
		if validateErr != nil {
			return model.CreateBatchInput{}, validateErr
		}
		if _, exists := seen[value]; exists {
			return model.CreateBatchInput{}, invalid("sample_ids must not repeat values")
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	sort.Strings(normalized)
	input.SampleIDs = normalized
	return input, nil
}
