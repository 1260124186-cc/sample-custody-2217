package httpapi

import (
	"net/http"
	"strconv"

	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/validation"
)

func (s *Server) summary(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, s.services.Queries.Summary())
}

func (s *Server) searchSamples(writer http.ResponseWriter, request *http.Request) {
	query, err := validation.Search(request.URL.Query().Get("q"), request.URL.Query().Get("limit"))
	if err != nil {
		writeError(writer, err)
		return
	}
	results := s.services.Queries.Search("", query.Limit)
	writeJSON(writer, http.StatusOK, map[string]any{
		"query":   query.Text,
		"results": results,
		"count":   len(results),
	})
}

func (s *Server) exportSamples(writer http.ResponseWriter, request *http.Request) {
	limit := 0
	if rawLimit := request.URL.Query().Get("limit"); rawLimit != "" {
		value, err := strconv.Atoi(rawLimit)
		if err != nil {
			writeError(writer, model.NewError(model.ErrorInvalid, "limit must be an integer"))
			return
		}
		limit = value
	}
	filter, err := validation.SampleFilter(request.URL.Query().Get("status"), request.URL.Query().Get("origin"), limit)
	if err != nil {
		writeError(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer.Header().Set("Content-Disposition", "attachment; filename=samples.csv")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(s.services.Queries.SampleExport(filter)))
}
