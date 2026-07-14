package document

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const tableDocument = "xy_business_document"

type RegisterInput struct {
	OpcEntityID, MemberID, AttachmentID uint64
	PeriodKey, FileHash                 string
}

type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) Register(ctx context.Context, in RegisterInput) (*BusinessDocument, error) {
	if in.OpcEntityID == 0 || in.MemberID == 0 || in.AttachmentID == 0 {
		return nil, errors.New("企业、会员和附件不能为空")
	}
	now := uint64(time.Now().Unix())
	result, err := g.DB().Model(tableDocument).Ctx(ctx).Data(g.Map{
		"opc_entity_id": in.OpcEntityID, "member_id": in.MemberID, "attachment_id": in.AttachmentID,
		"period_key": in.PeriodKey, "document_type": "unknown", "process_status": StatusManualReview,
		"file_hash": in.FileHash, "confidence": 0, "deleted": 0, "create_time": now, "update_time": now,
	}).Insert()
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, in.MemberID, uint64(id))
}

func (s *Service) Get(ctx context.Context, memberID, id uint64) (*BusinessDocument, error) {
	var item BusinessDocument
	var extracted string
	err := g.DB().Model(tableDocument).Ctx(ctx).Fields("*", "CAST(extracted_json AS CHAR) AS extracted_text").Where("id", id).Where("member_id", memberID).Where("deleted", 0).Scan(&item)
	_ = extracted
	if err != nil || item.ID == 0 {
		return nil, err
	}
	return &item, nil
}

func (s *Service) List(ctx context.Context, memberID uint64, periodKey string) ([]BusinessDocument, error) {
	model := g.DB().Model(tableDocument).Ctx(ctx).Where("member_id", memberID).Where("deleted", 0)
	if periodKey != "" {
		model = model.Where("period_key", periodKey)
	}
	var items []BusinessDocument
	if err := model.OrderDesc("create_time").Scan(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) Confirm(ctx context.Context, memberID, id uint64, documentType string, fields map[string]any) (*BusinessDocument, error) {
	item, err := s.Get(ctx, memberID, id)
	if err != nil || item == nil {
		return item, err
	}
	if err := Confirm(item, memberID, documentType, fields); err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}
	_, err = g.DB().Model(tableDocument).Ctx(ctx).Where("id", id).Where("member_id", memberID).Data(g.Map{
		"document_type": item.DocumentType, "process_status": item.ProcessStatus, "extracted_json": string(encoded),
		"confirmed_by": memberID, "confirmed_at": item.ConfirmedAt, "update_time": uint64(time.Now().Unix()),
	}).Update()
	if err != nil {
		return nil, err
	}
	return item, nil
}
