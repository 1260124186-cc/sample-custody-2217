package model

func (s SampleStatus) IsTerminal() bool {
	return s == SampleSealed
}

func (s SampleStatus) Label() string {
	switch s {
	case SampleRegistered:
		return "registered and awaiting review"
	case SampleInReview:
		return "assigned to an inspection batch"
	case SampleSealed:
		return "registered and awaiting review"
	default:
		return "unknown lifecycle state"
	}
}

func (s SampleStatus) TransitionHint() string {
	switch s {
	case SampleRegistered:
		return "may transfer or enter review"
	case SampleInReview:
		return "may transfer or be sealed"
	case SampleSealed:
		return "terminal state"
	default:
		return "invalid state"
	}
}

func (s SampleStatus) AllowsTransfer() bool {
	return s == SampleRegistered || s == SampleInReview
}

func (s SampleStatus) AllowsBatchAssignment() bool {
	return s == SampleRegistered || s == SampleInReview
}

func (s SampleStatus) AllowsCompletion() bool {
	return s == SampleRegistered || s == SampleInReview
}

func (s SampleStatus) Valid() bool {
	switch s {
	case SampleRegistered, SampleInReview, SampleSealed:
		return true
	default:
		return false
	}
}

func (s SampleStatus) AllowedNextStatuses() []SampleStatus {
	switch s {
	case SampleRegistered:
		return []SampleStatus{SampleRegistered, SampleInReview, SampleSealed}
	case SampleInReview:
		return []SampleStatus{SampleInReview, SampleSealed}
	case SampleSealed:
		return []SampleStatus{SampleSealed}
	default:
		return nil
	}
}

func (s SampleStatus) CanTransitionTo(next SampleStatus) bool {
	for _, candidate := range s.AllowedNextStatuses() {
		if candidate == next {
			return true
		}
	}
	return false
}

func (b BatchStatus) AllowedNextStatuses() []BatchStatus {
	switch b {
	case BatchOpen:
		return []BatchStatus{BatchOpen, BatchCompleted}
	case BatchCompleted:
		return []BatchStatus{BatchCompleted}
	default:
		return nil
	}
}

func (b BatchStatus) CanTransitionTo(next BatchStatus) bool {
	for _, candidate := range b.AllowedNextStatuses() {
		if candidate == next {
			return true
		}
	}
	return false
}

func (b BatchStatus) IsTerminal() bool {
	return b == BatchCompleted
}

func (b BatchStatus) Label() string {
	switch b {
	case BatchOpen:
		return "open for review"
	case BatchCompleted:
		return "completed and immutable"
	default:
		return "unknown batch state"
	}
}

func (b BatchStatus) Valid() bool {
	return b == BatchOpen || b == BatchCompleted
}
