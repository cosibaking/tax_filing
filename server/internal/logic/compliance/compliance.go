package compliance

import (
	_ "xygo/internal/logic/compliance/audit"
	_ "xygo/internal/logic/compliance/dashboard"
	"xygo/internal/logic/compliance/diagnosis"
	"xygo/internal/logic/compliance/ledger"
	"xygo/internal/logic/compliance/opc"
	"xygo/internal/logic/compliance/order"
	"xygo/internal/logic/compliance/profile"
	_ "xygo/internal/logic/compliance/social"
	_ "xygo/internal/logic/compliance/statement"
	_ "xygo/internal/logic/compliance/tax"
	"xygo/internal/service"
)

func init() {
	service.RegisterComplianceDiagnosis(diagnosis.New())
	service.RegisterComplianceOrder(order.New())
	service.RegisterComplianceOpc(opc.New())
	service.RegisterComplianceLedger(ledger.New())
	service.RegisterComplianceProfile(profile.New())
}
