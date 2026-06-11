package ocrspace

import (
	"fmt"
	"strings"
)

const (
	defaultPostEndpoint = "https://api.ocr.space/parse/image"
	defaultGetEndpoint  = "https://api.ocr.space/parse/imageurl"
)

// Language OCR 语言代码（3 字母，如 eng、chs、auto）
type Language string

const (
	LanguageEnglish            Language = "eng"
	LanguageChineseSimplified  Language = "chs"
	LanguageChineseTraditional Language = "cht"
	LanguageAuto               Language = "auto"
)

// FileType 文件类型
type FileType string

const (
	FileTypePDF FileType = "PDF"
	FileTypePNG FileType = "PNG"
	FileTypeJPG FileType = "JPG"
	FileTypeGIF FileType = "GIF"
	FileTypeTIF FileType = "TIF"
	FileTypeBMP FileType = "BMP"
)

// OCREngine OCR 引擎编号
type OCREngine int

const (
	Engine1 OCREngine = 1
	Engine2 OCREngine = 2
	Engine3 OCREngine = 3
)

// ParseOptions OCR 请求可选参数，对应 OCR.space POST 参数
type ParseOptions struct {
	Language                     Language
	IsOverlayRequired            bool
	Filetype                     FileType
	DetectOrientation            bool
	IsCreateSearchablePdf        bool
	IsSearchablePdfHideTextLayer bool
	Scale                        bool
	IsTable                      bool
	OCREngine                    OCREngine
}

// ParseResponse OCR.space API 响应
type ParseResponse struct {
	ParsedResults                []ParsedResult `json:"ParsedResults"`
	OCRExitCode                  int            `json:"OCRExitCode"`
	IsErroredOnProcessing        bool           `json:"IsErroredOnProcessing"`
	ErrorMessage                 *string        `json:"ErrorMessage"`
	ErrorDetails                 *string        `json:"ErrorDetails"`
	SearchablePDFURL             *string        `json:"SearchablePDFURL"`
	ProcessingTimeInMilliseconds string         `json:"ProcessingTimeInMilliseconds"`
}

// ParsedResult 单页/单图 OCR 结果
type ParsedResult struct {
	TextOverlay       *TextOverlay `json:"TextOverlay"`
	FileParseExitCode int          `json:"FileParseExitCode"`
	ParsedText        string       `json:"ParsedText"`
	ErrorMessage      *string      `json:"ErrorMessage"`
	ErrorDetails      *string      `json:"ErrorDetails"`
}

// TextOverlay 文字坐标 overlay
type TextOverlay struct {
	Lines      []TextLine `json:"Lines"`
	HasOverlay bool       `json:"HasOverlay"`
	Message    *string    `json:"Message"`
}

// TextLine overlay 中的一行
type TextLine struct {
	Words     []Word `json:"Words"`
	MaxHeight int    `json:"MaxHeight"`
	MinTop    int    `json:"MinTop"`
}

// Word overlay 中的单个词
type Word struct {
	WordText string `json:"WordText"`
	Left     int    `json:"Left"`
	Top      int    `json:"Top"`
	Height   int    `json:"Height"`
	Width    int    `json:"Width"`
}

// IsSuccess 判断 OCR 是否整体成功（OCRExitCode 1 或 2）
func (r *ParseResponse) IsSuccess() bool {
	if r == nil {
		return false
	}
	return r.OCRExitCode == 1 || r.OCRExitCode == 2
}

// CombinedText 合并所有页面的识别文本
func (r *ParseResponse) CombinedText() string {
	if r == nil || len(r.ParsedResults) == 0 {
		return ""
	}
	parts := make([]string, 0, len(r.ParsedResults))
	for _, page := range r.ParsedResults {
		text := strings.TrimSpace(page.ParsedText)
		if text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, "\n\n")
}

// Err 将 API 错误转为 Go error
func (r *ParseResponse) Err() error {
	if r == nil {
		return fmt.Errorf("ocrspace: empty response")
	}
	if r.IsSuccess() {
		return nil
	}
	msg := ""
	if r.ErrorMessage != nil {
		msg = *r.ErrorMessage
	}
	if msg == "" {
		msg = fmt.Sprintf("OCR failed with exit code %d", r.OCRExitCode)
	}
	if r.ErrorDetails != nil && *r.ErrorDetails != "" {
		msg = msg + ": " + *r.ErrorDetails
	}
	return fmt.Errorf("ocrspace: %s", msg)
}
