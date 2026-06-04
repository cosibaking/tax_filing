import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import {
  requireActiveOpc,
  listExpense,
  createExpense,
} from '@/lib/services/compliance/ledger/income-expense-service';

function parseYearMonth(searchParams: URLSearchParams) {
  const now = new Date();
  const year = parseInt(searchParams.get('year') ?? String(now.getFullYear()), 10);
  const month = parseInt(searchParams.get('month') ?? String(now.getMonth() + 1), 10);
  return { year, month };
}

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const { year, month } = parseYearMonth(new URL(request.url).searchParams);
    const items = await listExpense(opc.id, year, month);
    return jsonOk(items);
  } catch (error) {
    return handleApiError(error);
  }
}

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const body = await request.json();
    const { category, amount, invoiceType, attachmentId, description, occurredAt, forceConfirm } =
      body as {
        category?: string;
        amount?: number;
        invoiceType?: string;
        attachmentId?: string;
        description?: string;
        occurredAt?: string;
        forceConfirm?: boolean;
      };

    if (!category || amount == null || !invoiceType || !occurredAt) {
      return jsonFail(ErrorCodes.VALIDATION, '费用参数不完整');
    }

    try {
      const entry = await createExpense(opc.id, member.id, {
        category,
        amount,
        invoiceType,
        attachmentId: attachmentId ? BigInt(attachmentId) : undefined,
        description,
        occurredAt: new Date(occurredAt),
        forceConfirm,
      });
      return jsonOk(entry);
    } catch (err) {
      if (err instanceof Error && err.message.startsWith('COST_RATIO_WARNING:')) {
        return jsonFail(ErrorCodes.COST_RATIO_WARNING, err.message.replace('COST_RATIO_WARNING:', ''));
      }
      throw err;
    }
  } catch (error) {
    return handleApiError(error);
  }
}
