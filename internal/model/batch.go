package model

import "time"

type BatchStatus string

const (
	BatchOpen      BatchStatus = "open"
	BatchCompleted BatchStatus = "completed"
)

type Batch struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Purpose     string      `json:"purpose"`
	SampleIDs   []string    `json:"sample_ids"`
	Status      BatchStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}

func (b Batch) Clone() Batch {
	return b
}

func (b Batch) IsOpen() bool {
	return b.Status == BatchOpen
}

type CreateBatchInput struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Purpose   string   `json:"purpose"`
	SampleIDs []string `json:"sample_ids"`
}

type BatchCompletion struct {
	Batch       Batch `json:"batch"`
	SealedCount int   `json:"sealed_count"`
}
