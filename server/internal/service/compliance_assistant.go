package service

import (
	"context"

	compliancedocument "xygo/internal/logic/compliance/document"
	"xygo/internal/logic/compliance/profile"
	compliancereport "xygo/internal/logic/compliance/report"
	compliancerisk "xygo/internal/logic/compliance/risk"
	"xygo/internal/logic/compliance/ruleengine"
	compliancetask "xygo/internal/logic/compliance/task"
	complianceticket "xygo/internal/logic/compliance/ticket"
)

type IComplianceProfile interface {
	Save(ctx context.Context, in profile.SaveInput) (*profile.Profile, bool, error)
	GetForMember(ctx context.Context, memberID, opcID uint64) (*profile.Profile, error)
}

type IComplianceRuleVersion interface {
	Transition(ctx context.Context, in ruleengine.TransitionInput) (*ruleengine.RuleVersion, error)
	Simulate(ctx context.Context, versionID uint64, facts map[string]any) (bool, error)
}

type IComplianceTask interface {
	Generate(ctx context.Context, in compliancetask.GenerateInput) (*compliancetask.Task, bool, error)
	List(ctx context.Context, query compliancetask.Query) ([]compliancetask.Task, int, error)
	GetForMember(ctx context.Context, memberID, id uint64) (*compliancetask.Task, error)
	ApplyEvent(ctx context.Context, memberID, id uint64, event, note string) (*compliancetask.Task, error)
}

type IComplianceDocument interface {
	Register(ctx context.Context, in compliancedocument.RegisterInput) (*compliancedocument.BusinessDocument, error)
	List(ctx context.Context, memberID uint64, periodKey string) ([]compliancedocument.BusinessDocument, error)
	Confirm(ctx context.Context, memberID, id uint64, documentType string, fields map[string]any) (*compliancedocument.BusinessDocument, error)
}

type IComplianceRisk interface {
	ScanAndSave(ctx context.Context, opcEntityID, memberID uint64, periodKey string, facts compliancerisk.Facts) ([]compliancerisk.Event, error)
	List(ctx context.Context, memberID uint64, periodKey, status string) ([]compliancerisk.Event, error)
	Decide(ctx context.Context, memberID, id uint64, action, note string) (*compliancerisk.Event, error)
}
type IComplianceReport interface {
	Create(ctx context.Context, opcID, memberID uint64, in compliancereport.Input) (*compliancereport.MonthlyReport, error)
	List(ctx context.Context, memberID uint64, period string) ([]compliancereport.MonthlyReport, error)
	Publish(ctx context.Context, memberID, id uint64) (*compliancereport.MonthlyReport, error)
	ExportPDF(ctx context.Context, memberID, id uint64) (filename string, content []byte, err error)
}
type IComplianceTicket interface {
	Create(context.Context, complianceticket.CreateInput) (*complianceticket.Ticket, error)
	List(context.Context, uint64) ([]complianceticket.Ticket, error)
	Apply(context.Context, uint64, uint64, string, string) (*complianceticket.Ticket, error)
}

var localComplianceProfile IComplianceProfile
var localComplianceRuleVersion IComplianceRuleVersion
var localComplianceTask IComplianceTask
var localComplianceDocument IComplianceDocument
var localComplianceRisk IComplianceRisk
var localComplianceReport IComplianceReport
var localComplianceTicket IComplianceTicket

func ComplianceProfile() IComplianceProfile {
	if localComplianceProfile == nil {
		panic("implement not found for interface IComplianceProfile, forgot register?")
	}
	return localComplianceProfile
}

func RegisterComplianceProfile(i IComplianceProfile) {
	localComplianceProfile = i
}

func ComplianceRuleVersion() IComplianceRuleVersion {
	if localComplianceRuleVersion == nil {
		panic("implement not found for interface IComplianceRuleVersion, forgot register?")
	}
	return localComplianceRuleVersion
}

func RegisterComplianceRuleVersion(i IComplianceRuleVersion) {
	localComplianceRuleVersion = i
}

func ComplianceTask() IComplianceTask {
	if localComplianceTask == nil {
		panic("implement not found for interface IComplianceTask, forgot register?")
	}
	return localComplianceTask
}

func RegisterComplianceTask(i IComplianceTask) { localComplianceTask = i }

func ComplianceDocument() IComplianceDocument {
	if localComplianceDocument == nil {
		panic("IComplianceDocument not registered")
	}
	return localComplianceDocument
}
func RegisterComplianceDocument(i IComplianceDocument) { localComplianceDocument = i }

func ComplianceRisk() IComplianceRisk {
	if localComplianceRisk == nil {
		panic("IComplianceRisk not registered")
	}
	return localComplianceRisk
}
func RegisterComplianceRisk(i IComplianceRisk) { localComplianceRisk = i }
func ComplianceReport() IComplianceReport {
	if localComplianceReport == nil {
		panic("IComplianceReport not registered")
	}
	return localComplianceReport
}
func RegisterComplianceReport(i IComplianceReport) { localComplianceReport = i }
func ComplianceTicket() IComplianceTicket {
	if localComplianceTicket == nil {
		panic("IComplianceTicket not registered")
	}
	return localComplianceTicket
}
func RegisterComplianceTicket(i IComplianceTicket) { localComplianceTicket = i }
