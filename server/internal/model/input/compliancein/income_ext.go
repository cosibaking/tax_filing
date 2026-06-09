package compliancein

// IncomeOCRPreviewInp 平台流水 OCR 预览入参
type IncomeOCRPreviewInp struct {
	MemberId    uint64
	Platform    string
	OcrText     string
	ImageBytes  []byte
	Filename    string
	AttachmentId uint64
}

// IncomeOCRPreviewModel OCR 预览出参
type IncomeOCRPreviewModel struct {
	Platform       string                   `json:"platform"`
	OcrSource      string                   `json:"ocrSource"`
	RawTextPreview string                   `json:"rawTextPreview,omitempty"`
	Rows           []IncomeImportPreviewRow `json:"rows"`
	ValidCount     int                      `json:"validCount"`
	InvalidCount   int                      `json:"invalidCount"`
	TotalGross     float64                  `json:"totalGross"`
	TotalNet       float64                  `json:"totalNet"`
}

// IncomeOCRImportInp OCR 确认导入入参
type IncomeOCRImportInp struct {
	MemberId     uint64
	Platform     string
	OcrText      string
	ImageBytes   []byte
	Filename     string
	AttachmentId uint64
	Ip           string
}

// IncomeConsistencyInp 收入一致性比对入参
type IncomeConsistencyInp struct {
	MemberId uint64
	Month    string
}

// IncomeConsistencyPlatformItem 分平台比对项
type IncomeConsistencyPlatformItem struct {
	Platform         string  `json:"platform"`
	LedgerNet        float64 `json:"ledgerNet"`
	PlatformReported float64 `json:"platformReported"`
	Variance         float64 `json:"variance"`
	Status           string  `json:"status"`
}

// IncomeConsistencyModel 收入一致性比对出参
type IncomeConsistencyModel struct {
	Month                  string                          `json:"month"`
	LedgerGross            float64                         `json:"ledgerGross"`
	LedgerNet              float64                         `json:"ledgerNet"`
	BankInflow             float64                         `json:"bankInflow"`
	BankMatched            float64                         `json:"bankMatched"`
	BankUnmatched          float64                         `json:"bankUnmatched"`
	PlatformReported       float64                         `json:"platformReported"`
	VarianceLedgerBank     float64                         `json:"varianceLedgerBank"`
	VarianceLedgerPlatform float64                         `json:"varianceLedgerPlatform"`
	VarianceRate           float64                         `json:"varianceRate"`
	Status                 string                          `json:"status"`
	IncomeMatchPassed      bool                            `json:"incomeMatchPassed"`
	Hints                  []string                        `json:"hints,omitempty"`
	Platforms              []IncomeConsistencyPlatformItem `json:"platforms,omitempty"`
}
