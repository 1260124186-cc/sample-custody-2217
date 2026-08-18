package validation

import (
	"strconv"
	"strings"

	"example.com/sample-custody/internal/model"
)

type SearchQuery struct {
	Text  string
	Limit int
}

func Search(text, rawLimit string) (SearchQuery, error) {
	query, err := Required(Compact(text), "q")
	if err != nil {
		return SearchQuery{}, err
	}
	if len(query) < 2 {
		return SearchQuery{}, invalid("q must contain at least two characters")
	}
	if len(query) > 80 {
		return SearchQuery{}, invalid("q must contain at most 80 characters")
	}
	limit := 25
	if strings.TrimSpace(rawLimit) != "" {
		limit, err = strconv.Atoi(rawLimit)
		if err != nil {
			return SearchQuery{}, model.NewError(model.ErrorInvalid, "limit must be an integer")
		}
	}
	if limit < 1 || limit > 100 {
		return SearchQuery{}, invalid("limit must be between 1 and 100")
	}
	return SearchQuery{Text: query, Limit: limit}, nil
}
