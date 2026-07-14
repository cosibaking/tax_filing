package ruleengine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const tableRuleVersion = "xy_compliance_rule_version"

type DatabaseVersionStore struct{}

type ruleVersionRow struct {
	ID            uint64 `orm:"id"`
	RuleID        uint64 `orm:"rule_id"`
	Version       uint   `orm:"version"`
	RegionCode    string `orm:"region_code"`
	ConditionJSON string `orm:"condition_json"`
	ActionJSON    string `orm:"action_json"`
	Status        string `orm:"status"`
	EffectiveFrom uint64 `orm:"effective_from"`
	EffectiveTo   uint64 `orm:"effective_to"`
	CreatedBy     uint64 `orm:"created_by"`
	ReviewedBy    uint64 `orm:"reviewed_by"`
	PublishedAt   uint64 `orm:"published_at"`
}

func (DatabaseVersionStore) Get(ctx context.Context, id uint64) (*RuleVersion, error) {
	var row ruleVersionRow
	if err := g.DB().Model(tableRuleVersion).Ctx(ctx).Where("id", id).Scan(&row); err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	var expression Expression
	if err := json.Unmarshal([]byte(row.ConditionJSON), &expression); err != nil {
		return nil, err
	}
	var action map[string]any
	if err := json.Unmarshal([]byte(row.ActionJSON), &action); err != nil {
		return nil, err
	}
	return &RuleVersion{
		ID: row.ID, RuleID: row.RuleID, Version: row.Version, RegionCode: row.RegionCode,
		Expression: expression, Action: action, Status: row.Status,
		EffectiveFrom: row.EffectiveFrom, EffectiveTo: row.EffectiveTo,
		CreatedBy: row.CreatedBy, ReviewedBy: row.ReviewedBy, PublishedAt: row.PublishedAt,
	}, nil
}

func (DatabaseVersionStore) Save(ctx context.Context, item RuleVersion) error {
	_, err := g.DB().Model(tableRuleVersion).Ctx(ctx).Where("id", item.ID).Data(g.Map{
		"status": item.Status, "reviewed_by": item.ReviewedBy,
		"published_at": item.PublishedAt, "update_time": uint64(time.Now().Unix()),
	}).Update()
	return err
}

func (DatabaseVersionStore) RetirePublished(ctx context.Context, ruleID, exceptID uint64) error {
	_, err := g.DB().Model(tableRuleVersion).Ctx(ctx).
		Where("rule_id", ruleID).Where("status", RuleStatusPublished).WhereNot("id", exceptID).
		Data(g.Map{"status": RuleStatusRetired, "update_time": uint64(time.Now().Unix())}).Update()
	return err
}

func NewDatabaseVersionService() *VersionService {
	return NewVersionService(DatabaseVersionStore{})
}
