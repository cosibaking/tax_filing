package opc

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/consts"
	"xygo/internal/library/attachmentaccess"
	"xygo/internal/library/compliancecrypto"
	"xygo/internal/library/security"
	"xygo/internal/model/input/compliancein"
	"xygo/internal/service"
)

const (
	materialsSectionCompany     = "company"
	materialsSectionLegal       = "legal"
	materialsSectionAttachments = "attachments"
)

var allowedMaterialsSections = map[string]bool{
	materialsSectionCompany:     true,
	materialsSectionLegal:       true,
	materialsSectionAttachments: true,
}

// GetMaterialsOverview 已提交资料概览（按 OPC 实体卡片）
func (s *sComplianceOpc) GetMaterialsOverview(ctx context.Context, memberId uint64) (*compliancein.MaterialsOverviewModel, error) {
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !s.hasActiveOrder(ctx, memberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}

	rows, err := s.loadSubmittedOpcByMember(ctx, memberId)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &compliancein.MaterialsOverviewModel{Entities: []compliancein.MaterialsEntityOverview{}}, nil
	}

	entities := make([]compliancein.MaterialsEntityOverview, 0, len(rows))
	for _, row := range rows {
		status, statusLabel := materialsItemStatus(row)
		entities = append(entities, compliancein.MaterialsEntityOverview{
			OpcId:              row.Id,
			CompanyNameMasked:  companyNameMasked(row),
			LegalPersonSummary: legalOverviewSummary(row),
			Status:             status,
			StatusLabel:        statusLabel,
		})
	}
	return &compliancein.MaterialsOverviewModel{Entities: entities}, nil
}

// GetMaterialsDetail 单个 OPC 实体完整资料（脱敏）
func (s *sComplianceOpc) GetMaterialsDetail(ctx context.Context, memberId, opcId uint64) (*compliancein.MaterialsEntityDetailModel, error) {
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if opcId == 0 {
		return nil, gerror.New("无效的 OPC 主体")
	}
	if !s.hasActiveOrder(ctx, memberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}

	row, err := s.loadOpcByMemberAndId(ctx, memberId, opcId)
	if err != nil {
		return nil, err
	}
	if row == nil || row.MaterialsSubmittedAt == 0 {
		return nil, gerror.New("尚未提交注册资料")
	}

	status, statusLabel := materialsItemStatus(row)
	return &compliancein.MaterialsEntityDetailModel{
		OpcId:       row.Id,
		Status:      status,
		StatusLabel: statusLabel,
		Sections:    buildMaterialsDetailSections(ctx, row),
	}, nil
}

// GetMaterialsSection 资料分组详情（脱敏）
func (s *sComplianceOpc) GetMaterialsSection(ctx context.Context, memberId, opcId uint64, section string) (*compliancein.MaterialsSectionDetailModel, error) {
	if memberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !s.hasActiveOrder(ctx, memberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}
	section = strings.TrimSpace(section)
	if !allowedMaterialsSections[section] {
		return nil, gerror.New("无效的资料分组")
	}

	row, err := s.loadOpcForMaterials(ctx, memberId, opcId)
	if err != nil {
		return nil, err
	}
	if row == nil || row.MaterialsSubmittedAt == 0 {
		return nil, gerror.New("尚未提交注册资料")
	}

	status, statusLabel := materialsItemStatus(row)
	title := materialsSectionTitle(section)
	fields := buildMaskedSectionFields(ctx, row, section)

	return &compliancein.MaterialsSectionDetailModel{
		Key:         section,
		Title:       title,
		Status:      status,
		StatusLabel: statusLabel,
		Fields:      fields,
	}, nil
}

// RevealMaterialsAll 密码验证后返回完整资料（全部分组）
func (s *sComplianceOpc) RevealMaterialsAll(ctx context.Context, in *compliancein.MaterialsRevealInp) (*compliancein.MaterialsRevealAllModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !s.hasActiveOrder(ctx, in.MemberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}
	if strings.TrimSpace(in.Password) == "" {
		return nil, gerror.New("请输入密码")
	}
	if err := service.MemberUser().VerifyPassword(ctx, in.MemberId, in.Password); err != nil {
		return nil, err
	}

	row, err := s.loadOpcForMaterials(ctx, in.MemberId, in.OpcId)
	if err != nil {
		return nil, err
	}
	if row == nil || row.MaterialsSubmittedAt == 0 {
		return nil, gerror.New("尚未提交注册资料")
	}

	sections := []string{materialsSectionCompany, materialsSectionLegal, materialsSectionAttachments}
	out := make([]compliancein.MaterialsRevealSectionModel, 0, len(sections))
	for _, section := range sections {
		item := compliancein.MaterialsRevealSectionModel{
			Key:    section,
			Title:  materialsSectionTitle(section),
			Fields: buildPlainSectionFields(ctx, row, section),
		}
		if section == materialsSectionAttachments {
			attachments, err := buildAttachmentRevealItems(ctx, in.MemberId, row)
			if err != nil {
				return nil, err
			}
			item.Attachments = attachments
		}
		out = append(out, item)
	}
	return &compliancein.MaterialsRevealAllModel{
		OpcId:    row.Id,
		Sections: out,
	}, nil
}

// RevealMaterialsSection 密码验证后返回完整资料
func (s *sComplianceOpc) RevealMaterialsSection(ctx context.Context, in *compliancein.MaterialsRevealInp) (*compliancein.MaterialsRevealModel, error) {
	if in.MemberId == 0 {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先登录")
	}
	if !s.hasActiveOrder(ctx, in.MemberId) {
		return nil, gerror.NewCode(consts.CodeNoPermission, "请先完成签约")
	}
	section := strings.TrimSpace(in.Section)
	if !allowedMaterialsSections[section] {
		return nil, gerror.New("无效的资料分组")
	}
	if strings.TrimSpace(in.Password) == "" {
		return nil, gerror.New("请输入密码")
	}
	if err := service.MemberUser().VerifyPassword(ctx, in.MemberId, in.Password); err != nil {
		return nil, err
	}

	row, err := s.loadOpcForMaterials(ctx, in.MemberId, in.OpcId)
	if err != nil {
		return nil, err
	}
	if row == nil || row.MaterialsSubmittedAt == 0 {
		return nil, gerror.New("尚未提交注册资料")
	}

	fields := buildPlainSectionFields(ctx, row, section)
	result := &compliancein.MaterialsRevealModel{
		Key:    section,
		Fields: fields,
	}
	if section == materialsSectionAttachments {
		attachments, err := buildAttachmentRevealItems(ctx, in.MemberId, row)
		if err != nil {
			return nil, err
		}
		result.Attachments = attachments
	}
	return result, nil
}

func buildMaterialsDetailSections(ctx context.Context, row *opcRow) []compliancein.MaterialsDetailSection {
	sections := []string{materialsSectionCompany, materialsSectionLegal, materialsSectionAttachments}
	out := make([]compliancein.MaterialsDetailSection, 0, len(sections))
	for _, section := range sections {
		out = append(out, compliancein.MaterialsDetailSection{
			Key:    section,
			Title:  materialsSectionTitle(section),
			Fields: buildMaskedSectionFields(ctx, row, section),
		})
	}
	return out
}

func companyNameMasked(row *opcRow) string {
	if name := strings.TrimSpace(row.CompanyName); name != "" {
		return maskCompanyName(name)
	}
	var names []string
	if row.ProposedNames != "" {
		_ = gjson.DecodeTo(row.ProposedNames, &names)
	}
	if len(names) > 0 && strings.TrimSpace(names[0]) != "" {
		return maskCompanyName(names[0])
	}
	return "未填写"
}

func maskCompanyName(name string) string {
	runes := []rune(strings.TrimSpace(name))
	n := len(runes)
	if n <= 1 {
		return name
	}
	if n <= 3 {
		return string(runes[0]) + "*"
	}
	return string(runes[:2]) + "***" + string(runes[n-1])
}

func (s *sComplianceOpc) loadSubmittedOpcByMember(ctx context.Context, memberId uint64) ([]*opcRow, error) {
	var rows []*opcRow
	err := g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("member_id", memberId).
		Where("deleted", 0).
		Where("materials_submitted_at > ?", 0).
		OrderDesc("materials_submitted_at").
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "查询OPC主体失败")
	}
	return rows, nil
}

func (s *sComplianceOpc) loadOpcByMemberAndId(ctx context.Context, memberId, opcId uint64) (*opcRow, error) {
	var row opcRow
	err := g.DB().Model(tableOpcEntity).Ctx(ctx).
		Where("id", opcId).
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

func (s *sComplianceOpc) loadOpcForMaterials(ctx context.Context, memberId, opcId uint64) (*opcRow, error) {
	if opcId > 0 {
		return s.loadOpcByMemberAndId(ctx, memberId, opcId)
	}
	return s.loadOpcByMember(ctx, memberId)
}

func materialsSectionTitle(section string) string {
	switch section {
	case materialsSectionCompany:
		return "拟设公司信息"
	case materialsSectionLegal:
		return "法人身份信息"
	case materialsSectionAttachments:
		return "证件与证明材料"
	default:
		return section
	}
}

func materialsItemStatus(row *opcRow) (status, label string) {
	switch row.Status {
	case opcMaterials:
		return "rejected", "需补正"
	case opcMaterialsReview:
		return "reviewing", "审核中"
	default:
		if row.MaterialsSubmittedAt > 0 {
			return "approved", "已通过"
		}
		return "pending", "未提交"
	}
}

func legalOverviewSummary(row *opcRow) string {
	name := maskName(row.LegalPersonName)
	phone := security.MaskMobile(row.Phone)
	if name == "" && phone == "" {
		return "法人信息已提交"
	}
	return fmt.Sprintf("%s · %s", name, phone)
}

func buildMaskedSectionFields(ctx context.Context, row *opcRow, section string) []compliancein.MaterialsSectionField {
	switch section {
	case materialsSectionCompany:
		return maskedCompanyFields(row)
	case materialsSectionLegal:
		return maskedLegalFields(ctx, row)
	case materialsSectionAttachments:
		return maskedAttachmentFields(row)
	default:
		return nil
	}
}

func buildPlainSectionFields(ctx context.Context, row *opcRow, section string) []compliancein.MaterialsSectionField {
	switch section {
	case materialsSectionCompany:
		return plainCompanyFields(row)
	case materialsSectionLegal:
		return plainLegalFields(ctx, row)
	case materialsSectionAttachments:
		return plainAttachmentFields(row)
	default:
		return nil
	}
}

func maskedCompanyFields(row *opcRow) []compliancein.MaterialsSectionField {
	var names []string
	if row.ProposedNames != "" {
		_ = gjson.DecodeTo(row.ProposedNames, &names)
	}
	namesText := strings.Join(names, "、")
	if namesText == "" {
		namesText = "—"
	}
	termLabel := "长期"
	if row.BusinessTermType == "fixed" {
		termLabel = row.BusinessTermEnd
		if termLabel == "" {
			termLabel = "固定期限"
		}
	}
	region := strings.TrimSpace(row.RegisterProvince + " " + row.RegisterCity + " " + row.RegisterDistrict)
	return []compliancein.MaterialsSectionField{
		{Label: "备选公司名称", Value: namesText},
		{Label: "注册资本", Value: fmt.Sprintf("%.2f 万元", row.RegisteredCapital)},
		{Label: "认缴期限", Value: fmt.Sprintf("%d 年", row.CapitalTermYears)},
		{Label: "营业期限", Value: termLabel},
		{Label: "经营范围", Value: truncateMask(row.BusinessScope, 20)},
		{Label: "注册地区", Value: region},
		{Label: "详细地址", Value: security.MaskAddress(row.RegisterAddress)},
	}
}

func plainCompanyFields(row *opcRow) []compliancein.MaterialsSectionField {
	var names []string
	if row.ProposedNames != "" {
		_ = gjson.DecodeTo(row.ProposedNames, &names)
	}
	namesText := strings.Join(names, "、")
	if namesText == "" {
		namesText = "—"
	}
	termLabel := "长期"
	if row.BusinessTermType == "fixed" {
		termLabel = row.BusinessTermEnd
		if termLabel == "" {
			termLabel = "固定期限"
		}
	}
	region := strings.TrimSpace(row.RegisterProvince + " " + row.RegisterCity + " " + row.RegisterDistrict)
	return []compliancein.MaterialsSectionField{
		{Label: "备选公司名称", Value: namesText},
		{Label: "注册资本", Value: fmt.Sprintf("%.2f 万元", row.RegisteredCapital)},
		{Label: "认缴期限", Value: fmt.Sprintf("%d 年", row.CapitalTermYears)},
		{Label: "营业期限", Value: termLabel},
		{Label: "经营范围", Value: row.BusinessScope},
		{Label: "注册地区", Value: region},
		{Label: "详细地址", Value: row.RegisterAddress},
	}
}

func maskedLegalFields(ctx context.Context, row *opcRow) []compliancein.MaterialsSectionField {
	idCardMasked := ""
	if row.IdCardEncrypted != "" {
		if plain, decErr := compliancecrypto.Decrypt(ctx, row.IdCardEncrypted); decErr == nil {
			idCardMasked = compliancecrypto.MaskIdCard(plain)
		}
	}
	return []compliancein.MaterialsSectionField{
		{Label: "法人姓名", Value: maskName(row.LegalPersonName)},
		{Label: "身份证号", Value: idCardMasked},
		{Label: "身份证有效期", Value: row.IdCardValidFrom + " 至 " + row.IdCardValidTo},
		{Label: "民族", Value: row.Ethnicity},
		{Label: "户籍地址", Value: security.MaskAddress(row.HouseholdAddress)},
		{Label: "现居住地址", Value: security.MaskAddress(row.ResidentialAddress)},
		{Label: "手机号", Value: security.MaskMobile(row.Phone)},
		{Label: "邮箱", Value: security.MaskEmail(row.Email)},
		{Label: "电子签名授权", Value: esignLabel(row.EsignAuthorized)},
	}
}

func plainLegalFields(ctx context.Context, row *opcRow) []compliancein.MaterialsSectionField {
	idCard := ""
	if row.IdCardEncrypted != "" {
		if plain, decErr := compliancecrypto.Decrypt(ctx, row.IdCardEncrypted); decErr == nil {
			idCard = plain
		}
	}
	return []compliancein.MaterialsSectionField{
		{Label: "法人姓名", Value: row.LegalPersonName},
		{Label: "身份证号", Value: idCard},
		{Label: "身份证有效期", Value: row.IdCardValidFrom + " 至 " + row.IdCardValidTo},
		{Label: "民族", Value: row.Ethnicity},
		{Label: "户籍地址", Value: row.HouseholdAddress},
		{Label: "现居住地址", Value: row.ResidentialAddress},
		{Label: "手机号", Value: row.Phone},
		{Label: "邮箱", Value: row.Email},
		{Label: "电子签名授权", Value: esignLabel(row.EsignAuthorized)},
	}
}

func maskedAttachmentFields(row *opcRow) []compliancein.MaterialsSectionField {
	return []compliancein.MaterialsSectionField{
		{Label: "身份证正面", Value: attachmentStatusLabel(row.IdCardFrontFileId)},
		{Label: "身份证反面", Value: attachmentStatusLabel(row.IdCardBackFileId)},
		{Label: "地址证明", Value: attachmentStatusLabel(row.AddressProofFileId)},
	}
}

func plainAttachmentFields(row *opcRow) []compliancein.MaterialsSectionField {
	return []compliancein.MaterialsSectionField{
		{Label: "身份证正面", Value: attachmentFileLabel(row.IdCardFrontFileId)},
		{Label: "身份证反面", Value: attachmentFileLabel(row.IdCardBackFileId)},
		{Label: "地址证明", Value: attachmentFileLabel(row.AddressProofFileId)},
	}
}

func attachmentStatusLabel(fileId uint64) string {
	if fileId > 0 {
		return "已上传"
	}
	return "未上传"
}

func attachmentFileLabel(fileId uint64) string {
	if fileId > 0 {
		return fmt.Sprintf("附件 #%d", fileId)
	}
	return "未上传"
}

func buildAttachmentRevealItems(ctx context.Context, memberId uint64, row *opcRow) ([]compliancein.MaterialsAttachmentReveal, error) {
	type itemDef struct {
		label  string
		fileId uint64
	}
	items := []itemDef{
		{label: "身份证正面", fileId: row.IdCardFrontFileId},
		{label: "身份证反面", fileId: row.IdCardBackFileId},
		{label: "地址证明", fileId: row.AddressProofFileId},
	}
	out := make([]compliancein.MaterialsAttachmentReveal, 0, len(items))
	for _, item := range items {
		if item.fileId == 0 {
			continue
		}
		att, err := attachmentaccess.LoadMeta(ctx, item.fileId)
		if err != nil {
			return nil, err
		}
		if att.UserId != memberId {
			return nil, gerror.NewCode(consts.CodeNoPermission, "无权访问附件")
		}
		accessURL, err := attachmentaccess.BuildSignedURL(ctx, memberId, item.fileId, 0)
		if err != nil {
			return nil, err
		}
		out = append(out, compliancein.MaterialsAttachmentReveal{
			Label:     item.label,
			FileId:    item.fileId,
			FileName:  att.Name,
			MimeType:  att.Mimetype,
			AccessUrl: accessURL,
		})
	}
	return out, nil
}

func esignLabel(v int) string {
	if v == 1 {
		return "已授权"
	}
	return "未授权"
}

func truncateMask(text string, keep int) string {
	text = strings.TrimSpace(text)
	runes := []rune(text)
	if len(runes) <= keep {
		return text
	}
	return string(runes[:keep]) + "..."
}
