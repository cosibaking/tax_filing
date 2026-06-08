package complianceverify

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	envProvider = "COMPLIANCE_VERIFY_PROVIDER"
	cfgProvider = "compliance.verifyProvider"
)

// MaterialsVerifyInput OPC 资料真实性校验入参
type MaterialsVerifyInput struct {
	MemberId          uint64
	LegalPersonName   string
	IdCard            string
	Phone             string
	Email             string
	IdCardFrontFileId uint64
	IdCardBackFileId  uint64
}

// IsMockMode 是否使用 Mock 真实性校验（默认 mock，便于 MVP 联调）
func IsMockMode(ctx context.Context) bool {
	provider := strings.TrimSpace(os.Getenv(envProvider))
	if provider == "" {
		provider = g.Cfg().MustGet(ctx, cfgProvider, "mock").String()
	}
	provider = strings.ToLower(strings.TrimSpace(provider))
	return provider == "" || provider == "mock"
}

// VerifyMaterials 校验法人身份信息真实性（OCR / 三要素 / 手机实名等）
func VerifyMaterials(ctx context.Context, in *MaterialsVerifyInput) error {
	if in == nil {
		return gerror.New("校验参数无效")
	}
	if IsMockMode(ctx) {
		g.Log().Debugf(ctx, "[complianceverify] mock pass memberId=%d phone=%s", in.MemberId, in.Phone)
		return nil
	}
	// 生产环境接入 OCR、运营商三要素、电子签等第三方后再实现
	return gerror.New("真实性校验服务未配置，请设置 compliance.verifyProvider=mock 或接入第三方")
}
