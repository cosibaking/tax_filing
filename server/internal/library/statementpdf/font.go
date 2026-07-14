package statementpdf

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

var fontCandidates = []string{
	"resource/font/statement.ttf",
	"resource/captcha/fonts/SourceHanSansCN-Normal.ttf",
	"internal/library/statementpdf/fonts/statement.ttf",
	"/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
	"/usr/share/fonts/truetype/noto/NotoSansSC-Regular.ttf",
	"/usr/share/fonts/truetype/arphic/ukai.ttf",
	"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
}

func loadFontData() ([]byte, error) {
	ctx := context.Background()
	custom := ""
	if v, err := g.Cfg().Get(ctx, "compliance.statementPdfFont"); err == nil {
		custom = strings.TrimSpace(v.String())
	}
	if custom != "" {
		if data, err := os.ReadFile(custom); err == nil && len(data) > 0 {
			return data, nil
		}
	}
	for _, path := range fontCandidates {
		if !gfile.Exists(path) {
			continue
		}
		data, err := os.ReadFile(path)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}
	return nil, gerror.New("未找到 PDF 中文字体，请配置 compliance.statementPdfFont 或将字体放到 resource/font/statement.ttf")
}

// LoadChineseFontData returns the Chinese font used by statement PDFs.
// Callers must treat the returned bytes as read-only.
func LoadChineseFontData() ([]byte, error) {
	return loadFontData()
}

// ChineseFontData returns the Chinese font used by statement PDFs, or nil when
// no configured or bundled font can be found. Callers must treat it as read-only.
func ChineseFontData() []byte {
	data, _ := LoadChineseFontData()
	return data
}
