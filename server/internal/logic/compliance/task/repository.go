package task

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const tableTask = "xy_enterprise_compliance_task"

type Query struct {
	MemberID  uint64
	Status    string
	TaskType  string
	PeriodKey string
	Page      int
	PageSize  int
}

type DatabaseRepository struct{}

type taskRow struct {
	ID             uint64 `orm:"id"`
	OpcEntityID    uint64 `orm:"opc_entity_id"`
	MemberID       uint64 `orm:"member_id"`
	RuleVersionID  uint64 `orm:"rule_version_id"`
	PeriodKey      string `orm:"period_key"`
	TriggerKey     string `orm:"trigger_key"`
	TaskType       string `orm:"task_type"`
	Title          string `orm:"title"`
	Description    string `orm:"description"`
	SnapshotJSON   string `orm:"snapshot_json"`
	MaterialJSON   string `orm:"material_json"`
	Status         string `orm:"status"`
	Priority       string `orm:"priority"`
	DueAt          uint64 `orm:"due_at"`
	CompletedAt    uint64 `orm:"completed_at"`
	CompletionNote string `orm:"completion_note"`
}

func (DatabaseRepository) FindByKey(ctx context.Context, key IdempotencyKey) (*Task, error) {
	var row taskRow
	err := g.DB().Model(tableTask).Ctx(ctx).
		Where("opc_entity_id", key.OpcEntityID).Where("rule_version_id", key.RuleVersionID).
		Where("period_key", key.PeriodKey).Where("trigger_key", key.TriggerKey).
		Where("deleted", 0).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil, err
	}
	return row.task()
}

func (DatabaseRepository) Insert(ctx context.Context, item Task) (*Task, error) {
	snapshot, _ := json.Marshal(item.Snapshot)
	materials, _ := json.Marshal(item.Materials)
	result, err := g.DB().Model(tableTask).Ctx(ctx).Data(g.Map{
		"opc_entity_id": item.Key.OpcEntityID, "member_id": item.MemberID,
		"rule_version_id": item.Key.RuleVersionID, "period_key": item.Key.PeriodKey,
		"trigger_key": item.Key.TriggerKey, "task_type": item.TaskType,
		"title": item.Title, "description": item.Description,
		"snapshot_json": string(snapshot), "material_json": string(materials),
		"status": item.Status, "priority": item.Priority, "due_at": uint64(item.DueAt.Unix()),
		"deleted": 0, "create_time": uint64(time.Now().Unix()), "update_time": uint64(time.Now().Unix()),
	}).Insert()
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	item.ID = uint64(id)
	for _, reminder := range PlanReminders(item.ID, item.DueAt, []string{"in_app"}) {
		if _, err := g.DB().Model("xy_task_reminder").Ctx(ctx).Data(g.Map{
			"task_id": reminder.TaskID, "channel": reminder.Channel,
			"planned_at": uint64(reminder.PlannedAt.Unix()), "status": reminder.Status,
			"attempts": 0, "create_time": uint64(time.Now().Unix()), "update_time": uint64(time.Now().Unix()),
		}).InsertIgnore(); err != nil {
			return nil, err
		}
	}
	return &item, nil
}

func (DatabaseRepository) Save(ctx context.Context, item Task) error {
	completedAt := uint64(0)
	if !item.CompletedAt.IsZero() {
		completedAt = uint64(item.CompletedAt.Unix())
	}
	_, err := g.DB().Model(tableTask).Ctx(ctx).Where("id", item.ID).Data(g.Map{
		"status": item.Status, "completed_at": completedAt,
		"completion_note": item.CompletionNote, "update_time": uint64(time.Now().Unix()),
	}).Update()
	return err
}

func (DatabaseRepository) Get(ctx context.Context, id uint64) (*Task, error) {
	var row taskRow
	err := g.DB().Model(tableTask).Ctx(ctx).Where("id", id).Where("deleted", 0).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil, err
	}
	return row.task()
}

func (DatabaseRepository) List(ctx context.Context, query Query) ([]Task, int, error) {
	model := g.DB().Model(tableTask).Ctx(ctx).Where("member_id", query.MemberID).Where("deleted", 0)
	if query.Status != "" {
		model = model.Where("status", query.Status)
	}
	if query.TaskType != "" {
		model = model.Where("task_type", query.TaskType)
	}
	if query.PeriodKey != "" {
		model = model.Where("period_key", query.PeriodKey)
	}
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}
	page, size := query.Page, query.PageSize
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	var rows []taskRow
	if err := model.OrderAsc("due_at").Page(page, size).Scan(&rows); err != nil {
		return nil, 0, err
	}
	items := make([]Task, 0, len(rows))
	for _, row := range rows {
		item, err := row.task()
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, nil
}

func (row taskRow) task() (*Task, error) {
	var snapshot Snapshot
	if err := json.Unmarshal([]byte(row.SnapshotJSON), &snapshot); err != nil {
		return nil, err
	}
	var materials []string
	if err := json.Unmarshal([]byte(row.MaterialJSON), &materials); err != nil {
		return nil, err
	}
	item := &Task{
		ID: row.ID, Key: IdempotencyKey{OpcEntityID: row.OpcEntityID, RuleVersionID: row.RuleVersionID, PeriodKey: row.PeriodKey, TriggerKey: row.TriggerKey},
		MemberID: row.MemberID, TaskType: row.TaskType, Title: row.Title, Description: row.Description,
		Materials: materials, Snapshot: snapshot, Status: row.Status, Priority: row.Priority,
		DueAt: time.Unix(int64(row.DueAt), 0), CompletionNote: row.CompletionNote,
	}
	if row.CompletedAt > 0 {
		item.CompletedAt = time.Unix(int64(row.CompletedAt), 0)
	}
	return item, nil
}

type QueryRepository interface {
	Repository
	Get(ctx context.Context, id uint64) (*Task, error)
	List(ctx context.Context, query Query) ([]Task, int, error)
}

type Service struct {
	repository QueryRepository
	generator  *Generator
}

func NewService(repository QueryRepository) *Service {
	return &Service{repository: repository, generator: NewGenerator(repository)}
}

func NewDatabaseService() *Service { return NewService(DatabaseRepository{}) }

func (s *Service) Generate(ctx context.Context, in GenerateInput) (*Task, bool, error) {
	return s.generator.Generate(ctx, in)
}

func (s *Service) List(ctx context.Context, query Query) ([]Task, int, error) {
	return s.repository.List(ctx, query)
}

func (s *Service) GetForMember(ctx context.Context, memberID, id uint64) (*Task, error) {
	item, err := s.repository.Get(ctx, id)
	if err != nil || item == nil {
		return item, err
	}
	if item.MemberID != memberID {
		return nil, errors.New("无权访问该合规任务")
	}
	return item, nil
}

func (s *Service) ApplyEvent(ctx context.Context, memberID, id uint64, event, note string) (*Task, error) {
	item, err := s.GetForMember(ctx, memberID, id)
	if err != nil || item == nil {
		return item, err
	}
	next, err := Transition(item.Status, event)
	if err != nil {
		return nil, err
	}
	item.Status = next
	item.CompletionNote = note
	if next == StatusCompleted {
		item.CompletedAt = time.Now()
	}
	if err := s.repository.Save(ctx, *item); err != nil {
		return nil, err
	}
	return item, nil
}
