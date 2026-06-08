package audit

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/logic/compliance/shared"
	"xygo/internal/service"
)

func init() {
	service.RegisterComplianceAudit(New())
}

// ExportCSV 导出会员合规留痕（同意记录 + 审计日志）CSV
func (s *sComplianceAudit) ExportCSV(ctx context.Context, memberId uint64) ([]byte, string, error) {
	if memberId == 0 {
		return nil, "", gerror.New("请指定会员ID")
	}

	var consents []struct {
		Id              uint64 `json:"id"`
		OrderId         uint64 `json:"order_id"`
		Type            string `json:"type"`
		DocumentVersion string `json:"document_version"`
		Ip              string `json:"ip"`
		AgreedAt        uint64 `json:"agreed_at"`
	}
	err := g.DB().Model(shared.TableConsent).Ctx(ctx).
		Where("member_id", memberId).
		OrderAsc("id").
		Scan(&consents)
	if err != nil {
		return nil, "", gerror.Wrap(err, "查询同意留痕失败")
	}

	var logs []struct {
		Id           uint64 `json:"id"`
		EntityType   string `json:"entity_type"`
		EntityId     uint64 `json:"entity_id"`
		Action       string `json:"action"`
		OperatorId   uint64 `json:"operator_id"`
		OperatorType string `json:"operator_type"`
		Before       string `json:"before"`
		After        string `json:"after"`
		Ip           string `json:"ip"`
		CreateTime   uint64 `json:"create_time"`
	}
	err = g.DB().Model(shared.TableAuditLog).Ctx(ctx).
		Where("operator_id", memberId).
		Where("operator_type", "member").
		OrderAsc("id").
		Scan(&logs)
	if err != nil {
		return nil, "", gerror.Wrap(err, "查询审计日志失败")
	}

	// 补充与该会员 OPC/订单相关的审计
	var opcIds []uint64
	_ = g.DB().Model("xy_opc_entity").Ctx(ctx).
		Where("member_id", memberId).
		Where("deleted", 0).
		Fields("id").
		Scan(&opcIds)
	if len(opcIds) > 0 {
		var entityLogs []struct {
			Id           uint64 `json:"id"`
			EntityType   string `json:"entity_type"`
			EntityId     uint64 `json:"entity_id"`
			Action       string `json:"action"`
			OperatorId   uint64 `json:"operator_id"`
			OperatorType string `json:"operator_type"`
			Before       string `json:"before"`
			After        string `json:"after"`
			Ip           string `json:"ip"`
			CreateTime   uint64 `json:"create_time"`
		}
		_ = g.DB().Model(shared.TableAuditLog).Ctx(ctx).
			WhereIn("entity_id", opcIds).
			WhereIn("entity_type", []string{"opc_entity", "tax_filing_task"}).
			OrderAsc("id").
			Scan(&entityLogs)
		logs = append(logs, entityLogs...)
	}

	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(buf)

	_ = w.Write([]string{"section", "id", "ref_id", "type", "action", "operator_id", "operator_type", "detail", "ip", "time"})
	for _, c := range consents {
		_ = w.Write([]string{
			"consent",
			strconv.FormatUint(c.Id, 10),
			strconv.FormatUint(c.OrderId, 10),
			c.Type,
			"",
			strconv.FormatUint(memberId, 10),
			"member",
			c.DocumentVersion,
			c.Ip,
			strconv.FormatUint(c.AgreedAt, 10),
		})
	}
	for _, l := range logs {
		detail := l.After
		if detail == "" {
			detail = l.Before
		}
		_ = w.Write([]string{
			"audit",
			strconv.FormatUint(l.Id, 10),
			strconv.FormatUint(l.EntityId, 10),
			l.EntityType,
			l.Action,
			strconv.FormatUint(l.OperatorId, 10),
			l.OperatorType,
			compactJSON(detail),
			l.Ip,
			strconv.FormatUint(l.CreateTime, 10),
		})
	}
	w.Flush()
	if err = w.Error(); err != nil {
		return nil, "", gerror.Wrap(err, "生成CSV失败")
	}

	filename := fmt.Sprintf("compliance_audit_member_%d.csv", memberId)
	return buf.Bytes(), filename, nil
}

func compactJSON(raw string) string {
	if raw == "" {
		return ""
	}
	if gjson.Valid(raw) {
		return raw
	}
	return raw
}
