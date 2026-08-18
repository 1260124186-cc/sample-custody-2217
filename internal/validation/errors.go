package validation

import (
	"fmt"

	"example.com/sample-custody/internal/model"
)

func invalid(format string, args ...any) error {
	return model.NewError(model.ErrorInvalid, "%s", fmt.Sprintf(format, args...))
}
