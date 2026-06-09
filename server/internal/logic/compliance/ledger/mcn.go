package ledger

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	settlementPersonal  = "personal"
	settlementMcnPublic = "mcn_public"
)

type mcnCalcResult struct {
	GrossAmount      float64
	McnShareAmount   float64
	GrossBeforeSplit float64
	McnSplitRatio    float64
	SettlementType   string
	McnName          string
}

func applyMcnFields(in *mcnCalcResult, settlementType, mcnName string, splitRatio, grossBeforeSplit, grossAmount float64) error {
	settlementType = strings.ToLower(strings.TrimSpace(settlementType))
	if settlementType == "" {
		settlementType = settlementPersonal
	}
	if settlementType != settlementPersonal && settlementType != settlementMcnPublic {
		return gerror.New("结算方式无效")
	}
	in.SettlementType = settlementType
	in.McnName = strings.TrimSpace(mcnName)
	in.McnSplitRatio = splitRatio
	in.GrossBeforeSplit = roundMoney(grossBeforeSplit)

	if grossBeforeSplit > 0 && splitRatio > 0 {
		if splitRatio < 0 || splitRatio >= 1 {
			return gerror.New("MCN 分成比例须在 0-100% 之间")
		}
		if settlementType == settlementMcnPublic && mcnName == "" {
			return gerror.New("MCN 对公结算请填写 MCN 机构名称")
		}
		in.McnShareAmount = roundMoney(grossBeforeSplit * splitRatio)
		in.GrossAmount = roundMoney(grossBeforeSplit - in.McnShareAmount)
		return nil
	}
	if grossAmount > 0 {
		in.GrossAmount = roundMoney(grossAmount)
	}
	return nil
}

func mcnDataMap(calc mcnCalcResult) g.Map {
	return g.Map{
		"settlement_type":    calc.SettlementType,
		"mcn_name":           calc.McnName,
		"mcn_split_ratio":    calc.McnSplitRatio,
		"mcn_share_amount":   calc.McnShareAmount,
		"gross_before_split": calc.GrossBeforeSplit,
	}
}
