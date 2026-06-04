import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { requireActiveOpc } from '@/lib/services/compliance/ledger/income-expense-service';
import { getTaxCalendar } from '@/lib/services/compliance/ledger/ledger-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const params = new URL(request.url).searchParams;
    const now = new Date();
    const year = parseInt(params.get('year') ?? String(now.getFullYear()), 10);
    const month = parseInt(params.get('month') ?? String(now.getMonth() + 1), 10);

    const calendar = await getTaxCalendar(opc.id, year, month);
    return jsonOk(calendar);
  } catch (error) {
    return handleApiError(error);
  }
}
