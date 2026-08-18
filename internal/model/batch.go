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

// Clone 深拷贝 SampleIDs 切片，避免与外部共享底层数组导致历史数据被污染。
func (b Batch) Clone() Batch {
	clone := b
	if b.SampleIDs != nil {
		ids := make([]string, len(b.SampleIDs))
		copy(ids, b.SampleIDs)
		clone.SampleIDs = ids
	}
	return clone
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
