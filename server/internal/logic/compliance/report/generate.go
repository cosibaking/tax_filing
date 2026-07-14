package report

import (
	"context"
	"errors"

	"xygo/internal/library/reportpdf"
	"xygo/internal/logic/compliance/shared"
)

type Input struct {
	PeriodKey    string           `json:"periodKey"`
	Statistics   map[string]any   `json:"statistics"`
	Completeness map[string]any   `json:"completeness"`
	Risks        []map[string]any `json:"risks"`
	RuleVersions []uint64         `json:"ruleVersions,omitempty"`
}

type Result struct {
	Content          string           `json:"content"`
	AIModel          string           `json:"aiModel"`
	KnowledgeVersion string           `json:"knowledgeVersion"`
	Structured       StructuredReport `json:"structuredReport"`
}

func Generate(in Input) (Result, error) {
	structured, err := BuildStructured(in)
	if err != nil {
		return Result{}, err
	}
	return Result{Content: RenderContent(in.PeriodKey, structured), Structured: structured}, nil
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
	Get(context.Context, uint64, uint64) (*MonthlyReport, error)
	List(context.Context, uint64, string) ([]MonthlyReport, error)
	Publish(context.Context, uint64, uint64) (*MonthlyReport, error)
}
type Service struct {
	repository     Repository
	narrator       Narrator
	companyLoader  companyLoader
	pdfGenerator   func(reportpdf.Data) ([]byte, error)
	snapshotLoader SnapshotLoader
}

func NewService(repository Repository, narrator Narrator) *Service {
	return newServiceWithCompanyLoader(repository, narrator, shared.LoadOpcByMember)
}

type companyLoader func(context.Context, uint64) (*shared.OpcBrief, error)

func newServiceWithCompanyLoader(repository Repository, narrator Narrator, loader companyLoader) *Service {
	return &Service{repository: repository, narrator: narrator, companyLoader: loader, pdfGenerator: reportpdf.Generate, snapshotLoader: DatabaseSnapshotLoader{}}
}

func newServiceWithSnapshotLoader(repository Repository, narrator Narrator, loader SnapshotLoader) *Service {
	s := NewService(repository, narrator)
	s.snapshotLoader = loader
	return s
}
func (s *Service) Build(ctx context.Context, in Input) (Result, error) {
	base, err := Generate(in)
	if err != nil {
		return Result{}, err
	}
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
	ID               uint64            `json:"id" orm:"id"`
	OpcEntityID      uint64            `json:"opcEntityId" orm:"opc_entity_id"`
	MemberID         uint64            `json:"memberId" orm:"member_id"`
	PeriodKey        string            `json:"periodKey" orm:"period_key"`
	Version          uint              `json:"version" orm:"version"`
	Status           string            `json:"status" orm:"status"`
	Content          string            `json:"content" orm:"content"`
	AIModel          string            `json:"aiModel" orm:"ai_model"`
	KnowledgeVersion string            `json:"knowledgeVersion" orm:"knowledge_version"`
	StructuredJSON   string            `json:"-" orm:"structured_json"`
	StructuredReport *StructuredReport `json:"structuredReport,omitempty" orm:"-"`
	Legacy           bool              `json:"legacy" orm:"-"`
	PublishedAt      uint64            `json:"publishedAt" orm:"published_at"`
	CreatedAt        uint64            `json:"createdAt" orm:"create_time"`
}
