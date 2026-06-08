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
	g.Meta      `path:"/compliance/income" method:"post" tags:"会员合规" summary:"创建收入条目"`
	Platform    string  `json:"platform" v:"required#请选择平台"`
	Category    string  `json:"category" v:"required#请选择收入类型"`
	GrossAmount float64 `json:"grossAmount" v:"required|min:0.01#请填写含税收入|含税收入须大于0"`
	PlatformFee float64 `json:"platformFee"`
	OccurredAt  string  `json:"occurredAt" v:"required#请填写发生日期"`
	Remark      string  `json:"remark"`
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


