package compliance

import (
	"xygo/internal/logic/compliance/diagnosis"
	_ "xygo/internal/logic/compliance/audit"
	"xygo/internal/logic/compliance/ledger"
	"xygo/internal/logic/compliance/opc"
	"xygo/internal/logic/compliance/order"
	_ "xygo/internal/logic/compliance/dashboard"
	_ "xygo/internal/logic/compliance/statement"
	_ "xygo/internal/logic/compliance/tax"
	"xygo/internal/service"
)

func init() {
	service.RegisterComplianceDiagnosis(diagnosis.New())
	service.RegisterComplianceOrder(order.New())
	service.RegisterComplianceOpc(opc.New())
	service.RegisterComplianceLedger(ledger.New())
}
