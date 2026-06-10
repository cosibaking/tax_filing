package diagnosis

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/model/input/compliancein"
)

const (
	tableDiagnosis  = "xy_compliance_diagnosis"
	tableServicePlan = "xy_service_plan"
)

type sComplianceDiagnosis struct{}

func New() *sComplianceDiagnosis {
	return &sComplianceDiagnosis{}
}

// Preview 计算诊断结果但不落库（访客模式）
func (s *sComplianceDiagnosis) Preview(ctx context.Context, in *compliancein.DiagnosisSubmitInp) (*compliancein.DiagnosisSubmitModel, error) {
	prepared, err := s.prepareDiagnosis(in)
	if err != nil {
		return nil, err
	}
	return &compliancein.DiagnosisSubmitModel{
		Id:               0,
		RecommendedPlan:  prepared.recommendedPlan,
		TaxComparison:    prepared.comparison,
		Reasons:          prepared.reasons,
		AssumptionHints:  prepared.hints,
		RiskLevel:        prepared.riskLevel,
		ComplianceAlerts: prepared.complianceAlerts,
		McnGuidance:      prepared.mcnGuidance,
	}, nil
}

// Submit 提交合规诊断问卷并保存
func (s *sComplianceDiagnosis) Submit(ctx context.Context, in *compliancein.DiagnosisSubmitInp) (*compliancein.DiagnosisSubmitModel, error) {
	prepared, err := s.prepareDiagnosis(in)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	data := g.Map{
		"platforms":            prepared.platformsJson,
		"monthly_income_range": in.MonthlyIncomeRange,
		"annual_cost_estimate": prepared.annualCost,
		"existing_entity":      in.ExistingEntity,
		"has_filed_tax":        in.HasFiledTax,
		"tax_bureau_contact":   boolToInt(in.TaxBureauContact),
		"recommended_plan":     prepared.recommendedPlan,
		"tax_comparison":       prepared.comparisonJson,
		"deleted":              0,
		"create_time":          now,
		"update_time":          now,
	}
	if in.MemberId > 0 {
		data["member_id"] = in.MemberId
	}
	if in.Notes != "" {
		data["notes"] = in.Notes
	}
	if len(in.CostBreakdown) > 0 {
		costJson, encErr := gjson.Encode(in.CostBreakdown)
		if encErr != nil {
			return nil, gerror.Wrap(encErr, "成本明细编码失败")
		}
		data["cost_breakdown"] = string(costJson)
	}

	result, err := g.DB().Model(tableDiagnosis).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "保存诊断记录失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取诊断记录ID失败")
	}

	return &compliancein.DiagnosisSubmitModel{
		Id:               uint64(id),
		RecommendedPlan:  prepared.recommendedPlan,
		TaxComparison:    prepared.comparison,
		Reasons:          prepared.reasons,
		AssumptionHints:  prepared.hints,
		RiskLevel:        prepared.riskLevel,
		ComplianceAlerts: prepared.complianceAlerts,
		McnGuidance:      prepared.mcnGuidance,
	}, nil
}

type preparedDiagnosis struct {
	comparison      compliancein.TaxComparison
	recommendedPlan string
	reasons         []string
	hints           []string
	platformsJson   string
	comparisonJson   string
	annualCost       float64
	riskLevel        string
	complianceAlerts []string
	mcnGuidance      []string
}

func (s *sComplianceDiagnosis) prepareDiagnosis(in *compliancein.DiagnosisSubmitInp) (*preparedDiagnosis, error) {
	annualIncome := AnnualIncomeFromRange(in.MonthlyIncomeRange)
	if annualIncome <= 0 {
		return nil, gerror.New("无效的月收入区间")
	}

	annualCost := in.AnnualCostEstimate
	if annualCost < 0 {
		annualCost = 0
	}
	if annualCost > annualIncome {
		annualCost = annualIncome
	}

	comparison, recommendedPlan, reasons := CompareTaxSchemes(annualIncome, annualCost, in.TaxBureauContact)
	riskLevel, alerts, mcnGuidance := BuildComplianceInsights(in)

	platformsJson, err := gjson.Encode(in.Platforms)
	if err != nil {
		return nil, gerror.Wrap(err, "平台数据编码失败")
	}
	comparisonJson, err := gjson.Encode(comparison)
	if err != nil {
		return nil, gerror.Wrap(err, "税负对比数据编码失败")
	}

	return &preparedDiagnosis{
		comparison:       comparison,
		recommendedPlan:  recommendedPlan,
		reasons:          reasons,
		hints:            BuildAssumptionHints(annualIncome, annualCost),
		platformsJson:    string(platformsJson),
		comparisonJson:   string(comparisonJson),
		annualCost:       annualCost,
		riskLevel:        riskLevel,
		complianceAlerts: alerts,
		mcnGuidance:      mcnGuidance,
	}, nil
}

// SyncGuest 登录后将访客缓存的诊断批量写入数据库
func (s *sComplianceDiagnosis) SyncGuest(ctx context.Context, in *compliancein.DiagnosisSyncInp) (*compliancein.DiagnosisSyncModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.New("会员ID无效")
	}
	if len(in.Items) == 0 {
		return &compliancein.DiagnosisSyncModel{SavedIds: []uint64{}, Count: 0}, nil
	}

	savedIds := make([]uint64, 0, len(in.Items))
	for i := range in.Items {
		item := in.Items[i]
		item.MemberId = in.MemberId
		out, err := s.Submit(ctx, &item)
		if err != nil {
			return nil, gerror.Wrapf(err, "同步第 %d 条诊断失败", i+1)
		}
		savedIds = append(savedIds, out.Id)
	}

	return &compliancein.DiagnosisSyncModel{
		SavedIds: savedIds,
		Count:    len(savedIds),
	}, nil
}

// BindToMember 将匿名诊断记录绑定到当前会员
func (s *sComplianceDiagnosis) BindToMember(ctx context.Context, in *compliancein.DiagnosisBindInp) (*compliancein.DiagnosisBindModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.New("会员ID无效")
	}
	if len(in.DiagnosisIds) == 0 {
		return &compliancein.DiagnosisBindModel{BoundCount: 0}, nil
	}

	now := time.Now().Unix()
	result, err := g.DB().Model(tableDiagnosis).Ctx(ctx).
		WhereIn("id", in.DiagnosisIds).
		Where("member_id", 0).
		Where("deleted", 0).
		Data(g.Map{
			"member_id":   in.MemberId,
			"update_time": now,
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "绑定诊断记录失败")
	}

	affected, _ := result.RowsAffected()
	return &compliancein.DiagnosisBindModel{BoundCount: int(affected)}, nil
}

// GetDetail 获取会员诊断详情
func (s *sComplianceDiagnosis) GetDetail(ctx context.Context, in *compliancein.DiagnosisDetailInp) (*compliancein.DiagnosisDetailModel, error) {
	if in.MemberId == 0 || in.DiagnosisId == 0 {
		return nil, gerror.New("参数无效")
	}

	var row struct {
		MemberId           uint64  `json:"member_id"`
		RecommendedPlan    string  `json:"recommended_plan"`
		TaxComparison      string  `json:"tax_comparison"`
		TaxBureauContact   int     `json:"tax_bureau_contact"`
		MonthlyIncomeRange string  `json:"monthly_income_range"`
		AnnualCostEstimate float64 `json:"annual_cost_estimate"`
		ExistingEntity     string  `json:"existing_entity"`
		HasFiledTax        string  `json:"has_filed_tax"`
		Notes              string  `json:"notes"`
	}

	err := g.DB().Model(tableDiagnosis).Ctx(ctx).
		Where("id", in.DiagnosisId).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "查询诊断详情失败")
	}
	if row.RecommendedPlan == "" {
		return nil, gerror.New("诊断记录不存在")
	}
	if row.MemberId != in.MemberId {
		return nil, gerror.New("无权查看该诊断记录")
	}

	var comparison compliancein.TaxComparison
	if row.TaxComparison != "" {
		_ = gjson.DecodeTo(row.TaxComparison, &comparison)
	}

	annualIncome := comparison.AnnualIncome
	if annualIncome <= 0 {
		annualIncome = AnnualIncomeFromRange(row.MonthlyIncomeRange)
	}
	annualCost := comparison.AnnualCost
	if annualCost <= 0 {
		annualCost = row.AnnualCostEstimate
	}
	_, _, reasons := CompareTaxSchemes(annualIncome, annualCost, row.TaxBureauContact == 1)
	riskLevel, alerts, mcnGuidance := BuildComplianceInsights(&compliancein.DiagnosisSubmitInp{
		MonthlyIncomeRange: row.MonthlyIncomeRange,
		ExistingEntity:     row.ExistingEntity,
		HasFiledTax:        row.HasFiledTax,
		TaxBureauContact:   row.TaxBureauContact == 1,
		Notes:              row.Notes,
	})

	return &compliancein.DiagnosisDetailModel{
		Id:               in.DiagnosisId,
		RecommendedPlan:  row.RecommendedPlan,
		TaxComparison:    comparison,
		Reasons:          reasons,
		AssumptionHints:  BuildAssumptionHints(annualIncome, annualCost),
		RiskLevel:        riskLevel,
		ComplianceAlerts: alerts,
		McnGuidance:      mcnGuidance,
	}, nil
}

// Calculate 税负对比计算
func (s *sComplianceDiagnosis) Calculate(ctx context.Context, in *compliancein.TaxCalculatorInp) (*compliancein.TaxCalculatorModel, error) {
	if in.AnnualIncome <= 0 {
		return nil, gerror.New("年收入必须大于0")
	}

	annualCost := in.AnnualCost
	if annualCost < 0 {
		annualCost = 0
	}
	if annualCost > in.AnnualIncome {
		annualCost = in.AnnualIncome
	}

	taxBureauContact := in.TaxBureauContact
	if in.DiagnosisId > 0 && !taxBureauContact {
		var record struct {
			TaxBureauContact int `json:"tax_bureau_contact"`
		}
		err := g.DB().Model(tableDiagnosis).Ctx(ctx).
			Where("id", in.DiagnosisId).
			Where("deleted", 0).
			Scan(&record)
		if err == nil {
			taxBureauContact = record.TaxBureauContact == 1
		}
	}

	comparison, recommendedPlan, reasons := CompareTaxSchemes(in.AnnualIncome, annualCost, taxBureauContact)

	if in.DiagnosisId > 0 {
		comparisonJson, err := gjson.Encode(comparison)
		if err == nil {
			_, _ = g.DB().Model(tableDiagnosis).Ctx(ctx).
				Where("id", in.DiagnosisId).
				Where("deleted", 0).
				Data(g.Map{
					"recommended_plan": recommendedPlan,
					"tax_comparison":   string(comparisonJson),
				}).Update()
		}
	}

	return &compliancein.TaxCalculatorModel{
		TaxComparison:   comparison,
		RecommendedPlan: recommendedPlan,
		Reasons:         reasons,
		AssumptionHints: BuildAssumptionHints(in.AnnualIncome, annualCost),
	}, nil
}

// ListPlans 获取服务套餐列表
func (s *sComplianceDiagnosis) ListPlans(ctx context.Context) (*compliancein.ServicePlansModel, error) {
	var rows []struct {
		Id           uint64  `json:"id"`
		Name         string  `json:"name"`
		Tier         string  `json:"tier"`
		MonthlyPrice float64 `json:"monthly_price"`
		Features     string  `json:"features"`
		PriceDisplay string  `json:"price_display"`
	}

	err := g.DB().Model(tableServicePlan).Ctx(ctx).
		Where("status", 1).
		OrderAsc("sort").
		OrderAsc("id").
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询服务套餐失败")
	}

	list := make([]compliancein.ServicePlanItem, 0, len(rows))
	for _, row := range rows {
		item := compliancein.ServicePlanItem{
			Id:          row.Id,
			Name:        row.Name,
			Tier:        row.Tier,
			PriceLabel:  row.PriceDisplay,
			Recommended: row.Tier == "advanced",
		}
		if row.MonthlyPrice > 0 {
			price := row.MonthlyPrice
			item.MonthlyPrice = &price
		}
		if row.Features != "" {
			var features []string
			if decErr := gjson.DecodeTo(row.Features, &features); decErr == nil {
				item.Features = features
			}
		}
		if item.Features == nil {
			item.Features = []string{}
		}
		list = append(list, item)
	}

	return &compliancein.ServicePlansModel{List: list}, nil
}

// ListHistory 会员诊断历史
func (s *sComplianceDiagnosis) ListHistory(ctx context.Context, in *compliancein.DiagnosisHistoryInp) (*compliancein.DiagnosisHistoryModel, error) {
	page, pageSize := in.Page, in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	model := g.DB().Model(tableDiagnosis).Ctx(ctx).
		Where("member_id", in.MemberId).
		Where("deleted", 0)

	total, err := model.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询诊断历史失败")
	}

	var rows []struct {
		Id                 uint64      `json:"id"`
		Platforms          string      `json:"platforms"`
		MonthlyIncomeRange string      `json:"monthly_income_range"`
		AnnualCostEstimate float64     `json:"annual_cost_estimate"`
		ExistingEntity     string      `json:"existing_entity"`
		HasFiledTax        string      `json:"has_filed_tax"`
		TaxBureauContact   int         `json:"tax_bureau_contact"`
		RecommendedPlan    string      `json:"recommended_plan"`
		TaxComparison      string      `json:"tax_comparison"`
		CreateTime         int64 `json:"create_time"`
	}

	err = model.Page(page, pageSize).OrderDesc("id").Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询诊断历史失败")
	}

	list := make([]compliancein.DiagnosisHistoryItem, 0, len(rows))
	for _, row := range rows {
		item := compliancein.DiagnosisHistoryItem{
			Id:                 row.Id,
			MonthlyIncomeRange: row.MonthlyIncomeRange,
			AnnualCostEstimate: row.AnnualCostEstimate,
			ExistingEntity:     row.ExistingEntity,
			HasFiledTax:        row.HasFiledTax,
			TaxBureauContact:   row.TaxBureauContact == 1,
			RecommendedPlan:    row.RecommendedPlan,
			Platforms:          []string{},
		}
		if row.Platforms != "" {
			_ = gjson.DecodeTo(row.Platforms, &item.Platforms)
		}
		if row.TaxComparison != "" {
			_ = gjson.DecodeTo(row.TaxComparison, &item.TaxComparison)
		}
		if row.CreateTime > 0 {
			item.CreatedAt = time.Unix(row.CreateTime, 0).Format("2006-01-02 15:04:05")
		}
		list = append(list, item)
	}

	return &compliancein.DiagnosisHistoryModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
