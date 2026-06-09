package opc

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xygo/internal/consts"
	"xygo/internal/library/compliancecrypto"
	"xygo/internal/library/complianceverify"
	"xygo/internal/library/security"
	"xygo/internal/logic/compliance/audit"
	"xygo/internal/model/input/compliancein"
	"xygo/utility"
)

const (
	tableOpcEntity      = "xy_opc_entity"
	tableOpcProgressLog = "xy_opc_progress_log"
	tableServiceOrder   = "xy_service_order"
	tableServicePlan    = "xy_service_plan"
	tableMember         = "xy_member"

	opcPending         = "pending"
	opcMaterials       = "materials"
	opcMaterialsReview = "materials_review"
	opcRegistering     = "registering"
	opcTax             = "tax"
	opcBank            = "bank"
	opcActive          = "active"

	orderStatusActive = "active"
)

type sComplianceOpc struct{}

func New() *sComplianceOpc {
	return &sComplianceOpc{}
}

// GetSummary OPC 主体摘要
func (s *sComplianceOpc) GetSummary(ctx context.Context, memberId uint64) (*compliancein.OpcSummaryModel, error) {
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !s.hasActiveOrder(ctx, memberId) {
		return &compliancein.OpcSummaryModel{Status: "unsigned"}, nil
	}

	row, err := s.loadOpcByMember(ctx, memberId)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return &compliancein.OpcSummaryModel{Status: opcPending, IsActive: false}, nil
	}

	planTier := s.memberPlanTier(ctx, memberId)
	creditCode := row.CreditCode
	if creditCode != "" {
		creditCode = maskCreditCode(creditCode)
	}

	return &compliancein.OpcSummaryModel{
		OpcId:       row.Id,
		Status:      row.Status,
		CompanyName: row.CompanyName,
		CreditCode:  creditCode,
		PlanTier:    planTier,
		IsActive:    row.Status == opcActive,
	}, nil
}

// GetProgress 进度时间轴
func (s *sComplianceOpc) GetProgress(ctx context.Context, memberId uint64) (*compliancein.OpcProgressModel, error) {
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !s.hasActiveOrder(ctx, memberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}

	row, err := s.loadOpcByMember(ctx, memberId)
	if err != nil {
		return nil, err
	}
	if row == nil {
		planTier, planName, planAmount, signedAt := s.memberPlanInfo(ctx, memberId)
		return &compliancein.OpcProgressModel{
			OpcStatus:        opcPending,
			EstimatedSlaDays: 14,
			Steps:            buildProgressSteps(opcPending, 0),
			PlanTier:         planTier,
			PlanName:         planName,
			PlanAmount:       planAmount,
			SignedAt:         signedAt,
		}, nil
	}

	creditCode := row.CreditCode
	if creditCode != "" {
		creditCode = maskCreditCode(creditCode)
	}

	rejectNote := s.latestRejectNote(ctx, row.Id)
	readonly := &compliancein.MaterialsReadonlySummary{
		LegalPersonName: maskName(row.LegalPersonName),
		Phone:           security.MaskMobile(row.Phone),
	}
	if row.ProposedNames != "" {
		var names []string
		_ = gjson.DecodeTo(row.ProposedNames, &names)
		if len(names) > 0 {
			readonly.ProposedName = names[0]
		}
	}

	bankMasked := ""
	if row.BankAccountEnc != "" {
		if plain, decErr := compliancecrypto.Decrypt(ctx, row.BankAccountEnc); decErr == nil {
			bankMasked = compliancecrypto.MaskBankAccount(plain)
		}
	}

	planTier, planName, planAmount, signedAt := s.memberPlanInfo(ctx, memberId)

	return &compliancein.OpcProgressModel{
		OpcStatus:          row.Status,
		CompanyName:        row.CompanyName,
		CreditCode:         creditCode,
		EstimatedSlaDays:   14,
		Steps:              buildProgressSteps(row.Status, row.MaterialsSubmittedAt),
		MaterialsReadonly:  readonly,
		MaterialsSubmitted: row.MaterialsSubmittedAt > 0,
		RejectNote:         rejectNote,
		BankAccountMasked:  bankMasked,
		PlanTier:           planTier,
		PlanName:           planName,
		PlanAmount:         planAmount,
		SignedAt:           signedAt,
	}, nil
}

// SubmitMaterials 提交注册资料
func (s *sComplianceOpc) SubmitMaterials(ctx context.Context, in *compliancein.MaterialsSubmitInp) (*compliancein.MaterialsSubmitModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !s.hasActiveOrder(ctx, in.MemberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}
	if err := validateMaterials(ctx, in); err != nil {
		return nil, err
	}
	if err := complianceverify.VerifyMaterials(ctx, &complianceverify.MaterialsVerifyInput{
		MemberId:          in.MemberId,
		LegalPersonName:   in.LegalPersonName,
		IdCard:            in.IdCard,
		Phone:             in.Phone,
		Email:             in.Email,
		IdCardFrontFileId: in.IdCardFrontFileId,
		IdCardBackFileId:  in.IdCardBackFileId,
	}); err != nil {
		return nil, err
	}
	for _, fid := range []uint64{in.AddressProofFileId, in.IdCardFrontFileId, in.IdCardBackFileId} {
		if err := validateAttachment(ctx, fid); err != nil {
			return nil, err
		}
	}

	row, err := s.loadOpcByMember(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, gerror.New("OPC主体未创建，请重新签约")
	}
	if row.Status != opcPending && row.Status != opcMaterials {
		return nil, gerror.New("当前状态不允许提交资料")
	}

	encIdCard, err := compliancecrypto.Encrypt(ctx, strings.TrimSpace(in.IdCard))
	if err != nil {
		return nil, gerror.Wrap(err, "身份证加密失败")
	}

	namesJson, err := gjson.Encode(in.ProposedNames)
	if err != nil {
		return nil, gerror.Wrap(err, "公司名称编码失败")
	}

	now := uint64(utility.NowUnix())
	beforeStatus := row.Status
	data := g.Map{
		"proposed_names":         string(namesJson),
		"registered_capital":     in.RegisteredCapital,
		"capital_term_years":     in.CapitalTermYears,
		"business_term_type":     in.BusinessTermType,
		"business_scope":         strings.TrimSpace(in.BusinessScope),
		"register_province":      strings.TrimSpace(in.RegisterProvince),
		"register_city":          strings.TrimSpace(in.RegisterCity),
		"register_district":      strings.TrimSpace(in.RegisterDistrict),
		"register_address":       strings.TrimSpace(in.RegisterAddress),
		"address_proof_file_id":  in.AddressProofFileId,
		"legal_person_name":      strings.TrimSpace(in.LegalPersonName),
		"id_card_encrypted":      encIdCard,
		"id_card_valid_from":     in.IdCardValidFrom,
		"id_card_valid_to":       in.IdCardValidTo,
		"id_card_front_file_id":  in.IdCardFrontFileId,
		"id_card_back_file_id":   in.IdCardBackFileId,
		"ethnicity":              strings.TrimSpace(in.Ethnicity),
		"household_address":      strings.TrimSpace(in.HouseholdAddress),
		"residential_address":    strings.TrimSpace(in.ResidentialAddress),
		"phone":                  strings.TrimSpace(in.Phone),
		"email":                  strings.TrimSpace(in.Email),
		"esign_authorized":       boolToInt(in.EsignAuthorized),
		"status":                 opcMaterialsReview,
		"update_time":            now,
	}
	if in.BusinessTermType == "fixed" && in.BusinessTermEnd != "" {
		data["business_term_end"] = in.BusinessTermEnd
	}
	if row.MaterialsSubmittedAt == 0 {
		data["materials_submitted_at"] = now
	}

	_, err = g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("id", row.Id).
		Where("member_id", in.MemberId).
		Where("deleted", 0).
		Data(data).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "保存资料失败")
	}

	s.writeProgressLog(ctx, row.Id, "materials", "reviewing", "", 0)

	_ = audit.WriteAudit(ctx, "opc_entity", row.Id, "materials.submit", in.MemberId, "member",
		g.Map{"status": beforeStatus}, g.Map{"status": opcMaterialsReview}, in.Ip)

	return &compliancein.MaterialsSubmitModel{
		OpcId:  row.Id,
		Status: opcMaterialsReview,
	}, nil
}

// SubmitBankReceipt 上传开户回执
func (s *sComplianceOpc) SubmitBankReceipt(ctx context.Context, in *compliancein.BankReceiptInp) (*compliancein.BankReceiptModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if in.BankReceiptFileId == 0 {
		return nil, gerror.New("请上传开户回执")
	}
	if err := validateAttachment(ctx, in.BankReceiptFileId); err != nil {
		return nil, err
	}

	row, err := s.loadOpcByMember(ctx, in.MemberId)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}
	if row.Status != opcTax && row.Status != opcBank {
		return nil, gerror.New("当前状态不允许上传开户回执")
	}

	now := uint64(utility.NowUnix())
	update := g.Map{
		"bank_receipt_file_id": in.BankReceiptFileId,
		"update_time":          now,
	}
	if strings.TrimSpace(in.BankName) != "" {
		update["bank_name"] = strings.TrimSpace(in.BankName)
	}

	_, err = g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("id", row.Id).
		Where("member_id", in.MemberId).
		Data(update).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "保存开户回执失败")
	}

	_ = audit.WriteAudit(ctx, "opc_entity", row.Id, "bank.receipt.submit", in.MemberId, "member",
		nil, update, in.Ip)

	return &compliancein.BankReceiptModel{
		OpcId:  row.Id,
		Status: row.Status,
	}, nil
}

// ListTasks 顾问 OPC 任务列表
func (s *sComplianceOpc) ListTasks(ctx context.Context, in *compliancein.OpcTaskListInp) (*compliancein.OpcTaskListModel, error) {
	page, pageSize := in.Page, in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	model := g.DB().Model(tableOpcEntity+" oe").Ctx(ctx).
		LeftJoin(tableMember+" m", "m.id = oe.member_id").
		LeftJoin(tableServiceOrder+" so", "so.member_id = oe.member_id AND so.deleted = 0 AND so.status = '"+orderStatusActive+"'").
		LeftJoin(tableServicePlan+" sp", "sp.id = so.plan_id").
		Where("oe.deleted", 0).
		Where("oe.status != ?", opcPending)

	if in.Status != "" {
		model = model.Where("oe.status", in.Status)
	}
	if q := strings.TrimSpace(in.Query); q != "" {
		like := "%" + q + "%"
		model = model.Where(
			"oe.legal_person_name LIKE ? OR oe.company_name LIKE ? OR m.mobile LIKE ? OR m.nickname LIKE ?",
			like, like, like, like,
		)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务列表失败")
	}

	var rows []struct {
		Id                   uint64 `json:"id"`
		MemberId             uint64 `json:"member_id"`
		ProposedNames        string `json:"proposed_names"`
		LegalPersonName      string `json:"legal_person_name"`
		Status               string `json:"status"`
		MaterialsSubmittedAt uint64 `json:"materials_submitted_at"`
		Nickname             string `json:"nickname"`
		Mobile               string `json:"mobile"`
		Tier                 string `json:"tier"`
	}

	err = model.Fields(
		"oe.id", "oe.member_id", "oe.proposed_names", "oe.legal_person_name",
		"oe.status", "oe.materials_submitted_at", "m.nickname", "m.mobile", "sp.tier",
	).Page(page, pageSize).OrderDesc("oe.materials_submitted_at").OrderDesc("oe.id").Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务列表失败")
	}

	list := make([]compliancein.OpcTaskListItem, 0, len(rows))
	for _, row := range rows {
		proposedName := "未命名"
		var names []string
		if row.ProposedNames != "" {
			_ = gjson.DecodeTo(row.ProposedNames, &names)
			if len(names) > 0 {
				proposedName = names[0]
			}
		}
		memberName := row.Nickname
		if memberName == "" {
			memberName = row.LegalPersonName
		}
		item := compliancein.OpcTaskListItem{
			OpcId:             row.Id,
			MemberId:          row.MemberId,
			MemberName:        memberName,
			MemberPhoneMasked: security.MaskMobile(row.Mobile),
			ProposedName:      proposedName,
			LegalPersonName:   row.LegalPersonName,
			Status:            row.Status,
			StatusLabel:       statusLabel(row.Status),
			PlanTier:          row.Tier,
		}
		if row.MaterialsSubmittedAt > 0 {
			item.MaterialsSubmittedAt = formatUnix(row.MaterialsSubmittedAt)
			item.DaysSinceSubmit = daysSince(row.MaterialsSubmittedAt)
		}
		list = append(list, item)
	}

	return &compliancein.OpcTaskListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// GetTaskDetail 顾问任务详情（脱敏）
func (s *sComplianceOpc) GetTaskDetail(ctx context.Context, in *compliancein.OpcTaskDetailInp) (*compliancein.OpcTaskDetailModel, error) {
	if in.OpcId == 0 {
		return nil, gerror.New("请指定任务ID")
	}

	row, member, err := s.loadOpcDetail(ctx, in.OpcId)
	if err != nil {
		return nil, err
	}

	idCardMasked := ""
	if row.IdCardEncrypted != "" {
		if plain, decErr := compliancecrypto.Decrypt(ctx, row.IdCardEncrypted); decErr == nil {
			idCardMasked = compliancecrypto.MaskIdCard(plain)
		}
	}
	bankMasked := ""
	if row.BankAccountEnc != "" {
		if plain, decErr := compliancecrypto.Decrypt(ctx, row.BankAccountEnc); decErr == nil {
			bankMasked = compliancecrypto.MaskBankAccount(plain)
		}
	}

	var proposedNames []string
	if row.ProposedNames != "" {
		_ = gjson.DecodeTo(row.ProposedNames, &proposedNames)
	}

	logs, _ := s.loadProgressLogs(ctx, in.OpcId)
	planTier := s.memberPlanTier(ctx, row.MemberId)

	return &compliancein.OpcTaskDetailModel{
		OpcId:                row.Id,
		MemberId:             row.MemberId,
		MemberName:           member.Nickname,
		MemberPhone:          security.MaskMobile(member.Mobile),
		PlanTier:             planTier,
		Status:               row.Status,
		StatusLabel:          statusLabel(row.Status),
		RejectNote:           s.latestRejectNote(ctx, in.OpcId),
		MaterialsSubmittedAt: formatUnix(row.MaterialsSubmittedAt),
		Company: compliancein.OpcCompanyDetail{
			ProposedNames:      proposedNames,
			RegisteredCapital:  row.RegisteredCapital,
			CapitalTermYears:   row.CapitalTermYears,
			BusinessTermType:   row.BusinessTermType,
			BusinessTermEnd:    row.BusinessTermEnd,
			BusinessScope:      row.BusinessScope,
			RegisterProvince:   row.RegisterProvince,
			RegisterCity:       row.RegisterCity,
			RegisterDistrict:   row.RegisterDistrict,
			RegisterAddress:    row.RegisterAddress,
			AddressProofFileId: row.AddressProofFileId,
		},
		LegalPerson: compliancein.OpcLegalPersonDetail{
			LegalPersonName:    row.LegalPersonName,
			IdCardMasked:       idCardMasked,
			IdCardValidFrom:    row.IdCardValidFrom,
			IdCardValidTo:      row.IdCardValidTo,
			Ethnicity:          row.Ethnicity,
			HouseholdAddress:   row.HouseholdAddress,
			ResidentialAddress: row.ResidentialAddress,
			Phone:              security.MaskMobile(row.Phone),
			Email:              security.MaskEmail(row.Email),
			IdCardFrontFileId:  row.IdCardFrontFileId,
			IdCardBackFileId:   row.IdCardBackFileId,
			EsignAuthorized:    row.EsignAuthorized == 1,
		},
		Business: compliancein.OpcBusinessResult{
			CompanyName:   row.CompanyName,
			CreditCode:    row.CreditCode,
			EstablishedAt: row.EstablishedAt,
			LicenseFileId: row.LicenseFileId,
		},
		Tax: compliancein.OpcTaxDetail{
			TaxpayerType:   row.TaxpayerType,
			TaxActivatedAt: formatUnix(row.TaxActivatedAt),
		},
		Bank: compliancein.OpcBankDetail{
			BankName:          row.BankName,
			BankAccountMasked: bankMasked,
			BankReceiptFileId: row.BankReceiptFileId,
		},
		AllowedActions: allowedActions(row.Status),
		ProgressLogs:   logs,
	}, nil
}

// AdvanceTask 顾问推进 OPC 状态
func (s *sComplianceOpc) AdvanceTask(ctx context.Context, in *compliancein.OpcTaskActionInp) (*compliancein.OpcTaskActionModel, error) {
	if in.OpcId == 0 {
		return nil, gerror.New("请指定任务ID")
	}
	if in.AdminId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}

	row, _, err := s.loadOpcDetail(ctx, in.OpcId)
	if err != nil {
		return nil, err
	}
	beforeStatus := row.Status
	now := uint64(utility.NowUnix())

	switch in.Action {
	case "approve_materials":
		if row.Status != opcMaterialsReview {
			return nil, gerror.New("当前状态不允许通过资料审核")
		}
		_, err = g.DB().Model(tableOpcEntity).Ctx(ctx).Where("id", in.OpcId).Data(g.Map{
			"status":                opcRegistering,
			"materials_approved_at": now,
			"update_time":           now,
		}).Update()
		if err != nil {
			return nil, gerror.Wrap(err, "更新状态失败")
		}
		s.writeProgressLog(ctx, in.OpcId, "materials", "approved", in.Note, in.AdminId)
		_ = audit.WriteAudit(ctx, "opc_entity", in.OpcId, "materials.approve", in.AdminId, "admin",
			g.Map{"status": beforeStatus}, g.Map{"status": opcRegistering}, in.Ip)

	case "reject_materials":
		if row.Status != opcMaterialsReview {
			return nil, gerror.New("当前状态不允许驳回资料")
		}
		if len([]rune(strings.TrimSpace(in.Note))) < 10 {
			return nil, gerror.New("驳回备注至少10个字")
		}
		_, err = g.DB().Model(tableOpcEntity).Ctx(ctx).Where("id", in.OpcId).Data(g.Map{
			"status":      opcMaterials,
			"update_time": now,
		}).Update()
		if err != nil {
			return nil, gerror.Wrap(err, "更新状态失败")
		}
		s.writeProgressLog(ctx, in.OpcId, "materials", "rejected", in.Note, in.AdminId)
		_ = audit.WriteAudit(ctx, "opc_entity", in.OpcId, "materials.reject", in.AdminId, "admin",
			g.Map{"status": beforeStatus}, g.Map{"status": opcMaterials, "note": in.Note}, in.Ip)

	case "issue_license":
		if row.Status != opcRegistering {
			return nil, gerror.New("当前状态不允许下发执照")
		}
		if strings.TrimSpace(in.CompanyName) == "" {
			return nil, gerror.New("请填写核准公司名称")
		}
		if !validateCreditCode(in.CreditCode) {
			return nil, gerror.New("统一社会信用代码格式无效")
		}
		if in.LicenseFileId == 0 {
			return nil, gerror.New("请上传营业执照")
		}
		if err := validateAttachment(ctx, in.LicenseFileId); err != nil {
			return nil, err
		}
		update := g.Map{
			"status":          opcTax,
			"company_name":    strings.TrimSpace(in.CompanyName),
			"credit_code":     strings.ToUpper(strings.TrimSpace(in.CreditCode)),
			"license_file_id": in.LicenseFileId,
			"update_time":     now,
		}
		if in.EstablishedAt != "" {
			update["established_at"] = in.EstablishedAt
		}
		_, err = g.DB().Model(tableOpcEntity).Ctx(ctx).Where("id", in.OpcId).Data(update).Update()
		if err != nil {
			return nil, gerror.Wrap(err, "更新状态失败")
		}
		s.writeProgressLog(ctx, in.OpcId, "business", "license_issued", in.Note, in.AdminId)
		_ = audit.WriteAudit(ctx, "opc_entity", in.OpcId, "business.license_issued", in.AdminId, "admin",
			g.Map{"status": beforeStatus}, g.Map{"status": opcTax}, in.Ip)

	case "complete_tax":
		if row.Status != opcTax {
			return nil, gerror.New("当前状态不允许标记税务完成")
		}
		if strings.TrimSpace(in.TaxActivatedAt) == "" {
			return nil, gerror.New("请填写电子税务局激活日期")
		}
		taxTs := parseDateToUnix(in.TaxActivatedAt)
		_, err = g.DB().Model(tableOpcEntity).Ctx(ctx).Where("id", in.OpcId).Data(g.Map{
			"status":           opcBank,
			"tax_activated_at": taxTs,
			"update_time":      now,
		}).Update()
		if err != nil {
			return nil, gerror.Wrap(err, "更新状态失败")
		}
		s.writeProgressLog(ctx, in.OpcId, "tax", "completed", in.Note, in.AdminId)
		_ = audit.WriteAudit(ctx, "opc_entity", in.OpcId, "tax.completed", in.AdminId, "admin",
			g.Map{"status": beforeStatus}, g.Map{"status": opcBank}, in.Ip)

	case "complete_bank":
		if row.Status != opcBank {
			return nil, gerror.New("当前状态不允许标记银行开户完成")
		}
		if strings.TrimSpace(in.BankAccount) == "" {
			return nil, gerror.New("请填写对公账号")
		}
		receiptId := in.BankReceiptFileId
		if receiptId == 0 {
			receiptId = row.BankReceiptFileId
		}
		if receiptId == 0 {
			return nil, gerror.New("请上传开户回执")
		}
		encBank, encErr := compliancecrypto.Encrypt(ctx, strings.TrimSpace(in.BankAccount))
		if encErr != nil {
			return nil, gerror.Wrap(encErr, "银行账号加密失败")
		}
		_, err = g.DB().Model(tableOpcEntity).Ctx(ctx).Where("id", in.OpcId).Data(g.Map{
			"status":               opcActive,
			"bank_account_enc":     encBank,
			"bank_receipt_file_id": receiptId,
			"update_time":          now,
		}).Update()
		if err != nil {
			return nil, gerror.Wrap(err, "更新状态失败")
		}
		s.writeProgressLog(ctx, in.OpcId, "bank", "completed", in.Note, in.AdminId)
		s.writeProgressLog(ctx, in.OpcId, "complete", "completed", "", in.AdminId)
		_ = audit.WriteAudit(ctx, "opc_entity", in.OpcId, "bank.completed", in.AdminId, "admin",
			g.Map{"status": beforeStatus}, g.Map{"status": opcActive}, in.Ip)

	default:
		return nil, gerror.New("无效的操作类型")
	}

	newStatus := s.getOpcStatus(ctx, in.OpcId)
	return &compliancein.OpcTaskActionModel{
		OpcId:  in.OpcId,
		Status: newStatus,
	}, nil
}

// ListCustomers 合规客户简单列表
func (s *sComplianceOpc) ListCustomers(ctx context.Context, in *compliancein.ComplianceCustomerListInp) (*compliancein.ComplianceCustomerListModel, error) {
	page, pageSize := in.Page, in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	model := g.DB().Model(tableServiceOrder+" so").Ctx(ctx).
		LeftJoin(tableMember+" m", "m.id = so.member_id").
		LeftJoin(tableServicePlan+" sp", "sp.id = so.plan_id").
		LeftJoin(tableOpcEntity+" oe", "oe.member_id = so.member_id AND oe.deleted = 0").
		Where("so.deleted", 0).
		Where("so.status", orderStatusActive)

	if q := strings.TrimSpace(in.Query); q != "" {
		like := "%" + q + "%"
		model = model.Where("m.nickname LIKE ? OR m.mobile LIKE ?", like, like)
	}
	if in.OpcStatus != "" {
		model = model.Where("oe.status", in.OpcStatus)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询客户列表失败")
	}

	var rows []struct {
		MemberId   uint64 `json:"member_id"`
		Nickname   string `json:"nickname"`
		Mobile     string `json:"mobile"`
		Tier       string `json:"tier"`
		SignedAt   uint64 `json:"signed_at"`
		OpcId      uint64 `json:"opc_id"`
		OpcStatus  string `json:"opc_status"`
		CompanyName string `json:"company_name"`
	}

	err = model.Fields(
		"so.member_id", "m.nickname", "m.mobile", "sp.tier", "so.signed_at",
		"oe.id AS opc_id", "oe.status AS opc_status", "oe.company_name",
	).Page(page, pageSize).OrderDesc("so.signed_at").Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询客户列表失败")
	}

	list := make([]compliancein.ComplianceCustomerItem, 0, len(rows))
	for _, row := range rows {
		item := compliancein.ComplianceCustomerItem{
			MemberId:          row.MemberId,
			MemberName:        row.Nickname,
			MemberPhoneMasked: security.MaskMobile(row.Mobile),
			PlanTier:          row.Tier,
			OrderStatus:       orderStatusActive,
			OpcId:             row.OpcId,
			OpcStatus:         row.OpcStatus,
			OpcStatusLabel:    statusLabel(row.OpcStatus),
			OpcCompanyName:    row.CompanyName,
			SignedAt:          formatUnix(row.SignedAt),
		}
		list = append(list, item)
	}

	return &compliancein.ComplianceCustomerListModel{
		List:     list,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// --- helpers ---

type opcRow struct {
	Id                   uint64  `json:"id"`
	MemberId             uint64  `json:"member_id"`
	ProposedNames        string  `json:"proposed_names"`
	RegisteredCapital    float64 `json:"registered_capital"`
	CapitalTermYears     int     `json:"capital_term_years"`
	BusinessTermType     string  `json:"business_term_type"`
	BusinessTermEnd      string  `json:"business_term_end"`
	BusinessScope        string  `json:"business_scope"`
	RegisterProvince     string  `json:"register_province"`
	RegisterCity         string  `json:"register_city"`
	RegisterDistrict     string  `json:"register_district"`
	RegisterAddress      string  `json:"register_address"`
	AddressProofFileId   uint64  `json:"address_proof_file_id"`
	LegalPersonName      string  `json:"legal_person_name"`
	IdCardEncrypted      string  `json:"id_card_encrypted"`
	IdCardValidFrom      string  `json:"id_card_valid_from"`
	IdCardValidTo        string  `json:"id_card_valid_to"`
	IdCardFrontFileId    uint64  `json:"id_card_front_file_id"`
	IdCardBackFileId     uint64  `json:"id_card_back_file_id"`
	Ethnicity            string  `json:"ethnicity"`
	HouseholdAddress     string  `json:"household_address"`
	ResidentialAddress   string  `json:"residential_address"`
	Phone                string  `json:"phone"`
	Email                string  `json:"email"`
	EsignAuthorized      int     `json:"esign_authorized"`
	CompanyName          string  `json:"company_name"`
	CreditCode           string  `json:"credit_code"`
	EstablishedAt        string  `json:"established_at"`
	LicenseFileId        uint64  `json:"license_file_id"`
	TaxpayerType         string  `json:"taxpayer_type"`
	TaxActivatedAt       uint64  `json:"tax_activated_at"`
	BankName             string  `json:"bank_name"`
	BankAccountEnc       string  `json:"bank_account_enc"`
	BankReceiptFileId    uint64  `json:"bank_receipt_file_id"`
	Status               string  `json:"status"`
	MaterialsSubmittedAt uint64  `json:"materials_submitted_at"`
}

type memberBrief struct {
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
}

func (s *sComplianceOpc) hasActiveOrder(ctx context.Context, memberId uint64) bool {
	count, err := g.DB().Model(tableServiceOrder).Ctx(ctx).
		Where("member_id", memberId).
		Where("status", orderStatusActive).
		Where("deleted", 0).
		Count()
	return err == nil && count > 0
}

func (s *sComplianceOpc) loadOpcByMember(ctx context.Context, memberId uint64) (*opcRow, error) {
	var row opcRow
	err := g.DB().Model(tableOpcEntity).Ctx(ctx).
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

func (s *sComplianceOpc) loadOpcDetail(ctx context.Context, opcId uint64) (*opcRow, memberBrief, error) {
	var row opcRow
	err := g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("id", opcId).
		Where("deleted", 0).
		Scan(&row)
	if err != nil {
		return nil, memberBrief{}, gerror.Wrap(err, "查询OPC主体失败")
	}
	if row.Id == 0 {
		return nil, memberBrief{}, gerror.New("任务不存在")
	}
	var member memberBrief
	_ = g.DB().Model(tableMember).Ctx(ctx).Where("id", row.MemberId).Scan(&member)
	return &row, member, nil
}

func (s *sComplianceOpc) getOpcStatus(ctx context.Context, opcId uint64) string {
	var row struct {
		Status string `json:"status"`
	}
	_ = g.DB().Model(tableOpcEntity).Ctx(ctx).Where("id", opcId).Scan(&row)
	return row.Status
}

func (s *sComplianceOpc) memberPlanTier(ctx context.Context, memberId uint64) string {
	tier, _, _, _ := s.memberPlanInfo(ctx, memberId)
	return tier
}

func (s *sComplianceOpc) memberPlanInfo(ctx context.Context, memberId uint64) (tier, name string, amount float64, signedAt string) {
	var row struct {
		Tier      string  `json:"tier"`
		Name      string  `json:"name"`
		Amount    float64 `json:"amount"`
		SignedAt  uint64  `json:"signed_at"`
	}
	_ = g.DB().Model(tableServiceOrder+" so").Ctx(ctx).
		LeftJoin(tableServicePlan+" sp", "sp.id = so.plan_id").
		Where("so.member_id", memberId).
		Where("so.status", orderStatusActive).
		Where("so.deleted", 0).
		Fields("sp.tier, sp.name, so.amount, so.signed_at").
		OrderDesc("so.id").
		Limit(1).
		Scan(&row)
	return row.Tier, row.Name, row.Amount, formatUnix(row.SignedAt)
}

func (s *sComplianceOpc) writeProgressLog(ctx context.Context, opcId uint64, step, status, note string, operatedBy uint64) {
	_, _ = g.DB().Model(tableOpcProgressLog).Ctx(ctx).Data(g.Map{
		"opc_id":       opcId,
		"step":         step,
		"status":       status,
		"note":         note,
		"operated_by":  operatedBy,
		"create_time":  utility.NowUnix(),
	}).Insert()
}

func (s *sComplianceOpc) loadProgressLogs(ctx context.Context, opcId uint64) ([]compliancein.OpcProgressLogItem, error) {
	var rows []struct {
		Step       string `json:"step"`
		Status     string `json:"status"`
		Note       string `json:"note"`
		OperatedBy uint64 `json:"operated_by"`
		CreateTime int64  `json:"create_time"`
	}
	err := g.DB().Model(tableOpcProgressLog).Ctx(ctx).
		Where("opc_id", opcId).
		OrderAsc("id").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	list := make([]compliancein.OpcProgressLogItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, compliancein.OpcProgressLogItem{
			Step:       r.Step,
			Status:     r.Status,
			Note:       r.Note,
			OperatedBy: r.OperatedBy,
			CreatedAt:  formatUnix(uint64(r.CreateTime)),
		})
	}
	return list, nil
}

func (s *sComplianceOpc) latestRejectNote(ctx context.Context, opcId uint64) string {
	var row struct {
		Note string `json:"note"`
	}
	_ = g.DB().Model(tableOpcProgressLog).Ctx(ctx).
		Where("opc_id", opcId).
		Where("step", "materials").
		Where("status", "rejected").
		OrderDesc("id").
		Limit(1).
		Scan(&row)
	return row.Note
}

func buildProgressSteps(status string, submittedAt uint64) []compliancein.ProgressStepItem {
	steps := []compliancein.ProgressStepItem{
		{Key: "materials", Label: "资料提交"},
		{Key: "business", Label: "工商注册"},
		{Key: "tax", Label: "税务登记"},
		{Key: "bank", Label: "银行开户"},
	}
	statusOrder := map[string]int{
		opcPending: 0, opcMaterials: 0, opcMaterialsReview: 0,
		opcRegistering: 1, opcTax: 2, opcBank: 3, opcActive: 4,
	}
	current := statusOrder[status]
	submittedDate := ""
	if submittedAt > 0 {
		submittedDate = formatUnixDate(submittedAt)
	}
	for i := range steps {
		switch {
		case i < current:
			steps[i].Status = "done"
			if i == 0 && submittedDate != "" {
				steps[i].Date = submittedDate
			}
		case i == current && status != opcActive:
			steps[i].Status = "current"
		default:
			steps[i].Status = "pending"
		}
		if status == opcActive {
			steps[i].Status = "done"
		}
	}
	return steps
}

func allowedActions(status string) []string {
	switch status {
	case opcMaterialsReview:
		return []string{"approve_materials", "reject_materials"}
	case opcRegistering:
		return []string{"issue_license"}
	case opcTax:
		return []string{"complete_tax"}
	case opcBank:
		return []string{"complete_bank"}
	default:
		return []string{}
	}
}

func statusLabel(status string) string {
	labels := map[string]string{
		opcPending:         "待提交资料",
		opcMaterials:       "资料待补正",
		opcMaterialsReview: "资料审核中",
		opcRegistering:     "工商注册中",
		opcTax:             "税务登记中",
		opcBank:            "银行开户中",
		opcActive:          "已激活",
	}
	if l, ok := labels[status]; ok {
		return l
	}
	return status
}

func maskCreditCode(code string) string {
	if len(code) < 10 {
		return code
	}
	return code[:6] + strings.Repeat("*", len(code)-10) + code[len(code)-4:]
}

func formatUnix(ts uint64) string {
	if ts == 0 {
		return ""
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m-d H:i:s")
}

func formatUnixDate(ts uint64) string {
	if ts == 0 {
		return ""
	}
	return utility.UnixToGTime(int64(ts)).Format("Y-m-d")
}

func daysSince(ts uint64) int {
	if ts == 0 {
		return 0
	}
	now := utility.NowUnix()
	diff := now - int64(ts)
	if diff < 0 {
		return 0
	}
	return int(diff / 86400)
}

func parseDateToUnix(dateStr string) uint64 {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return 0
	}
	t, err := gtime.StrToTime(dateStr)
	if err != nil {
		return 0
	}
	return uint64(t.Unix())
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
