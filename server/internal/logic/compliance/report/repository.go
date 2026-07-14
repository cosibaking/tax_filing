package report

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

const tableReport = "xy_compliance_monthly_report"

type DatabaseRepository struct{}

func (DatabaseRepository) SaveDraft(ctx context.Context, opcID, memberID uint64, in Input, result Result) (*MonthlyReport, error) {
	stats, _ := json.Marshal(in.Statistics)
	complete, _ := json.Marshal(in.Completeness)
	risks, _ := json.Marshal(in.Risks)
	now := uint64(time.Now().Unix())
	max, err := g.DB().Model(tableReport).Ctx(ctx).Where("opc_entity_id", opcID).Where("period_key", in.PeriodKey).Max("version")
	if err != nil {
		return nil, err
	}
	version := uint(max) + 1
	res, err := g.DB().Model(tableReport).Ctx(ctx).Data(g.Map{"opc_entity_id": opcID, "member_id": memberID, "period_key": in.PeriodKey, "version": version, "status": "draft", "statistics_json": string(stats), "completeness_json": string(complete), "risk_snapshot_json": string(risks), "rule_versions_json": "[]", "content": result.Content, "ai_model": result.AIModel, "knowledge_version": result.KnowledgeVersion, "create_time": now, "update_time": now}).Insert()
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &MonthlyReport{ID: uint64(id), OpcEntityID: opcID, MemberID: memberID, PeriodKey: in.PeriodKey, Version: version, Status: "draft", Content: result.Content, AIModel: result.AIModel}, nil
}
func (DatabaseRepository) List(ctx context.Context, memberID uint64, period string) ([]MonthlyReport, error) {
	model := g.DB().Model(tableReport).Ctx(ctx).Where("member_id", memberID)
	if period != "" {
		model = model.Where("period_key", period)
	}
	var items []MonthlyReport
	err := model.OrderDesc("period_key,version").Scan(&items)
	return items, err
}
func (DatabaseRepository) Publish(ctx context.Context, memberID, id uint64) (*MonthlyReport, error) {
	var item MonthlyReport
	if err := g.DB().Model(tableReport).Ctx(ctx).Where("id", id).Where("member_id", memberID).Scan(&item); err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, errors.New("报告不存在或无权访问")
	}
	if err := ValidateMutation(item.Status); err != nil {
		return nil, err
	}
	now := uint64(time.Now().Unix())
	_, err := g.DB().Model(tableReport).Ctx(ctx).Where("id", id).Data(g.Map{"status": "published", "published_by": memberID, "published_at": now, "update_time": now}).Update()
	item.Status = "published"
	item.PublishedAt = now
	return &item, err
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
