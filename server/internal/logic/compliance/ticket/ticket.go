package ticket

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type Ticket struct {
	ID           uint64 `json:"id" orm:"id"`
	OpcEntityID  uint64 `json:"opcEntityId" orm:"opc_entity_id"`
	MemberID     uint64 `json:"memberId" orm:"member_id"`
	TicketType   string `json:"ticketType" orm:"ticket_type"`
	Status       string `json:"status" orm:"status"`
	Title        string `json:"title" orm:"title"`
	Description  string `json:"description" orm:"description"`
	TaskID       uint64 `json:"taskId" orm:"task_id"`
	RiskEventID  uint64 `json:"riskEventId" orm:"risk_event_id"`
	AssigneeID   uint64 `json:"assigneeId" orm:"assignee_id"`
	ProviderName string `json:"providerName" orm:"provider_name"`
	CompletedAt  uint64 `json:"completedAt" orm:"completed_at"`
}
type CreateInput struct {
	OpcEntityID, MemberID, TaskID, RiskEventID   uint64
	TicketType, Title, Description, ProviderName string
}

var transitions = map[string]map[string]string{"submitted": {"accept": "accepted", "cancel": "cancelled"}, "accepted": {"start": "processing", "cancel": "cancelled"}, "processing": {"complete": "completed", "request_info": "waiting_member"}, "waiting_member": {"resume": "processing", "cancel": "cancelled"}}

func Transition(status, event string) (string, error) {
	if next := transitions[status][event]; next != "" {
		return next, nil
	}
	return "", errors.New("工单状态流转无效")
}
func ValidateProvider(ticketType, provider string) error {
	if (ticketType == "tax_filing" || ticketType == "bookkeeping") && provider == "" {
		return errors.New("代理记账或申报服务必须指定有资质服务机构")
	}
	return nil
}

type Service struct{}

func New() *Service { return &Service{} }
func (s *Service) Create(ctx context.Context, in CreateInput) (*Ticket, error) {
	if in.OpcEntityID == 0 || in.MemberID == 0 {
		return nil, errors.New("企业和会员不能为空")
	}
	if err := ValidateProvider(in.TicketType, in.ProviderName); err != nil {
		return nil, err
	}
	now := uint64(time.Now().Unix())
	res, err := g.DB().Model("xy_service_ticket").Ctx(ctx).Data(g.Map{"opc_entity_id": in.OpcEntityID, "member_id": in.MemberID, "ticket_type": in.TicketType, "status": "submitted", "title": in.Title, "description": in.Description, "task_id": in.TaskID, "risk_event_id": in.RiskEventID, "provider_name": in.ProviderName, "create_time": now, "update_time": now}).Insert()
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return s.Get(ctx, in.MemberID, uint64(id))
}
func (s *Service) Get(ctx context.Context, memberID, id uint64) (*Ticket, error) {
	var item Ticket
	err := g.DB().Model("xy_service_ticket").Ctx(ctx).Where("id", id).Where("member_id", memberID).Where("deleted", 0).Scan(&item)
	if err != nil || item.ID == 0 {
		return nil, err
	}
	return &item, nil
}
func (s *Service) List(ctx context.Context, memberID uint64) ([]Ticket, error) {
	var items []Ticket
	err := g.DB().Model("xy_service_ticket").Ctx(ctx).Where("member_id", memberID).Where("deleted", 0).OrderDesc("create_time").Scan(&items)
	return items, err
}
func (s *Service) Apply(ctx context.Context, memberID, id uint64, event, note string) (*Ticket, error) {
	item, err := s.Get(ctx, memberID, id)
	if err != nil || item == nil {
		return item, err
	}
	next, err := Transition(item.Status, event)
	if err != nil {
		return nil, err
	}
	now := uint64(time.Now().Unix())
	completed := uint64(0)
	if next == "completed" {
		completed = now
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model("xy_service_ticket").Where("id", id).Data(g.Map{"status": next, "completed_at": completed, "update_time": now}).Update(); e != nil {
			return e
		}
		_, e := tx.Model("xy_service_ticket_record").Data(g.Map{"ticket_id": id, "from_status": item.Status, "to_status": next, "content": note, "operator_type": "member", "operator_id": memberID, "create_time": now}).Insert()
		return e
	})
	if err != nil {
		return nil, err
	}
	item.Status = next
	item.CompletedAt = completed
	return item, nil
}
