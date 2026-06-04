import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { requireActiveOpc } from '@/lib/services/compliance/ledger/income-expense-service';
import { getMemberTaxChecklist } from '@/lib/services/compliance/ledger/tax-calendar-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const data = await getMemberTaxChecklist(opc.id);
    return jsonOk(data);
  } catch (error) {
    return handleApiError(error);
  }
}
