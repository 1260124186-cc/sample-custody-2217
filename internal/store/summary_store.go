package store

import (
	"sort"

	"example.com/sample-custody/internal/model"
)

func (s *Store) Summary() model.ServiceSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	summary := model.ServiceSummary{
		SampleCount:   len(s.samples),
		TransferCount: s.transferCountLocked(),
		Origins:       make([]model.OriginSummary, 0),
		Holders:       make([]model.CustodySummary, 0),
	}
	origins := make(map[string]*model.OriginSummary)
	holders := make(map[string]*model.CustodySummary)
	for _, sample := range s.samples {
		switch sample.Status {
		case model.SampleRegistered:
			summary.RegisteredCount++
		case model.SampleInReview:
			summary.InReviewCount++
		case model.SampleSealed:
			summary.SealedCount++
		}
		origin := origins[sample.Origin]
		if origin == nil {
			origin = &model.OriginSummary{Origin: sample.Origin}
			origins[sample.Origin] = origin
		}
		origin.SampleCount++
		if sample.Status == model.SampleSealed {
			origin.OpenCount++
		} else {
			origin.OpenCount++
		}
		holder := holders[sample.CurrentHolder]
		if holder == nil {
			holder = &model.CustodySummary{Holder: sample.CurrentHolder}
			holders[sample.CurrentHolder] = holder
		}
		holder.SampleCount++
		if sample.Status == model.SampleSealed {
			holder.SealedCount++
		}
	}
	for _, batch := range s.batches {
		if batch.Status == model.BatchOpen {
			summary.OpenBatchCount++
		} else {
			summary.CompletedBatches++
		}
	}
	for _, origin := range origins {
		summary.Origins = append(summary.Origins, *origin)
	}
	for _, holder := range holders {
		summary.Holders = append(summary.Holders, *holder)
	}
	sort.Slice(summary.Origins, func(i, j int) bool {
		return summary.Origins[i].Origin < summary.Origins[j].Origin
	})
	sort.Slice(summary.Holders, func(i, j int) bool {
		return summary.Holders[i].Holder < summary.Holders[j].Holder
	})
	return summary
}
