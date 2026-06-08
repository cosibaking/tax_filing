package order

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/utility"
	"xygo/internal/consts"
	"xygo/internal/logic/compliance/audit"
	"xygo/internal/model/input/compliancein"
)

const (
	tableServiceOrder   = "xy_service_order"
	tableServicePlan    = "xy_service_plan"
	tableConsent        = "xy_compliance_consent"
	tableOpcEntity      = "xy_opc_entity"

	orderStatusPending   = "pending"
	orderStatusActive    = "active"
	orderStatusCancelled = "cancelled"

	consentRiskDisclosure = "risk_disclosure"
	consentPlanConfirm    = "plan_confirm"
	consentContractSign   = "contract_sign"

	opcStatusPending = "pending"
)

type sComplianceOrder struct{}

func New() *sComplianceOrder {
	return &sComplianceOrder{}
}

// CreateOrder 创建服务订单（status=pending）
func (s *sComplianceOrder) CreateOrder(ctx context.Context, in *compliancein.OrderCreateInp) (*compliancein.OrderCreateModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if in.PlanId == 0 {
		return nil, gerror.New("请选择服务套餐")
	}

	var plan struct {
		Id           uint64  `json:"id"`
		MonthlyPrice float64 `json:"monthly_price"`
		Status       int     `json:"status"`
	}
	err := g.DB().Model(tableServicePlan).Ctx(ctx).
		Where("id", in.PlanId).
		Where("status", 1).
		Scan(&plan)
	if err != nil {
		return nil, gerror.Wrap(err, "查询套餐失败")
	}
	if plan.Id == 0 {
		return nil, gerror.New("套餐不存在或已下架")
	}

	// 已有待处理或生效订单则返回已有 pending 订单
	var existing struct {
		Id     uint64  `json:"id"`
		Status string  `json:"status"`
		Amount float64 `json:"amount"`
		PlanId uint64  `json:"plan_id"`
	}
	_ = g.DB().Model(tableServiceOrder).Ctx(ctx).
		Where("member_id", in.MemberId).
		Where("deleted", 0).
		WhereIn("status", []string{orderStatusPending, orderStatusActive}).
		OrderDesc("id").
		Scan(&existing)
	if existing.Id > 0 {
		if existing.Status == orderStatusActive {
			return nil, gerror.New("您已有生效中的服务订单")
		}
		return &compliancein.OrderCreateModel{
			OrderId: existing.Id,
			Status:  existing.Status,
			Amount:  existing.Amount,
			PlanId:  existing.PlanId,
		}, nil
	}

	now := uint64(utility.NowUnix())
	data := g.Map{
		"member_id":    in.MemberId,
		"plan_id":      in.PlanId,
		"status":       orderStatusPending,
		"amount":       plan.MonthlyPrice,
		"deleted":      0,
		"create_time":  now,
		"update_time":  now,
	}
	if in.DiagnosisId > 0 {
		data["diagnosis_id"] = in.DiagnosisId
	}

	result, err := g.DB().Model(tableServiceOrder).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "创建订单失败")
	}
	orderId, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.Wrap(err, "获取订单ID失败")
	}

	_ = audit.WriteAudit(ctx, "service_order", uint64(orderId), "order.create", in.MemberId, "member", nil, g.Map{
		"planId": in.PlanId,
		"status": orderStatusPending,
	}, "")

	return &compliancein.OrderCreateModel{
		OrderId: uint64(orderId),
		Status:  orderStatusPending,
		Amount:  plan.MonthlyPrice,
		PlanId:  in.PlanId,
	}, nil
}

// RecordConsent 风险告知 / 方案确认留痕
func (s *sComplianceOrder) RecordConsent(ctx context.Context, in *compliancein.ConsentInp) (*compliancein.ConsentModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if in.OrderId == 0 {
		return nil, gerror.New("请指定订单")
	}

	order, err := s.loadMemberOrder(ctx, in.MemberId, in.OrderId)
	if err != nil {
		return nil, err
	}
	if order.Status != orderStatusPending {
		return nil, gerror.New("当前订单状态不允许此操作")
	}

	consentType := strings.TrimSpace(in.Type)
	switch consentType {
	case consentRiskDisclosure:
		if in.Acknowledgments == nil ||
			!in.Acknowledgments.ComplianceNotEvasion ||
			!in.Acknowledgments.NoAuditGuarantee ||
			!in.Acknowledgments.NoFakeInvoice ||
			!in.Acknowledgments.TaxDependsOnReality {
			return nil, gerror.New("请勾选全部风险告知条款")
		}
	case consentPlanConfirm:
		if !in.PlanConfirmed {
			return nil, gerror.New("请确认服务方案")
		}
	default:
		return nil, gerror.New("无效的同意类型")
	}

	// 幂等：已存在则直接返回
	var existing struct {
		Id uint64 `json:"id"`
	}
	_ = g.DB().Model(tableConsent).Ctx(ctx).
		Where("order_id", in.OrderId).
		Where("member_id", in.MemberId).
		Where("type", consentType).
		Scan(&existing)
	if existing.Id > 0 {
		return &compliancein.ConsentModel{
			ConsentId: existing.Id,
			Type:      consentType,
			OrderId:   in.OrderId,
		}, nil
	}

	now := uint64(utility.NowUnix())
	docVersion := in.DocumentVersion
	if docVersion == "" {
		docVersion = "v1.0"
	}
	result, err := g.DB().Model(tableConsent).Ctx(ctx).Data(g.Map{
		"member_id":         in.MemberId,
		"order_id":          in.OrderId,
		"type":              consentType,
		"document_version":  docVersion,
		"ip":                in.Ip,
		"user_agent":        in.UserAgent,
		"agreed_at":         now,
		"create_time":       now,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "保存同意记录失败")
	}
	consentId, _ := result.LastInsertId()

	after := g.Map{"type": consentType, "documentVersion": docVersion}
	if consentType == consentRiskDisclosure && in.Acknowledgments != nil {
		after["acknowledgments"] = in.Acknowledgments
	}
	if consentType == consentPlanConfirm {
		after["planConfirmed"] = in.PlanConfirmed
	}
	_ = audit.WriteAudit(ctx, "service_order", in.OrderId, "consent."+consentType, in.MemberId, "member", nil, after, in.Ip)

	return &compliancein.ConsentModel{
		ConsentId: uint64(consentId),
		Type:      consentType,
		OrderId:   in.OrderId,
	}, nil
}

// SignContract 电子签约（姓名确认 MVP），激活订单并创建 OpcEntity
func (s *sComplianceOrder) SignContract(ctx context.Context, in *compliancein.SignInp) (*compliancein.SignModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if in.OrderId == 0 {
		return nil, gerror.New("请指定订单")
	}
	legalName := strings.TrimSpace(in.LegalName)
	if len(legalName) < 2 || len(legalName) > 20 {
		return nil, gerror.New("签约姓名须为2-20个字符")
	}

	order, err := s.loadMemberOrder(ctx, in.MemberId, in.OrderId)
	if err != nil {
		return nil, err
	}
	if order.Status == orderStatusActive {
		opcId, opcStatus := s.findOpcByMember(ctx, in.MemberId)
		return &compliancein.SignModel{
			OrderId:   in.OrderId,
			Status:    orderStatusActive,
			SignedAt:  formatUnixTime(order.SignedAt),
			OpcId:     opcId,
			OpcStatus: opcStatus,
		}, nil
	}
	if order.Status != orderStatusPending {
		return nil, gerror.New("当前订单状态不允许签约")
	}

	if !s.hasConsent(ctx, in.OrderId, consentRiskDisclosure) {
		return nil, gerror.New("请先完成风险告知确认")
	}
	if !s.hasConsent(ctx, in.OrderId, consentPlanConfirm) {
		return nil, gerror.New("请先完成方案确认")
	}

	now := uint64(utility.NowUnix())

	// 签约留痕
	var signConsent struct{ Id uint64 `json:"id"` }
	_ = g.DB().Model(tableConsent).Ctx(ctx).
		Where("order_id", in.OrderId).
		Where("type", consentContractSign).
		Scan(&signConsent)
	if signConsent.Id == 0 {
		_, _ = g.DB().Model(tableConsent).Ctx(ctx).Data(g.Map{
			"member_id":        in.MemberId,
			"order_id":         in.OrderId,
			"type":             consentContractSign,
			"document_version": "v1.0",
			"ip":               in.Ip,
			"user_agent":       in.UserAgent,
			"agreed_at":        now,
			"create_time":      now,
		}).Insert()
		_ = audit.WriteAudit(ctx, "service_order", in.OrderId, "consent."+consentContractSign, in.MemberId, "member", nil,
			g.Map{"type": consentContractSign, "legalName": legalName}, in.Ip)
	}

	_, err = g.DB().Model(tableServiceOrder).Ctx(ctx).
		Where("id", in.OrderId).
		Where("member_id", in.MemberId).
		Where("deleted", 0).
		Data(g.Map{
			"status":      orderStatusActive,
			"signed_at":   now,
			"legal_name":  legalName,
			"update_time": now,
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新订单状态失败")
	}

	opcId, err := s.ensureOpcEntity(ctx, in.MemberId, now)
	if err != nil {
		return nil, err
	}

	_ = audit.WriteAudit(ctx, "service_order", in.OrderId, "contract.sign", in.MemberId, "member", g.Map{
		"status": orderStatusPending,
	}, g.Map{
		"status":    orderStatusActive,
		"legalName": legalName,
		"opcId":     opcId,
	}, in.Ip)

	return &compliancein.SignModel{
		OrderId:   in.OrderId,
		Status:    orderStatusActive,
		SignedAt:  formatUnixTime(now),
		OpcId:     opcId,
		OpcStatus: opcStatusPending,
	}, nil
}

func (s *sComplianceOrder) loadMemberOrder(ctx context.Context, memberId, orderId uint64) (*struct {
	Id       uint64 `json:"id"`
	Status   string `json:"status"`
	SignedAt uint64 `json:"signed_at"`
}, error) {
	var order struct {
		Id       uint64 `json:"id"`
		Status   string `json:"status"`
		SignedAt uint64 `json:"signed_at"`
	}
	err := g.DB().Model(tableServiceOrder).Ctx(ctx).
		Where("id", orderId).
		Where("member_id", memberId).
		Where("deleted", 0).
		Scan(&order)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}
	if order.Id == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "订单不存在")
	}
	return &order, nil
}

func (s *sComplianceOrder) hasConsent(ctx context.Context, orderId uint64, consentType string) bool {
	count, err := g.DB().Model(tableConsent).Ctx(ctx).
		Where("order_id", orderId).
		Where("type", consentType).
		Count()
	return err == nil && count > 0
}

func (s *sComplianceOrder) ensureOpcEntity(ctx context.Context, memberId uint64, now uint64) (uint64, error) {
	var existing struct {
		Id     uint64 `json:"id"`
		Status string `json:"status"`
	}
	_ = g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("member_id", memberId).
		Where("deleted", 0).
		Scan(&existing)
	if existing.Id > 0 {
		return existing.Id, nil
	}

	result, err := g.DB().Model(tableOpcEntity).Ctx(ctx).Data(g.Map{
		"member_id":    memberId,
		"status":       opcStatusPending,
		"taxpayer_type": "small_scale",
		"deleted":      0,
		"create_time":  now,
		"update_time":  now,
	}).Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "创建OPC主体失败")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "获取OPC主体ID失败")
	}
	return uint64(id), nil
}

func (s *sComplianceOrder) findOpcByMember(ctx context.Context, memberId uint64) (uint64, string) {
	var row struct {
		Id     uint64 `json:"id"`
		Status string `json:"status"`
	}
	_ = g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("member_id", memberId).
		Where("deleted", 0).
		Scan(&row)
	return row.Id, row.Status
}

// GetActiveOrder 获取当前会员的待处理或生效订单
func (s *sComplianceOrder) GetActiveOrder(ctx context.Context, memberId uint64) (*compliancein.ActiveOrderModel, error) {
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}

	var row struct {
		Id          uint64  `json:"id"`
		PlanId      uint64  `json:"plan_id"`
		DiagnosisId uint64  `json:"diagnosis_id"`
		Status      string  `json:"status"`
		Amount      float64 `json:"amount"`
		SignedAt    uint64  `json:"signed_at"`
		LegalName   string  `json:"legal_name"`
		CreateTime  uint64  `json:"create_time"`
		PlanName    string  `json:"plan_name"`
		PlanTier    string  `json:"plan_tier"`
	}
	err := g.DB().Model(tableServiceOrder+" o").Ctx(ctx).
		LeftJoin(tableServicePlan+" p", "p.id = o.plan_id").
		Fields("o.id", "o.plan_id", "o.diagnosis_id", "o.status", "o.amount", "o.signed_at", "o.legal_name", "o.create_time", "p.name AS plan_name", "p.tier AS plan_tier").
		Where("o.member_id", memberId).
		Where("o.deleted", 0).
		WhereIn("o.status", []string{orderStatusPending, orderStatusActive}).
		OrderDesc("o.id").
		Limit(1).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}
	if row.Id == 0 {
		return nil, nil
	}

	out := &compliancein.ActiveOrderModel{
		OrderId:     row.Id,
		PlanId:      row.PlanId,
		DiagnosisId: row.DiagnosisId,
		Status:      row.Status,
		Amount:      row.Amount,
		PlanName:    row.PlanName,
		PlanTier:    row.PlanTier,
		LegalName:   row.LegalName,
	}
	if row.SignedAt > 0 {
		out.SignedAt = formatUnixTime(row.SignedAt)
	}
	if row.CreateTime > 0 {
		out.CreatedAt = formatUnixTime(row.CreateTime)
	}
	return out, nil
}

func formatUnixTime(ts uint64) string {
	if ts == 0 {
		return ""
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m-d H:i:s")
}

