// ================================================================================
// Compliance service interfaces.
// ================================================================================

package service

import (
	"context"

	"xygo/internal/model/input/compliancein"
)

type (
	IComplianceDiagnosis interface {
		Preview(ctx context.Context, in *compliancein.DiagnosisSubmitInp) (*compliancein.DiagnosisSubmitModel, error)
		Submit(ctx context.Context, in *compliancein.DiagnosisSubmitInp) (*compliancein.DiagnosisSubmitModel, error)
		SyncGuest(ctx context.Context, in *compliancein.DiagnosisSyncInp) (*compliancein.DiagnosisSyncModel, error)
		BindToMember(ctx context.Context, in *compliancein.DiagnosisBindInp) (*compliancein.DiagnosisBindModel, error)
		GetDetail(ctx context.Context, in *compliancein.DiagnosisDetailInp) (*compliancein.DiagnosisDetailModel, error)
		Calculate(ctx context.Context, in *compliancein.TaxCalculatorInp) (*compliancein.TaxCalculatorModel, error)
		ListPlans(ctx context.Context) (*compliancein.ServicePlansModel, error)
		ListHistory(ctx context.Context, in *compliancein.DiagnosisHistoryInp) (*compliancein.DiagnosisHistoryModel, error)
	}

	IComplianceOrder interface {
		CreateOrder(ctx context.Context, in *compliancein.OrderCreateInp) (*compliancein.OrderCreateModel, error)
		GetActiveOrder(ctx context.Context, memberId uint64) (*compliancein.ActiveOrderModel, error)
		RecordConsent(ctx context.Context, in *compliancein.ConsentInp) (*compliancein.ConsentModel, error)
		SignContract(ctx context.Context, in *compliancein.SignInp) (*compliancein.SignModel, error)
	}

	IComplianceOpc interface {
		GetSummary(ctx context.Context, memberId uint64) (*compliancein.OpcSummaryModel, error)
		GetProgress(ctx context.Context, memberId uint64) (*compliancein.OpcProgressModel, error)
		GetMaterialsOverview(ctx context.Context, memberId uint64) (*compliancein.MaterialsOverviewModel, error)
		GetMaterialsDetail(ctx context.Context, memberId, opcId uint64) (*compliancein.MaterialsEntityDetailModel, error)
		GetMaterialsSection(ctx context.Context, memberId, opcId uint64, section string) (*compliancein.MaterialsSectionDetailModel, error)
		RevealMaterialsSection(ctx context.Context, in *compliancein.MaterialsRevealInp) (*compliancein.MaterialsRevealModel, error)
		RevealMaterialsAll(ctx context.Context, in *compliancein.MaterialsRevealInp) (*compliancein.MaterialsRevealAllModel, error)
		SubmitMaterials(ctx context.Context, in *compliancein.MaterialsSubmitInp) (*compliancein.MaterialsSubmitModel, error)
		SubmitBankReceipt(ctx context.Context, in *compliancein.BankReceiptInp) (*compliancein.BankReceiptModel, error)
		ListTasks(ctx context.Context, in *compliancein.OpcTaskListInp) (*compliancein.OpcTaskListModel, error)
		GetTaskDetail(ctx context.Context, in *compliancein.OpcTaskDetailInp) (*compliancein.OpcTaskDetailModel, error)
		AdvanceTask(ctx context.Context, in *compliancein.OpcTaskActionInp) (*compliancein.OpcTaskActionModel, error)
		ListCustomers(ctx context.Context, in *compliancein.ComplianceCustomerListInp) (*compliancein.ComplianceCustomerListModel, error)
	}

	IComplianceLedger interface {
		ListIncome(ctx context.Context, in *compliancein.IncomeListInp) (*compliancein.IncomeListModel, error)
		CreateIncome(ctx context.Context, in *compliancein.IncomeCreateInp) (*compliancein.IncomeCreateModel, error)
		DeleteIncome(ctx context.Context, in *compliancein.IncomeDeleteInp) error
		ImportIncomeCSV(ctx context.Context, in *compliancein.IncomeImportInp) (*compliancein.IncomeImportModel, error)
		PreviewIncomeImport(ctx context.Context, in *compliancein.IncomeImportInp) (*compliancein.IncomeImportPreviewModel, error)
		ListExpense(ctx context.Context, in *compliancein.ExpenseListInp) (*compliancein.ExpenseListModel, error)
		CreateExpense(ctx context.Context, in *compliancein.ExpenseCreateInp) (*compliancein.ExpenseCreateModel, error)
		DeleteExpense(ctx context.Context, in *compliancein.ExpenseDeleteInp) error
		ListExpenseTypes(ctx context.Context) (*compliancein.ExpenseTypesModel, error)
		ImportBankStatement(ctx context.Context, in *compliancein.BankImportInp) (*compliancein.BankImportModel, error)
		ListBankUnmatched(ctx context.Context, in *compliancein.BankUnmatchedInp) (*compliancein.BankUnmatchedModel, error)
		ListLedgerVouchers(ctx context.Context, in *compliancein.LedgerVoucherListInp) (*compliancein.LedgerVoucherListModel, error)
		GetProfit(ctx context.Context, in *compliancein.ProfitQueryInp) (*compliancein.ProfitSummaryModel, error)
	}

	IComplianceTax interface {
		GetCalendar(ctx context.Context, in *compliancein.TaxCalendarInp) (*compliancein.TaxCalendarModel, error)
		GetChecklist(ctx context.Context, in *compliancein.TaxChecklistInp) (*compliancein.TaxChecklistModel, error)
		GetTaskDetail(ctx context.Context, in *compliancein.TaxTaskDetailInp) (*compliancein.TaxTaskDetailModel, error)
	}

	IComplianceFiling interface {
		ListTasks(ctx context.Context, in *compliancein.FilingListInp) (*compliancein.FilingListModel, error)
		MarkFiled(ctx context.Context, in *compliancein.FilingMarkInp) (*compliancein.FilingMarkModel, error)
	}

	IComplianceStatement interface {
		ListStatements(ctx context.Context, in *compliancein.StatementListInp) (*compliancein.StatementListModel, error)
		ListMemberStatements(ctx context.Context, in *compliancein.MemberStatementListInp) (*compliancein.MemberStatementListModel, error)
		GenerateStatements(ctx context.Context, in *compliancein.StatementGenerateInp) (*compliancein.StatementGenerateModel, error)
		NotifyStatements(ctx context.Context, in *compliancein.StatementNotifyInp) (*compliancein.StatementNotifyModel, error)
		SendStatement(ctx context.Context, in *compliancein.StatementSendInp) (*compliancein.StatementSendModel, error)
	}

	IComplianceAudit interface {
		WriteAudit(ctx context.Context, entityType string, entityId uint64, action string, operatorId uint64, operatorType string, before, after interface{}, ip string) error
		ExportCSV(ctx context.Context, memberId uint64) ([]byte, string, error)
	}

	IComplianceDashboard interface {
		GetOverview(ctx context.Context) (*compliancein.DashboardOverviewModel, error)
	}

	IComplianceSocial interface {
		GetEmployment(ctx context.Context, memberId uint64) (*compliancein.EmploymentStatusModel, error)
		SetEmployment(ctx context.Context, in *compliancein.EmploymentStatusInp) (*compliancein.EmploymentStatusModel, error)
		ListGuides(ctx context.Context) (*compliancein.SocialGuideListModel, error)
		GetGuide(ctx context.Context, in *compliancein.SocialGuideDetailInp) (*compliancein.SocialGuideDetailModel, error)
		CreateConsult(ctx context.Context, in *compliancein.SocialConsultCreateInp) (*compliancein.SocialConsultCreateModel, error)
		ListMemberConsults(ctx context.Context, in *compliancein.SocialConsultListInp) (*compliancein.SocialConsultListModel, error)
		ListAdminConsults(ctx context.Context, in *compliancein.AdminSocialConsultListInp) (*compliancein.AdminSocialConsultListModel, error)
		ReplyConsult(ctx context.Context, in *compliancein.AdminSocialConsultReplyInp) (*compliancein.AdminSocialConsultReplyModel, error)
	}
)

var (
	localComplianceDiagnosis  IComplianceDiagnosis
	localComplianceOrder      IComplianceOrder
	localComplianceOpc        IComplianceOpc
	localComplianceLedger     IComplianceLedger
	localComplianceTax        IComplianceTax
	localComplianceFiling     IComplianceFiling
	localComplianceStatement  IComplianceStatement
	localComplianceAudit      IComplianceAudit
	localComplianceDashboard  IComplianceDashboard
	localComplianceSocial     IComplianceSocial
)

func ComplianceDiagnosis() IComplianceDiagnosis {
	if localComplianceDiagnosis == nil {
		panic("implement not found for interface IComplianceDiagnosis, forgot register?")
	}
	return localComplianceDiagnosis
}

func RegisterComplianceDiagnosis(i IComplianceDiagnosis) {
	localComplianceDiagnosis = i
}

func ComplianceOrder() IComplianceOrder {
	if localComplianceOrder == nil {
		panic("implement not found for interface IComplianceOrder, forgot register?")
	}
	return localComplianceOrder
}

func RegisterComplianceOrder(i IComplianceOrder) {
	localComplianceOrder = i
}

func ComplianceOpc() IComplianceOpc {
	if localComplianceOpc == nil {
		panic("implement not found for interface IComplianceOpc, forgot register?")
	}
	return localComplianceOpc
}

func RegisterComplianceOpc(i IComplianceOpc) {
	localComplianceOpc = i
}

func ComplianceLedger() IComplianceLedger {
	if localComplianceLedger == nil {
		panic("implement not found for interface IComplianceLedger, forgot register?")
	}
	return localComplianceLedger
}

func RegisterComplianceLedger(i IComplianceLedger) {
	localComplianceLedger = i
}

func ComplianceTax() IComplianceTax {
	if localComplianceTax == nil {
		panic("implement not found for interface IComplianceTax, forgot register?")
	}
	return localComplianceTax
}

func RegisterComplianceTax(i IComplianceTax) {
	localComplianceTax = i
}

func ComplianceFiling() IComplianceFiling {
	if localComplianceFiling == nil {
		panic("implement not found for interface IComplianceFiling, forgot register?")
	}
	return localComplianceFiling
}

func RegisterComplianceFiling(i IComplianceFiling) {
	localComplianceFiling = i
}

func ComplianceStatement() IComplianceStatement {
	if localComplianceStatement == nil {
		panic("implement not found for interface IComplianceStatement, forgot register?")
	}
	return localComplianceStatement
}

func RegisterComplianceStatement(i IComplianceStatement) {
	localComplianceStatement = i
}

func ComplianceAudit() IComplianceAudit {
	if localComplianceAudit == nil {
		panic("implement not found for interface IComplianceAudit, forgot register?")
	}
	return localComplianceAudit
}

func RegisterComplianceAudit(i IComplianceAudit) {
	localComplianceAudit = i
}

func ComplianceDashboard() IComplianceDashboard {
	if localComplianceDashboard == nil {
		panic("implement not found for interface IComplianceDashboard, forgot register?")
	}
	return localComplianceDashboard
}

func RegisterComplianceDashboard(i IComplianceDashboard) {
	localComplianceDashboard = i
}

func ComplianceSocial() IComplianceSocial {
	if localComplianceSocial == nil {
		panic("implement not found for interface IComplianceSocial, forgot register?")
	}
	return localComplianceSocial
}

func RegisterComplianceSocial(i IComplianceSocial) {
	localComplianceSocial = i
}
