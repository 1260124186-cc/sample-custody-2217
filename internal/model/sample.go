package model

import "time"

type SampleStatus string

const (
	SampleRegistered SampleStatus = "registered"
	SampleInReview   SampleStatus = "in_review"
	SampleSealed     SampleStatus = "sealed"
)

type Sample struct {
	ID            string       `json:"id"`
	Code          string       `json:"code"`
	Material      string       `json:"material"`
	Origin        string       `json:"origin"`
	Quantity      float64      `json:"quantity"`
	Unit          string       `json:"unit"`
	Status        SampleStatus `json:"status"`
	CurrentHolder string       `json:"current_holder"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (s Sample) IsTransferable() bool {
	return s.Status != SampleSealed
}

func (s Sample) IsReviewable() bool {
	return s.Status == SampleRegistered || s.Status == SampleInReview
}

func (s Sample) Clone() Sample {
	return s
}

type RegisterSampleInput struct {
	ID       string  `json:"id"`
	Code     string  `json:"code"`
	Material string  `json:"material"`
	Origin   string  `json:"origin"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type SampleFilter struct {
	Status SampleStatus
	Origin string
	Limit  int
}
