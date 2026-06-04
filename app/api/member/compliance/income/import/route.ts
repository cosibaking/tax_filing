import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import {
  requireActiveOpc,
  importIncomeCsv,
} from '@/lib/services/compliance/ledger/income-expense-service';

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const body = await request.json();
    const { rows } = body as {
      rows?: Array<{
        platform: string;
        category: string;
        grossAmount: number;
        platformFee: number;
        occurredAt: string;
      }>;
    };

    if (!rows?.length) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供导入数据');
    }

    const result = await importIncomeCsv(opc.id, member.id, rows);
    return jsonOk(result);
  } catch (error) {
    return handleApiError(error);
  }
}
