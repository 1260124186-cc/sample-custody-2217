package audit

import (
	"sort"

	"example.com/sample-custody/internal/model"
)

type ActionSummary struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
}

func (l *Log) SummaryByAction() []ActionSummary {
	l.mu.RLock()
	defer l.mu.RUnlock()
	counts := make(map[string]int)
	for _, event := range l.events {
		counts[event.Action]++
	}
	summary := make([]ActionSummary, 0, len(counts))
	for action, count := range counts {
		summary = append(summary, ActionSummary{Action: action, Count: count})
	}
	sort.Slice(summary, func(i, j int) bool {
		if summary[i].Count == summary[j].Count {
			return summary[i].Action < summary[j].Action
		}
		return summary[i].Count > summary[j].Count
	})
	return summary
}

func (l *Log) LatestForEntity(entityID string, limit int) []model.AuditEvent {
	events := l.ByEntity(entityID)
	if limit <= 0 || len(events) <= limit {
		return events
	}
	start := len(events) - limit
	return append([]model.AuditEvent(nil), events[start:]...)
}
