package validation

import "example.com/sample-custody/internal/model"

func Transfer(input model.TransferInput) (model.TransferInput, error) {
	var err error
	input.From, err = Required(Compact(input.From), "from")
	if err != nil {
		return model.TransferInput{}, err
	}
	input.To, err = Required(Compact(input.To), "to")
	if err != nil {
		return model.TransferInput{}, err
	}
	if input.From == input.To {
		return model.TransferInput{}, invalid("from and to must be different")
	}
	input.Location, err = Required(Compact(input.Location), "location")
	if err != nil {
		return model.TransferInput{}, err
	}
	input.Operator, err = Required(Compact(input.Operator), "operator")
	if err != nil {
		return model.TransferInput{}, err
	}
	input.Note = Compact(input.Note)
	if len(input.Note) > 500 {
		return model.TransferInput{}, invalid("note must be at most 500 characters")
	}
	return input, nil
}
