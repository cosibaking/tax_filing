import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { requireActiveOpc } from '@/lib/services/compliance/ledger/income-expense-service';
import { loadTaxCalendar } from '@/lib/services/compliance/ledger/tax-calendar-service';
import { serializeBigInt } from '@/lib/db';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const params = new URL(request.url).searchParams;
    const now = new Date();
    const year = parseInt(params.get('year') ?? String(now.getFullYear()), 10);
    const month = parseInt(params.get('month') ?? String(now.getMonth() + 1), 10);

    const calendar = await loadTaxCalendar(opc.id, year, month);
    return jsonOk(serializeBigInt(calendar));
  } catch (error) {
    return handleApiError(error);
  }
}
