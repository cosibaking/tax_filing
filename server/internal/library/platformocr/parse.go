package platformocr

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParsedRow OCR/文本解析出的收入行
type ParsedRow struct {
	OccurredAt  string
	GrossAmount float64
	PlatformFee float64
	Category    string
	Remark      string
}

var (
	datePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(\d{4})[年/.-](\d{1,2})[月/.-](\d{1,2})`),
		regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`),
		regexp.MustCompile(`(\d{4})/(\d{2})/(\d{2})`),
	}
	amountPattern = regexp.MustCompile(`[¥￥]?\s*([\d,]+(?:\.\d{1,2})?)\s*元?`)
	lineKeywords  = []string{"结算", "收入", "提现", "到账", "分成", "佣金", "打赏", "广告", "服务费", "销售额"}
)

// ParseStatementText 从平台流水 OCR 文本中解析收入行
func ParseStatementText(text, platform string) []ParsedRow {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	var rows []ParsedRow
	seen := map[string]bool{}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if !lineLooksLikeIncome(line) {
			continue
		}
		date := extractDate(line)
		amounts := extractAmounts(line)
		if date == "" || len(amounts) == 0 {
			continue
		}
		gross := amounts[0]
		fee := 0.0
		if len(amounts) > 1 {
			fee = amounts[1]
			if fee > gross {
				fee, gross = gross, fee
			}
		}
		key := date + "|" + strconv.FormatFloat(gross, 'f', 2, 64)
		if seen[key] {
			continue
		}
		seen[key] = true
		rows = append(rows, ParsedRow{
			OccurredAt:  date,
			GrossAmount: gross,
			PlatformFee: fee,
			Category:    guessCategory(line),
			Remark:      truncate(line, 120),
		})
	}

	// 整段文本未按行命中时，尝试日期+金额邻近匹配
	if len(rows) == 0 {
		rows = parseBlockPairs(text)
	}

	_ = platform
	return rows
}

func lineLooksLikeIncome(line string) bool {
	for _, kw := range lineKeywords {
		if strings.Contains(line, kw) {
			return true
		}
	}
	if extractDate(line) != "" && len(extractAmounts(line)) > 0 {
		return true
	}
	return false
}

func extractDate(s string) string {
	for _, re := range datePatterns {
		m := re.FindStringSubmatch(s)
		if len(m) >= 4 {
			y, _ := strconv.Atoi(m[1])
			mo, _ := strconv.Atoi(m[2])
			d, _ := strconv.Atoi(m[3])
			if y >= 2000 && mo >= 1 && mo <= 12 && d >= 1 && d <= 31 {
				return time.Date(y, time.Month(mo), d, 0, 0, 0, 0, time.Local).Format("2006-01-02")
			}
		}
	}
	return ""
}

func extractAmounts(s string) []float64 {
	matches := amountPattern.FindAllStringSubmatch(s, -1)
	out := make([]float64, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		v, err := strconv.ParseFloat(strings.ReplaceAll(m[1], ",", ""), 64)
		if err != nil || v <= 0 || v > 1e9 {
			continue
		}
		out = append(out, v)
	}
	return out
}

func guessCategory(line string) string {
	switch {
	case strings.Contains(line, "佣金") || strings.Contains(line, "带货"):
		return "commission"
	case strings.Contains(line, "广告"):
		return "ad"
	case strings.Contains(line, "服务费"):
		return "service_fee"
	case strings.Contains(line, "打赏"):
		return "tip"
	default:
		return "other"
	}
}

func parseBlockPairs(text string) []ParsedRow {
	dates := regexp.MustCompile(`\d{4}[年/.-]\d{1,2}[月/.-]\d{1,2}`).FindAllString(text, -1)
	amounts := extractAmounts(text)
	if len(dates) == 0 || len(amounts) == 0 {
		return nil
	}
	n := len(dates)
	if len(amounts) < n {
		n = len(amounts)
	}
	var rows []ParsedRow
	for i := 0; i < n; i++ {
		d := extractDate(dates[i])
		if d == "" {
			continue
		}
		rows = append(rows, ParsedRow{
			OccurredAt:  d,
			GrossAmount: amounts[i],
			Category:    "other",
			Remark:      "OCR 块解析",
		})
	}
	return rows
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
