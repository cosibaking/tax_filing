package task

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type IdempotencyKey struct {
	OpcEntityID   uint64
	RuleVersionID uint64
	PeriodKey     string
	TriggerKey    string
}

func (k IdempotencyKey) String() string {
	return fmt.Sprintf("%d:%d:%s:%s", k.OpcEntityID, k.RuleVersionID, k.PeriodKey, k.TriggerKey)
}

type Snapshot struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Materials   []string `json:"materials"`
	Priority    string   `json:"priority"`
	DueAt       int64    `json:"dueAt"`
}

type Task struct {
	ID             uint64         `json:"id"`
	Key            IdempotencyKey `json:"key"`
	MemberID       uint64         `json:"memberId"`
	TaskType       string         `json:"taskType"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Materials      []string       `json:"materials"`
	Snapshot       Snapshot       `json:"snapshot"`
	Status         string         `json:"status"`
	Priority       string         `json:"priority"`
	DueAt          time.Time      `json:"dueAt"`
	CompletedAt    time.Time      `json:"completedAt"`
	CompletionNote string         `json:"completionNote"`
}

type GenerateInput struct {
	OpcEntityID   uint64
	MemberID      uint64
	RuleVersionID uint64
	PeriodKey     string
	TriggerKey    string
	TaskType      string
	Title         string
	Description   string
	Materials     []string
	Priority      string
	DueAt         time.Time
}

type Repository interface {
	FindByKey(ctx context.Context, key IdempotencyKey) (*Task, error)
	Insert(ctx context.Context, item Task) (*Task, error)
	Save(ctx context.Context, item Task) error
}

type Generator struct {
	repository Repository
}

func NewGenerator(repository Repository) *Generator {
	return &Generator{repository: repository}
}

func (g *Generator) Generate(ctx context.Context, in GenerateInput) (*Task, bool, error) {
	if err := validateGenerateInput(in); err != nil {
		return nil, false, err
	}
	key := IdempotencyKey{
		OpcEntityID: in.OpcEntityID, RuleVersionID: in.RuleVersionID,
		PeriodKey: in.PeriodKey, TriggerKey: in.TriggerKey,
	}
	existing, err := g.repository.FindByKey(ctx, key)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		return existing, false, nil
	}
	materials := append([]string(nil), in.Materials...)
	item := Task{
		Key: key, MemberID: in.MemberID, TaskType: in.TaskType,
		Title: in.Title, Description: in.Description, Materials: materials,
		Snapshot: Snapshot{
			Title: in.Title, Description: in.Description, Materials: append([]string(nil), materials...),
			Priority: in.Priority, DueAt: in.DueAt.Unix(),
		},
		Status: StatusNotStarted, Priority: in.Priority, DueAt: in.DueAt,
	}
	saved, err := g.repository.Insert(ctx, item)
	if err != nil {
		// 数据库唯一键是并发幂等的最终防线；冲突后返回已存在任务。
		if existing, lookupErr := g.repository.FindByKey(ctx, key); lookupErr == nil && existing != nil {
			return existing, false, nil
		}
		return nil, false, err
	}
	return saved, true, nil
}

func MarkOverdue(item *Task, now time.Time) {
	if item == nil || item.DueAt.IsZero() || !now.After(item.DueAt) {
		return
	}
	if item.Status == StatusNotStarted || item.Status == StatusPreparing {
		item.Status = StatusOverdue
	}
}

func validateGenerateInput(in GenerateInput) error {
	if in.OpcEntityID == 0 || in.MemberID == 0 || in.RuleVersionID == 0 {
		return errors.New("企业、会员和规则版本不能为空")
	}
	if strings.TrimSpace(in.PeriodKey) == "" || strings.TrimSpace(in.TriggerKey) == "" {
		return errors.New("任务期间和触发键不能为空")
	}
	if strings.TrimSpace(in.Title) == "" || in.DueAt.IsZero() {
		return errors.New("任务标题和截止时间不能为空")
	}
	if in.Priority != "low" && in.Priority != "medium" && in.Priority != "high" {
		return errors.New("任务优先级无效")
	}
	return nil
}
