package statement

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"

	"xygo/internal/library/statementpdf"
	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
	"xygo/utility"
)

const statementPdfDir = "resource/public/attachment/upload/statements"

// GetStatementPdfUrl 会员获取对账单 PDF 地址（不存在则按需生成）
func (s *sComplianceStatement) GetStatementPdfUrl(ctx context.Context, in *compliancein.StatementPdfInp) (*compliancein.StatementPdfModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	row, err := loadStatementRow(ctx, in.StatementId)
	if err != nil {
		return nil, err
	}
	if row.OpcId != opc.Id {
		return nil, gerror.New("无权访问该对账单")
	}
	url, err := s.ensureStatementPdf(ctx, row)
	if err != nil {
		return nil, err
	}
	return &compliancein.StatementPdfModel{Url: url}, nil
}

// AdminStatementPdfUrl 管理端获取对账单 PDF 地址
func (s *sComplianceStatement) AdminStatementPdfUrl(ctx context.Context, statementId uint64) (*compliancein.StatementPdfModel, error) {
	row, err := loadStatementRow(ctx, statementId)
	if err != nil {
		return nil, err
	}
	url, err := s.ensureStatementPdf(ctx, row)
	if err != nil {
		return nil, err
	}
	return &compliancein.StatementPdfModel{Url: url}, nil
}

type statementRow struct {
	Id          uint64
	OpcId       uint64
	Year        int
	Month       int
	Summary     string
	SentAt      uint64
	CompanyName string
}

func loadStatementRow(ctx context.Context, statementId uint64) (*statementRow, error) {
	var row statementRow
	err := g.DB().Model(shared.TableMonthlyStmt+" s").Ctx(ctx).
		LeftJoin(shared.TableOpcEntity+" o", "o.id = s.opc_id AND o.deleted = 0").
		Fields("s.id", "s.opc_id", "s.year", "s.month", "s.summary", "s.sent_at", "o.company_name").
		Where("s.id", statementId).
		Where("s.deleted", 0).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "查询对账单失败")
	}
	if row.Id == 0 {
		return nil, gerror.New("对账单不存在")
	}
	return &row, nil
}

func (s *sComplianceStatement) ensureStatementPdf(ctx context.Context, row *statementRow) (string, error) {
	relPath := statementPdfRelPath(row.Id)
	physical := statementPdfPhysical(row.Id)
	if gfile.Exists(physical) {
		return relPath, nil
	}

	data, err := buildPdfData(ctx, row)
	if err != nil {
		return "", err
	}
	pdfBytes, err := statementpdf.Generate(*data)
	if err != nil {
		return "", gerror.Wrap(err, "生成 PDF 失败")
	}
	if err = gfile.Mkdir(statementPdfDir); err != nil {
		return "", gerror.Wrap(err, "创建 PDF 目录失败")
	}
	if err = gfile.PutBytes(physical, pdfBytes); err != nil {
		return "", gerror.Wrap(err, "保存 PDF 失败")
	}
	return relPath, nil
}

func buildPdfData(ctx context.Context, row *statementRow) (*statementpdf.Data, error) {
	period := fmt.Sprintf("%04d-%02d", row.Year, row.Month)
	revenue, cost, profit, prepaidTax, cumulative, filingStatus := parseMemberSummary(row.Summary)
	incomeEntries, expenseEntries, pendingFilings, generatedAt := parseAdminSummaryExtras(row.Summary)

	data := &statementpdf.Data{
		CompanyName:      row.CompanyName,
		Period:           period,
		Revenue:          revenue,
		Cost:             cost,
		Profit:           profit,
		PrepaidTax:       prepaidTax,
		CumulativeProfit: cumulative,
		FilingStatus:     filingStatus,
		IncomeEntries:    incomeEntries,
		ExpenseEntries:   expenseEntries,
		PendingFilings:   pendingFilings,
		GeneratedAt:      generatedAt,
	}
	if data.GeneratedAt == "" {
		data.GeneratedAt = utility.UnixToGTime(utility.NowUnix()).Format("Y-m-d H:i:s")
	}
	if row.SentAt > 0 {
		data.SentAt = utility.UnixToGTime(int64(row.SentAt)).Format("Y-m-d H:i:s")
	}
	return data, nil
}

func parseAdminSummaryExtras(raw string) (incomeEntries, expenseEntries, pendingFilings int, generatedAt string) {
	if raw == "" {
		return
	}
	var summary struct {
		IncomeEntries  int    `json:"incomeEntries"`
		ExpenseEntries int    `json:"expenseEntries"`
		PendingFilings int    `json:"pendingFilings"`
		GeneratedAt    string `json:"generatedAt"`
	}
	if err := gjson.DecodeTo(raw, &summary); err != nil {
		return
	}
	return summary.IncomeEntries, summary.ExpenseEntries, summary.PendingFilings, summary.GeneratedAt
}

func statementPdfRelPath(statementId uint64) string {
	return "/attachment/upload/statements/" + strconv.FormatUint(statementId, 10) + ".pdf"
}

func statementPdfPhysical(statementId uint64) string {
	return gfile.Join(statementPdfDir, strconv.FormatUint(statementId, 10)+".pdf")
}

func statementPdfUrlIfExists(statementId uint64) string {
	if gfile.Exists(statementPdfPhysical(statementId)) {
		return statementPdfRelPath(statementId)
	}
	return ""
}
