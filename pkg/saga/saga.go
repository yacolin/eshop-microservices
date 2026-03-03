package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"eshop-microservices/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Status Saga 状态
type Status string

const (
	StatusPending    Status = "pending"    // 待执行
	StatusRunning    Status = "running"    // 执行中
	StatusSucceeded  Status = "succeeded"  // 成功
	StatusFailed     Status = "failed"     // 失败
	StatusCompensating Status = "compensating" // 补偿中
	StatusCompensated Status = "compensated"   // 已补偿
)

// Step Saga 步骤
type Step struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Action        func(ctx context.Context) error `json:"-"`
	Compensate    func(ctx context.Context) error `json:"-"`
	Status        Status          `json:"status"`
	Error         string          `json:"error,omitempty"`
	StartedAt     *time.Time      `json:"started_at,omitempty"`
	CompletedAt   *time.Time      `json:"completed_at,omitempty"`
	CompensatedAt *time.Time      `json:"compensated_at,omitempty"`
}

// Saga 分布式事务 Saga
type Saga struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Steps       []*Step   `json:"steps"`
	Status      Status    `json:"status"`
	CurrentStep int       `json:"current_step"`
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// SagaLog Saga 日志存储接口
type SagaLog interface {
	Save(ctx context.Context, saga *Saga) error
	Get(ctx context.Context, sagaID string) (*Saga, error)
	UpdateStep(ctx context.Context, sagaID string, stepIndex int, step *Step) error
}

// NewSaga 创建新的 Saga
func NewSaga(name string) *Saga {
	return &Saga{
		ID:          uuid.New().String(),
		Name:        name,
		Steps:       make([]*Step, 0),
		Status:      StatusPending,
		CurrentStep: -1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Data:        make(map[string]interface{}),
	}
}

// AddStep 添加步骤
func (s *Saga) AddStep(name string, action, compensate func(ctx context.Context) error) *Saga {
	step := &Step{
		ID:         uuid.New().String(),
		Name:       name,
		Action:     action,
		Compensate: compensate,
		Status:     StatusPending,
	}
	s.Steps = append(s.Steps, step)
	return s
}

// SetData 设置 Saga 数据
func (s *Saga) SetData(key string, value interface{}) *Saga {
	s.Data[key] = value
	return s
}

// GetData 获取 Saga 数据
func (s *Saga) GetData(key string) (interface{}, bool) {
	val, ok := s.Data[key]
	return val, ok
}

// Coordinator Saga 协调器
type Coordinator struct {
	log SagaLog
}

// NewCoordinator 创建 Saga 协调器
func NewCoordinator(log SagaLog) *Coordinator {
	return &Coordinator{log: log}
}

// Execute 执行 Saga
func (c *Coordinator) Execute(ctx context.Context, saga *Saga) error {
	logger.Info("starting saga execution",
		zap.String("saga_id", saga.ID),
		zap.String("saga_name", saga.Name))

	saga.Status = StatusRunning
	saga.UpdatedAt = time.Now()

	// 保存初始状态
	if err := c.log.Save(ctx, saga); err != nil {
		return fmt.Errorf("failed to save saga: %w", err)
	}

	// 执行每个步骤
	for i, step := range saga.Steps {
		saga.CurrentStep = i
		now := time.Now()
		step.StartedAt = &now
		step.Status = StatusRunning

		logger.Info("executing saga step",
			zap.String("saga_id", saga.ID),
			zap.String("step_id", step.ID),
			zap.String("step_name", step.Name),
			zap.Int("step_index", i))

		// 执行动作
		if err := step.Action(ctx); err != nil {
			step.Status = StatusFailed
			step.Error = err.Error()
			completedAt := time.Now()
			step.CompletedAt = &completedAt
			saga.Error = fmt.Sprintf("step %d (%s) failed: %v", i, step.Name, err)
			saga.Status = StatusFailed
			saga.UpdatedAt = time.Now()

			// 保存失败状态
			if saveErr := c.log.Save(ctx, saga); saveErr != nil {
				logger.Error("failed to save saga failure state",
					zap.String("saga_id", saga.ID),
					zap.Error(saveErr))
			}

			// 触发补偿
			logger.Warn("saga step failed, starting compensation",
				zap.String("saga_id", saga.ID),
				zap.String("step_name", step.Name),
				zap.Error(err))

			if compErr := c.compensate(ctx, saga, i); compErr != nil {
				logger.Error("saga compensation failed",
					zap.String("saga_id", saga.ID),
					zap.Error(compErr))
				return fmt.Errorf("saga failed and compensation also failed: %v, compensation error: %w", err, compErr)
			}

			return fmt.Errorf("saga failed and compensated: %w", err)
		}

		// 步骤成功
		step.Status = StatusSucceeded
		completedAt := time.Now()
		step.CompletedAt = &completedAt
		saga.UpdatedAt = time.Now()

		// 更新步骤状态
		if err := c.log.UpdateStep(ctx, saga.ID, i, step); err != nil {
			logger.Error("failed to update step status",
				zap.String("saga_id", saga.ID),
				zap.String("step_id", step.ID),
				zap.Error(err))
		}

		logger.Info("saga step completed",
			zap.String("saga_id", saga.ID),
			zap.String("step_name", step.Name))
	}

	// Saga 成功完成
	saga.Status = StatusSucceeded
	now := time.Now()
	saga.CompletedAt = &now
	saga.UpdatedAt = now

	if err := c.log.Save(ctx, saga); err != nil {
		logger.Error("failed to save saga completion",
			zap.String("saga_id", saga.ID),
			zap.Error(err))
	}

	logger.Info("saga completed successfully",
		zap.String("saga_id", saga.ID),
		zap.String("saga_name", saga.Name))

	return nil
}

// compensate 执行补偿
func (c *Coordinator) compensate(ctx context.Context, saga *Saga, failedStepIndex int) error {
	saga.Status = StatusCompensating
	saga.UpdatedAt = time.Now()

	logger.Info("starting saga compensation",
		zap.String("saga_id", saga.ID),
		zap.Int("failed_step", failedStepIndex),
		zap.Int("steps_to_compensate", failedStepIndex))

	// 逆序执行补偿
	for i := failedStepIndex - 1; i >= 0; i-- {
		step := saga.Steps[i]

		if step.Compensate == nil {
			logger.Warn("step has no compensation function, skipping",
				zap.String("saga_id", saga.ID),
				zap.String("step_name", step.Name))
			continue
		}

		logger.Info("compensating saga step",
			zap.String("saga_id", saga.ID),
			zap.String("step_name", step.Name),
			zap.Int("step_index", i))

		if err := step.Compensate(ctx); err != nil {
			step.Error = fmt.Sprintf("compensation failed: %v", err)
			logger.Error("step compensation failed",
				zap.String("saga_id", saga.ID),
				zap.String("step_name", step.Name),
				zap.Error(err))
			// 继续补偿其他步骤
		} else {
			compensatedAt := time.Now()
			step.CompensatedAt = &compensatedAt
			logger.Info("step compensated successfully",
				zap.String("saga_id", saga.ID),
				zap.String("step_name", step.Name))
		}

		// 更新步骤状态
		if err := c.log.UpdateStep(ctx, saga.ID, i, step); err != nil {
			logger.Error("failed to update compensation status",
				zap.String("saga_id", saga.ID),
				zap.String("step_id", step.ID),
				zap.Error(err))
		}
	}

	saga.Status = StatusCompensated
	saga.UpdatedAt = time.Now()

	if err := c.log.Save(ctx, saga); err != nil {
		logger.Error("failed to save compensated saga",
			zap.String("saga_id", saga.ID),
			zap.Error(err))
	}

	logger.Info("saga compensation completed",
		zap.String("saga_id", saga.ID))

	return nil
}

// GetSaga 获取 Saga 状态
func (c *Coordinator) GetSaga(ctx context.Context, sagaID string) (*Saga, error) {
	return c.log.Get(ctx, sagaID)
}

// ToJSON 序列化为 JSON
func (s *Saga) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON 从 JSON 反序列化
func (s *Saga) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}
