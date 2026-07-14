package document

import (
	"errors"
	"strings"
	"time"
)

const (
	StatusUploaded            = "uploaded"
	StatusManualReview        = "manual_review"
	StatusPendingConfirmation = "pending_confirmation"
	StatusConfirmed           = "confirmed"
)

type BusinessDocument struct {
	ID            uint64         `json:"id"`
	OpcEntityID   uint64         `json:"opcEntityId"`
	MemberID      uint64         `json:"memberId"`
	AttachmentID  uint64         `json:"attachmentId"`
	PeriodKey     string         `json:"periodKey"`
	DocumentType  string         `json:"documentType"`
	ProcessStatus string         `json:"processStatus"`
	FileHash      string         `json:"fileHash"`
	Extracted     map[string]any `json:"extracted"`
	Confidence    float64        `json:"confidence"`
	ConfirmedBy   uint64         `json:"confirmedBy"`
	ConfirmedAt   uint64         `json:"confirmedAt"`
	Deleted       bool           `json:"deleted"`
	DuplicateOfID uint64         `json:"duplicateOfId,omitempty"`
}

func ResolveProcessingStatus(ocrSucceeded bool, _ float64) string {
	if !ocrSucceeded {
		return StatusManualReview
	}
	// 无论置信度高低，AI 字段都必须由用户或顾问确认。
	return StatusPendingConfirmation
}

func FindDuplicate(items []BusinessDocument, opcID uint64, hash string) *BusinessDocument {
	if strings.TrimSpace(hash) == "" {
		return nil
	}
	for i := range items {
		if items[i].OpcEntityID == opcID && items[i].FileHash == hash && !items[i].Deleted {
			copy := items[i]
			return &copy
		}
	}
	return nil
}

func Confirm(item *BusinessDocument, memberID uint64, documentType string, fields map[string]any) error {
	if item == nil || item.MemberID != memberID {
		return errors.New("无权确认该经营资料")
	}
	if item.ProcessStatus != StatusPendingConfirmation && item.ProcessStatus != StatusManualReview {
		return errors.New("当前资料状态不可确认")
	}
	if strings.TrimSpace(documentType) == "" {
		return errors.New("资料类型不能为空")
	}
	item.DocumentType = documentType
	item.Extracted = fields
	item.ProcessStatus = StatusConfirmed
	item.ConfirmedBy = memberID
	item.ConfirmedAt = uint64(time.Now().Unix())
	return nil
}
