import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { requireActiveOpc } from '@/lib/services/compliance/ledger/income-expense-service';
import { getTaxTaskDetail } from '@/lib/services/compliance/ledger/tax-calendar-service';
import { serializeBigInt } from '@/lib/db';

type RouteParams = { params: Promise<{ id: string }> };

export async function GET(request: Request, { params }: RouteParams) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const { id } = await params;
    const detail = await getTaxTaskDetail(opc.id, BigInt(id));
    if (!detail) {
      return jsonFail(ErrorCodes.NOT_FOUND, '申报任务不存在', 404);
    }
    return jsonOk(serializeBigInt(detail));
  } catch (error) {
    return handleApiError(error);
  }
}
