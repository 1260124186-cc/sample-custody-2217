package model

type OriginSummary struct {
	Origin      string `json:"origin"`
	SampleCount int    `json:"sample_count"`
	SealedCount int    `json:"sealed_count"`
	OpenCount   int    `json:"open_count"`
}

type CustodySummary struct {
	Holder      string `json:"holder"`
	SampleCount int    `json:"sample_count"`
	SealedCount int    `json:"sealed_count"`
}

type ServiceSummary struct {
	SampleCount      int              `json:"sample_count"`
	RegisteredCount  int              `json:"registered_count"`
	InReviewCount    int              `json:"in_review_count"`
	SealedCount      int              `json:"sealed_count"`
	OpenBatchCount   int              `json:"open_batch_count"`
	CompletedBatches int              `json:"completed_batch_count"`
	TransferCount    int              `json:"transfer_count"`
	Origins          []OriginSummary  `json:"origins"`
	Holders          []CustodySummary `json:"holders"`
}

type SearchResult struct {
	ID            string       `json:"id"`
	Code          string       `json:"code"`
	Material      string       `json:"material"`
	Origin        string       `json:"origin"`
	Status        SampleStatus `json:"status"`
	CurrentHolder string       `json:"current_holder"`
	MatchFields   []string     `json:"match_fields"`
}
