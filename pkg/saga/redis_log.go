package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// SagaKeyPrefix Redis key 前缀
	SagaKeyPrefix = "saga:"
	// SagaTTL Saga 数据过期时间（7天）
	SagaTTL = 7 * 24 * time.Hour
)

// RedisLog Redis 实现的 Saga 日志存储
type RedisLog struct {
	client *redis.Client
	prefix string
}

// NewRedisLog 创建 Redis 日志存储
func NewRedisLog(client *redis.Client) *RedisLog {
	return &RedisLog{
		client: client,
		prefix: SagaKeyPrefix,
	}
}

// NewRedisLogWithPrefix 创建带自定义前缀的 Redis 日志存储
func NewRedisLogWithPrefix(client *redis.Client, prefix string) *RedisLog {
	return &RedisLog{
		client: client,
		prefix: prefix,
	}
}

// getKey 生成 Redis key
func (r *RedisLog) getKey(sagaID string) string {
	return r.prefix + sagaID
}

// Save 保存 Saga
func (r *RedisLog) Save(ctx context.Context, saga *Saga) error {
	key := r.getKey(saga.ID)

	// 序列化 Saga
	data, err := saga.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal saga: %w", err)
	}

	// 保存到 Redis，设置过期时间
	if err := r.client.SetEx(ctx, key, data, SagaTTL).Err(); err != nil {
		return fmt.Errorf("failed to save saga to redis: %w", err)
	}

	return nil
}

// Get 获取 Saga
func (r *RedisLog) Get(ctx context.Context, sagaID string) (*Saga, error) {
	key := r.getKey(sagaID)

	// 从 Redis 获取
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("saga not found: %s", sagaID)
		}
		return nil, fmt.Errorf("failed to get saga from redis: %w", err)
	}

	// 反序列化
	saga := &Saga{}
	if err := saga.FromJSON(data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal saga: %w", err)
	}

	return saga, nil
}

// UpdateStep 更新步骤状态
func (r *RedisLog) UpdateStep(ctx context.Context, sagaID string, stepIndex int, step *Step) error {
	// 先获取完整的 Saga
	saga, err := r.Get(ctx, sagaID)
	if err != nil {
		return err
	}

	// 检查步骤索引
	if stepIndex < 0 || stepIndex >= len(saga.Steps) {
		return fmt.Errorf("invalid step index: %d", stepIndex)
	}

	// 更新步骤（保留 Action 和 Compensate 函数）
	// 注意：从 Redis 获取的 Saga 不包含函数，这里只更新状态字段
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

	// 更新 Saga 元数据
	saga.UpdatedAt = time.Now()

	// 重新保存
	return r.Save(ctx, saga)
}

// Delete 删除 Saga
func (r *RedisLog) Delete(ctx context.Context, sagaID string) error {
	key := r.getKey(sagaID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete saga from redis: %w", err)
	}
	return nil
}

// List 列出所有 Saga（用于调试）
func (r *RedisLog) List(ctx context.Context) ([]*Saga, error) {
	pattern := r.prefix + "*"
	var cursor uint64
	var sagaIDs []string

	// 扫描所有 key
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan saga keys: %w", err)
		}

		for _, key := range keys {
			// 去掉前缀获取 sagaID
			sagaID := key[len(r.prefix):]
			sagaIDs = append(sagaIDs, sagaID)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	// 获取所有 Saga
	sagas := make([]*Saga, 0, len(sagaIDs))
	for _, sagaID := range sagaIDs {
		saga, err := r.Get(ctx, sagaID)
		if err != nil {
			// 跳过无法获取的 Saga
			continue
		}
		sagas = append(sagas, saga)
	}

	return sagas, nil
}

// ListByStatus 按状态列出 Saga
func (r *RedisLog) ListByStatus(ctx context.Context, status Status) ([]*Saga, error) {
	allSagas, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*Saga
	for _, saga := range allSagas {
		if saga.Status == status {
			filtered = append(filtered, saga)
		}
	}

	return filtered, nil
}

// Exists 检查 Saga 是否存在
func (r *RedisLog) Exists(ctx context.Context, sagaID string) (bool, error) {
	key := r.getKey(sagaID)
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check saga existence: %w", err)
	}
	return n > 0, nil
}

// GetStep 获取特定步骤
func (r *RedisLog) GetStep(ctx context.Context, sagaID string, stepIndex int) (*Step, error) {
	saga, err := r.Get(ctx, sagaID)
	if err != nil {
		return nil, err
	}

	if stepIndex < 0 || stepIndex >= len(saga.Steps) {
		return nil, fmt.Errorf("invalid step index: %d", stepIndex)
	}

	return saga.Steps[stepIndex], nil
}

// RedisSaga Redis 存储的 Saga 结构（不含函数）
type RedisSaga struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Steps       []*RedisStep           `json:"steps"`
	Status      Status                 `json:"status"`
	CurrentStep int                    `json:"current_step"`
	Error       string                 `json:"error,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// RedisStep Redis 存储的步骤结构（不含函数）
type RedisStep struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Status        Status     `json:"status"`
	Error         string     `json:"error,omitempty"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CompensatedAt *time.Time `json:"compensated_at,omitempty"`
}

// ToRedisSaga 转换为 Redis 存储格式
func (s *Saga) ToRedisSaga() *RedisSaga {
	rs := &RedisSaga{
		ID:          s.ID,
		Name:        s.Name,
		Status:      s.Status,
		CurrentStep: s.CurrentStep,
		Error:       s.Error,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Data:        s.Data,
	}

	if s.CompletedAt != nil {
		completedAt := *s.CompletedAt
		rs.CompletedAt = &completedAt
	}

	rs.Steps = make([]*RedisStep, len(s.Steps))
	for i, step := range s.Steps {
		rs.Steps[i] = &RedisStep{
			ID:     step.ID,
			Name:   step.Name,
			Status: step.Status,
			Error:  step.Error,
		}
		if step.StartedAt != nil {
			startedAt := *step.StartedAt
			rs.Steps[i].StartedAt = &startedAt
		}
		if step.CompletedAt != nil {
			completedAt := *step.CompletedAt
			rs.Steps[i].CompletedAt = &completedAt
		}
		if step.CompensatedAt != nil {
			compensatedAt := *step.CompensatedAt
			rs.Steps[i].CompensatedAt = &compensatedAt
		}
	}

	return rs
}

// Serialize 序列化为 JSON
func (rs *RedisSaga) Serialize() ([]byte, error) {
	return json.Marshal(rs)
}
