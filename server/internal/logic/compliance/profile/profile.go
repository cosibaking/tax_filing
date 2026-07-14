package profile

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/shared"
)

const (
	StatusActive   = "active"
	StatusArchived = "archived"
	SourceUser     = "user"
	SourceAdvisor  = "advisor"
	SourceImport   = "import"

	tableProfile = "xy_enterprise_profile"
)

type Data struct {
	Region               string `json:"region"`
	EntityType           string `json:"entityType"`
	TaxpayerType         string `json:"taxpayerType"`
	VATPeriod            string `json:"vatPeriod"`
	EmployeeCount        int    `json:"employeeCount"`
	InvoiceEnabled       bool   `json:"invoiceEnabled"`
	HasRevenue           bool   `json:"hasRevenue"`
	HasPublicBankAccount bool   `json:"hasPublicBankAccount"`
	Complexity           string `json:"complexity"`
}

type Profile struct {
	ID           uint64            `json:"id"`
	OpcEntityID  uint64            `json:"opcEntityId"`
	Version      uint              `json:"version"`
	Data         Data              `json:"data"`
	Labels       map[string]string `json:"labels"`
	Source       string            `json:"source"`
	Status       string            `json:"status"`
	ChangeReason string            `json:"changeReason"`
	ConfirmedBy  uint64            `json:"confirmedBy"`
	ConfirmedAt  uint64            `json:"confirmedAt"`
}

type SaveInput struct {
	MemberID     uint64 `json:"memberId"`
	OpcID        uint64 `json:"opcId"`
	Source       string `json:"source"`
	ChangeReason string `json:"changeReason"`
	ConfirmedBy  uint64 `json:"confirmedBy"`
	Data         Data   `json:"data"`
}

type Repository interface {
	Active(ctx context.Context, opcID uint64) (*Profile, error)
	SaveVersion(ctx context.Context, item Profile) (*Profile, error)
}

type Ownership interface {
	MemberOwnsOpc(ctx context.Context, memberID, opcID uint64) (bool, error)
}

type Service struct {
	repository Repository
	ownership  Ownership
}

func NewService(repository Repository, ownership Ownership) *Service {
	return &Service{repository: repository, ownership: ownership}
}

func New() *Service {
	return NewService(databaseRepository{}, opcOwnership{})
}

func (s *Service) Save(ctx context.Context, in SaveInput) (*Profile, bool, error) {
	if in.Source == "" {
		in.Source = SourceUser
	}
	if err := validate(in); err != nil {
		return nil, false, err
	}
	owned, err := s.ownership.MemberOwnsOpc(ctx, in.MemberID, in.OpcID)
	if err != nil {
		return nil, false, err
	}
	if !owned {
		return nil, false, errors.New("无权修改该企业画像")
	}

	active, err := s.active(ctx, in.OpcID)
	if err != nil {
		return nil, false, err
	}
	if active != nil && reflect.DeepEqual(active.Data, in.Data) {
		return active, false, nil
	}

	version := uint(1)
	if active != nil {
		version = active.Version + 1
	}
	confirmedBy := in.ConfirmedBy
	if confirmedBy == 0 {
		confirmedBy = in.MemberID
	}
	item := Profile{
		OpcEntityID:  in.OpcID,
		Version:      version,
		Data:         in.Data,
		Labels:       labels(in.Data),
		Source:       in.Source,
		Status:       StatusActive,
		ChangeReason: strings.TrimSpace(in.ChangeReason),
		ConfirmedBy:  confirmedBy,
		ConfirmedAt:  uint64(time.Now().Unix()),
	}
	saved, err := s.repository.SaveVersion(ctx, item)
	if err != nil {
		return nil, false, err
	}
	return saved, true, nil
}

func (s *Service) GetForMember(ctx context.Context, memberID, opcID uint64) (*Profile, error) {
	owned, err := s.ownership.MemberOwnsOpc(ctx, memberID, opcID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, errors.New("无权访问该企业画像")
	}
	return s.active(ctx, opcID)
}

func (s *Service) active(ctx context.Context, opcID uint64) (*Profile, error) {
	item, err := s.repository.Active(ctx, opcID)
	if err == nil {
		return item, nil
	}
	if errors.Is(err, sql.ErrNoRows) || errors.Is(gerror.Cause(err), sql.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}

func validate(in SaveInput) error {
	if in.MemberID == 0 || in.OpcID == 0 {
		return errors.New("会员和企业不能为空")
	}
	if in.Source != SourceUser && in.Source != SourceAdvisor && in.Source != SourceImport {
		return errors.New("画像来源无效")
	}
	if in.Source == SourceAdvisor && strings.TrimSpace(in.ChangeReason) == "" {
		return errors.New("顾问修改企业画像必须填写原因")
	}
	if in.Data.Region != "CN-BJ" {
		return errors.New("首期仅支持北京企业")
	}
	if in.Data.EntityType != "one_person_limited_company" || in.Data.TaxpayerType != "small_scale" {
		return errors.New("首期仅支持一人有限公司小规模纳税人")
	}
	if in.Data.VATPeriod != "monthly" && in.Data.VATPeriod != "quarterly" {
		return errors.New("增值税申报周期无效")
	}
	if in.Data.EmployeeCount < 0 || in.Data.EmployeeCount > 5 {
		return errors.New("首期仅支持0至5名员工")
	}
	if in.Data.Complexity != "low" {
		return errors.New("首期仅支持低复杂度经营")
	}
	return nil
}

func labels(data Data) map[string]string {
	employeeBand := "1-5"
	if data.EmployeeCount == 0 {
		employeeBand = "0"
	}
	return map[string]string{
		"region":               data.Region,
		"entityType":           data.EntityType,
		"taxpayerType":         data.TaxpayerType,
		"vatPeriod":            data.VATPeriod,
		"employeeBand":         employeeBand,
		"invoiceEnabled":       fmt.Sprint(data.InvoiceEnabled),
		"hasRevenue":           fmt.Sprint(data.HasRevenue),
		"hasPublicBankAccount": fmt.Sprint(data.HasPublicBankAccount),
		"complexity":           data.Complexity,
	}
}

type opcOwnership struct{}

func (opcOwnership) MemberOwnsOpc(ctx context.Context, memberID, opcID uint64) (bool, error) {
	opc, err := shared.LoadOpcById(ctx, opcID)
	if err != nil {
		return false, err
	}
	return opc.MemberId == memberID, nil
}

type databaseRepository struct{}

type profileRow struct {
	ID           uint64 `orm:"id"`
	OpcEntityID  uint64 `orm:"opc_entity_id"`
	Version      uint   `orm:"version"`
	ProfileJSON  string `orm:"profile_json"`
	Source       string `orm:"source"`
	Status       string `orm:"status"`
	ChangeReason string `orm:"change_reason"`
	ConfirmedBy  uint64 `orm:"confirmed_by"`
	ConfirmedAt  uint64 `orm:"confirmed_at"`
}

func (databaseRepository) Active(ctx context.Context, opcID uint64) (*Profile, error) {
	var row profileRow
	if err := g.DB().Model(tableProfile).Ctx(ctx).
		Where("opc_entity_id", opcID).
		Where("status", StatusActive).
		Where("deleted", 0).
		OrderDesc("version").
		Limit(1).
		Scan(&row); err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}
	return row.profile()
}

func (databaseRepository) SaveVersion(ctx context.Context, item Profile) (*Profile, error) {
	dataJSON, err := json.Marshal(item.Data)
	if err != nil {
		return nil, err
	}
	labelsJSON, err := json.Marshal(item.Labels)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{
		"opc_entity_id": item.OpcEntityID,
		"version":       item.Version,
		"profile_json":  mergeProfileJSON(dataJSON, labelsJSON),
		"source":        item.Source,
		"status":        StatusActive,
		"change_reason": item.ChangeReason,
		"confirmed_by":  item.ConfirmedBy,
		"confirmed_at":  item.ConfirmedAt,
		"deleted":       0,
		"create_time":   uint64(time.Now().Unix()),
		"update_time":   uint64(time.Now().Unix()),
	}
	var id int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(tableProfile).Where("opc_entity_id", item.OpcEntityID).Where("status", StatusActive).
			Data(g.Map{"status": StatusArchived, "update_time": uint64(time.Now().Unix())}).Update(); err != nil {
			return err
		}
		result, err := tx.Model(tableProfile).Data(payload).Insert()
		if err != nil {
			return err
		}
		id, err = result.LastInsertId()
		return err
	})
	if err != nil {
		return nil, err
	}
	item.ID = uint64(id)
	return &item, nil
}

func mergeProfileJSON(dataJSON, labelsJSON []byte) string {
	return fmt.Sprintf(`{"data":%s,"labels":%s}`, dataJSON, labelsJSON)
}

func (row profileRow) profile() (*Profile, error) {
	var envelope struct {
		Data   Data              `json:"data"`
		Labels map[string]string `json:"labels"`
	}
	if err := json.Unmarshal([]byte(row.ProfileJSON), &envelope); err != nil {
		return nil, err
	}
	return &Profile{
		ID: row.ID, OpcEntityID: row.OpcEntityID, Version: row.Version,
		Data: envelope.Data, Labels: envelope.Labels, Source: row.Source, Status: row.Status,
		ChangeReason: row.ChangeReason, ConfirmedBy: row.ConfirmedBy, ConfirmedAt: row.ConfirmedAt,
	}, nil
}
