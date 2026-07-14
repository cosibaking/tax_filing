package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

const tableReport = "xy_compliance_monthly_report"

type DatabaseRepository struct{}

func (DatabaseRepository) SaveDraft(ctx context.Context, opcID, memberID uint64, in Input, result Result) (*MonthlyReport, error) {
	stats, err := json.Marshal(in.Statistics)
	if err != nil {
		return nil, fmt.Errorf("序列化经营统计快照失败: %w", err)
	}
	complete, err := json.Marshal(in.Completeness)
	if err != nil {
		return nil, fmt.Errorf("序列化资料完整度快照失败: %w", err)
	}
	risks, err := json.Marshal(in.Risks)
	if err != nil {
		return nil, fmt.Errorf("序列化风险快照失败: %w", err)
	}
	structured, err := json.Marshal(result.Structured)
	if err != nil {
		return nil, fmt.Errorf("序列化结构化报告失败: %w", err)
	}
	now := uint64(time.Now().Unix())
	var id int64
	var version uint
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// Existing rows serialize concurrent version generation. For the first row,
		// the unique constraint remains the cross-database concurrency backstop.
		if _, err := tx.Model(tableReport).Ctx(ctx).
			Where("opc_entity_id", opcID).Where("period_key", in.PeriodKey).
			Fields("id").OrderAsc("id").LockUpdate().Array(); err != nil {
			return err
		}
		// Use a fresh model so Fields/Order/Lock state cannot leak into MAX SQL.
		max, err := tx.Model(tableReport).Ctx(ctx).
			Where("opc_entity_id", opcID).Where("period_key", in.PeriodKey).
			Max("version")
		if err != nil {
			return err
		}
		version = uint(max) + 1
		res, err := tx.Model(tableReport).Ctx(ctx).Data(g.Map{"opc_entity_id": opcID, "member_id": memberID, "period_key": in.PeriodKey, "version": version, "status": "draft", "statistics_json": string(stats), "completeness_json": string(complete), "risk_snapshot_json": string(risks), "structured_json": string(structured), "rule_versions_json": "[]", "content": result.Content, "ai_model": result.AIModel, "knowledge_version": result.KnowledgeVersion, "create_time": now, "update_time": now}).Insert()
		if err != nil {
			return normalizeVersionConflict(err)
		}
		id, err = res.LastInsertId()
		return err
	})
	if err != nil {
		return nil, err
	}
	structuredReport := result.Structured
	return &MonthlyReport{ID: uint64(id), OpcEntityID: opcID, MemberID: memberID, PeriodKey: in.PeriodKey, Version: version, Status: "draft", Content: result.Content, AIModel: result.AIModel, KnowledgeVersion: result.KnowledgeVersion, StructuredJSON: string(structured), StructuredReport: &structuredReport, Legacy: false, CreatedAt: now}, nil
}

func (DatabaseRepository) Get(ctx context.Context, memberID, id uint64) (*MonthlyReport, error) {
	var item MonthlyReport
	if err := g.DB().Model(tableReport).Ctx(ctx).Where("id", id).Where("member_id", memberID).Scan(&item); err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, errors.New("报告不存在或无权访问")
	}
	if err := item.ApplyStructuredJSON(item.StructuredJSON); err != nil {
		return nil, fmt.Errorf("报告 %d 结构化快照损坏: %w", item.ID, err)
	}
	return &item, nil
}

func (DatabaseRepository) List(ctx context.Context, memberID uint64, period string) ([]MonthlyReport, error) {
	model := g.DB().Model(tableReport).Ctx(ctx).Where("member_id", memberID)
	if period != "" {
		model = model.Where("period_key", period)
	}
	var items []MonthlyReport
	if err := model.OrderDesc("period_key,version").Scan(&items); err != nil {
		return nil, err
	}
	for i := range items {
		if err := items[i].ApplyStructuredJSON(items[i].StructuredJSON); err != nil {
			g.Log().Warningf(ctx, "月度报告结构化快照解析失败，报告ID=%d: %v", items[i].ID, err)
		}
	}
	return items, nil
}
func (DatabaseRepository) Publish(ctx context.Context, memberID, id uint64) (*MonthlyReport, error) {
	item, err := (DatabaseRepository{}).Get(ctx, memberID, id)
	if err != nil {
		return nil, err
	}
	if err := ValidateMutation(item.Status); err != nil {
		return nil, err
	}
	now := uint64(time.Now().Unix())
	_, err = g.DB().Model(tableReport).Ctx(ctx).Where("id", id).Data(g.Map{"status": "published", "published_by": memberID, "published_at": now, "update_time": now}).Update()
	item.Status = "published"
	item.PublishedAt = now
	return item, err
}

func (item *MonthlyReport) ApplyStructuredJSON(raw string) error {
	item.StructuredReport = nil
	item.Legacy = true
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var structured StructuredReport
	if err := json.Unmarshal([]byte(raw), &structured); err != nil {
		return fmt.Errorf("解析结构化报告 JSON 失败: %w", err)
	}
	if structured.SchemaVersion <= 0 {
		return errors.New("结构化报告 schemaVersion 必须大于 0")
	}
	item.StructuredReport = &structured
	item.Legacy = false
	return nil
}

func normalizeVersionConflict(err error) error {
	if err == nil {
		return nil
	}
	constraintMatches := strings.Contains(strings.ToLower(err.Error()), "uk_report_version")
	if !constraintMatches {
		return err
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return fmt.Errorf("报告版本生成冲突，请重试: %w", err)
	}
	var stateErr interface{ SQLState() string }
	if errors.As(err, &stateErr) && stateErr.SQLState() == "23505" {
		return fmt.Errorf("报告版本生成冲突，请重试: %w", err)
	}
	return err
}
func NewDatabaseService() *Service { return NewService(DatabaseRepository{}, nil) }
func (s *Service) Create(ctx context.Context, opcID, memberID uint64, in Input) (*MonthlyReport, error) {
	result, err := s.Build(ctx, in)
	if err != nil {
		return nil, err
	}
	return s.repository.SaveDraft(ctx, opcID, memberID, in, result)
}
func (s *Service) List(ctx context.Context, memberID uint64, period string) ([]MonthlyReport, error) {
	return s.repository.List(ctx, memberID, period)
}
func (s *Service) Publish(ctx context.Context, memberID, id uint64) (*MonthlyReport, error) {
	return s.repository.Publish(ctx, memberID, id)
}
