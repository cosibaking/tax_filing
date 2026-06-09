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

	"github.com/gogf/gf/v2/frame/g"



	"xygo/internal/model/input/compliancein"

)



// ComplianceDiagnosisHistoryReq 会员诊断历史

type ComplianceDiagnosisHistoryReq struct {

	g.Meta   `path:"/compliance/diagnosis" method:"get" tags:"会员合规" summary:"诊断历史列表"`

	Page     int `p:"page" json:"page" d:"1"`

	PageSize int `p:"pageSize" json:"pageSize" d:"20"`

}



// ComplianceDiagnosisHistoryRes 会员诊断历史响应

type ComplianceDiagnosisHistoryRes struct {

	List     []compliancein.DiagnosisHistoryItem `json:"list"`

	Page     int                                 `json:"page"`

	PageSize int                                 `json:"pageSize"`

	Total    int                                 `json:"total"`

}



// ComplianceDiagnosisSyncItem 访客诊断同步项

type ComplianceDiagnosisSyncItem struct {

	Platforms          []string           `json:"platforms" v:"required#请选择直播平台"`

	MonthlyIncomeRange string             `json:"monthlyIncomeRange" v:"required#请选择月收入区间"`

	AnnualCostEstimate float64            `json:"annualCostEstimate"`

	ExistingEntity     string             `json:"existingEntity" v:"required#请选择现有主体"`

	HasFiledTax        string             `json:"hasFiledTax" v:"required#请选择报税状态"`

	TaxBureauContact   bool               `json:"taxBureauContact"`

	Notes              string             `json:"notes"`

	CostBreakdown      map[string]float64 `json:"costBreakdown"`

}



// ComplianceDiagnosisSyncReq 登录后同步访客诊断

type ComplianceDiagnosisSyncReq struct {

	g.Meta `path:"/compliance/diagnosis/sync" method:"post" tags:"会员合规" summary:"同步访客诊断到账户"`

	Items  []ComplianceDiagnosisSyncItem `json:"items" v:"required#请提供诊断数据"`

}



// ComplianceDiagnosisSyncRes 同步响应

type ComplianceDiagnosisSyncRes struct {

	*compliancein.DiagnosisSyncModel

}



// ComplianceDiagnosisBindReq 绑定匿名诊断记录

type ComplianceDiagnosisBindReq struct {

	g.Meta       `path:"/compliance/diagnosis/bind" method:"post" tags:"会员合规" summary:"绑定匿名诊断记录"`

	DiagnosisIds []uint64 `json:"diagnosisIds" v:"required#请提供诊断ID"`

}



// ComplianceDiagnosisBindRes 绑定响应

type ComplianceDiagnosisBindRes struct {

	*compliancein.DiagnosisBindModel

}



// ComplianceDiagnosisDetailReq 诊断详情

type ComplianceDiagnosisDetailReq struct {

	g.Meta      `path:"/compliance/diagnosis/{id}" method:"get" tags:"会员合规" summary:"诊断详情"`

	DiagnosisId uint64 `p:"id" json:"id" v:"required|min:1#请提供诊断ID|诊断ID无效"`

}



// ComplianceDiagnosisDetailRes 诊断详情响应

type ComplianceDiagnosisDetailRes struct {

	*compliancein.DiagnosisDetailModel

}



// ComplianceOrderCreateReq 创建服务订单

type ComplianceOrderCreateReq struct {

	g.Meta      `path:"/compliance/order" method:"post" tags:"会员合规" summary:"创建服务订单"`

	PlanId      uint64 `json:"planId" v:"required|min:1#请选择套餐|套餐无效"`

	DiagnosisId uint64 `json:"diagnosisId"`

}



// ComplianceOrderCreateRes 创建订单响应

type ComplianceOrderCreateRes struct {

	*compliancein.OrderCreateModel

}



// ComplianceConsentReq 风险告知 / 方案确认

type ComplianceConsentReq struct {

	g.Meta          `path:"/compliance/consent" method:"post" tags:"会员合规" summary:"合规同意留痕"`

	OrderId         uint64                           `json:"orderId" v:"required|min:1#请指定订单"`

	Type            string                           `json:"type" v:"required|in:risk_disclosure,plan_confirm#请指定类型|类型无效"`

	DocumentVersion string                           `json:"documentVersion"`

	Acknowledgments *compliancein.RiskAcknowledgments `json:"acknowledgments"`

	PlanConfirmed   bool                             `json:"planConfirmed"`

}



// ComplianceConsentRes 同意留痕响应

type ComplianceConsentRes struct {

	*compliancein.ConsentModel

}



// ComplianceSignReq 电子签约

type ComplianceSignReq struct {

	g.Meta    `path:"/compliance/sign" method:"post" tags:"会员合规" summary:"电子签约"`

	OrderId   uint64 `json:"orderId" v:"required|min:1#请指定订单"`

	LegalName string `json:"legalName" v:"required|length:2,20#请填写签约姓名|签约姓名须为2-20字"`

}



// ComplianceSignRes 签约响应

type ComplianceSignRes struct {

	*compliancein.SignModel

}



// ComplianceOpcSummaryReq OPC 主体摘要

type ComplianceOpcSummaryReq struct {

	g.Meta `path:"/compliance/opc" method:"get" tags:"会员合规" summary:"OPC主体摘要"`

}



// ComplianceOpcSummaryRes OPC 主体摘要响应

type ComplianceOpcSummaryRes struct {

	*compliancein.OpcSummaryModel

}



// ComplianceOpcProgressReq OPC 进度时间轴

type ComplianceOpcProgressReq struct {

	g.Meta `path:"/compliance/opc/progress" method:"get" tags:"会员合规" summary:"OPC进度时间轴"`

}



// ComplianceOpcProgressRes OPC 进度响应

type ComplianceOpcProgressRes struct {

	*compliancein.OpcProgressModel

}



// ComplianceOpcMaterialsReq 提交 OPC 注册资料

type ComplianceOpcMaterialsReq struct {

	g.Meta             `path:"/compliance/opc/materials" method:"post" tags:"会员合规" summary:"提交OPC注册资料"`

	ProposedNames      []string                         `json:"proposedNames" v:"required#请填写备选公司名称"`

	RegisteredCapital  float64                          `json:"registeredCapital" v:"required|min:0.01|max:1000#请填写注册资本"`

	CapitalTermYears   int                              `json:"capitalTermYears" v:"required|in:5,10,20,30#请选择认缴期限"`

	BusinessTermType   string                           `json:"businessTermType" v:"required|in:long_term,fixed#请选择营业期限"`

	BusinessTermEnd    string                           `json:"businessTermEnd"`

	BusinessScope      string                           `json:"businessScope" v:"required#请填写经营范围"`

	RegisterProvince   string                           `json:"registerProvince" v:"required#请填写注册省份"`

	RegisterCity       string                           `json:"registerCity" v:"required#请填写注册城市"`

	RegisterDistrict   string                           `json:"registerDistrict" v:"required#请填写注册区县"`

	RegisterAddress    string                           `json:"registerAddress" v:"required#请填写详细地址"`

	AddressProofFileId uint64                           `json:"addressProofFileId" v:"required|min:1#请上传地址证明"`

	LegalPersonName    string                           `json:"legalPersonName" v:"required#请填写法人姓名"`

	IdCard             string                           `json:"idCard" v:"required|length:18,18#请填写身份证号"`

	IdCardValidFrom    string                           `json:"idCardValidFrom" v:"required#请填写身份证有效期起"`

	IdCardValidTo      string                           `json:"idCardValidTo" v:"required#请填写身份证有效期止"`

	Ethnicity          string                           `json:"ethnicity"`

	HouseholdAddress   string                           `json:"householdAddress" v:"required#请填写户籍地址"`

	ResidentialAddress string                           `json:"residentialAddress" v:"required#请填写现居住地址"`

	Phone              string                           `json:"phone" v:"required|phone#请填写手机号"`

	Email              string                           `json:"email" v:"required|email#请填写邮箱"`

	IdCardFrontFileId  uint64                           `json:"idCardFrontFileId" v:"required|min:1#请上传身份证正面"`

	IdCardBackFileId   uint64                           `json:"idCardBackFileId" v:"required|min:1#请上传身份证反面"`

	EsignAuthorized    bool                             `json:"esignAuthorized"`

	Confirmations      compliancein.MaterialsConfirmations `json:"confirmations"`

}



// ComplianceOpcMaterialsRes 资料提交响应

type ComplianceOpcMaterialsRes struct {

	*compliancein.MaterialsSubmitModel

}



// ComplianceOpcBankReceiptReq 上传开户回执

type ComplianceOpcBankReceiptReq struct {

	g.Meta            `path:"/compliance/opc/bank-receipt" method:"post" tags:"会员合规" summary:"上传开户回执"`

	BankName          string `json:"bankName"`

	BankReceiptFileId uint64 `json:"bankReceiptFileId" v:"required|min:1#请上传开户回执"`

}



// ComplianceOpcBankReceiptRes 开户回执响应

type ComplianceOpcBankReceiptRes struct {

	*compliancein.BankReceiptModel

}

// ComplianceOpcMaterialsOverviewReq 已提交资料概览
type ComplianceOpcMaterialsOverviewReq struct {
	g.Meta `path:"/compliance/opc/materials/overview" method:"get" tags:"会员合规" summary:"OPC已提交资料概览"`
}

// ComplianceOpcMaterialsOverviewRes 资料概览响应
type ComplianceOpcMaterialsOverviewRes struct {
	*compliancein.MaterialsOverviewModel
}

// ComplianceOpcMaterialsDetailReq 单个 OPC 实体完整资料（脱敏）
type ComplianceOpcMaterialsDetailReq struct {
	g.Meta `path:"/compliance/opc/materials/detail" method:"get" tags:"会员合规" summary:"OPC已提交资料详情"`
	OpcId  uint64 `p:"opcId" json:"opcId" v:"required|min:1#请指定OPC主体"`
}

// ComplianceOpcMaterialsDetailRes 资料详情响应
type ComplianceOpcMaterialsDetailRes struct {
	*compliancein.MaterialsEntityDetailModel
}

// ComplianceOpcMaterialsSectionReq 资料分组详情（脱敏）
type ComplianceOpcMaterialsSectionReq struct {
	g.Meta  `path:"/compliance/opc/materials/section" method:"get" tags:"会员合规" summary:"OPC资料分组详情"`
	OpcId   uint64 `p:"opcId" json:"opcId"`
	Section string `p:"section" json:"section" v:"required#请指定资料分组"`
}

// ComplianceOpcMaterialsSectionRes 资料分组详情响应
type ComplianceOpcMaterialsSectionRes struct {
	*compliancein.MaterialsSectionDetailModel
}

// ComplianceOpcMaterialsRevealReq 密码验证查看完整资料
type ComplianceOpcMaterialsRevealReq struct {
	g.Meta   `path:"/compliance/opc/materials/reveal" method:"post" tags:"会员合规" summary:"密码验证查看完整OPC资料"`
	OpcId    uint64 `json:"opcId" v:"required|min:1#请指定OPC主体"`
	Section  string `json:"section"`
	Password string `json:"password" v:"required#请输入密码"`
}

// ComplianceOpcMaterialsRevealRes 完整资料响应
type ComplianceOpcMaterialsRevealRes struct {
	*compliancein.MaterialsRevealModel
}

// ComplianceOpcMaterialsRevealAllReq 密码验证查看全部完整资料
type ComplianceOpcMaterialsRevealAllReq struct {
	g.Meta   `path:"/compliance/opc/materials/reveal-all" method:"post" tags:"会员合规" summary:"密码验证查看全部OPC资料"`
	OpcId    uint64 `json:"opcId" v:"required|min:1#请指定OPC主体"`
	Password string `json:"password" v:"required#请输入密码"`
}

// ComplianceOpcMaterialsRevealAllRes 全部完整资料响应
type ComplianceOpcMaterialsRevealAllRes struct {
	*compliancein.MaterialsRevealAllModel
}

// ComplianceIncomeListReq 收入台账列表
type ComplianceIncomeListReq struct {
	g.Meta   `path:"/compliance/income" method:"get" tags:"会员合规" summary:"收入台账列表"`
	Page     int    `p:"page" json:"page" d:"1"`
	PageSize int    `p:"pageSize" json:"pageSize" d:"20"`
	Month    string `p:"month" json:"month"`
	Platform string `p:"platform" json:"platform"`
	Category string `p:"category" json:"category"`
}

type ComplianceIncomeListRes struct {
	*compliancein.IncomeListModel
}

type ComplianceIncomeCreateReq struct {
	g.Meta           `path:"/compliance/income" method:"post" tags:"会员合规" summary:"创建收入条目"`
	Platform         string  `json:"platform" v:"required#请选择平台"`
	Category         string  `json:"category" v:"required#请选择收入类型"`
	GrossAmount      float64 `json:"grossAmount"`
	PlatformFee      float64 `json:"platformFee"`
	OccurredAt       string  `json:"occurredAt" v:"required#请填写发生日期"`
	SettlementType   string  `json:"settlementType"`
	McnName          string  `json:"mcnName"`
	McnSplitRatio    float64 `json:"mcnSplitRatio"`
	GrossBeforeSplit float64 `json:"grossBeforeSplit"`
	Remark           string  `json:"remark"`
}

type ComplianceIncomeCreateRes struct {
	*compliancein.IncomeCreateModel
}

type ComplianceIncomeImportReq struct {
	g.Meta     `path:"/compliance/income/import" method:"post" mime:"multipart/form-data" tags:"会员合规" summary:"CSV导入收入"`
	CsvContent string `json:"csvContent"`
}

type ComplianceIncomeImportRes struct {
	*compliancein.IncomeImportModel
}

type ComplianceExpenseListReq struct {
	g.Meta   `path:"/compliance/expense" method:"get" tags:"会员合规" summary:"费用台账列表"`
	Page     int    `p:"page" json:"page" d:"1"`
	PageSize int    `p:"pageSize" json:"pageSize" d:"20"`
	Month    string `p:"month" json:"month"`
	Category string `p:"category" json:"category"`
}

type ComplianceExpenseListRes struct {
	*compliancein.ExpenseListModel
}

type ComplianceExpenseCreateReq struct {
	g.Meta         `path:"/compliance/expense" method:"post" tags:"会员合规" summary:"创建费用条目"`
	Category       string  `json:"category" v:"required#请选择费用类型"`
	Amount         float64 `json:"amount" v:"required|min:0.01#请填写金额|金额须大于0"`
	InvoiceType    string  `json:"invoiceType"`
	AttachmentId   uint64  `json:"attachmentId"`
	Description    string  `json:"description"`
	OccurredAt     string  `json:"occurredAt" v:"required#请填写发生日期"`
	ConfirmWarning bool    `json:"confirmWarning"`
}

type ComplianceExpenseCreateRes struct {
	*compliancein.ExpenseCreateModel
}

type ComplianceBankImportReq struct {
	g.Meta     `path:"/compliance/bank/import" method:"post" mime:"multipart/form-data" tags:"会员合规" summary:"银行流水导入与对账"`
	CsvContent string `json:"csvContent"`
}

type ComplianceBankImportRes struct {
	*compliancein.BankImportModel
}

type ComplianceExpenseTypesReq struct {
	g.Meta `path:"/compliance/expense-types" method:"get" tags:"会员合规" summary:"费用类型库"`
}

type ComplianceExpenseTypesRes struct {
	*compliancein.ExpenseTypesModel
}

// ComplianceOrderActiveReq 当前活跃订单
type ComplianceOrderActiveReq struct {
	g.Meta `path:"/compliance/order/active" method:"get" tags:"会员合规" summary:"当前活跃订单"`
}

type ComplianceOrderActiveRes struct {
	*compliancein.ActiveOrderModel
}

// ComplianceLedgerProfitReq 利润表
type ComplianceLedgerProfitReq struct {
	g.Meta     `path:"/compliance/ledger/profit" method:"get" tags:"会员合规" summary:"利润表"`
	Period     string `p:"period" json:"period"`
	PeriodType string `p:"periodType" json:"periodType"`
	Year       int    `p:"year" json:"year"`
	Month      int    `p:"month" json:"month"`
	Quarter    int    `p:"quarter" json:"quarter"`
}

type ComplianceLedgerProfitRes struct {
	*compliancein.ProfitSummaryModel
}

type ComplianceIncomeDeleteReq struct {
	g.Meta `path:"/compliance/income" method:"delete" tags:"会员合规" summary:"删除收入条目"`
	Id     uint64 `p:"id" json:"id" v:"required|min:1#请指定条目|条目无效"`
}

type ComplianceIncomeDeleteRes struct{}

type ComplianceIncomeImportPreviewReq struct {
	g.Meta     `path:"/compliance/income/import/preview" method:"post" mime:"multipart/form-data" tags:"会员合规" summary:"CSV导入预览"`
	CsvContent string `json:"csvContent"`
}

type ComplianceIncomeImportPreviewRes struct {
	*compliancein.IncomeImportPreviewModel
}

type ComplianceIncomeOCRPreviewReq struct {
	g.Meta     `path:"/compliance/income/ocr/preview" method:"post" mime:"multipart/form-data" tags:"会员合规" summary:"平台流水OCR预览"`
	Platform   string `p:"platform" json:"platform"`
	OcrText    string `p:"ocrText" json:"ocrText"`
	AttachmentId uint64 `p:"attachmentId" json:"attachmentId"`
}

type ComplianceIncomeOCRPreviewRes struct {
	*compliancein.IncomeOCRPreviewModel
}

type ComplianceIncomeOCRImportReq struct {
	g.Meta       `path:"/compliance/income/ocr/import" method:"post" mime:"multipart/form-data" tags:"会员合规" summary:"平台流水OCR导入"`
	Platform     string `p:"platform" json:"platform"`
	OcrText      string `p:"ocrText" json:"ocrText"`
	AttachmentId uint64 `p:"attachmentId" json:"attachmentId"`
}

type ComplianceIncomeOCRImportRes struct {
	*compliancein.IncomeImportModel
}

type ComplianceIncomeConsistencyReq struct {
	g.Meta `path:"/compliance/income/consistency" method:"get" tags:"会员合规" summary:"收入一致性比对"`
	Month  string `p:"month" json:"month"`
}

type ComplianceIncomeConsistencyRes struct {
	*compliancein.IncomeConsistencyModel
}

type ComplianceExpenseDeleteReq struct {
	g.Meta `path:"/compliance/expense" method:"delete" tags:"会员合规" summary:"删除费用条目"`
	Id     uint64 `p:"id" json:"id" v:"required|min:1#请指定条目|条目无效"`
}

type ComplianceExpenseDeleteRes struct{}

type ComplianceExpenseCategoriesReq struct {
	g.Meta `path:"/compliance/expense/categories" method:"get" tags:"会员合规" summary:"费用类型库"`
}

type ComplianceExpenseCategoriesRes struct {
	List []compliancein.ExpenseTypeItem `json:"list"`
}

type ComplianceBankUnmatchedReq struct {
	g.Meta `path:"/compliance/bank/unmatched" method:"get" tags:"会员合规" summary:"未匹配银行流水"`
	Month  string `p:"month" json:"month"`
}

type ComplianceBankUnmatchedRes struct {
	*compliancein.BankUnmatchedModel
}

type ComplianceLedgerVouchersReq struct {
	g.Meta `path:"/compliance/ledger/vouchers" method:"get" tags:"会员合规" summary:"会计分录列表"`
	Period string `p:"period" json:"period"`
}

type ComplianceLedgerVouchersRes struct {
	*compliancein.LedgerVoucherListModel
}

// ComplianceTaxCalendarReq 申报日历
type ComplianceTaxCalendarReq struct {
	g.Meta `path:"/compliance/tax/calendar" method:"get" tags:"会员合规" summary:"申报日历"`
	Year   int `p:"year" json:"year"`
	Month  int `p:"month" json:"month"`
}

type ComplianceTaxCalendarRes struct {
	*compliancein.TaxCalendarModel
}

// ComplianceTaxChecklistReq 报税前自查清单
type ComplianceTaxChecklistReq struct {
	g.Meta `path:"/compliance/tax/checklist" method:"get" tags:"会员合规" summary:"报税前自查清单"`
	TaskId uint64 `p:"taskId" json:"taskId"`
	Period string `p:"period" json:"period"`
}

type ComplianceTaxChecklistRes struct {
	*compliancein.TaxChecklistModel
}

type ComplianceTaxTaskDetailReq struct {
	g.Meta `path:"/compliance/tax/tasks/detail" method:"get" tags:"会员合规" summary:"申报任务详情"`
	Id     uint64 `p:"id" json:"id" v:"required|min:1#请指定任务|任务无效"`
}

type ComplianceTaxTaskDetailRes struct {
	*compliancein.TaxTaskDetailModel
}

type ComplianceMemberStatementsReq struct {
	g.Meta   `path:"/compliance/statements" method:"get" tags:"会员合规" summary:"月度对账单列表"`
	Page     int `p:"page" json:"page" d:"1"`
	PageSize int `p:"pageSize" json:"pageSize" d:"20"`
}

type ComplianceMemberStatementsRes struct {
	*compliancein.MemberStatementListModel
}

type ComplianceStatementPdfReq struct {
	g.Meta `path:"/compliance/statements/{id}/pdf" method:"get" tags:"会员合规" summary:"对账单PDF（MVP未实现）"`
	Id     uint64 `p:"id" in:"path" json:"id" v:"required|min:1#请指定对账单|对账单无效"`
}

type ComplianceStatementPdfRes struct {
	Url string `json:"url"`
}

// ComplianceEmploymentGetReq 获取用工状态
type ComplianceEmploymentGetReq struct {
	g.Meta `path:"/compliance/employment" method:"get" tags:"会员合规" summary:"获取用工状态"`
}

type ComplianceEmploymentGetRes struct {
	*compliancein.EmploymentStatusModel
}

// ComplianceEmploymentSetReq 设置用工状态
type ComplianceEmploymentSetReq struct {
	g.Meta           `path:"/compliance/employment" method:"put" tags:"会员合规" summary:"设置用工状态"`
	EmploymentStatus string `json:"employmentStatus" v:"required|in:no_employee,has_employee#请选择用工状态|用工状态无效"`
}

type ComplianceEmploymentSetRes struct {
	*compliancein.EmploymentStatusModel
}

// ComplianceSocialGuidesReq 社保指引列表
type ComplianceSocialGuidesReq struct {
	g.Meta `path:"/compliance/social/guides" method:"get" tags:"会员合规" summary:"社保指引列表"`
}

type ComplianceSocialGuidesRes struct {
	*compliancein.SocialGuideListModel
}

// ComplianceSocialGuideDetailReq 社保指引详情
type ComplianceSocialGuideDetailReq struct {
	g.Meta `path:"/compliance/social/guides/{slug}" method:"get" tags:"会员合规" summary:"社保指引详情"`
	Slug   string `p:"slug" in:"path" json:"slug" v:"required#请指定指引文章"`
}

type ComplianceSocialGuideDetailRes struct {
	*compliancein.SocialGuideDetailModel
}

// ComplianceSocialConsultCreateReq 提交社保咨询
type ComplianceSocialConsultCreateReq struct {
	g.Meta     `path:"/compliance/social/consults" method:"post" tags:"会员合规" summary:"提交社保咨询"`
	Category   string `json:"category" d:"founder"`
	Question   string `json:"question" v:"required|min-length:5#请描述问题|问题至少5个字"`
	RegionCode string `json:"regionCode"`
}

type ComplianceSocialConsultCreateRes struct {
	*compliancein.SocialConsultCreateModel
}

// ComplianceSocialConsultListReq 社保咨询列表
type ComplianceSocialConsultListReq struct {
	g.Meta   `path:"/compliance/social/consults" method:"get" tags:"会员合规" summary:"社保咨询列表"`
	Page     int `p:"page" json:"page" d:"1"`
	PageSize int `p:"pageSize" json:"pageSize" d:"20"`
}

type ComplianceSocialConsultListRes struct {
	*compliancein.SocialConsultListModel
}

// ComplianceTaxFilingTemplateReq 下载报税 Excel 模板
type ComplianceTaxFilingTemplateReq struct {
	g.Meta `path:"/compliance/tax/filing/template" method:"get" tags:"会员合规" summary:"下载报税Excel模板"`
}

type ComplianceTaxFilingTemplateRes struct{}

// ComplianceTaxFilingImportPreviewReq 报税 Excel 导入预览
type ComplianceTaxFilingImportPreviewReq struct {
	g.Meta `path:"/compliance/tax/filing/import/preview" method:"post" mime:"multipart/form-data" tags:"会员合规" summary:"报税Excel导入预览"`
}

type ComplianceTaxFilingImportPreviewRes struct {
	*compliancein.TaxFilingImportPreviewModel
}

// ComplianceTaxFilingImportReq 报税 Excel 导入
type ComplianceTaxFilingImportReq struct {
	g.Meta        `path:"/compliance/tax/filing/import" method:"post" mime:"multipart/form-data" tags:"会员合规" summary:"报税Excel导入"`
	ExcelFileId   uint64 `json:"excelFileId"`
}

type ComplianceTaxFilingImportRes struct {
	*compliancein.TaxFilingImportModel
}

// ComplianceTaxFilingSubmissionsReq 报税提交历史
type ComplianceTaxFilingSubmissionsReq struct {
	g.Meta   `path:"/compliance/tax/filing/submissions" method:"get" tags:"会员合规" summary:"报税提交历史"`
	Page     int `p:"page" json:"page" d:"1"`
	PageSize int `p:"pageSize" json:"pageSize" d:"20"`
}

type ComplianceTaxFilingSubmissionsRes struct {
	*compliancein.TaxFilingSubmissionListModel
}


