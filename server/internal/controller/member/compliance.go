// +----------------------------------------------------------------------

// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]

// +----------------------------------------------------------------------

// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.

// +----------------------------------------------------------------------

// | Licensed ( https://opensource.org/licenses/MIT )

// +----------------------------------------------------------------------

// | Author: 喜羊羊 <751300685@qq.com>

// +----------------------------------------------------------------------



package member



import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/v2/net/ghttp"



	"xygo/api/member"

	"xygo/internal/consts"

	"xygo/internal/library/contexts"

	"xygo/internal/model/input/compliancein"

	"xygo/internal/service"

)



func requireMemberId(ctx context.Context) (uint64, error) {

	memberId := contexts.GetMemberId(ctx)

	if memberId == 0 {

		return 0, gerror.NewCode(consts.CodeNotAuthorized, "请先登录")

	}

	return memberId, nil

}



func requestMeta(ctx context.Context) (ip, ua string) {

	r := ghttp.RequestFromCtx(ctx)

	if r == nil {

		return "", ""

	}

	return r.GetClientIp(), r.Header.Get("User-Agent")

}



// ComplianceDiagnosisHistory 会员诊断历史

func (c *ControllerV1) ComplianceDiagnosisHistory(ctx context.Context, req *member.ComplianceDiagnosisHistoryReq) (res *member.ComplianceDiagnosisHistoryRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}



	out, err := service.ComplianceDiagnosis().ListHistory(ctx, &compliancein.DiagnosisHistoryInp{

		MemberId: memberId,

		Page:     req.Page,

		PageSize: req.PageSize,

	})

	if err != nil {

		return nil, err

	}



	return &member.ComplianceDiagnosisHistoryRes{

		List:     out.List,

		Page:     out.Page,

		PageSize: out.PageSize,

		Total:    out.Total,

	}, nil

}



// ComplianceDiagnosisSync 登录后同步访客诊断缓存

func (c *ControllerV1) ComplianceDiagnosisSync(ctx context.Context, req *member.ComplianceDiagnosisSyncReq) (res *member.ComplianceDiagnosisSyncRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}



	items := make([]compliancein.DiagnosisSubmitInp, 0, len(req.Items))

	for _, item := range req.Items {

		items = append(items, compliancein.DiagnosisSubmitInp{

			Platforms:          item.Platforms,

			MonthlyIncomeRange: item.MonthlyIncomeRange,

			AnnualCostEstimate: item.AnnualCostEstimate,

			ExistingEntity:     item.ExistingEntity,

			HasFiledTax:        item.HasFiledTax,

			TaxBureauContact:   item.TaxBureauContact,

			Notes:              item.Notes,

			CostBreakdown:      item.CostBreakdown,

		})

	}



	out, err := service.ComplianceDiagnosis().SyncGuest(ctx, &compliancein.DiagnosisSyncInp{

		MemberId: memberId,

		Items:    items,

	})

	if err != nil {

		return nil, err

	}



	return &member.ComplianceDiagnosisSyncRes{DiagnosisSyncModel: out}, nil

}



// ComplianceDiagnosisBind 绑定匿名诊断记录到当前会员

func (c *ControllerV1) ComplianceDiagnosisBind(ctx context.Context, req *member.ComplianceDiagnosisBindReq) (res *member.ComplianceDiagnosisBindRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}



	out, err := service.ComplianceDiagnosis().BindToMember(ctx, &compliancein.DiagnosisBindInp{

		MemberId:     memberId,

		DiagnosisIds: req.DiagnosisIds,

	})

	if err != nil {

		return nil, err

	}



	return &member.ComplianceDiagnosisBindRes{DiagnosisBindModel: out}, nil

}



// ComplianceDiagnosisDetail 诊断详情

func (c *ControllerV1) ComplianceDiagnosisDetail(ctx context.Context, req *member.ComplianceDiagnosisDetailReq) (res *member.ComplianceDiagnosisDetailRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}



	out, err := service.ComplianceDiagnosis().GetDetail(ctx, &compliancein.DiagnosisDetailInp{

		MemberId:    memberId,

		DiagnosisId: req.DiagnosisId,

	})

	if err != nil {

		return nil, err

	}



	return &member.ComplianceDiagnosisDetailRes{DiagnosisDetailModel: out}, nil

}



// ComplianceOrderCreate 创建服务订单

func (c *ControllerV1) ComplianceOrderCreate(ctx context.Context, req *member.ComplianceOrderCreateReq) (res *member.ComplianceOrderCreateRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}



	out, err := service.ComplianceOrder().CreateOrder(ctx, &compliancein.OrderCreateInp{

		MemberId:    memberId,

		PlanId:      req.PlanId,

		DiagnosisId: req.DiagnosisId,

	})

	if err != nil {

		return nil, err

	}

	return &member.ComplianceOrderCreateRes{OrderCreateModel: out}, nil

}



// ComplianceConsent 风险告知 / 方案确认

func (c *ControllerV1) ComplianceConsent(ctx context.Context, req *member.ComplianceConsentReq) (res *member.ComplianceConsentRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}

	ip, ua := requestMeta(ctx)



	out, err := service.ComplianceOrder().RecordConsent(ctx, &compliancein.ConsentInp{

		MemberId:        memberId,

		OrderId:         req.OrderId,

		Type:            req.Type,

		DocumentVersion: req.DocumentVersion,

		Acknowledgments: req.Acknowledgments,

		PlanConfirmed:   req.PlanConfirmed,

		Ip:              ip,

		UserAgent:       ua,

	})

	if err != nil {

		return nil, err

	}

	return &member.ComplianceConsentRes{ConsentModel: out}, nil

}



// ComplianceSign 电子签约

func (c *ControllerV1) ComplianceSign(ctx context.Context, req *member.ComplianceSignReq) (res *member.ComplianceSignRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}

	ip, ua := requestMeta(ctx)



	out, err := service.ComplianceOrder().SignContract(ctx, &compliancein.SignInp{

		MemberId:  memberId,

		OrderId:   req.OrderId,

		LegalName: req.LegalName,

		Ip:        ip,

		UserAgent: ua,

	})

	if err != nil {

		return nil, err

	}

	return &member.ComplianceSignRes{SignModel: out}, nil

}



// ComplianceOpcSummary OPC 主体摘要

func (c *ControllerV1) ComplianceOpcSummary(ctx context.Context, req *member.ComplianceOpcSummaryReq) (res *member.ComplianceOpcSummaryRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}



	out, err := service.ComplianceOpc().GetSummary(ctx, memberId)

	if err != nil {

		return nil, err

	}

	return &member.ComplianceOpcSummaryRes{OpcSummaryModel: out}, nil

}



// ComplianceOpcProgress OPC 进度时间轴

func (c *ControllerV1) ComplianceOpcProgress(ctx context.Context, req *member.ComplianceOpcProgressReq) (res *member.ComplianceOpcProgressRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}



	out, err := service.ComplianceOpc().GetProgress(ctx, memberId)

	if err != nil {

		return nil, err

	}

	return &member.ComplianceOpcProgressRes{OpcProgressModel: out}, nil

}



// ComplianceOpcMaterials 提交 OPC 注册资料

func (c *ControllerV1) ComplianceOpcMaterials(ctx context.Context, req *member.ComplianceOpcMaterialsReq) (res *member.ComplianceOpcMaterialsRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}

	ip, _ := requestMeta(ctx)



	out, err := service.ComplianceOpc().SubmitMaterials(ctx, &compliancein.MaterialsSubmitInp{

		MemberId:           memberId,

		ProposedNames:      req.ProposedNames,

		RegisteredCapital:  req.RegisteredCapital,

		CapitalTermYears:   req.CapitalTermYears,

		BusinessTermType:   req.BusinessTermType,

		BusinessTermEnd:    req.BusinessTermEnd,

		BusinessScope:      req.BusinessScope,

		RegisterProvince:   req.RegisterProvince,

		RegisterCity:       req.RegisterCity,

		RegisterDistrict:   req.RegisterDistrict,

		RegisterAddress:    req.RegisterAddress,

		AddressProofFileId: req.AddressProofFileId,

		LegalPersonName:    req.LegalPersonName,

		IdCard:             req.IdCard,

		IdCardValidFrom:    req.IdCardValidFrom,

		IdCardValidTo:      req.IdCardValidTo,

		Ethnicity:          req.Ethnicity,

		HouseholdAddress:   req.HouseholdAddress,

		ResidentialAddress: req.ResidentialAddress,

		Phone:              req.Phone,

		Email:              req.Email,

		IdCardFrontFileId:  req.IdCardFrontFileId,

		IdCardBackFileId:   req.IdCardBackFileId,

		EsignAuthorized:    req.EsignAuthorized,

		Confirmations:      req.Confirmations,

		Ip:                 ip,

	})

	if err != nil {

		return nil, err

	}

	return &member.ComplianceOpcMaterialsRes{MaterialsSubmitModel: out}, nil

}



// ComplianceOpcBankReceipt 上传开户回执

func (c *ControllerV1) ComplianceOpcBankReceipt(ctx context.Context, req *member.ComplianceOpcBankReceiptReq) (res *member.ComplianceOpcBankReceiptRes, err error) {

	memberId, err := requireMemberId(ctx)

	if err != nil {

		return nil, err

	}

	ip, _ := requestMeta(ctx)



	out, err := service.ComplianceOpc().SubmitBankReceipt(ctx, &compliancein.BankReceiptInp{

		MemberId:          memberId,

		BankName:          req.BankName,

		BankReceiptFileId: req.BankReceiptFileId,

		Ip:                ip,

	})

	if err != nil {

		return nil, err

	}

	return &member.ComplianceOpcBankReceiptRes{BankReceiptModel: out}, nil
}

// ComplianceOpcMaterialsOverview 已提交资料概览
func (c *ControllerV1) ComplianceOpcMaterialsOverview(ctx context.Context, req *member.ComplianceOpcMaterialsOverviewReq) (res *member.ComplianceOpcMaterialsOverviewRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceOpc().GetMaterialsOverview(ctx, memberId)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceOpcMaterialsOverviewRes{MaterialsOverviewModel: out}, nil
}

// ComplianceOpcMaterialsDetail 单个 OPC 实体完整资料（脱敏）
func (c *ControllerV1) ComplianceOpcMaterialsDetail(ctx context.Context, req *member.ComplianceOpcMaterialsDetailReq) (res *member.ComplianceOpcMaterialsDetailRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceOpc().GetMaterialsDetail(ctx, memberId, req.OpcId)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceOpcMaterialsDetailRes{MaterialsEntityDetailModel: out}, nil
}

// ComplianceOpcMaterialsSection 资料分组详情（脱敏）
func (c *ControllerV1) ComplianceOpcMaterialsSection(ctx context.Context, req *member.ComplianceOpcMaterialsSectionReq) (res *member.ComplianceOpcMaterialsSectionRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceOpc().GetMaterialsSection(ctx, memberId, req.OpcId, req.Section)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceOpcMaterialsSectionRes{MaterialsSectionDetailModel: out}, nil
}

// ComplianceOpcMaterialsReveal 密码验证查看完整资料
func (c *ControllerV1) ComplianceOpcMaterialsReveal(ctx context.Context, req *member.ComplianceOpcMaterialsRevealReq) (res *member.ComplianceOpcMaterialsRevealRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceOpc().RevealMaterialsSection(ctx, &compliancein.MaterialsRevealInp{
		MemberId: memberId,
		OpcId:    req.OpcId,
		Section:  req.Section,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceOpcMaterialsRevealRes{MaterialsRevealModel: out}, nil
}

// ComplianceOpcMaterialsRevealAll 密码验证查看全部完整资料
func (c *ControllerV1) ComplianceOpcMaterialsRevealAll(ctx context.Context, req *member.ComplianceOpcMaterialsRevealAllReq) (res *member.ComplianceOpcMaterialsRevealAllRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceOpc().RevealMaterialsAll(ctx, &compliancein.MaterialsRevealInp{
		MemberId: memberId,
		OpcId:    req.OpcId,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceOpcMaterialsRevealAllRes{MaterialsRevealAllModel: out}, nil
}

func readUploadCSV(ctx context.Context, field string, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return fallback
	}
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return ""
	}
	up := r.GetUploadFile(field)
	if up == nil {
		return strings.TrimSpace(r.Get("csvContent").String())
	}
	f, err := up.Open()
	if err != nil {
		return ""
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(b)
}

func readUploadBytes(ctx context.Context, field string) []byte {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil
	}
	up := r.GetUploadFile(field)
	if up == nil {
		return nil
	}
	f, err := up.Open()
	if err != nil {
		return nil
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil
	}
	return b
}

// ComplianceIncomeList 收入台账列表
func (c *ControllerV1) ComplianceIncomeList(ctx context.Context, req *member.ComplianceIncomeListReq) (res *member.ComplianceIncomeListRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceLedger().ListIncome(ctx, &compliancein.IncomeListInp{
		MemberId: memberId,
		Page:     req.Page,
		PageSize: req.PageSize,
		Month:    req.Month,
		Platform: req.Platform,
		Category: req.Category,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceIncomeListRes{IncomeListModel: out}, nil
}

// ComplianceIncomeCreate 创建收入
func (c *ControllerV1) ComplianceIncomeCreate(ctx context.Context, req *member.ComplianceIncomeCreateReq) (res *member.ComplianceIncomeCreateRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	out, err := service.ComplianceLedger().CreateIncome(ctx, &compliancein.IncomeCreateInp{
		MemberId:    memberId,
		Platform:    req.Platform,
		Category:    req.Category,
		GrossAmount: req.GrossAmount,
		PlatformFee: req.PlatformFee,
		OccurredAt:  req.OccurredAt,
		Remark:      req.Remark,
		Ip:          ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceIncomeCreateRes{IncomeCreateModel: out}, nil
}

// ComplianceIncomeImport CSV 导入收入
func (c *ControllerV1) ComplianceIncomeImport(ctx context.Context, req *member.ComplianceIncomeImportReq) (res *member.ComplianceIncomeImportRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	csvContent := readUploadCSV(ctx, "file", req.CsvContent)
	out, err := service.ComplianceLedger().ImportIncomeCSV(ctx, &compliancein.IncomeImportInp{
		MemberId:   memberId,
		CsvContent: csvContent,
		Ip:         ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceIncomeImportRes{IncomeImportModel: out}, nil
}

// ComplianceExpenseList 费用台账列表
func (c *ControllerV1) ComplianceExpenseList(ctx context.Context, req *member.ComplianceExpenseListReq) (res *member.ComplianceExpenseListRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceLedger().ListExpense(ctx, &compliancein.ExpenseListInp{
		MemberId: memberId,
		Page:     req.Page,
		PageSize: req.PageSize,
		Month:    req.Month,
		Category: req.Category,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceExpenseListRes{ExpenseListModel: out}, nil
}

// ComplianceExpenseCreate 创建费用
func (c *ControllerV1) ComplianceExpenseCreate(ctx context.Context, req *member.ComplianceExpenseCreateReq) (res *member.ComplianceExpenseCreateRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	out, err := service.ComplianceLedger().CreateExpense(ctx, &compliancein.ExpenseCreateInp{
		MemberId:       memberId,
		Category:       req.Category,
		Amount:         req.Amount,
		InvoiceType:    req.InvoiceType,
		AttachmentId:   req.AttachmentId,
		Description:    req.Description,
		OccurredAt:     req.OccurredAt,
		ConfirmWarning: req.ConfirmWarning,
		Ip:             ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceExpenseCreateRes{ExpenseCreateModel: out}, nil
}

// ComplianceBankImport 银行流水导入
func (c *ControllerV1) ComplianceBankImport(ctx context.Context, req *member.ComplianceBankImportReq) (res *member.ComplianceBankImportRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	csvContent := readUploadCSV(ctx, "file", req.CsvContent)
	out, err := service.ComplianceLedger().ImportBankStatement(ctx, &compliancein.BankImportInp{
		MemberId:   memberId,
		CsvContent: csvContent,
		Ip:         ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceBankImportRes{BankImportModel: out}, nil
}

// ComplianceExpenseTypes 费用类型库
func (c *ControllerV1) ComplianceExpenseTypes(ctx context.Context, req *member.ComplianceExpenseTypesReq) (res *member.ComplianceExpenseTypesRes, err error) {
	if _, err = requireMemberId(ctx); err != nil {
		return nil, err
	}
	out, err := service.ComplianceLedger().ListExpenseTypes(ctx)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceExpenseTypesRes{ExpenseTypesModel: out}, nil
}

// ComplianceLedgerProfit 利润表
func (c *ControllerV1) ComplianceLedgerProfit(ctx context.Context, req *member.ComplianceLedgerProfitReq) (res *member.ComplianceLedgerProfitRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	profitIn := buildProfitQuery(req)
	profitIn.MemberId = memberId
	out, err := service.ComplianceLedger().GetProfit(ctx, profitIn)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceLedgerProfitRes{ProfitSummaryModel: out}, nil
}

// ComplianceTaxCalendar 申报日历
func (c *ControllerV1) ComplianceTaxCalendar(ctx context.Context, req *member.ComplianceTaxCalendarReq) (res *member.ComplianceTaxCalendarRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceTax().GetCalendar(ctx, &compliancein.TaxCalendarInp{
		MemberId: memberId,
		Year:     req.Year,
		Month:    req.Month,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceTaxCalendarRes{TaxCalendarModel: out}, nil
}

// ComplianceTaxChecklist 报税前自查清单
func (c *ControllerV1) ComplianceTaxChecklist(ctx context.Context, req *member.ComplianceTaxChecklistReq) (res *member.ComplianceTaxChecklistRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceTax().GetChecklist(ctx, &compliancein.TaxChecklistInp{
		MemberId: memberId,
		TaskId:   req.TaskId,
		Period:   req.Period,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceTaxChecklistRes{TaxChecklistModel: out}, nil
}

// ComplianceOrderActive 当前活跃订单
func (c *ControllerV1) ComplianceOrderActive(ctx context.Context, req *member.ComplianceOrderActiveReq) (res *member.ComplianceOrderActiveRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceOrder().GetActiveOrder(ctx, memberId)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceOrderActiveRes{ActiveOrderModel: out}, nil
}

// ComplianceIncomeDelete 删除收入
func (c *ControllerV1) ComplianceIncomeDelete(ctx context.Context, req *member.ComplianceIncomeDeleteReq) (res *member.ComplianceIncomeDeleteRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	err = service.ComplianceLedger().DeleteIncome(ctx, &compliancein.IncomeDeleteInp{
		MemberId: memberId,
		Id:       req.Id,
		Ip:       ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceIncomeDeleteRes{}, nil
}

// ComplianceIncomeImportPreview CSV 导入预览
func (c *ControllerV1) ComplianceIncomeImportPreview(ctx context.Context, req *member.ComplianceIncomeImportPreviewReq) (res *member.ComplianceIncomeImportPreviewRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	csvContent := readUploadCSV(ctx, "file", req.CsvContent)
	out, err := service.ComplianceLedger().PreviewIncomeImport(ctx, &compliancein.IncomeImportInp{
		MemberId:   memberId,
		CsvContent: csvContent,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceIncomeImportPreviewRes{IncomeImportPreviewModel: out}, nil
}

// ComplianceExpenseDelete 删除费用
func (c *ControllerV1) ComplianceExpenseDelete(ctx context.Context, req *member.ComplianceExpenseDeleteReq) (res *member.ComplianceExpenseDeleteRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	err = service.ComplianceLedger().DeleteExpense(ctx, &compliancein.ExpenseDeleteInp{
		MemberId: memberId,
		Id:       req.Id,
		Ip:       ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceExpenseDeleteRes{}, nil
}

// ComplianceExpenseCategories 费用类型库（前端路径别名）
func (c *ControllerV1) ComplianceExpenseCategories(ctx context.Context, req *member.ComplianceExpenseCategoriesReq) (res *member.ComplianceExpenseCategoriesRes, err error) {
	if _, err = requireMemberId(ctx); err != nil {
		return nil, err
	}
	out, err := service.ComplianceLedger().ListExpenseTypes(ctx)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceExpenseCategoriesRes{List: out.List}, nil
}

// ComplianceBankUnmatched 未匹配银行流水
func (c *ControllerV1) ComplianceBankUnmatched(ctx context.Context, req *member.ComplianceBankUnmatchedReq) (res *member.ComplianceBankUnmatchedRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceLedger().ListBankUnmatched(ctx, &compliancein.BankUnmatchedInp{
		MemberId: memberId,
		Month:    req.Month,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceBankUnmatchedRes{BankUnmatchedModel: out}, nil
}

// ComplianceLedgerVouchers 会计分录列表
func (c *ControllerV1) ComplianceLedgerVouchers(ctx context.Context, req *member.ComplianceLedgerVouchersReq) (res *member.ComplianceLedgerVouchersRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceLedger().ListLedgerVouchers(ctx, &compliancein.LedgerVoucherListInp{
		MemberId: memberId,
		Period:   req.Period,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceLedgerVouchersRes{LedgerVoucherListModel: out}, nil
}

// ComplianceTaxTaskDetail 申报任务详情
func (c *ControllerV1) ComplianceTaxTaskDetail(ctx context.Context, req *member.ComplianceTaxTaskDetailReq) (res *member.ComplianceTaxTaskDetailRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceTax().GetTaskDetail(ctx, &compliancein.TaxTaskDetailInp{
		MemberId: memberId,
		TaskId:   req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceTaxTaskDetailRes{TaxTaskDetailModel: out}, nil
}

// ComplianceMemberStatements 会员对账单列表
func (c *ControllerV1) ComplianceMemberStatements(ctx context.Context, req *member.ComplianceMemberStatementsReq) (res *member.ComplianceMemberStatementsRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceStatement().ListMemberStatements(ctx, &compliancein.MemberStatementListInp{
		MemberId: memberId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceMemberStatementsRes{MemberStatementListModel: out}, nil
}

// ComplianceStatementPdf 对账单 PDF（MVP 以 JSON 摘要为准，无 PDF 文件）
func (c *ControllerV1) ComplianceStatementPdf(ctx context.Context, req *member.ComplianceStatementPdfReq) (res *member.ComplianceStatementPdfRes, err error) {
	if _, err = requireMemberId(ctx); err != nil {
		return nil, err
	}
	return &member.ComplianceStatementPdfRes{Url: ""}, nil
}

// ComplianceEmploymentGet 获取用工状态
func (c *ControllerV1) ComplianceEmploymentGet(ctx context.Context, req *member.ComplianceEmploymentGetReq) (res *member.ComplianceEmploymentGetRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceSocial().GetEmployment(ctx, memberId)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceEmploymentGetRes{EmploymentStatusModel: out}, nil
}

// ComplianceEmploymentSet 设置用工状态
func (c *ControllerV1) ComplianceEmploymentSet(ctx context.Context, req *member.ComplianceEmploymentSetReq) (res *member.ComplianceEmploymentSetRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	out, err := service.ComplianceSocial().SetEmployment(ctx, &compliancein.EmploymentStatusInp{
		MemberId:         memberId,
		EmploymentStatus: req.EmploymentStatus,
		Ip:               ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceEmploymentSetRes{EmploymentStatusModel: out}, nil
}

// ComplianceSocialGuides 社保指引列表
func (c *ControllerV1) ComplianceSocialGuides(ctx context.Context, req *member.ComplianceSocialGuidesReq) (res *member.ComplianceSocialGuidesRes, err error) {
	if _, err = requireMemberId(ctx); err != nil {
		return nil, err
	}
	out, err := service.ComplianceSocial().ListGuides(ctx)
	if err != nil {
		return nil, err
	}
	return &member.ComplianceSocialGuidesRes{SocialGuideListModel: out}, nil
}

// ComplianceSocialGuideDetail 社保指引详情
func (c *ControllerV1) ComplianceSocialGuideDetail(ctx context.Context, req *member.ComplianceSocialGuideDetailReq) (res *member.ComplianceSocialGuideDetailRes, err error) {
	if _, err = requireMemberId(ctx); err != nil {
		return nil, err
	}
	out, err := service.ComplianceSocial().GetGuide(ctx, &compliancein.SocialGuideDetailInp{Slug: req.Slug})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceSocialGuideDetailRes{SocialGuideDetailModel: out}, nil
}

// ComplianceSocialConsultCreate 提交社保咨询
func (c *ControllerV1) ComplianceSocialConsultCreate(ctx context.Context, req *member.ComplianceSocialConsultCreateReq) (res *member.ComplianceSocialConsultCreateRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	out, err := service.ComplianceSocial().CreateConsult(ctx, &compliancein.SocialConsultCreateInp{
		MemberId:   memberId,
		Category:   req.Category,
		Question:   req.Question,
		RegionCode: req.RegionCode,
		Ip:         ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceSocialConsultCreateRes{SocialConsultCreateModel: out}, nil
}

// ComplianceSocialConsultList 社保咨询列表
func (c *ControllerV1) ComplianceSocialConsultList(ctx context.Context, req *member.ComplianceSocialConsultListReq) (res *member.ComplianceSocialConsultListRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceSocial().ListMemberConsults(ctx, &compliancein.SocialConsultListInp{
		MemberId: memberId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceSocialConsultListRes{SocialConsultListModel: out}, nil
}

func buildProfitQuery(req *member.ComplianceLedgerProfitReq) *compliancein.ProfitQueryInp {
	in := &compliancein.ProfitQueryInp{}
	switch strings.ToLower(strings.TrimSpace(req.PeriodType)) {
	case "quarter":
		in.Period = "quarterly"
	case "year":
		in.Period = "yearly"
	default:
		in.Period = "monthly"
	}
	if req.Year > 0 {
		in.Year = req.Year
	}
	if req.Month > 0 {
		in.Month = req.Month
	}
	if req.Quarter > 0 {
		in.Quarter = req.Quarter
	}
	period := strings.TrimSpace(req.Period)
	if period == "" {
		return in
	}
	parts := strings.Split(period, "-")
	if len(parts) == 2 {
		if y, err := strconv.Atoi(parts[0]); err == nil && y > 0 {
			in.Year = y
		}
		if m, err := strconv.Atoi(parts[1]); err == nil && m > 0 {
			in.Month = m
		}
	} else if len(parts) == 1 {
		if y, err := strconv.Atoi(parts[0]); err == nil && y > 0 {
			in.Year = y
		}
	}
	return in
}

// ComplianceTaxFilingTemplate 下载报税 Excel 模板
func (c *ControllerV1) ComplianceTaxFilingTemplate(ctx context.Context, req *member.ComplianceTaxFilingTemplateReq) (res *member.ComplianceTaxFilingTemplateRes, err error) {
	if _, err = requireMemberId(ctx); err != nil {
		return nil, err
	}
	data, err := service.ComplianceTax().GenerateFilingTemplate(ctx)
	if err != nil {
		return nil, err
	}
	if r := ghttp.RequestFromCtx(ctx); r != nil {
		r.Response.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		r.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "tax-filing-import.xlsx"))
		r.Response.Write(data)
	}
	return &member.ComplianceTaxFilingTemplateRes{}, nil
}

// ComplianceTaxFilingImportPreview 报税 Excel 导入预览
func (c *ControllerV1) ComplianceTaxFilingImportPreview(ctx context.Context, req *member.ComplianceTaxFilingImportPreviewReq) (res *member.ComplianceTaxFilingImportPreviewRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	fileContent := readUploadBytes(ctx, "file")
	out, err := service.ComplianceTax().PreviewTaxFilingImport(ctx, &compliancein.TaxFilingImportInp{
		MemberId:    memberId,
		FileContent: fileContent,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceTaxFilingImportPreviewRes{TaxFilingImportPreviewModel: out}, nil
}

// ComplianceTaxFilingImport 报税 Excel 导入
func (c *ControllerV1) ComplianceTaxFilingImport(ctx context.Context, req *member.ComplianceTaxFilingImportReq) (res *member.ComplianceTaxFilingImportRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	ip, _ := requestMeta(ctx)
	fileContent := readUploadBytes(ctx, "file")
	out, err := service.ComplianceTax().ImportTaxFilingExcel(ctx, &compliancein.TaxFilingImportInp{
		MemberId:    memberId,
		FileContent: fileContent,
		ExcelFileId: req.ExcelFileId,
		Ip:          ip,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceTaxFilingImportRes{TaxFilingImportModel: out}, nil
}

// ComplianceTaxFilingSubmissions 报税提交历史
func (c *ControllerV1) ComplianceTaxFilingSubmissions(ctx context.Context, req *member.ComplianceTaxFilingSubmissionsReq) (res *member.ComplianceTaxFilingSubmissionsRes, err error) {
	memberId, err := requireMemberId(ctx)
	if err != nil {
		return nil, err
	}
	out, err := service.ComplianceTax().ListTaxFilingSubmissions(ctx, &compliancein.TaxFilingSubmissionListInp{
		MemberId: memberId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &member.ComplianceTaxFilingSubmissionsRes{TaxFilingSubmissionListModel: out}, nil
}

