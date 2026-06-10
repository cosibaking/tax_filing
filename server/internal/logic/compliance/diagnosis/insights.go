package diagnosis

import (
	"strings"

	"xygo/internal/model/input/compliancein"
)

// BuildComplianceInsights 合规风险评级、漏报提示与 MCN 指引（PRD F-05/F-06/F-70）
func BuildComplianceInsights(in *compliancein.DiagnosisSubmitInp) (riskLevel string, alerts []string, mcnGuidance []string) {
	alerts = make([]string, 0, 4)
	mcnGuidance = make([]string, 0, 4)
	score := 0

	switch in.HasFiledTax {
	case "no":
		score += 3
		alerts = append(alerts,
			"您选择了存在漏报/逾期。在未收到税务机关正式通知前，不建议私下主动补报，应先与专业顾问确认处理路径。",
			"平台收入报送已上线，历史未申报数据存在被比对风险，建议尽快建立规范申报节奏。",
		)
	case "unsure":
		score += 2
		alerts = append(alerts,
			"您不确定历史申报情况。建议梳理近三年的平台结算与个税申报记录，避免「以为平台代扣即等于已报税」的误区。",
		)
	}

	if in.TaxBureauContact {
		score += 3
		alerts = append(alerts,
			"您已收到税务机关或平台的合规通知，建议优先配合提供真实资料，在顾问指导下完成规范申报。",
		)
	}

	if in.ExistingEntity == "none" && in.MonthlyIncomeRange == "15万+" {
		score++
		alerts = append(alerts,
			"高收入且暂无经营主体，个人账户大额流水与平台报送数据比对风险较高，建议优先考虑 OPC 合规落地。",
		)
	}

	notes := strings.ToLower(in.Notes)
	isMcn := strings.Contains(notes, "mcn签约主播") ||
		strings.Contains(notes, "mcn对公") ||
		strings.Contains(notes, "mcn_public") ||
		strings.Contains(notes, "混合结算")

	if isMcn {
		mcnGuidance = append(mcnGuidance,
			"MCN 签约主播须区分「机构分成」与「个人实收」，台账按归属主播部分入账，留存 MCN 协议与结算单。",
			"若仍以个人卡提现为主，OPC 主体可能名存实亡，建议评估结算主体变更与对公合同签署。",
			"签约前请审阅 MCN 合同中的独家条款、主体限制与代扣代缴约定，避免违约或重复计税。",
		)
		if strings.Contains(notes, "mcn对公") || strings.Contains(notes, "mcn_public") {
			mcnGuidance = append(mcnGuidance,
				"当前为 MCN 对公结算模式，OPC 设立后需协调平台/MCN 将结算主体切换至公司对公账户。",
			)
		}
		if strings.Contains(notes, "个人提现") {
			score++
			mcnGuidance = append(mcnGuidance,
				"个人提现模式下收入易与 OPC 台账脱节，务必按月归集平台结算并勾稽银行流水。",
			)
		}
	}

	switch {
	case score >= 3:
		riskLevel = "red"
	case score >= 1:
		riskLevel = "yellow"
	default:
		riskLevel = "green"
	}
	return riskLevel, alerts, mcnGuidance
}
