package opc

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"xygo/internal/library/complianceverify"
	"xygo/internal/model/input/compliancein"
)

var (
	emailRe      = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	mobileRe      = regexp.MustCompile(`^1[3-9]\d{9}$`)
	creditCodeRe  = regexp.MustCompile(`^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$`)
	allowedMime   = map[string]bool{
		"image/jpeg":       true,
		"image/png":        true,
		"application/pdf":  true,
	}
	idCardWeights = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	idCardCheckMap = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
)

func validateMaterials(ctx context.Context, in *compliancein.MaterialsSubmitInp) error {
	if len(in.ProposedNames) < 1 || len(in.ProposedNames) > 3 {
		return gerror.New("备选公司名称须为1-3条")
	}
	for _, name := range in.ProposedNames {
		name = strings.TrimSpace(name)
		if len([]rune(name)) < 2 || len([]rune(name)) > 30 {
			return gerror.New("公司名称每条须为2-30字")
		}
		if isAllDigits(name) {
			return gerror.New("公司名称不能为纯数字")
		}
	}
	if in.RegisteredCapital <= 0 || in.RegisteredCapital > 1000 {
		return gerror.New("注册资本须在0-1000万元之间")
	}
	if in.CapitalTermYears != 5 && in.CapitalTermYears != 10 && in.CapitalTermYears != 20 && in.CapitalTermYears != 30 {
		return gerror.New("请选择有效的认缴期限")
	}
	if in.BusinessTermType != "long_term" && in.BusinessTermType != "fixed" {
		return gerror.New("请选择营业期限类型")
	}
	if in.BusinessTermType == "fixed" && strings.TrimSpace(in.BusinessTermEnd) == "" {
		return gerror.New("固定营业期限须填写截止日期")
	}
	scope := strings.TrimSpace(in.BusinessScope)
	if len([]rune(scope)) < 10 || len([]rune(scope)) > 2000 {
		return gerror.New("经营范围须为10-2000字")
	}
	if !containsScopeKeyword(scope) {
		return gerror.New("经营范围须包含直播/文化/信息技术相关表述")
	}
	if strings.TrimSpace(in.RegisterProvince) == "" || strings.TrimSpace(in.RegisterCity) == "" ||
		strings.TrimSpace(in.RegisterDistrict) == "" {
		return gerror.New("请完整填写注册地址省市区")
	}
	if len([]rune(strings.TrimSpace(in.RegisterAddress))) < 5 {
		return gerror.New("详细注册地址至少5个字")
	}
	if in.AddressProofFileId == 0 {
		return gerror.New("请上传地址证明附件")
	}
	if len([]rune(strings.TrimSpace(in.LegalPersonName))) < 2 || len([]rune(strings.TrimSpace(in.LegalPersonName))) > 20 {
		return gerror.New("法人姓名须为2-20个中文")
	}
	if !validateIdCardForMode(ctx, in.IdCard) {
		return gerror.New("身份证号格式无效")
	}
	if strings.TrimSpace(in.IdCardValidFrom) == "" || strings.TrimSpace(in.IdCardValidTo) == "" {
		return gerror.New("请填写身份证有效期")
	}
	if len([]rune(strings.TrimSpace(in.HouseholdAddress))) < 5 {
		return gerror.New("户籍地址至少5个字")
	}
	if len([]rune(strings.TrimSpace(in.ResidentialAddress))) < 5 {
		return gerror.New("现居住地址至少5个字")
	}
	if !mobileRe.MatchString(strings.TrimSpace(in.Phone)) {
		return gerror.New("手机号格式无效")
	}
	if !emailRe.MatchString(strings.TrimSpace(in.Email)) {
		return gerror.New("邮箱格式无效")
	}
	if in.IdCardFrontFileId == 0 {
		return gerror.New("请上传身份证正面")
	}
	if in.IdCardBackFileId == 0 {
		return gerror.New("请上传身份证反面")
	}
	if !in.EsignAuthorized {
		return gerror.New("请勾选电子签名授权")
	}
	if !in.Confirmations.Truthful || !in.Confirmations.UsageConsent || !in.Confirmations.OpcLimitAck {
		return gerror.New("请勾选全部授权确认项")
	}
	return nil
}

var idCardFormatRe = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`)

func validateIdCardForMode(ctx context.Context, id string) bool {
	if complianceverify.IsMockMode(ctx) {
		return validateIdCardFormat(id)
	}
	return validateIdCard(id)
}

func validateIdCardFormat(id string) bool {
	return idCardFormatRe.MatchString(strings.ToUpper(strings.TrimSpace(id)))
}

func validateIdCard(id string) bool {
	id = strings.ToUpper(strings.TrimSpace(id))
	if !validateIdCardFormat(id) {
		return false
	}
	sum := 0
	for i := 0; i < 17; i++ {
		sum += int(id[i]-'0') * idCardWeights[i]
	}
	return idCardCheckMap[sum%11] == id[17]
}

func validateCreditCode(code string) bool {
	return creditCodeRe.MatchString(strings.ToUpper(strings.TrimSpace(code)))
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func containsScopeKeyword(scope string) bool {
	keywords := []string{"直播", "文化", "信息技术", "技术", "互联网", "文艺", "传媒", "数字内容"}
	for _, kw := range keywords {
		if strings.Contains(scope, kw) {
			return true
		}
	}
	return false
}

func validateAttachment(ctx context.Context, fileId uint64) error {
	if fileId == 0 {
		return gerror.New("附件ID无效")
	}
	var att struct {
		Id       uint64 `json:"id"`
		Mimetype string `json:"mimetype"`
	}
	err := g.DB().Model("xy_sys_attachment").Ctx(ctx).Where("id", fileId).Scan(&att)
	if err != nil {
		return gerror.Wrap(err, "查询附件失败")
	}
	if att.Id == 0 {
		return gerror.New("附件不存在")
	}
	if !allowedMime[att.Mimetype] {
		return gerror.New("附件格式不支持，仅允许 jpg/png/pdf")
	}
	return nil
}

func maskName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return name
	}
	return string(runes[0]) + "*"
}
