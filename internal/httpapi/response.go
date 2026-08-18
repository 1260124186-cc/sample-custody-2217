package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"example.com/sample-custody/internal/model"
)

func decodeJSON(request *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return model.NewError(model.ErrorInvalid, "request body must be valid JSON: %v", err)
	}
	if decoder.More() {
		return model.NewError(model.ErrorInvalid, "request body must contain one JSON value")
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeError(writer http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	kind := model.ErrorInternal
	var domain *model.DomainError
	if errors.As(err, &domain) {
		kind = domain.Kind
		switch kind {
		case model.ErrorInvalid:
			status = http.StatusBadRequest
		case model.ErrorNotFound:
			status = http.StatusNotFound
		case model.ErrorConflict:
			status = http.StatusConflict
		}
	}
	writeJSON(writer, status, map[string]string{
		"error":   string(kind),
		"message": err.Error(),
	})
}
