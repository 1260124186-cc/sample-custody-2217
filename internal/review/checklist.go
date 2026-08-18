package review

import (
	"fmt"
	"sort"
	"strings"

	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/policy"
)

type Severity string

const (
	SeverityInfo  Severity = "info"
	SeverityWarn  Severity = "warning"
	SeverityBlock Severity = "blocking"
)

type Finding struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	SampleID string   `json:"sample_id,omitempty"`
}

type Checklist struct {
	Allowed  bool
	Findings []Finding
}

func (c Checklist) BlockingCount() int {
	count := 0
	for _, finding := range c.Findings {
		if finding.Severity == SeverityBlock {
			count++
		}
	}
	return count
}

func (c Checklist) WarningCount() int {
	count := 0
	for _, finding := range c.Findings {
		if finding.Severity == SeverityWarn {
			count++
		}
	}
	return count
}

func (c Checklist) Codes() []string {
	codes := make([]string, 0, len(c.Findings))
	for _, finding := range c.Findings {
		codes = append(codes, finding.Code)
	}
	sort.Strings(codes)
	return codes
}

func ForRegistration(input model.RegisterSampleInput) Checklist {
	findings := make([]Finding, 0)
	if strings.TrimSpace(input.Material) == "" {
		findings = append(findings, Finding{
			Code:     "missing-material",
			Severity: SeverityBlock,
			Message:  "material must be recorded before intake",
		})
	}
	if strings.TrimSpace(input.Origin) == "" {
		findings = append(findings, Finding{
			Code:     "missing-origin",
			Severity: SeverityBlock,
			Message:  "origin must be recorded before intake",
		})
	}
	quantity := policy.AssessQuantity(input.Quantity, input.Unit)
	if !quantity.Acceptable {
		findings = append(findings, Finding{
			Code:     "invalid-quantity",
			Severity: SeverityBlock,
			Message:  quantity.Reason,
		})
	}
	return Checklist{Allowed: !hasBlocking(findings), Findings: findings}
}

func ForTransfer(sample model.Sample, from, to string) Checklist {
	findings := make([]Finding, 0, 3)
	if !sample.IsTransferable() {
		findings = append(findings, Finding{
			Code:     "sealed-sample",
			Severity: SeverityBlock,
			Message:  "sealed specimens cannot move",
			SampleID: sample.ID,
		})
	}
	if !policy.HolderMatches(sample.CurrentHolder, from) {
		findings = append(findings, Finding{
			Code:     "stale-holder",
			Severity: SeverityBlock,
			Message:  fmt.Sprintf("expected holder %q", sample.CurrentHolder),
			SampleID: sample.ID,
		})
	}
	if strings.EqualFold(strings.TrimSpace(from), strings.TrimSpace(to)) {
		findings = append(findings, Finding{
			Code:     "same-holder",
			Severity: SeverityBlock,
			Message:  "a transfer needs a new holder",
			SampleID: sample.ID,
		})
	}
	return Checklist{Allowed: !hasBlocking(findings), Findings: sortFindings(findings)}
}

func ForBatch(samples []model.Sample) Checklist {
	findings := make([]Finding, 0)
	if len(samples) == 0 {
		findings = append(findings, Finding{
			Code:     "empty-batch",
			Severity: SeverityBlock,
			Message:  "an inspection batch must contain at least one specimen",
		})
	}
	seenCodes := make(map[string]string, len(samples))
	for _, sample := range samples {
		if !policy.CanEnterBatch(sample) {
			findings = append(findings, Finding{
				Code:     "sample-not-reviewable",
				Severity: SeverityBlock,
				Message:  fmt.Sprintf("sample %q is not eligible for review", sample.ID),
				SampleID: sample.ID,
			})
		}
		// 编码去重大小写不敏感，仅大小写不同的编码视为同一编码
		code := strings.ToLower(strings.TrimSpace(sample.Code))
		if previous, exists := seenCodes[code]; exists {
			findings = append(findings, Finding{
				Code:     "duplicate-code",
				Severity: SeverityBlock,
				Message:  fmt.Sprintf("samples %q and %q share the same code", previous, sample.ID),
				SampleID: sample.ID,
			})
		}
		seenCodes[code] = sample.ID
	}
	return Checklist{Allowed: !hasBlocking(findings), Findings: sortFindings(findings)}
}

func ForCompletion(samples []model.Sample) Checklist {
	findings := make([]Finding, 0)
	for _, sample := range samples {
		if !policy.CanSeal(sample) {
			findings = append(findings, Finding{
				Code:     "sample-cannot-seal",
				Severity: SeverityBlock,
				Message:  fmt.Sprintf("sample %q cannot be sealed from status %q", sample.ID, sample.Status),
				SampleID: sample.ID,
			})
		}
		if strings.TrimSpace(sample.CurrentHolder) == "" {
			findings = append(findings, Finding{
				Code:     "holder-missing",
				Severity: SeverityWarn,
				Message:  "sample has no current holder",
				SampleID: sample.ID,
			})
		}
	}
	return Checklist{Allowed: !hasBlocking(findings), Findings: sortFindings(findings)}
}

func Explain(checklist Checklist) string {
	if checklist.Allowed && len(checklist.Findings) == 0 {
		return "review passed"
	}
	parts := make([]string, 0, len(checklist.Findings))
	for _, finding := range checklist.Findings {
		parts = append(parts, finding.Code+": "+finding.Message)
	}
	sort.Strings(parts)
	prefix := fmt.Sprintf(
		"%d blocking, %d warning (%s)",
		checklist.BlockingCount(),
		checklist.WarningCount(),
		strings.Join(checklist.Codes(), ","),
	)
	return prefix + ": " + strings.Join(parts, "; ")
}

func hasBlocking(findings []Finding) bool {
	for _, finding := range findings {
		if finding.Severity == SeverityBlock {
			return true
		}
	}
	return false
}

func sortFindings(findings []Finding) []Finding {
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].Code == findings[j].Code {
			return findings[i].SampleID < findings[j].SampleID
		}
		return findings[i].Code < findings[j].Code
	})
	return findings
}
