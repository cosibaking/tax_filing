package statementpdf

import (
	"bytes"
	"fmt"

	"github.com/signintech/gopdf"
)

// Data 对账单 PDF 渲染数据
type Data struct {
	CompanyName      string
	Period           string
	Revenue          float64
	Cost             float64
	Profit           float64
	PrepaidTax       float64
	CumulativeProfit float64
	FilingStatus     string
	IncomeEntries    int
	ExpenseEntries   int
	PendingFilings   int
	GeneratedAt      string
	SentAt           string
}

// Generate 生成对账单 PDF 字节流
func Generate(data Data) ([]byte, error) {
	fontData, err := loadFontData()
	if err != nil {
		return nil, err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	if err = pdf.AddTTFFontData("zh", fontData); err != nil {
		return nil, fmt.Errorf("加载字体失败: %w", err)
	}

	pageW := 595.28
	margin := 40.0
	contentW := pageW - margin*2
	y := margin

	y = drawHeader(&pdf, data.Period, margin, y, contentW)
	y = drawCompany(&pdf, data.CompanyName, margin, y, contentW)
	y = drawSummaryTable(&pdf, data, margin, y, contentW)
	y = drawServiceSection(&pdf, data, margin, y, contentW)
	drawFooter(&pdf, data.GeneratedAt, data.SentAt, margin, y, contentW)

	var buf bytes.Buffer
	if _, err = pdf.WriteTo(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawHeader(pdf *gopdf.GoPdf, period string, x, y, w float64) float64 {
	pdf.SetFillColor(90, 141, 238)
	pdf.RectFromUpperLeftWithStyle(0, 0, 595.28, 72, "F")
	pdf.SetTextColor(255, 255, 255)
	setFont(pdf, 20, "")
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, "月度对账单")
	y += 28
	setFont(pdf, 11, "")
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, fmt.Sprintf("账期：%s", period))
	pdf.SetTextColor(50, 50, 93)
	return y + 28
}

func drawCompany(pdf *gopdf.GoPdf, company string, x, y, w float64) float64 {
	if company == "" {
		company = "—"
	}
	setFont(pdf, 11, "")
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, fmt.Sprintf("OPC 经营主体：%s", company))
	return y + 24
}

func drawSummaryTable(pdf *gopdf.GoPdf, data Data, x, y, w float64) float64 {
	setFont(pdf, 13, "")
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, "财务摘要")
	y += 22

	rows := [][]string{
		{"本月收入", money(data.Revenue)},
		{"本月成本", money(data.Cost)},
		{"本月利润", money(data.Profit)},
		{"本月预缴税额", money(data.PrepaidTax)},
		{"累计年度利润", money(data.CumulativeProfit)},
		{"申报状态", filingStatusLabel(data.FilingStatus)},
	}

	setFont(pdf, 10, "")
	for _, row := range rows {
		pdf.SetX(x)
		pdf.SetY(y)
		pdf.Cell(nil, row[0])
		pdf.SetX(x + w*0.55)
		pdf.SetY(y)
		pdf.Cell(nil, row[1]+" 元")
		y += 18
	}
	return y + 12
}

func drawServiceSection(pdf *gopdf.GoPdf, data Data, x, y, w float64) float64 {
	setFont(pdf, 13, "")
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, "服务完成摘要")
	y += 22

	items := []string{
		fmt.Sprintf("收入流水笔数：%d 笔", data.IncomeEntries),
		fmt.Sprintf("成本费用笔数：%d 笔", data.ExpenseEntries),
		fmt.Sprintf("待完成申报：%d 项", data.PendingFilings),
	}
	setFont(pdf, 10, "")
	for _, item := range items {
		pdf.SetX(x)
		pdf.SetY(y)
		pdf.Cell(nil, "• "+item)
		y += 18
	}
	return y + 8
}

func drawFooter(pdf *gopdf.GoPdf, generatedAt, sentAt string, x, y, w float64) {
	pdf.SetTextColor(136, 152, 170)
	setFont(pdf, 9, "")
	if generatedAt != "" {
		pdf.SetX(x)
		pdf.SetY(y)
		pdf.Cell(nil, fmt.Sprintf("生成时间：%s", generatedAt))
		y += 14
	}
	if sentAt != "" {
		pdf.SetX(x)
		pdf.SetY(y)
		pdf.Cell(nil, fmt.Sprintf("发送时间：%s", sentAt))
		y += 14
	}
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.MultiCellWithOption(&gopdf.Rect{W: w, H: 40}, "本对账单由税务合规服务平台自动生成，仅供参考。如有疑问请联系您的专属顾问。", gopdf.CellOption{})
}

func setFont(pdf *gopdf.GoPdf, size float64, _ string) {
	_ = pdf.SetFont("zh", "", size)
}

func money(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func filingStatusLabel(status string) string {
	switch status {
	case "filed":
		return "已申报"
	case "pending":
		return "待申报"
	default:
		if status == "" {
			return "已申报完成"
		}
		return status
	}
}
