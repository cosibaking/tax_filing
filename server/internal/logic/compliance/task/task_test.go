package task

import (
	"context"
	"testing"
	"time"
)

type memoryTaskRepository struct {
	nextID uint64
	items  map[string]Task
}

func (m *memoryTaskRepository) FindByKey(_ context.Context, key IdempotencyKey) (*Task, error) {
	item, ok := m.items[key.String()]
	if !ok {
		return nil, nil
	}
	return &item, nil
}

func (m *memoryTaskRepository) Insert(_ context.Context, item Task) (*Task, error) {
	m.nextID++
	item.ID = m.nextID
	m.items[item.Key.String()] = item
	return &item, nil
}

func (m *memoryTaskRepository) Save(_ context.Context, item Task) error {
	m.items[item.Key.String()] = item
	return nil
}

func validGenerateInput() GenerateInput {
	return GenerateInput{
		OpcEntityID: 10, MemberID: 1, RuleVersionID: 100, PeriodKey: "2026-Q3", TriggerKey: "periodic",
		TaskType: "tax", Title: "第三季度申报资料准备", Description: "准备申报资料",
		Materials: []string{"银行流水", "销项发票"}, Priority: "high", DueAt: time.Date(2026, 10, 20, 23, 59, 59, 0, time.Local),
	}
}

func TestGenerateIsIdempotent(t *testing.T) {
	repo := &memoryTaskRepository{items: map[string]Task{}}
	generator := NewGenerator(repo)

	first, created, err := generator.Generate(context.Background(), validGenerateInput())
	if err != nil || !created {
		t.Fatalf("first generate: created=%v err=%v", created, err)
	}
	second, created, err := generator.Generate(context.Background(), validGenerateInput())
	if err != nil || created || first.ID != second.ID || len(repo.items) != 1 {
		t.Fatalf("duplicate generate: created=%v first=%d second=%d items=%d err=%v", created, first.ID, second.ID, len(repo.items), err)
	}
}

func TestGeneratedTaskKeepsRuleSnapshot(t *testing.T) {
	repo := &memoryTaskRepository{items: map[string]Task{}}
	generator := NewGenerator(repo)
	in := validGenerateInput()
	got, _, _ := generator.Generate(context.Background(), in)
	in.Title = "规则更新后的标题"

	if got.Snapshot.Title != "第三季度申报资料准备" {
		t.Fatalf("snapshot changed: %+v", got.Snapshot)
	}
}

func TestTaskStateTransitions(t *testing.T) {
	if got, err := Transition(StatusNotStarted, EventStart); err != nil || got != StatusPreparing {
		t.Fatalf("start got=%s err=%v", got, err)
	}
	if got, err := Transition(StatusPreparing, EventSubmit); err != nil || got != StatusPendingConfirmation {
		t.Fatalf("submit got=%s err=%v", got, err)
	}
	if got, err := Transition(StatusPendingConfirmation, EventReject); err != nil || got != StatusPreparing {
		t.Fatalf("reject got=%s err=%v", got, err)
	}
	if got, err := Transition(StatusPendingConfirmation, EventComplete); err != nil || got != StatusCompleted {
		t.Fatalf("complete got=%s err=%v", got, err)
	}
}

func TestInvalidTaskTransitionFails(t *testing.T) {
	if _, err := Transition(StatusCompleted, EventStart); err == nil {
		t.Fatal("expected completed task transition to fail")
	}
}

func TestMarkOverdueAndCompleteLater(t *testing.T) {
	task := Task{Status: StatusPreparing, DueAt: time.Now().Add(-time.Hour)}
	MarkOverdue(&task, time.Now())
	if task.Status != StatusOverdue {
		t.Fatalf("status=%s want=%s", task.Status, StatusOverdue)
	}
	got, err := Transition(task.Status, EventComplete)
	if err != nil || got != StatusCompleted {
		t.Fatalf("overdue complete got=%s err=%v", got, err)
	}
}
