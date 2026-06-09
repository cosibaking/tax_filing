package shared

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	EmploymentUnknown    = "unknown"
	EmploymentNoEmployee = "no_employee"
	EmploymentHasEmployee = "has_employee"
)

type EmploymentInfo struct {
	Status      string
	ConfirmedAt uint64
}

func LoadEmploymentStatus(ctx context.Context, opcId uint64) (EmploymentInfo, error) {
	var row struct {
		EmploymentStatus      string `json:"employment_status"`
		EmploymentConfirmedAt uint64 `json:"employment_confirmed_at"`
	}
	err := g.DB().Model(TableOpcEntity).Ctx(ctx).
		Fields("employment_status", "employment_confirmed_at").
		Where("id", opcId).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return EmploymentInfo{}, gerror.Wrap(err, "查询用工状态失败")
	}
	status := row.EmploymentStatus
	if status == "" {
		status = EmploymentUnknown
	}
	return EmploymentInfo{Status: status, ConfirmedAt: row.EmploymentConfirmedAt}, nil
}
