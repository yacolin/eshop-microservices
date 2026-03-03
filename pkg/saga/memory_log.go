package saga

import (
	"context"
	"fmt"
	"sync"
)

// MemoryLog 内存实现的 Saga 日志存储（用于测试）
type MemoryLog struct {
	mu    sync.RWMutex
	sagas map[string]*Saga
}

// NewMemoryLog 创建内存日志存储
func NewMemoryLog() *MemoryLog {
	return &MemoryLog{
		sagas: make(map[string]*Saga),
	}
}

// Save 保存 Saga
func (m *MemoryLog) Save(ctx context.Context, saga *Saga) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 深拷贝
	sagaCopy := &Saga{
		ID:          saga.ID,
		Name:        saga.Name,
		Status:      saga.Status,
		CurrentStep: saga.CurrentStep,
		Error:       saga.Error,
		CreatedAt:   saga.CreatedAt,
		UpdatedAt:   saga.UpdatedAt,
		Data:        make(map[string]interface{}),
	}

	if saga.CompletedAt != nil {
		completedAt := *saga.CompletedAt
		sagaCopy.CompletedAt = &completedAt
	}

	// 拷贝 steps
	sagaCopy.Steps = make([]*Step, len(saga.Steps))
	for i, step := range saga.Steps {
		stepCopy := &Step{
			ID:     step.ID,
			Name:   step.Name,
			Status: step.Status,
			Error:  step.Error,
		}

		if step.StartedAt != nil {
			startedAt := *step.StartedAt
			stepCopy.StartedAt = &startedAt
		}
		if step.CompletedAt != nil {
			completedAt := *step.CompletedAt
			stepCopy.CompletedAt = &completedAt
		}
		if step.CompensatedAt != nil {
			compensatedAt := *step.CompensatedAt
			stepCopy.CompensatedAt = &compensatedAt
		}

		sagaCopy.Steps[i] = stepCopy
	}

	// 拷贝 data
	for k, v := range saga.Data {
		sagaCopy.Data[k] = v
	}

	m.sagas[saga.ID] = sagaCopy
	return nil
}

// Get 获取 Saga
func (m *MemoryLog) Get(ctx context.Context, sagaID string) (*Saga, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	saga, ok := m.sagas[sagaID]
	if !ok {
		return nil, fmt.Errorf("saga not found: %s", sagaID)
	}

	return saga, nil
}

// UpdateStep 更新步骤状态
func (m *MemoryLog) UpdateStep(ctx context.Context, sagaID string, stepIndex int, step *Step) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	saga, ok := m.sagas[sagaID]
	if !ok {
		return fmt.Errorf("saga not found: %s", sagaID)
	}

	if stepIndex < 0 || stepIndex >= len(saga.Steps) {
		return fmt.Errorf("invalid step index: %d", stepIndex)
	}

	// 更新步骤
	saga.Steps[stepIndex].Status = step.Status
	saga.Steps[stepIndex].Error = step.Error

	if step.StartedAt != nil {
		startedAt := *step.StartedAt
		saga.Steps[stepIndex].StartedAt = &startedAt
	}
	if step.CompletedAt != nil {
		completedAt := *step.CompletedAt
		saga.Steps[stepIndex].CompletedAt = &completedAt
	}
	if step.CompensatedAt != nil {
		compensatedAt := *step.CompensatedAt
		saga.Steps[stepIndex].CompensatedAt = &compensatedAt
	}

	return nil
}

// List 列出所有 Saga（用于调试）
func (m *MemoryLog) List() []*Saga {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*Saga, 0, len(m.sagas))
	for _, saga := range m.sagas {
		list = append(list, saga)
	}
	return list
}

// Clear 清空所有 Saga（用于测试）
func (m *MemoryLog) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sagas = make(map[string]*Saga)
}
