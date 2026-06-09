package platformocr

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	envOCRProvider = "COMPLIANCE_OCR_PROVIDER"
	cfgOCRProvider = "compliance.ocrProvider"
)

// ExtractText 从截图或粘贴文本提取 OCR 内容
func ExtractText(ctx context.Context, imageBytes []byte, filename, ocrText string) (string, string, error) {
	ocrText = strings.TrimSpace(ocrText)
	if ocrText != "" {
		return ocrText, "paste", nil
	}
	if len(imageBytes) == 0 {
		return "", "", gerror.New("请上传平台流水截图或粘贴 OCR 识别文本")
	}

	provider := strings.ToLower(strings.TrimSpace(os.Getenv(envOCRProvider)))
	if provider == "" {
		provider = strings.ToLower(strings.TrimSpace(g.Cfg().MustGet(ctx, cfgOCRProvider, "auto").String()))
	}

	switch provider {
	case "mock":
		return mockExtract(filename), "mock", nil
	case "tesseract", "auto":
		text, err := tesseractExtract(ctx, imageBytes)
		if err == nil && strings.TrimSpace(text) != "" {
			return strings.TrimSpace(text), "tesseract", nil
		}
		if provider == "tesseract" {
			return "", "", gerror.Wrap(err, "OCR 识别失败，请粘贴识别文本后重试")
		}
		// auto 降级：提示用户粘贴文本
		return "", "need_text", gerror.New("未能自动识别截图文字，请在下方粘贴 OCR 识别文本后重新预览")
	default:
		return mockExtract(filename), "mock", nil
	}
}

func mockExtract(filename string) string {
	lower := strings.ToLower(filename)
	platform := "douyin"
	switch {
	case strings.Contains(lower, "kuaishou") || strings.Contains(lower, "快手"):
		platform = "kuaishou"
	case strings.Contains(lower, "bilibili") || strings.Contains(lower, "b站"):
		platform = "bilibili"
	}
	return "2025-05-15 平台结算收入 ¥12,800.00 元\n2025-05-22 带货佣金 ¥3,200.50 元\n#platform:" + platform
}

func tesseractExtract(ctx context.Context, imageBytes []byte) (string, error) {
	if _, err := exec.LookPath("tesseract"); err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "platform-ocr-*.png")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err = tmp.Write(imageBytes); err != nil {
		tmp.Close()
		return "", err
	}
	tmp.Close()

	outBase := tmpPath + "_out"
	defer os.Remove(outBase + ".txt")

	cmd := exec.CommandContext(ctx, "tesseract", tmpPath, outBase, "-l", "chi_sim+eng")
	if err = cmd.Run(); err != nil {
		return "", err
	}
	b, err := os.ReadFile(outBase + ".txt")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GuessPlatformFromText 从 OCR 文本猜测平台
func GuessPlatformFromText(text, hint string) string {
	hint = strings.ToLower(strings.TrimSpace(hint))
	if hint != "" && hint != "other" {
		return hint
	}
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "抖音") || strings.Contains(lower, "#platform:douyin"):
		return "douyin"
	case strings.Contains(lower, "快手") || strings.Contains(lower, "#platform:kuaishou"):
		return "kuaishou"
	case strings.Contains(lower, "b站") || strings.Contains(lower, "bilibili"):
		return "bilibili"
	case strings.Contains(lower, "视频号"):
		return "channels"
	case strings.Contains(lower, "小红书"):
		return "xiaohongshu"
	default:
		return "other"
	}
}

// IsImageFilename 判断是否图片文件
func IsImageFilename(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".bmp":
		return true
	default:
		return false
	}
}
