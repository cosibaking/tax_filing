package social

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/consts"
	"xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/shared"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
	"xygo/utility"
)

const (
	tableSocialConsult = "xy_compliance_social_consult"
	tableSocialGuide   = "xy_compliance_social_guide"
	tableServiceOrder  = "xy_service_order"
	tableServicePlan   = "xy_service_plan"

	consultStatusOpen    = "open"
	consultStatusReplied = "replied"
	consultStatusClosed  = "closed"

	consultCategoryFounder  = "founder"
	consultCategoryEmployee = "employee"
	consultCategoryOther    = "other"
)

type sComplianceSocial struct{}

func init() {
	service.RegisterComplianceSocial(New())
}

func New() *sComplianceSocial {
	return &sComplianceSocial{}
}

// GetEmployment 获取用工状态
func (s *sComplianceSocial) GetEmployment(ctx context.Context, memberId uint64) (*compliancein.EmploymentStatusModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, memberId)
	if err != nil {
		return nil, err
	}
	info, err := shared.LoadEmploymentStatus(ctx, opc.Id)
	if err != nil {
		return nil, err
	}
	return &compliancein.EmploymentStatusModel{
		EmploymentStatus:      info.Status,
		EmploymentConfirmedAt: formatUnix(info.ConfirmedAt),
		CanChangeToNoEmployee: true,
	}, nil
}

// SetEmployment 设置用工状态
func (s *sComplianceSocial) SetEmployment(ctx context.Context, in *compliancein.EmploymentStatusInp) (*compliancein.EmploymentStatusModel, error) {
	status := strings.TrimSpace(in.EmploymentStatus)
	if status != shared.EmploymentNoEmployee && status != shared.EmploymentHasEmployee {
		return nil, gerror.New("请选择有效的用工状态")
	}

	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}

	before, _ := shared.LoadEmploymentStatus(ctx, opc.Id)
	if before.Status == shared.EmploymentNoEmployee && status == shared.EmploymentHasEmployee {
		// allowed
	}
	if before.Status == shared.EmploymentHasEmployee && status == shared.EmploymentNoEmployee {
		return nil, gerror.New("已登记有雇员，暂不支持改为无雇员（请先联系顾问）")
	}

	now := utility.NowUnix()
	_, err = g.DB().Model(shared.TableOpcEntity).Ctx(ctx).
		Where("id", opc.Id).
		Data(g.Map{
			"employment_status":       status,
			"employment_confirmed_at": now,
			"update_time":             now,
		}).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新用工状态失败")
	}

	_ = audit.WriteAudit(ctx, "opc_entity", opc.Id, "employment_status_change", in.MemberId, "member",
		g.Map{"employment_status": before.Status},
		g.Map{"employment_status": status},
		in.Ip)

	return &compliancein.EmploymentStatusModel{
		EmploymentStatus:      status,
		EmploymentConfirmedAt: formatUnix(uint64(now)),
		CanChangeToNoEmployee: status != shared.EmploymentHasEmployee,
	}, nil
}

// ListGuides 社保指引列表
func (s *sComplianceSocial) ListGuides(ctx context.Context) (*compliancein.SocialGuideListModel, error) {
	var rows []struct {
		Slug    string `json:"slug"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
	}
	err := g.DB().Model(tableSocialGuide).Ctx(ctx).
		Fields("slug", "title", "summary").
		Where("enabled", 1).
		OrderAsc("sort").
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询社保指引失败")
	}
	list := make([]compliancein.SocialGuideItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, compliancein.SocialGuideItem{
			Slug:    row.Slug,
			Title:   row.Title,
			Summary: row.Summary,
		})
	}
	return &compliancein.SocialGuideListModel{List: list}, nil
}

// GetGuide 社保指引详情
func (s *sComplianceSocial) GetGuide(ctx context.Context, in *compliancein.SocialGuideDetailInp) (*compliancein.SocialGuideDetailModel, error) {
	slug := strings.TrimSpace(in.Slug)
	if slug == "" {
		return nil, gerror.New("请指定指引文章")
	}
	var row struct {
		Slug      string `json:"slug"`
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		ContentMd string `json:"content_md"`
		Audience  string `json:"audience"`
	}
	err := g.DB().Model(tableSocialGuide).Ctx(ctx).
		Where("slug", slug).
		Where("enabled", 1).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "查询社保指引失败")
	}
	if row.Slug == "" {
		return nil, gerror.NewCode(consts.CodeDataNotFound, "指引文章不存在")
	}
	return &compliancein.SocialGuideDetailModel{
		Slug:      row.Slug,
		Title:     row.Title,
		Summary:   row.Summary,
		ContentMd: row.ContentMd,
		Audience:  row.Audience,
	}, nil
}

// CreateConsult 提交社保咨询
func (s *sComplianceSocial) CreateConsult(ctx context.Context, in *compliancein.SocialConsultCreateInp) (*compliancein.SocialConsultCreateModel, error) {
	opc, err := shared.RequireActiveOpc(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	if !s.hasConsultPrivilege(ctx, in.MemberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "社保咨询为进阶版及以上套餐权益，请升级套餐或联系顾问")
	}

	question := strings.TrimSpace(in.Question)
	if len(question) < 5 {
		return nil, gerror.New("请描述您的问题（至少5个字）")
	}
	category := strings.TrimSpace(in.Category)
	if category == "" {
		category = consultCategoryFounder
	}

	now := utility.NowUnix()
	result, err := g.DB().Model(tableSocialConsult).Ctx(ctx).Data(g.Map{
		"member_id":   in.MemberId,
		"opc_id":      opc.Id,
		"category":    category,
		"question":    question,
		"region_code": strings.TrimSpace(in.RegionCode),
		"status":      consultStatusOpen,
		"create_time": now,
		"update_time": now,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "提交咨询失败")
	}
	id, _ := result.LastInsertId()

	_ = audit.WriteAudit(ctx, "social_consult", uint64(id), "social_consult_create", in.MemberId, "member",
		nil, g.Map{"category": category, "question": question}, in.Ip)

	return &compliancein.SocialConsultCreateModel{Id: uint64(id)}, nil
}

// ListMemberConsults 会员咨询列表
func (s *sComplianceSocial) ListMemberConsults(ctx context.Context, in *compliancein.SocialConsultListInp) (*compliancein.SocialConsultListModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	m := g.DB().Model(tableSocialConsult).Ctx(ctx).
		Where("member_id", in.MemberId).
		Where("deleted", 0)

	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计咨询失败")
	}

	var rows []consultRow
	err = m.OrderDesc("create_time").Page(page, pageSize).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询咨询失败")
	}

	list := make([]compliancein.SocialConsultItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, toConsultItem(row))
	}
	return &compliancein.SocialConsultListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// ListAdminConsults 管理端咨询列表
func (s *sComplianceSocial) ListAdminConsults(ctx context.Context, in *compliancein.AdminSocialConsultListInp) (*compliancein.AdminSocialConsultListModel, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	m := g.DB().Model(tableSocialConsult+" c").Ctx(ctx).
		LeftJoin(shared.TableOpcEntity+" o", "o.id = c.opc_id AND o.deleted = 0").
		LeftJoin(shared.TableMember+" m", "m.id = c.member_id").
		Where("c.deleted", 0)

	if status := strings.TrimSpace(in.Status); status != "" {
		m = m.Where("c.status", status)
	}
	if q := strings.TrimSpace(in.Query); q != "" {
		like := "%" + q + "%"
		m = m.Where("(c.question LIKE ? OR m.nickname LIKE ? OR m.mobile LIKE ? OR o.company_name LIKE ?)", like, like, like, like)
	}

	total, err := m.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计咨询失败")
	}

	var rows []adminConsultRow
	err = m.Fields(
		"c.id AS id",
		"c.member_id AS member_id",
		"c.opc_id AS opc_id",
		"c.category AS category",
		"c.question AS question",
		"c.region_code AS region_code",
		"c.status AS status",
		"c.reply AS reply",
		"c.replied_at AS replied_at",
		"c.create_time AS create_time",
		"m.nickname AS member_name",
		"o.company_name AS company_name",
	).OrderDesc("c.create_time").Page(page, pageSize).Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询咨询失败")
	}

	list := make([]compliancein.AdminSocialConsultItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, compliancein.AdminSocialConsultItem{
			Id:          row.Id,
			MemberId:    row.MemberId,
			MemberName:  row.MemberName,
			OpcId:       row.OpcId,
			CompanyName: row.CompanyName,
			Category:    row.Category,
			Question:    row.Question,
			RegionCode:  row.RegionCode,
			Status:      row.Status,
			Reply:       row.Reply,
			RepliedAt:   formatUnix(row.RepliedAt),
			CreatedAt:   formatUnix(row.CreateTime),
		})
	}
	return &compliancein.AdminSocialConsultListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// ReplyConsult 顾问回复或关闭咨询
func (s *sComplianceSocial) ReplyConsult(ctx context.Context, in *compliancein.AdminSocialConsultReplyInp) (*compliancein.AdminSocialConsultReplyModel, error) {
	if in.Id == 0 {
		return nil, gerror.New("请指定咨询工单")
	}

	var row consultRow
	err := g.DB().Model(tableSocialConsult).Ctx(ctx).
		Where("id", in.Id).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "查询咨询失败")
	}
	if row.Id == 0 {
		return nil, gerror.New("咨询工单不存在")
	}
	if row.Status == consultStatusClosed {
		return nil, gerror.New("该咨询已关闭")
	}

	action := strings.TrimSpace(in.Action)
	if action == "" {
		action = "reply"
	}

	now := utility.NowUnix()
	update := g.Map{"update_time": now}
	status := row.Status

	switch action {
	case "close":
		update["status"] = consultStatusClosed
		update["closed_at"] = now
		status = consultStatusClosed
	case "reply":
		reply := strings.TrimSpace(in.Reply)
		if len(reply) < 2 {
			return nil, gerror.New("请填写回复内容")
		}
		update["reply"] = reply
		update["replied_by"] = in.AdminId
		update["replied_at"] = now
		update["status"] = consultStatusReplied
		status = consultStatusReplied
	default:
		return nil, gerror.New("无效操作")
	}

	_, err = g.DB().Model(tableSocialConsult).Ctx(ctx).
		Where("id", in.Id).
		Data(update).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "更新咨询失败")
	}

	_ = audit.WriteAudit(ctx, "social_consult", in.Id, "social_consult_"+action, in.AdminId, "admin",
		g.Map{"status": row.Status},
		update,
		in.Ip)

	return &compliancein.AdminSocialConsultReplyModel{Id: in.Id, Status: status}, nil
}

func (s *sComplianceSocial) hasConsultPrivilege(ctx context.Context, memberId uint64) bool {
	var row struct {
		Tier string `json:"tier"`
	}
	err := g.DB().Model(tableServiceOrder+" o").Ctx(ctx).
		LeftJoin(tableServicePlan+" p", "p.id = o.plan_id").
		Fields("p.tier").
		Where("o.member_id", memberId).
		Where("o.status", shared.OrderStatusActive).
		Where("o.deleted", 0).
		Limit(1).
		Scan(&row)
	if err != nil {
		return false
	}
	return row.Tier == "advanced" || row.Tier == "premium"
}

type consultRow struct {
	Id         uint64 `json:"id" orm:"id"`
	MemberId   uint64 `json:"member_id" orm:"member_id"`
	OpcId      uint64 `json:"opc_id" orm:"opc_id"`
	Category   string `json:"category" orm:"category"`
	Question   string `json:"question" orm:"question"`
	RegionCode string `json:"region_code" orm:"region_code"`
	Status     string `json:"status" orm:"status"`
	Reply      string `json:"reply" orm:"reply"`
	RepliedAt  uint64 `json:"replied_at" orm:"replied_at"`
	CreateTime uint64 `json:"create_time" orm:"create_time"`
}

// adminConsultRow 管理端联表查询行（勿嵌入 consultRow，否则 GoFrame Scan 无法填充内层字段）
type adminConsultRow struct {
	Id          uint64 `json:"id" orm:"id"`
	MemberId    uint64 `json:"member_id" orm:"member_id"`
	OpcId       uint64 `json:"opc_id" orm:"opc_id"`
	Category    string `json:"category" orm:"category"`
	Question    string `json:"question" orm:"question"`
	RegionCode  string `json:"region_code" orm:"region_code"`
	Status      string `json:"status" orm:"status"`
	Reply       string `json:"reply" orm:"reply"`
	RepliedAt   uint64 `json:"replied_at" orm:"replied_at"`
	CreateTime  uint64 `json:"create_time" orm:"create_time"`
	MemberName  string `json:"member_name" orm:"member_name"`
	CompanyName string `json:"company_name" orm:"company_name"`
}

func toConsultItem(row consultRow) compliancein.SocialConsultItem {
	return compliancein.SocialConsultItem{
		Id:         row.Id,
		Category:   row.Category,
		Question:   row.Question,
		RegionCode: row.RegionCode,
		Status:     row.Status,
		Reply:      row.Reply,
		RepliedAt:  formatUnix(row.RepliedAt),
		CreatedAt:  formatUnix(row.CreateTime),
	}
}

func formatUnix(ts uint64) string {
	if ts == 0 {
		return ""
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m-d H:i:s")
}
