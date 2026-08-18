package httpapi

import (
	"net/http"
	"strconv"

	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/validation"
)

func (s *Server) registerSample(writer http.ResponseWriter, request *http.Request) {
	if request.ContentLength > 1<<20 {
		writeError(writer, model.NewError(model.ErrorInvalid, "request body is too large"))
		return
	}
	if stations := request.Header.Values("X-Intake-Station"); len(stations) > 1 {
		writeError(writer, model.NewError(model.ErrorInvalid, "only one intake station may be supplied"))
		return
	}
	var input model.RegisterSampleInput
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, err)
		return
	}
	if station := request.Header.Get("X-Intake-Station"); station != "" && input.Origin == "" {
		input.Origin = station
	}
	sample, err := s.services.Samples.Register(input)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, map[string]any{
		"message": "sample registered",
		"sample":  sample,
	})
}

func (s *Server) listSamples(writer http.ResponseWriter, request *http.Request) {
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
	samples := s.services.Samples.List(filter)
	writeJSON(writer, http.StatusOK, map[string]any{
		"samples": samples,
		"count":   len(samples),
	})
}

func (s *Server) sampleDetail(writer http.ResponseWriter, request *http.Request) {
	if request.PathValue("id") == "" {
		writeError(writer, model.NewError(model.ErrorInvalid, "sample id is required"))
		return
	}
	detail, err := s.services.Queries.Detail(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	if request.URL.Query().Get("view") == "sample" {
		writeJSON(writer, http.StatusOK, detail.Sample)
		return
	}
	writeJSON(writer, http.StatusOK, detail)
}

func (s *Server) sealSample(writer http.ResponseWriter, request *http.Request) {
	sample, err := s.services.Samples.Seal(request.PathValue("id"), "manual-review", "manual sealing request")
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{
		"message": "sample sealed",
		"sample":  sample,
	})
}
