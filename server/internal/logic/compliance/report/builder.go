package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

var categoryDefinitions = []struct{ code, name string }{
	{"business", "经营数据"},
	{"documents", "资料完整度"},
	{"tax", "税务与申报"},
	{"funds", "资金与股东往来"},
	{"employment", "用工与社保"},
	{"annual", "工商与年度事项"},
}

func BuildStructured(in Input) (StructuredReport, error) {
	if !validPeriod(in.PeriodKey) {
		return StructuredReport{}, errors.New("periodKey 必须是真实有效的 YYYY-MM")
	}
	statisticsMissing := len(in.Statistics) == 0
	completenessMissing := len(in.Completeness) == 0
	if trusted, _ := in.Statistics["trustedSnapshot"].(bool); trusted {
		available, _ := in.Statistics["sourceAvailable"].(bool)
		statisticsMissing = !available
	}
	if trusted, _ := in.Completeness["trustedSnapshot"].(bool); trusted {
		count, ok := number(in.Completeness["confirmedDocumentCount"])
		completenessMissing = !ok || count == 0
	}
	rate, rateOK := number(in.Completeness["rate"])
	if math.IsNaN(rate) || math.IsInf(rate, 0) {
		rate, rateOK = 0, false
	}
	if rate < 0 {
		rate = 0
	} else if rate > 100 {
		rate = 100
	}
	if completenessMissing || !rateOK {
		rate = 0
	}

	report := StructuredReport{SchemaVersion: StructuredSchemaVersion, Anomalies: make([]ReportAnomaly, 0)}
	for i, raw := range in.Risks {
		report.Anomalies = append(report.Anomalies, normalizeRisk(raw, i+1))
	}
	sort.SliceStable(report.Anomalies, func(i, j int) bool {
		a, b := report.Anomalies[i], report.Anomalies[j]
		if severityRank(a.Severity) != severityRank(b.Severity) {
			return severityRank(a.Severity) < severityRank(b.Severity)
		}
		if a.DueDate == "" {
			return false
		}
		if b.DueDate == "" {
			return true
		}
		return a.DueDate < b.DueDate
	})
	for _, risk := range report.Anomalies {
		switch risk.Severity {
		case "high":
			report.Summary.HighCount++
		case "low":
			report.Summary.LowCount++
		default:
			report.Summary.MediumCount++
		}
	}

	insufficient := statisticsMissing || completenessMissing || !rateOK
	report.Summary.CompletenessRate = rate
	switch {
	case insufficient:
		report.Summary.DataNotice = "当前经营统计或资料完整度数据不足，结论仅供初步参考，请补充资料后重新体检。"
	case rate < 100:
		report.Summary.DataNotice = "资料尚未完全齐备，请根据缺失清单及时补充。"
	default:
		report.Summary.DataNotice = "资料完整度数据已提供，仍建议结合原始凭证人工复核。"
	}
	switch {
	case report.Summary.HighCount > 0:
		report.Summary.Conclusion = ConclusionUrgent
	case report.Summary.MediumCount > 0 || insufficient:
		report.Summary.Conclusion = ConclusionAttention
	default:
		report.Summary.Conclusion = ConclusionNormal
	}
	report.Categories = buildCategories(report.Anomalies, statisticsMissing, completenessMissing || !rateOK, statisticsMissing && completenessMissing && len(report.Anomalies) == 0)
	return report, nil
}

func validPeriod(period string) bool {
	if len(period) != 7 || period[4] != '-' {
		return false
	}
	t, err := time.Parse("2006-01", period)
	return err == nil && t.Format("2006-01") == period
}

func number(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func normalizeRisk(raw map[string]any, index int) ReportAnomaly {
	severity := strings.ToLower(text(raw["severity"]))
	if severity != "high" && severity != "medium" && severity != "low" {
		severity = "medium"
	}
	category := text(raw["categoryCode"])
	if !validCategory(category) {
		category = "documents"
	}
	code := text(raw["code"])
	if code == "" {
		code = fmt.Sprintf("risk-%d", index)
	}
	title := text(raw["title"])
	if title == "" {
		title = "待核实的合规风险"
	}
	manual, explicit := raw["requiresManualReview"].(bool)
	if !explicit {
		manual = severity == "high"
	}
	return ReportAnomaly{
		Code: code, CategoryCode: category, Title: title, Severity: severity,
		Facts:             fallback(text(raw["facts"]), "现有数据提示可能存在异常，具体事实需结合原始资料核实。"),
		Basis:             fallback(text(raw["basis"]), "依据当前体检规则进行初步识别，适用口径需人工确认。"),
		Impact:            fallback(text(raw["impact"]), "如未及时核实，可能影响申报准确性或合规管理。"),
		Recommendation:    fallback(text(raw["recommendation"]), "建议尽快整理相关凭证，并由专业人员复核处理。"),
		RequiredMaterials: materials(raw["requiredMaterials"]), DueDate: text(raw["dueDate"]),
		RequiresManualReview: manual, RuleVersion: fallback(text(raw["ruleVersion"]), "unknown"),
	}
}

func buildCategories(anomalies []ReportAnomaly, statisticsMissing, documentsMissing, allDataMissing bool) []ReportCategory {
	result := make([]ReportCategory, 0, len(categoryDefinitions))
	for _, def := range categoryDefinitions {
		missing := allDataMissing || def.code == "business" && statisticsMissing || def.code == "documents" && documentsMissing
		count, status := 0, ConclusionNormal
		for _, anomaly := range anomalies {
			if anomaly.CategoryCode == def.code {
				count++
				if severityRank(anomaly.Severity) < severityRank(status) {
					status = anomaly.Severity
				}
			}
		}
		summary := "基于当前数据未识别到异常，仍需以实际资料和人工复核为准。"
		message := summary
		if missing {
			status, summary = "insufficient", "数据不足，暂不能判断该分类是否正常。"
			message = "请补充该分类所需的基础资料后重新检查。"
		} else if count > 0 {
			summary = fmt.Sprintf("识别到 %d 项待处理风险，请按异常清单核实。", count)
			message = summary
		}
		result = append(result, ReportCategory{Code: def.code, Name: def.name, Status: status, Summary: summary, Checks: []ReportCheck{{Code: def.code + "-overview", Name: def.name + "检查", Status: status, Message: message}}})
	}
	return result
}

func severityRank(s string) int {
	switch s {
	case "high", ConclusionUrgent:
		return 0
	case "medium", ConclusionAttention:
		return 1
	case "low":
		return 2
	default:
		return 3
	}
}

func validCategory(code string) bool {
	for _, def := range categoryDefinitions {
		if def.code == code {
			return true
		}
	}
	return false
}

func text(v any) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func fallback(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func materials(v any) []string {
	result := make([]string, 0)
	switch values := v.(type) {
	case []string:
		for _, value := range values {
			if value = strings.TrimSpace(value); value != "" {
				result = append(result, value)
			}
		}
	case []any:
		for _, value := range values {
			if item := text(value); item != "" {
				result = append(result, item)
			}
		}
	}
	return result
}

func RenderContent(period string, structured StructuredReport) string {
	conclusions := map[string]string{ConclusionNormal: "正常", ConclusionAttention: "需关注", ConclusionUrgent: "紧急处理"}
	return fmt.Sprintf("%s 月度经营合规体检\n总体结论：%s\n资料完整度：%.1f%%\n风险分布：高风险 %d 项、中风险 %d 项、低风险 %d 项\n资料提示：%s\n免责声明：本报告基于现有数据和规则自动生成，仅供经营合规参考，重要结论需由专业人员结合原始资料人工复核。",
		period, conclusions[structured.Summary.Conclusion], structured.Summary.CompletenessRate, structured.Summary.HighCount, structured.Summary.MediumCount, structured.Summary.LowCount, structured.Summary.DataNotice)
}
