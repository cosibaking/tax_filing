package report

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// SnapshotLoader is the trust boundary for monthly reports. Implementations
// must scope every source by the authenticated member, OPC and requested period.
type SnapshotLoader interface {
	Load(ctx context.Context, opcID, memberID uint64, periodKey string) (Input, error)
}

type DatabaseSnapshotLoader struct{}

type documentSnapshotRow struct {
	DocumentType string `orm:"document_type"`
}

type transactionSnapshotRow struct {
	TransactionCount int     `orm:"transaction_count"`
	MatchedCount     int     `orm:"matched_count"`
	TotalAmount      float64 `orm:"total_amount"`
}

type riskSnapshotRow struct {
	RiskCode      string `orm:"risk_code"`
	Severity      string `orm:"severity"`
	Summary       string `orm:"summary"`
	EvidenceJSON  string `orm:"evidence_json"`
	RuleVersionID uint64 `orm:"rule_version_id"`
}

func (DatabaseSnapshotLoader) Load(ctx context.Context, opcID, memberID uint64, periodKey string) (Input, error) {
	if opcID == 0 || memberID == 0 || !validPeriod(periodKey) {
		return Input{}, fmt.Errorf("企业、会员和有效账期不能为空")
	}

	var docs []documentSnapshotRow
	err := g.DB().Model("xy_business_document").Ctx(ctx).
		Fields("document_type").
		Where("opc_entity_id", opcID).Where("member_id", memberID).
		Where("period_key", periodKey).Where("process_status", "confirmed").Where("deleted", 0).
		Scan(&docs)
	if err != nil {
		return Input{}, fmt.Errorf("读取已确认经营资料失败: %w", err)
	}

	confirmedTypes := make(map[string]bool)
	for _, item := range docs {
		if item.DocumentType != "" && item.DocumentType != "unknown" {
			confirmedTypes[item.DocumentType] = true
		}
	}
	expected := []string{"bank_statement", "sales_invoice", "expense_invoice", "contract", "tax_receipt"}
	missing := make([]string, 0, len(expected))
	for _, documentType := range expected {
		if !confirmedTypes[documentType] {
			missing = append(missing, documentType)
		}
	}
	confirmed := make([]string, 0, len(confirmedTypes))
	for documentType := range confirmedTypes {
		confirmed = append(confirmed, documentType)
	}
	sort.Strings(confirmed)
	rate := float64(len(expected)-len(missing)) * 100 / float64(len(expected))

	start, end, err := periodUnixRange(periodKey)
	if err != nil {
		return Input{}, err
	}
	var tx transactionSnapshotRow
	err = g.DB().Model("xy_bank_transaction").Ctx(ctx).
		Fields("COUNT(*) AS transaction_count", "COALESCE(SUM(CASE WHEN matched = 1 THEN 1 ELSE 0 END),0) AS matched_count", "COALESCE(SUM(amount),0) AS total_amount").
		Where("opc_id", opcID).Where("deleted", 0).
		WhereGTE("occurred_at", start).WhereLT("occurred_at", end).
		Scan(&tx)
	if err != nil {
		return Input{}, fmt.Errorf("读取银行流水统计失败: %w", err)
	}
	statistics := map[string]any{
		"transactionCount": tx.TransactionCount,
		"matchedCount":     tx.MatchedCount,
		"unmatchedCount":   tx.TransactionCount - tx.MatchedCount,
		"totalAmount":      tx.TotalAmount,
		"sourceAvailable":  tx.TransactionCount > 0 || confirmedTypes["bank_statement"], "trustedSnapshot": true,
	}
	completeness := map[string]any{
		"rate": rate, "confirmedDocumentTypes": confirmed, "missingDocumentTypes": missing,
		"confirmedDocumentCount": len(docs), "expectedDocumentTypes": expected, "trustedSnapshot": true,
	}

	var riskRows []riskSnapshotRow
	err = g.DB().Model("xy_risk_event").Ctx(ctx).
		Fields("risk_code", "severity", "summary", "evidence_json", "rule_version_id").
		Where("opc_entity_id", opcID).Where("member_id", memberID).Where("period_key", periodKey).
		WhereIn("status", []string{"open", "confirmed"}).Where("deleted", 0).
		OrderAsc("risk_code").Scan(&riskRows)
	if err != nil {
		return Input{}, fmt.Errorf("读取合规风险事件失败: %w", err)
	}
	risks := make([]map[string]any, 0, len(riskRows))
	ruleIDs := make(map[uint64]bool)
	for _, row := range riskRows {
		var evidence map[string]any
		if strings.TrimSpace(row.EvidenceJSON) != "" {
			_ = json.Unmarshal([]byte(row.EvidenceJSON), &evidence)
		}
		risks = append(risks, trustedRisk(row, evidence))
		if row.RuleVersionID > 0 {
			ruleIDs[row.RuleVersionID] = true
		}
	}
	if len(missing) > 0 {
		risks = append(risks, missingDocumentsRisk(missing))
	}
	if tx.TransactionCount > tx.MatchedCount {
		risks = append(risks, map[string]any{
			"code": "BANK_TRANSACTION_UNMATCHED", "categoryCode": "business", "title": "存在未匹配银行流水", "severity": "medium",
			"facts": fmt.Sprintf("当前账期共识别 %d 笔银行流水，其中 %d 笔尚未匹配业务记录。", tx.TransactionCount, tx.TransactionCount-tx.MatchedCount),
			"basis": "服务端银行流水的 matched 状态统计。", "impact": "未匹配流水可能导致收入、支出或往来款项归类不完整。",
			"recommendation": "请逐笔核对未匹配流水并补充对应发票、合同或款项说明。", "requiredMaterials": []string{"银行流水", "发票", "业务合同"},
			"requiresManualReview": false, "ruleVersion": "report-snapshot-v1",
		})
	}
	ruleVersions := make([]uint64, 0, len(ruleIDs))
	for id := range ruleIDs {
		ruleVersions = append(ruleVersions, id)
	}
	sort.Slice(ruleVersions, func(i, j int) bool { return ruleVersions[i] < ruleVersions[j] })

	return Input{PeriodKey: periodKey, Statistics: statistics, Completeness: completeness, Risks: risks, RuleVersions: ruleVersions}, nil
}

func missingDocumentsRisk(missing []string) map[string]any {
	labels := map[string]string{"bank_statement": "银行流水", "sales_invoice": "销项发票", "expense_invoice": "费用凭证", "contract": "业务合同", "tax_receipt": "完税证明"}
	names := make([]string, 0, len(missing))
	for _, code := range missing {
		if name := labels[code]; name != "" {
			names = append(names, name)
		}
	}
	return map[string]any{
		"code": "DOCUMENTS_INCOMPLETE", "categoryCode": "documents", "title": "本期经营资料不完整", "severity": "medium",
		"facts": "当前账期尚未确认以下资料：" + strings.Join(names, "、") + "。", "basis": "服务端按本账期已确认资料类型与基础资料清单逐项比对。",
		"impact":         "资料缺失可能导致经营统计或风险判断不完整，本报告不会将未知情况判定为正常。",
		"recommendation": "请在经营资料库上传并确认缺失资料后重新生成体检报告。", "requiredMaterials": names,
		"requiresManualReview": false, "ruleVersion": "report-snapshot-v1",
	}
}

func periodUnixRange(period string) (uint64, uint64, error) {
	t, err := time.ParseInLocation("2006-01", period, time.Local)
	if err != nil || t.Format("2006-01") != period {
		return 0, 0, fmt.Errorf("periodKey 必须是真实有效的 YYYY-MM")
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	return uint64(start.Unix()), uint64(start.AddDate(0, 1, 0).Unix()), nil
}

type riskNarrative struct {
	category, facts, basis, impact, recommendation string
	materials                                      []string
}

var riskNarratives = map[string]riskNarrative{
	"REV_NO_SALES_DOC":          {"business", "已记录经营收入，但当前账期未发现对应销项或未开票收入记录。", "收入记录与销项资料匹配检查触发。", "可能导致收入与开票资料不一致，影响后续申报资料核对。", "请核对收入来源并补充销项发票、未开票收入台账或业务合同。", []string{"销项发票", "收入台账", "业务合同"}},
	"INVOICE_RECEIPT_GAP":       {"business", "当前账期开票金额与已确认回款金额存在较大差异。", "开票金额与回款金额差异比例超过体检规则阈值。", "可能存在回款未匹配、跨期收款或资料遗漏。", "请逐笔核对发票、银行流水和应收账款记录。", []string{"销项发票", "银行流水", "应收账款明细"}},
	"ZERO_WITH_CASHFLOW":        {"tax", "零申报标记期间存在经营流水。", "零申报状态与经营流水交叉检查触发。", "可能影响申报信息一致性，需尽快由专业人员复核。", "请暂停直接下结论，整理流水和收入凭证后申请人工复核。", []string{"银行流水", "收入台账", "申报表"}},
	"COMPANY_TO_PERSONAL":       {"funds", "公司向个人账户转账次数较多或金额较大。", "个人转账次数或金额达到体检规则阈值。", "可能存在款项性质不清或公司与个人资金混同风险。", "请逐笔标注款项用途并补充合同、报销单或借款资料。", []string{"银行流水", "报销单", "借款协议"}},
	"COST_DOC_MISSING":          {"documents", "成本费用金额高于当前已确认的成本凭证金额。", "费用记录与成本凭证覆盖检查触发。", "可能影响成本费用的完整归集和后续申报核对。", "请补充费用发票、合同及付款凭证。", []string{"费用发票", "付款凭证", "业务合同"}},
	"TASK_DUE_SOON":             {"tax", "当前存在临近截止日的高优先级合规任务。", "任务截止日期进入规则预警窗口。", "如未及时处理，可能发生事项逾期。", "请查看合规日历并优先完成对应任务。", []string{"任务所需材料"}},
	"ANNUAL_REPORT_UNCONFIRMED": {"annual", "当前处于工商年报办理窗口，但年报事项尚未确认。", "年度事项窗口与完成状态检查触发。", "如未在规定期限内处理，可能影响企业公示状态。", "请核对年报资料并通过官方渠道办理，重要信息请人工复核。", []string{"工商年报资料", "财务数据"}},
	"ADDRESS_CHANGED":           {"annual", "经营地址发生变化，但企业档案尚未同步更新。", "地址变化记录与企业档案一致性检查触发。", "可能导致登记信息不一致或联系异常。", "请核实实际经营地址并咨询是否需要办理信息变更。", []string{"地址证明", "企业登记资料"}},
	"EMPLOYEE_PROFILE_GAP":      {"employment", "已有员工记录，但用工资料尚不完整。", "员工数量与用工档案完整性检查触发。", "可能影响工资、个税、社保或公积金事项核对。", "请补充员工、工资及社保资料，并由专业人员复核。", []string{"员工花名册", "工资表", "社保缴费记录"}},
	"SALES_THRESHOLD_NEAR":      {"tax", "滚动销售额接近当前纳税人类型的体检预警阈值。", "滚动销售额达到规则配置的预警比例。", "可能需要关注纳税人登记及后续申报口径变化。", "请核对滚动销售额，并咨询主管税务机关或专业人员。", []string{"销售台账", "纳税申报表"}},
}

func trustedRisk(row riskSnapshotRow, evidence map[string]any) map[string]any {
	n, ok := riskNarratives[row.RiskCode]
	if !ok {
		n = riskNarrative{"documents", "服务端已记录一项待核实的合规风险事件。", "依据已持久化的合规风险事件生成。", "如不核实，可能影响经营资料完整性或合规事项判断。", "请结合原始资料申请人工复核。", []string{"相关原始资料"}}
	}
	facts := n.facts
	if len(evidence) > 0 {
		facts += " 已确认检测值：" + compactEvidence(evidence) + "。"
	}
	ruleVersion := "builtin-risk-v1"
	if row.RuleVersionID > 0 {
		ruleVersion = fmt.Sprintf("rule-version-%d", row.RuleVersionID)
	}
	title := strings.TrimSpace(row.Summary)
	if title == "" {
		title = "待核实的合规风险"
	}
	return map[string]any{"code": row.RiskCode, "categoryCode": n.category, "title": title,
		"severity": row.Severity, "facts": facts, "basis": n.basis, "impact": n.impact,
		"recommendation": n.recommendation, "requiredMaterials": n.materials,
		"requiresManualReview": row.Severity == "high", "ruleVersion": ruleVersion}
}

func compactEvidence(evidence map[string]any) string {
	keys := make([]string, 0, len(evidence))
	for key := range evidence {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", key, evidence[key]))
	}
	return strings.Join(parts, "，")
}
