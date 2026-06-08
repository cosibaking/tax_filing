package shared

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/consts"
)

const (
	TableOpcEntity    = "xy_opc_entity"
	TableServiceOrder = "xy_service_order"
	TableMember       = "xy_member"

	OpcStatusActive   = "active"
	OrderStatusActive = "active"
)

type OpcBrief struct {
	Id          uint64 `json:"id"`
	MemberId    uint64 `json:"member_id"`
	CompanyName string `json:"company_name"`
	Status      string `json:"status"`
}

func HasActiveOrder(ctx context.Context, memberId uint64) bool {
	count, err := g.DB().Model(TableServiceOrder).Ctx(ctx).
		Where("member_id", memberId).
		Where("status", OrderStatusActive).
		Where("deleted", 0).
		Count()
	return err == nil && count > 0
}

func LoadOpcByMember(ctx context.Context, memberId uint64) (*OpcBrief, error) {
	var row OpcBrief
	err := g.DB().Model(TableOpcEntity).Ctx(ctx).
		Where("member_id", memberId).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "查询OPC主体失败")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

func RequireActiveOpc(ctx context.Context, memberId uint64) (*OpcBrief, error) {
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !HasActiveOrder(ctx, memberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}
	row, err := LoadOpcByMember(ctx, memberId)
	if err != nil {
		return nil, err
	}
	if row == nil || row.Status != OpcStatusActive {
		return nil, gerror.NewCode(consts.CodeNoPermission, "OPC主体尚未激活")
	}
	return row, nil
}

func LoadOpcById(ctx context.Context, opcId uint64) (*OpcBrief, error) {
	var row OpcBrief
	err := g.DB().Model(TableOpcEntity).Ctx(ctx).
		Where("id", opcId).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "查询OPC主体失败")
	}
	if row.Id == 0 {
		return nil, gerror.New("OPC主体不存在")
	}
	return &row, nil
}
