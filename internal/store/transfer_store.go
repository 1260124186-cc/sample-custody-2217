package store

import (
	"example.com/sample-custody/internal/model"
	"example.com/sample-custody/internal/policy"
)

func (s *Store) TransfersForSample(sampleID string) []model.Transfer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := s.transfers[sampleID]
	result := make([]model.Transfer, len(entries))
	for index, entry := range entries {
		result[index] = entry.Clone()
	}
	return result
}

func (s *Store) UpdateSampleAndAppendTransfer(sample model.Sample, transfer model.Transfer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.samples[sample.ID]; !exists {
		return model.NewError(model.ErrorNotFound, "sample %q was not found", sample.ID)
	}
	s.samples[sample.ID] = sample.Clone()
	s.transfers[sample.ID] = append(s.transfers[sample.ID], transfer.Clone())
	return nil
}

// UpdateSampleAndAppendTransferGuarded 原子地校验交接前置条件：在持有写锁的状态下
// 确认样本当前保管人仍与 guard 记录的 ExpectedHolder 一致。这是消除 check-then-act
// 竞态的关键——Transfer 在读锁下读取样本快照后释放锁，两个并发请求可能读到同一旧值，
// 仅当真正写入时再次校验，才能保证同一时刻只有一次有效交接。
func (s *Store) UpdateSampleAndAppendTransferGuarded(sample model.Sample, transfer model.Transfer, guard model.TransferGuard) error {
	if !guard.Valid() {
		return model.NewError(model.ErrorInvalid, "transfer guard is incomplete")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, exists := s.samples[sample.ID]
	if !exists {
		return model.NewError(model.ErrorNotFound, "sample %q was not found", sample.ID)
	}
	if !policy.HolderMatches(current.CurrentHolder, guard.ExpectedHolder) {
		return model.NewError(
			model.ErrorConflict,
			"sample %q holder changed from %q to %q before transfer committed",
			sample.ID,
			guard.ExpectedHolder,
			current.CurrentHolder,
		)
	}
	s.samples[sample.ID] = sample.Clone()
	s.transfers[sample.ID] = append(s.transfers[sample.ID], transfer.Clone())
	return nil
}
