package risk

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const tableRiskEvent = "xy_risk_event"

type Event struct {
	ID             uint64         `json:"id" orm:"id"`
	OpcEntityID    uint64         `json:"opcEntityId" orm:"opc_entity_id"`
	MemberID       uint64         `json:"memberId" orm:"member_id"`
	RiskCode       string         `json:"riskCode" orm:"risk_code"`
	PeriodKey      string         `json:"periodKey" orm:"period_key"`
	Severity       string         `json:"severity" orm:"severity"`
	Status         string         `json:"status" orm:"status"`
	Summary        string         `json:"summary" orm:"summary"`
	Evidence       map[string]any `json:"evidence" orm:"-"`
	ResolutionNote string         `json:"resolutionNote" orm:"resolution_note"`
	FirstHitAt     uint64         `json:"firstHitAt" orm:"first_hit_at"`
	LastHitAt      uint64         `json:"lastHitAt" orm:"last_hit_at"`
	ResolvedAt     uint64         `json:"resolvedAt" orm:"resolved_at"`
}

type eventRow struct {
	Event
	EvidenceJSON string `orm:"evidence_json"`
}

type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) ScanAndSave(ctx context.Context, opcEntityID, memberID uint64, periodKey string, facts Facts) ([]Event, error) {
	if opcEntityID == 0 || memberID == 0 || periodKey == "" {
		return nil, errors.New("企业、会员和期间不能为空")
	}
	now := uint64(time.Now().Unix())
	for _, finding := range Scan(facts) {
		evidence, err := json.Marshal(finding.Evidence)
		if err != nil {
			return nil, err
		}
		_, err = g.DB().Model(tableRiskEvent).Ctx(ctx).Data(g.Map{
			"opc_entity_id": opcEntityID, "member_id": memberID, "risk_code": finding.Code,
			"period_key": periodKey, "severity": finding.Severity, "status": "open",
			"summary": finding.Summary, "evidence_json": string(evidence),
			"first_hit_at": now, "last_hit_at": now, "create_time": now, "update_time": now,
		}).OnDuplicate("severity=VALUES(severity), summary=VALUES(summary), evidence_json=VALUES(evidence_json), last_hit_at=VALUES(last_hit_at), update_time=VALUES(update_time)").Save()
		if err != nil {
			return nil, err
		}
	}
	return s.List(ctx, memberID, periodKey, "")
}

func (s *Service) List(ctx context.Context, memberID uint64, periodKey, status string) ([]Event, error) {
	model := g.DB().Model(tableRiskEvent).Ctx(ctx).Where("member_id", memberID).Where("deleted", 0)
	if periodKey != "" {
		model = model.Where("period_key", periodKey)
	}
	if status != "" {
		model = model.Where("status", status)
	}
	var rows []eventRow
	if err := model.OrderDesc("last_hit_at").Scan(&rows); err != nil {
		return nil, err
	}
	items := make([]Event, 0, len(rows))
	for _, row := range rows {
		item := row.Event
		if row.EvidenceJSON != "" {
			_ = json.Unmarshal([]byte(row.EvidenceJSON), &item.Evidence)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Service) Decide(ctx context.Context, memberID, id uint64, action, note string) (*Event, error) {
	status := map[string]string{"confirm": "confirmed", "dismiss": "dismissed", "resolve": "resolved"}[action]
	if status == "" {
		return nil, errors.New("风险处置操作无效")
	}
	now := uint64(time.Now().Unix())
	resolvedAt := uint64(0)
	if status == "resolved" || status == "dismissed" {
		resolvedAt = now
	}
	result, err := g.DB().Model(tableRiskEvent).Ctx(ctx).Where("id", id).Where("member_id", memberID).Where("deleted", 0).Data(g.Map{
		"status": status, "resolution_note": note, "resolved_at": resolvedAt, "update_time": now,
	}).Update()
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, errors.New("风险事件不存在或无权访问")
	}
	var row eventRow
	if err := g.DB().Model(tableRiskEvent).Ctx(ctx).Where("id", id).Where("member_id", memberID).Scan(&row); err != nil {
		return nil, err
	}
	item := row.Event
	_ = json.Unmarshal([]byte(row.EvidenceJSON), &item.Evidence)
	return &item, nil
}
