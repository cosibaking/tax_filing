import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import {
  requireActiveOpc,
  listIncome,
  createIncome,
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
    const items = await listIncome(opc.id, year, month);
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
    const { platform, category, grossAmount, platformFee, occurredAt, note } = body as {
      platform?: string;
      category?: string;
      grossAmount?: number;
      platformFee?: number;
      occurredAt?: string;
      note?: string;
    };

    if (!platform || !category || grossAmount == null || platformFee == null || !occurredAt) {
      return jsonFail(ErrorCodes.VALIDATION, '收入参数不完整');
    }

    const entry = await createIncome(opc.id, member.id, {
      platform,
      category,
      grossAmount,
      platformFee,
      occurredAt: new Date(occurredAt),
      note,
    });
    return jsonOk(entry);
  } catch (error) {
    return handleApiError(error);
  }
}
