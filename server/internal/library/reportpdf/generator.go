package reportpdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/signintech/gopdf"

	"xygo/internal/library/statementpdf"
)

const (
	pageWidth    = 595.28
	pageHeight   = 841.89
	margin       = 42.0
	contentWidth = pageWidth - margin*2
	bottomLimit  = pageHeight - 52
	lineHeight   = 15.0
)

type Data struct {
	CompanyName, PeriodKey string
	Version                uint
	Status, GeneratedAt    string
	LegacyContent          string
	Summary                Summary
	Categories             []Category
	Anomalies              []Anomaly
}

type Summary struct {
	Conclusion       string
	CompletenessRate float64
	HighCount        int
	MediumCount      int
	LowCount         int
	DataNotice       string
}

type Category struct {
	Code, Name, Status, Summary string
	Checks                      []Check
}

type Check struct{ Code, Name, Status, Message string }

type Anomaly struct {
	Code, CategoryCode, Title, Severity  string
	Facts, Basis, Impact, Recommendation string
	RequiredMaterials                    []string
	DueDate                              string
	RequiresManualReview                 bool
	RuleVersion                          string
}

func Generate(data Data) ([]byte, error) {
	fontData, err := statementpdf.LoadChineseFontData()
	if err != nil {
		return nil, fmt.Errorf("加载中文字体失败: %w", err)
	}

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	if err := pdf.AddTTFFontData("zh", fontData); err != nil {
		return nil, fmt.Errorf("加载中文字体失败: %w", err)
	}
	r := renderer{pdf: pdf}
	if err := r.addPage(); err != nil {
		return nil, err
	}
	if err := r.render(data); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := pdf.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("写入 PDF 失败: %w", err)
	}
	return buf.Bytes(), nil
}

type renderer struct {
	pdf *gopdf.GoPdf
	y   float64
}

func (r *renderer) addPage() error {
	r.pdf.AddPage()
	r.y = margin
	return r.setFont(10)
}

func (r *renderer) setFont(size float64) error {
	if err := r.pdf.SetFont("zh", "", size); err != nil {
		return fmt.Errorf("设置中文字体失败: %w", err)
	}
	return nil
}

func (r *renderer) ensureSpace(height float64) error {
	if r.y+height <= bottomLimit {
		return nil
	}
	return r.addPage()
}

func (r *renderer) render(data Data) error {
	if err := r.heading("月度经营体检报告", 20); err != nil {
		return err
	}
	meta := []string{
		"企业：" + fallback(data.CompanyName, "—"), "账期：" + fallback(data.PeriodKey, "—"),
		fmt.Sprintf("版本：%d", data.Version), "状态：" + statusLabel(data.Status), "生成时间：" + fallback(data.GeneratedAt, "—"),
	}
	for _, item := range meta {
		if err := r.paragraph(item, 10); err != nil {
			return err
		}
	}

	if isStructuredEmpty(data) && strings.TrimSpace(data.LegacyContent) != "" {
		if err := r.heading("历史报告摘要", 15); err != nil {
			return err
		}
		if err := r.paragraph(data.LegacyContent, 10); err != nil {
			return err
		}
		return r.disclaimer()
	}
	if err := r.heading("总体结论", 15); err != nil {
		return err
	}
	summary := []string{
		"结论：" + conclusionLabel(data.Summary.Conclusion),
		fmt.Sprintf("资料完整度：%.1f%%", data.Summary.CompletenessRate),
		fmt.Sprintf("风险数量：高风险 %d 项、中风险 %d 项、低风险 %d 项", data.Summary.HighCount, data.Summary.MediumCount, data.Summary.LowCount),
		"资料提示：" + fallback(data.Summary.DataNotice, "暂无资料提示。"),
	}
	for _, item := range summary {
		if err := r.paragraph(item, 10); err != nil {
			return err
		}
	}

	if err := r.heading("六类检查结果", 15); err != nil {
		return err
	}
	for i, category := range data.Categories {
		if err := r.subheading(fmt.Sprintf("%d. %s（%s）", i+1, fallback(category.Name, category.Code), statusLabel(category.Status))); err != nil {
			return err
		}
		if err := r.paragraph("分类结论："+fallback(category.Summary, "—"), 10); err != nil {
			return err
		}
		for _, check := range category.Checks {
			if err := r.paragraph(fmt.Sprintf("检查项：%s（%s） %s", fallback(check.Name, check.Code), statusLabel(check.Status), fallback(check.Message, "—")), 9); err != nil {
				return err
			}
		}
	}

	if err := r.heading("异常详情", 15); err != nil {
		return err
	}
	if len(data.Anomalies) == 0 {
		if err := r.paragraph("当前未列出异常事项。", 10); err != nil {
			return err
		}
	}
	for i, anomaly := range data.Anomalies {
		if err := r.subheading(fmt.Sprintf("异常 %d：%s", i+1, fallback(anomaly.Title, "待核实异常"))); err != nil {
			return err
		}
		fields := []string{
			"等级：" + severityLabel(anomaly.Severity), "事实：" + fallback(anomaly.Facts, "—"), "判断依据：" + fallback(anomaly.Basis, "—"),
			"可能影响：" + fallback(anomaly.Impact, "—"), "处理建议：" + fallback(anomaly.Recommendation, "—"),
			"所需材料：" + fallback(strings.Join(anomaly.RequiredMaterials, "、"), "—"), "处理期限：" + fallback(anomaly.DueDate, "—"),
			"是否人工复核：" + yesNo(anomaly.RequiresManualReview),
		}
		for _, field := range fields {
			if err := r.paragraph(field, 9); err != nil {
				return err
			}
		}
	}
	return r.disclaimer()
}

func (r *renderer) heading(text string, size float64) error {
	if err := r.ensureSpace(size + 18); err != nil {
		return err
	}
	if err := r.setFont(size); err != nil {
		return err
	}
	r.pdf.SetX(margin)
	r.pdf.SetY(r.y)
	r.pdf.Cell(nil, text)
	r.y += size + 14
	return nil
}

func (r *renderer) subheading(text string) error { return r.paragraph(text, 11) }

func (r *renderer) paragraph(text string, size float64) error {
	if err := r.setFont(size); err != nil {
		return err
	}
	for _, line := range wrapText(strings.TrimSpace(text), charsPerLine(size)) {
		if err := r.ensureSpace(lineHeight); err != nil {
			return err
		}
		if err := r.setFont(size); err != nil {
			return err
		}
		r.pdf.SetX(margin)
		r.pdf.SetY(r.y)
		r.pdf.Cell(nil, line)
		r.y += lineHeight
	}
	r.y += 3
	return nil
}

func (r *renderer) disclaimer() error {
	if err := r.heading("免责声明", 13); err != nil {
		return err
	}
	return r.paragraph("本报告基于当前提供的数据和自动化规则生成，仅作为经营事项提醒，不替代会计、税务或法律专业意见；重要事项请结合原始资料由专业人员人工复核。", 9)
}

func wrapText(text string, width int) []string {
	if text == "" {
		return []string{"—"}
	}
	var result []string
	for _, raw := range strings.Split(text, "\n") {
		runes := []rune(raw)
		if len(runes) == 0 {
			result = append(result, " ")
			continue
		}
		for len(runes) > width {
			result = append(result, string(runes[:width]))
			runes = runes[width:]
		}
		result = append(result, string(runes))
	}
	return result
}

func charsPerLine(size float64) int {
	n := int(contentWidth / (size * 1.05))
	if n < 12 {
		return 12
	}
	return n
}

func isStructuredEmpty(data Data) bool {
	return data.Summary == (Summary{}) && len(data.Categories) == 0 && len(data.Anomalies) == 0
}

func fallback(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return strings.TrimSpace(value)
}
func yesNo(v bool) string {
	if v {
		return "是"
	}
	return "否"
}
func conclusionLabel(v string) string {
	switch strings.ToLower(v) {
	case "normal":
		return "正常"
	case "attention":
		return "需关注"
	case "urgent":
		return "紧急处理"
	default:
		return fallback(v, "未知")
	}
}
func severityLabel(v string) string {
	switch strings.ToLower(v) {
	case "high":
		return "高风险"
	case "medium":
		return "中风险"
	case "low":
		return "低风险"
	default:
		return fallback(v, "未知")
	}
}
func statusLabel(v string) string {
	switch strings.ToLower(v) {
	case "generated":
		return "已生成"
	case "normal":
		return "正常"
	case "attention":
		return "需关注"
	case "urgent":
		return "紧急处理"
	case "insufficient":
		return "数据不足"
	default:
		return fallback(v, "未知")
	}
}
