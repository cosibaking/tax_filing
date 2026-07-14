package report

import (
	"context"
	"errors"
	"fmt"
)

type Input struct {
	PeriodKey    string           `json:"periodKey"`
	Statistics   map[string]any   `json:"statistics"`
	Completeness map[string]any   `json:"completeness"`
	Risks        []map[string]any `json:"risks"`
}

type Result struct {
	Content          string `json:"content"`
	AIModel          string `json:"aiModel"`
	KnowledgeVersion string `json:"knowledgeVersion"`
}

func Generate(in Input) Result {
	content := fmt.Sprintf("%s 月度经营合规体检\n", in.PeriodKey)
	if len(in.Statistics) == 0 {
		content += "经营统计：数据不足，请补充银行流水、发票或合同。\n"
	} else {
		content += "经营统计：已根据确认数据生成。\n"
	}
	if len(in.Completeness) == 0 {
		content += "资料完整度：数据不足。\n"
	} else {
		content += "资料完整度：请按缺失清单补充。\n"
	}
	content += fmt.Sprintf("风险提示：%d 项。重要结论需人工复核。", len(in.Risks))
	return Result{Content: content}
}

func ValidateMutation(status string) error {
	if status == "published" {
		return errors.New("已发布报告不可修改，请创建新版本")
	}
	return nil
}

type Narrator interface {
	Narrate(context.Context, Input) (content, model string, err error)
}
type Repository interface {
	SaveDraft(context.Context, uint64, uint64, Input, Result) (*MonthlyReport, error)
	List(context.Context, uint64, string) ([]MonthlyReport, error)
	Publish(context.Context, uint64, uint64) (*MonthlyReport, error)
}
type Service struct {
	repository Repository
	narrator   Narrator
}

func NewService(repository Repository, narrator Narrator) *Service {
	return &Service{repository: repository, narrator: narrator}
}
func (s *Service) Build(ctx context.Context, in Input) (Result, error) {
	base := Generate(in)
	if s.narrator == nil {
		return base, nil
	}
	content, model, err := s.narrator.Narrate(ctx, in)
	if err != nil || content == "" {
		return base, nil
	}
	base.Content, base.AIModel = content, model
	return base, nil
}

type MonthlyReport struct {
	ID               uint64 `json:"id" orm:"id"`
	OpcEntityID      uint64 `json:"opcEntityId" orm:"opc_entity_id"`
	MemberID         uint64 `json:"memberId" orm:"member_id"`
	PeriodKey        string `json:"periodKey" orm:"period_key"`
	Version          uint   `json:"version" orm:"version"`
	Status           string `json:"status" orm:"status"`
	Content          string `json:"content" orm:"content"`
	AIModel          string `json:"aiModel" orm:"ai_model"`
	KnowledgeVersion string `json:"knowledgeVersion" orm:"knowledge_version"`
	PublishedAt      uint64 `json:"publishedAt" orm:"published_at"`
}
