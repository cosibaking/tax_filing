package rules

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

// PeriodRange 解析期间字符串，返回 Unix 秒级起止时间（含起止日）
func PeriodRange(period string) (start, end uint64, err error) {
	return periodRange(period)
}

func periodRange(period string) (start, end uint64, err error) {
	period = strings.TrimSpace(period)
	if period == "" {
		return 0, 0, gerror.New("申报期间无效")
	}

	if strings.HasPrefix(period, "Q") && len(period) >= 6 {
		return quarterRange(period)
	}
	if len(period) == 4 {
		y, e := strconv.Atoi(period)
		if e != nil {
			return 0, 0, gerror.New("年度期间格式无效")
		}
		startT := time.Date(y, 1, 1, 0, 0, 0, 0, time.Local)
		endT := time.Date(y, 12, 31, 23, 59, 59, 0, time.Local)
		return uint64(startT.Unix()), uint64(endT.Unix()), nil
	}
	if len(period) == 7 && period[4] == '-' {
		return monthRange(period)
	}
	return 0, 0, gerror.Newf("不支持的期间格式: %s", period)
}

func monthRange(period string) (uint64, uint64, error) {
	t, err := time.ParseInLocation("2006-01", period, time.Local)
	if err != nil {
		return 0, 0, gerror.Wrap(err, "月份期间格式无效")
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return uint64(start.Unix()), uint64(end.Unix()), nil
}

func quarterRange(period string) (uint64, uint64, error) {
	var y, q int
	if _, err := fmt.Sscanf(period, "Q%d-%d", &q, &y); err != nil {
		if _, err = fmt.Sscanf(period, "%d-Q%d", &y, &q); err != nil {
			return 0, 0, gerror.New("季度期间格式无效，应为 Q1-2026")
		}
	}
	if q < 1 || q > 4 {
		return 0, 0, gerror.New("季度序号须在 1-4 之间")
	}
	startMonth := time.Month((q-1)*3 + 1)
	start := time.Date(y, startMonth, 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 3, 0).Add(-time.Second)
	return uint64(start.Unix()), uint64(end.Unix()), nil
}
