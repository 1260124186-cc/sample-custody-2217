package service

import (
	"example.com/sample-custody/internal/audit"
	"example.com/sample-custody/internal/export"
	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/search"
	"example.com/sample-custody/internal/store"
)

type QueryService struct {
	store *store.Store
	log   *audit.Log
}

type SampleDetail struct {
	Sample    model.Sample       `json:"sample"`
	Transfers []model.Transfer   `json:"transfers"`
	Events    []model.AuditEvent `json:"events"`
}

func (s *QueryService) Detail(id string) (SampleDetail, error) {
	sample, err := s.store.GetSample(id)
	if err != nil {
		return SampleDetail{}, err
	}
	return SampleDetail{
		Sample:    sample,
		Transfers: s.store.TransfersForSample(id),
		Events:    s.log.LatestForEntity(id, 50),
	}, nil
}

func (s *QueryService) Manifest(batchID string) (string, error) {
	batch, err := s.store.GetBatch(batchID)
	if err != nil {
		return "", err
	}
	samples, err := s.store.SamplesByIDs(batch.SampleIDs)
	if err != nil {
		return "", err
	}
	return export.Manifest(batch, samples), nil
}

func (s *QueryService) Summary() model.ServiceSummary {
	return s.store.Summary()
}

func (s *QueryService) Search(query string, limit int) []model.SearchResult {
	return search.Samples(s.store.ListSamples(model.SampleFilter{Status: model.SampleRegistered}), query, limit)
}

func (s *QueryService) SampleExport(filter model.SampleFilter) string {
	return export.SampleCSV(s.store.ListSamples(filter))
}
