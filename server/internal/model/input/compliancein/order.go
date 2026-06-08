package compliancein

// OrderCreateInp 创建服务订单入参
type OrderCreateInp struct {
	MemberId    uint64
	PlanId      uint64
	DiagnosisId uint64
}

// OrderCreateModel 创建订单出参
type OrderCreateModel struct {
	OrderId uint64  `json:"orderId"`
	Status  string  `json:"status"`
	Amount  float64 `json:"amount"`
	PlanId  uint64  `json:"planId"`
}

// ConsentInp 合规同意留痕入参
type ConsentInp struct {
	MemberId           uint64
	OrderId            uint64
	Type               string
	DocumentVersion    string
	Acknowledgments    *RiskAcknowledgments
	PlanConfirmed      bool
	Ip                 string
	UserAgent          string
}

// RiskAcknowledgments 风险告知四条确认
type RiskAcknowledgments struct {
	ComplianceNotEvasion bool `json:"complianceNotEvasion"`
	NoAuditGuarantee     bool `json:"noAuditGuarantee"`
	NoFakeInvoice        bool `json:"noFakeInvoice"`
	TaxDependsOnReality  bool `json:"taxDependsOnReality"`
}

// ConsentModel 同意留痕出参
type ConsentModel struct {
	ConsentId uint64 `json:"consentId"`
	Type      string `json:"type"`
	OrderId   uint64 `json:"orderId"`
}

// SignInp 电子签约入参
type SignInp struct {
	MemberId    uint64
	OrderId     uint64
	LegalName   string
	Ip          string
	UserAgent   string
}

// SignModel 签约出参
type SignModel struct {
	OrderId   uint64 `json:"orderId"`
	Status    string `json:"status"`
	SignedAt  string `json:"signedAt"`
	OpcId     uint64 `json:"opcId,omitempty"`
	OpcStatus string `json:"opcStatus,omitempty"`
}

// ActiveOrderModel 当前活跃订单
type ActiveOrderModel struct {
	OrderId     uint64  `json:"orderId"`
	PlanId      uint64  `json:"planId"`
	PlanName    string  `json:"planName,omitempty"`
	PlanTier    string  `json:"planTier,omitempty"`
	DiagnosisId uint64  `json:"diagnosisId,omitempty"`
	Status      string  `json:"status"`
	Amount      float64 `json:"amount,omitempty"`
	SignedAt    string  `json:"signedAt,omitempty"`
	LegalName   string  `json:"legalName,omitempty"`
	CreatedAt   string  `json:"createdAt,omitempty"`
}
