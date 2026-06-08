// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package site

import (
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/model/input/compliancein"
)

// ComplianceDiagnosisReq 提交合规诊断问卷
type ComplianceDiagnosisReq struct {
	g.Meta             `path:"/site/compliance/diagnosis" method:"post" tags:"SiteCompliance" summary:"提交合规诊断问卷"`
	Platforms          []string           `json:"platforms" v:"required#请选择直播平台"`
	MonthlyIncomeRange string             `json:"monthlyIncomeRange" v:"required|in:0-2万,2-5万,5-15万,15万+#请选择月收入区间|月收入区间无效"`
	AnnualCostEstimate float64            `json:"annualCostEstimate"`
	ExistingEntity     string             `json:"existingEntity" v:"required|in:none,individual,company,opc,other#请选择现有主体|现有主体类型无效"`
	HasFiledTax        string             `json:"hasFiledTax" v:"required|in:yes,no,unsure#请选择报税状态|报税状态无效"`
	TaxBureauContact   bool               `json:"taxBureauContact"`
	Notes              string             `json:"notes"`
	CostBreakdown      map[string]float64 `json:"costBreakdown"`
}

// ComplianceDiagnosisRes 诊断提交响应
type ComplianceDiagnosisRes struct {
	Id              uint64                    `json:"id"`
	RecommendedPlan string                    `json:"recommendedPlan"`
	TaxComparison   compliancein.TaxComparison `json:"taxComparison"`
	Reasons         []string                  `json:"reasons,omitempty"`
	AssumptionHints []string                  `json:"assumptionHints,omitempty"`
}

// ComplianceCalculatorReq 税负对比计算器
type ComplianceCalculatorReq struct {
	g.Meta       `path:"/site/compliance/calculator" method:"post" tags:"SiteCompliance" summary:"税负对比计算"`
	AnnualIncome float64 `json:"annualIncome" v:"required|min:1#请输入年收入|年收入必须大于0"`
	AnnualCost   float64 `json:"annualCost"`
	DiagnosisId  uint64  `json:"diagnosisId"`
}

// ComplianceCalculatorRes 税负计算响应
type ComplianceCalculatorRes struct {
	TaxComparison   compliancein.TaxComparison `json:"taxComparison"`
	RecommendedPlan string                     `json:"recommendedPlan,omitempty"`
	Reasons         []string                   `json:"reasons,omitempty"`
	AssumptionHints []string                   `json:"assumptionHints,omitempty"`
}

// CompliancePlansReq 服务套餐列表
type CompliancePlansReq struct {
	g.Meta `path:"/site/compliance/plans" method:"get" tags:"SiteCompliance" summary:"服务套餐列表"`
}

// CompliancePlansRes 服务套餐列表响应
type CompliancePlansRes struct {
	List []compliancein.ServicePlanItem `json:"list"`
}
