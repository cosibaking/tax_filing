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
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	api "xygo/api/site"
	"xygo/internal/library/contexts"
	"xygo/internal/library/token"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
)

func optionalMemberId(ctx context.Context) uint64 {
	if memberId := contexts.GetMemberId(ctx); memberId > 0 {
		return memberId
	}
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	tokenStr := r.Header.Get("Xy-User-Token")
	if tokenStr == "" {
		return 0
	}
	memberUser, err := token.ParseMember(ctx, tokenStr)
	if err != nil || memberUser == nil {
		return 0
	}
	return memberUser.Id
}

// ComplianceDiagnosis 提交合规诊断问卷
func (c *ControllerV1) ComplianceDiagnosis(ctx context.Context, req *api.ComplianceDiagnosisReq) (res *api.ComplianceDiagnosisRes, err error) {
	out, err := service.ComplianceDiagnosis().Submit(ctx, &compliancein.DiagnosisSubmitInp{
		MemberId:           optionalMemberId(ctx),
		Platforms:          req.Platforms,
		MonthlyIncomeRange: req.MonthlyIncomeRange,
		AnnualCostEstimate: req.AnnualCostEstimate,
		ExistingEntity:     req.ExistingEntity,
		HasFiledTax:        req.HasFiledTax,
		TaxBureauContact:   req.TaxBureauContact,
		Notes:              req.Notes,
		CostBreakdown:      req.CostBreakdown,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceDiagnosisRes{
		Id:              out.Id,
		RecommendedPlan: out.RecommendedPlan,
		TaxComparison:   out.TaxComparison,
		Reasons:         out.Reasons,
		AssumptionHints: out.AssumptionHints,
	}, nil
}

// ComplianceCalculator 税负对比计算
func (c *ControllerV1) ComplianceCalculator(ctx context.Context, req *api.ComplianceCalculatorReq) (res *api.ComplianceCalculatorRes, err error) {
	out, err := service.ComplianceDiagnosis().Calculate(ctx, &compliancein.TaxCalculatorInp{
		AnnualIncome: req.AnnualIncome,
		AnnualCost:   req.AnnualCost,
		DiagnosisId:  req.DiagnosisId,
	})
	if err != nil {
		return nil, err
	}
	return &api.ComplianceCalculatorRes{
		TaxComparison:   out.TaxComparison,
		RecommendedPlan: out.RecommendedPlan,
		Reasons:         out.Reasons,
		AssumptionHints: out.AssumptionHints,
	}, nil
}

// CompliancePlans 服务套餐列表
func (c *ControllerV1) CompliancePlans(ctx context.Context, req *api.CompliancePlansReq) (res *api.CompliancePlansRes, err error) {
	out, err := service.ComplianceDiagnosis().ListPlans(ctx)
	if err != nil {
		return nil, err
	}
	return &api.CompliancePlansRes{List: out.List}, nil
}
