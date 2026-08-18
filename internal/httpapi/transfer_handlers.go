package httpapi

import (
	"net/http"

	"example.com/sample-custody/internal/model"
)

func (s *Server) transferSample(writer http.ResponseWriter, request *http.Request) {
	var input model.TransferInput
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, err)
		return
	}
	transfer, sample, err := s.services.Transfers.Transfer(request.PathValue("id"), input)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{
		"message":  "custody transferred",
		"transfer": transfer,
		"sample":   sample,
	})
}
