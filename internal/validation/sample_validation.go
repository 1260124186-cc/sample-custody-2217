package validation

import (
	"strings"

	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/policy"
)

func RegisterSample(input model.RegisterSampleInput) (model.RegisterSampleInput, error) {
	var err error
	if input.ID != "" {
		input.ID, err = Identifier(input.ID, "id")
		if err != nil {
			return model.RegisterSampleInput{}, err
		}
	}
	input.Code, err = Identifier(strings.ToUpper(input.Code), "code")
	if err != nil {
		return model.RegisterSampleInput{}, err
	}
	input.Material, err = Required(Compact(input.Material), "material")
	if err != nil {
		return model.RegisterSampleInput{}, err
	}
	input.Origin, err = Required(Compact(input.Origin), "origin")
	if err != nil {
		return model.RegisterSampleInput{}, err
	}
	input.Unit, err = Required(strings.ToLower(Compact(input.Unit)), "unit")
	if err != nil {
		return model.RegisterSampleInput{}, err
	}
	assessment := policy.AssessQuantity(input.Quantity, input.Unit)
	if !assessment.Acceptable {
		return model.RegisterSampleInput{}, invalid("%s", assessment.Reason)
	}
	return input, nil
}

func SampleFilter(status, origin string, limit int) (model.SampleFilter, error) {
	filter := model.SampleFilter{Origin: Compact(origin), Limit: limit}
	if status != "" {
		switch model.SampleStatus(strings.TrimSpace(status)) {
		case model.SampleRegistered, model.SampleInReview, model.SampleSealed:
			filter.Status = model.SampleStatus(strings.TrimSpace(status))
		default:
			return model.SampleFilter{}, invalid("status is not recognized")
		}
	}
	if filter.Limit < 0 || filter.Limit > 200 {
		return model.SampleFilter{}, invalid("limit must be between 0 and 200")
	}
	return filter, nil
}
