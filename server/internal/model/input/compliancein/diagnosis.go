package compliancein

// TaxComparisonItem 单方案税负对比项
type TaxComparisonItem struct {
	Plan        string  `json:"plan"`
	Label       string  `json:"label"`
	AnnualTax   float64 `json:"annualTax"`
	TaxRate     float64 `json:"taxRate"`
	Recommended bool    `json:"recommended,omitempty"`
	Warning     bool    `json:"warning,omitempty"`
	Note        string  `json:"note,omitempty"`
}

// TaxComparison 税负对比结果
type TaxComparison struct {
	AnnualIncome float64             `json:"annualIncome"`
	AnnualCost   float64             `json:"annualCost"`
	Items        []TaxComparisonItem `json:"items"`
}

// DiagnosisSubmitInp 诊断问卷提交入参
type DiagnosisSubmitInp struct {
	MemberId            uint64
	Platforms           []string
	MonthlyIncomeRange  string
	AnnualCostEstimate  float64
	ExistingEntity      string
	HasFiledTax         string
	TaxBureauContact    bool
	Notes               string
	CostBreakdown       map[string]float64
}

// DiagnosisSubmitModel 诊断提交出参
type DiagnosisSubmitModel struct {
	Id              uint64        `json:"id"`
	RecommendedPlan string        `json:"recommendedPlan"`
	TaxComparison   TaxComparison `json:"taxComparison"`
	Reasons         []string      `json:"reasons,omitempty"`
	AssumptionHints []string      `json:"assumptionHints,omitempty"`
}

// TaxCalculatorInp 税负计算器入参
type TaxCalculatorInp struct {
	AnnualIncome float64
	AnnualCost   float64
	DiagnosisId  uint64
	TaxBureauContact bool
}

// TaxCalculatorModel 税负计算器出参
type TaxCalculatorModel struct {
	TaxComparison   TaxComparison `json:"taxComparison"`
	RecommendedPlan string        `json:"recommendedPlan,omitempty"`
	Reasons         []string      `json:"reasons,omitempty"`
	AssumptionHints []string      `json:"assumptionHints,omitempty"`
}

// ServicePlanItem 服务套餐项
type ServicePlanItem struct {
	Id            uint64   `json:"id"`
	Name          string   `json:"name"`
	Tier          string   `json:"tier"`
	MonthlyPrice  *float64 `json:"monthlyPrice"`
	PriceLabel    string   `json:"priceLabel,omitempty"`
	Features      []string `json:"features"`
	Recommended   bool     `json:"recommended,omitempty"`
}

// ServicePlansModel 套餐列表出参
type ServicePlansModel struct {
	List []ServicePlanItem `json:"list"`
}

// DiagnosisHistoryInp 诊断历史查询入参
type DiagnosisHistoryInp struct {
	MemberId uint64
	Page     int
	PageSize int
}

// DiagnosisHistoryItem 诊断历史项
type DiagnosisHistoryItem struct {
	Id                 uint64        `json:"id"`
	Platforms          []string      `json:"platforms"`
	MonthlyIncomeRange string        `json:"monthlyIncomeRange"`
	AnnualCostEstimate float64       `json:"annualCostEstimate"`
	ExistingEntity     string        `json:"existingEntity"`
	HasFiledTax        string        `json:"hasFiledTax"`
	TaxBureauContact   bool          `json:"taxBureauContact"`
	RecommendedPlan    string        `json:"recommendedPlan"`
	TaxComparison      TaxComparison `json:"taxComparison"`
	CreatedAt          string        `json:"createdAt"`
}

// DiagnosisHistoryModel 诊断历史出参
type DiagnosisHistoryModel struct {
	List     []DiagnosisHistoryItem `json:"list"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
	Total    int                    `json:"total"`
}
