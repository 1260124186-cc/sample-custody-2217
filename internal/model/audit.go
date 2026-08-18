package model

import "time"

type AuditEvent struct {
	ID        string    `json:"id"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (e AuditEvent) Clone() AuditEvent {
	return e
}
