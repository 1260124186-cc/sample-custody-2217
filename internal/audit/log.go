package audit

import (
	"sort"
	"sync"

	"example.com/sample-custody/internal/model"
)

// Log is an append-only in-memory audit trail shared by all business workflows.
type Log struct {
	mu     sync.RWMutex
	events []model.AuditEvent
}

func NewLog() *Log {
	return &Log{events: make([]model.AuditEvent, 0)}
}

func (l *Log) Append(event model.AuditEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, event.Clone())
}

func (l *Log) ByEntity(entityID string) []model.AuditEvent {
	l.mu.RLock()
	defer l.mu.RUnlock()
	result := make([]model.AuditEvent, 0)
	for _, event := range l.events {
		if event.EntityID == entityID {
			result = append(result, event.Clone())
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result
}
