package model

import "time"

type Transfer struct {
	ID        string    `json:"id"`
	SampleID  string    `json:"sample_id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Location  string    `json:"location"`
	Operator  string    `json:"operator"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

func (t Transfer) Clone() Transfer {
	return t
}

type TransferInput struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Location string `json:"location"`
	Operator string `json:"operator"`
	Note     string `json:"note"`
}

type TransferGuard struct {
	SampleID       string
	ExpectedHolder string
}

func (g TransferGuard) Valid() bool {
	return g.SampleID != "" && g.ExpectedHolder != ""
}
