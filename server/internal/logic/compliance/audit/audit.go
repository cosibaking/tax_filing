// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

package audit

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const auditTable = "xy_compliance_audit_log"

type sComplianceAudit struct{}

func New() *sComplianceAudit {
	return &sComplianceAudit{}
}

// WriteAudit 写入合规审计日志（dao 生成前直接使用 g.DB().Model）
func WriteAudit(ctx context.Context, entityType string, entityId uint64, action string, operatorId uint64, operatorType string, before, after interface{}, ip string) error {
	return writeAuditLog(ctx, entityType, entityId, action, operatorId, operatorType, before, after, ip)
}

// WriteAudit 实现 service.IComplianceAudit
func (s *sComplianceAudit) WriteAudit(ctx context.Context, entityType string, entityId uint64, action string, operatorId uint64, operatorType string, before, after interface{}, ip string) error {
	return writeAuditLog(ctx, entityType, entityId, action, operatorId, operatorType, before, after, ip)
}

func writeAuditLog(ctx context.Context, entityType string, entityId uint64, action string, operatorId uint64, operatorType string, before, after interface{}, ip string) error {
	_, err := g.DB().Model(auditTable).Ctx(ctx).Data(g.Map{
		"entity_type":   entityType,
		"entity_id":     entityId,
		"action":        action,
		"operator_id":   operatorId,
		"operator_type": operatorType,
		"before":        before,
		"after":         after,
		"ip":            ip,
	}).Insert()
	if err != nil {
		return gerror.Wrap(err, "write compliance audit log failed")
	}
	return nil
}
