package ruleengine

import (
	"context"
	"errors"
	"time"
)

const (
	RuleStatusDraft     = "draft"
	RuleStatusReviewing = "reviewing"
	RuleStatusPublished = "published"
	RuleStatusRetired   = "retired"

	ActionSubmitReview = "submit_review"
	ActionPublish      = "publish"
	ActionRetire       = "retire"
)

type RuleVersion struct {
	ID            uint64         `json:"id"`
	RuleID        uint64         `json:"ruleId"`
	Version       uint           `json:"version"`
	RegionCode    string         `json:"regionCode"`
	Expression    Expression     `json:"expression"`
	Action        map[string]any `json:"action"`
	Status        string         `json:"status"`
	EffectiveFrom uint64         `json:"effectiveFrom"`
	EffectiveTo   uint64         `json:"effectiveTo"`
	CreatedBy     uint64         `json:"createdBy"`
	ReviewedBy    uint64         `json:"reviewedBy"`
	PublishedAt   uint64         `json:"publishedAt"`
}

type TransitionInput struct {
	VersionID  uint64 `json:"versionId"`
	Action     string `json:"action"`
	OperatorID uint64 `json:"operatorId"`
}

type VersionStore interface {
	Get(ctx context.Context, id uint64) (*RuleVersion, error)
	Save(ctx context.Context, item RuleVersion) error
	RetirePublished(ctx context.Context, ruleID, exceptID uint64) error
}

type VersionService struct {
	store VersionStore
}

func NewVersionService(store VersionStore) *VersionService {
	return &VersionService{store: store}
}

func (s *VersionService) Transition(ctx context.Context, in TransitionInput) (*RuleVersion, error) {
	if in.VersionID == 0 || in.OperatorID == 0 {
		return nil, errors.New("规则版本和操作人不能为空")
	}
	version, err := s.store.Get(ctx, in.VersionID)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("规则版本不存在")
	}

	switch in.Action {
	case ActionSubmitReview:
		if version.Status != RuleStatusDraft {
			return nil, errors.New("只有草稿规则可以提交审核")
		}
		if _, err := Evaluate(version.Expression, map[string]any{
			"profile": map[string]any{
				"region": "CN-BJ", "entityType": "one_person_limited_company",
				"taxpayerType": "small_scale", "vatPeriod": "quarterly",
				"employeeCount": float64(0), "invoiceEnabled": false,
				"hasRevenue": false, "hasPublicBankAccount": false, "complexity": "low",
			},
			"period": map[string]any{
				"revenue": float64(0), "expense": float64(0), "invoiceAmount": float64(0),
				"receiptAmount": float64(0), "rollingSales": float64(0),
			},
			"documents": map[string]any{"completeness": float64(0), "missingTypes": []any{}},
		}); err != nil {
			return nil, err
		}
		version.Status = RuleStatusReviewing
	case ActionPublish:
		if version.Status != RuleStatusReviewing {
			return nil, errors.New("只有待审核规则可以发布")
		}
		if version.CreatedBy == in.OperatorID {
			return nil, errors.New("规则创建人不能审核并发布自己的版本")
		}
		if err := s.store.RetirePublished(ctx, version.RuleID, version.ID); err != nil {
			return nil, err
		}
		version.Status = RuleStatusPublished
		version.ReviewedBy = in.OperatorID
		version.PublishedAt = uint64(time.Now().Unix())
	case ActionRetire:
		if version.Status != RuleStatusPublished {
			return nil, errors.New("只有已发布规则可以停用")
		}
		version.Status = RuleStatusRetired
	default:
		return nil, errors.New("规则状态操作无效")
	}
	if err := s.store.Save(ctx, *version); err != nil {
		return nil, err
	}
	return version, nil
}

func (s *VersionService) Simulate(ctx context.Context, versionID uint64, facts map[string]any) (bool, error) {
	version, err := s.store.Get(ctx, versionID)
	if err != nil {
		return false, err
	}
	if version == nil {
		return false, errors.New("规则版本不存在")
	}
	return Evaluate(version.Expression, facts)
}
