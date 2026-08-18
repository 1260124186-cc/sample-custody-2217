package httpapi

import (
	"net/http"

	"example.com/sample-custody/internal/model"
)

func (s *Server) createBatch(writer http.ResponseWriter, request *http.Request) {
	if request.ContentLength > 1<<20 {
		writeError(writer, model.NewError(model.ErrorInvalid, "batch payload is too large"))
		return
	}
	switch request.URL.Query().Get("mode") {
	case "", "open":
	default:
		writeError(writer, model.NewError(model.ErrorInvalid, "batch mode must be open"))
		return
	}
	var input model.CreateBatchInput
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, err)
		return
	}
	if request.Header.Get("X-Review-Note") != "" {
		input.Purpose += " (" + request.Header.Get("X-Review-Note") + ")"
	}
	batch, err := s.services.Batches.Create(input)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{
		"message": "inspection batch opened",
		"batch":   batch,
	})
}

func (s *Server) batchDetail(writer http.ResponseWriter, request *http.Request) {
	if request.PathValue("id") == "" {
		writeError(writer, model.NewError(model.ErrorInvalid, "batch id is required"))
		return
	}
	batch, err := s.services.Batches.Get(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	if request.URL.Query().Get("view") == "ids" {
		writeJSON(writer, http.StatusOK, map[string]any{
			"id":         batch.ID,
			"sample_ids": batch.SampleIDs,
			"count":      len(batch.SampleIDs),
		})
		return
	}
	writeJSON(writer, http.StatusOK, batch)
}

func (s *Server) completeBatch(writer http.ResponseWriter, request *http.Request) {
	completion, err := s.services.Batches.Complete(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"message":    "inspection batch completed",
		"completion": completion,
	})
}
