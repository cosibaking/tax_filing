import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import {
  ErrorCodes,
  isValidIncomeCategory,
  isValidIncomePlatform,
  parseIncomeSort,
} from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import {
  requireActiveOpc,
  listIncome,
  createIncome,
} from '@/lib/services/compliance/ledger/income-expense-service';
import { parseDateOnly } from '@/lib/services/compliance/opc/opc-materials-validation';

function parsePagination(searchParams: URLSearchParams) {
  const page = Math.max(1, parseInt(searchParams.get('page') ?? '1', 10) || 1);
  const rawSize = parseInt(searchParams.get('pageSize') ?? '10', 10) || 10;
  const pageSize = Math.min(50, Math.max(1, rawSize));
  return { page, pageSize };
}

function endOfDay(d: Date): Date {
  const end = new Date(d);
  end.setHours(23, 59, 59, 999);
  return end;
}

function startOfDay(d: Date): Date {
  const start = new Date(d);
  start.setHours(0, 0, 0, 0);
  return start;
}

function parseIncomeFilters(searchParams: URLSearchParams) {
  const dateFromStr = searchParams.get('dateFrom')?.trim() ?? '';
  const dateToStr = searchParams.get('dateTo')?.trim() ?? '';

  let dateFrom: Date | undefined;
  let dateTo: Date | undefined;

  if (dateFromStr || dateToStr) {
    if (!dateFromStr || !dateToStr) {
      return { error: '请同时选择开始与结束日期' as const };
    }
    try {
      dateFrom = startOfDay(parseDateOnly(dateFromStr));
      dateTo = endOfDay(parseDateOnly(dateToStr));
    } catch {
      return { error: '日期格式无效' as const };
    }
    if (dateFrom > dateTo) {
      return { error: '开始日期不能晚于结束日期' as const };
    }
  }

  const platformRaw = searchParams.get('platform')?.trim() ?? '';
  const categoryRaw = searchParams.get('category')?.trim() ?? '';
  if (platformRaw && !isValidIncomePlatform(platformRaw)) {
    return { error: '无效平台' as const };
  }
  if (categoryRaw && !isValidIncomeCategory(categoryRaw)) {
    return { error: '无效类别' as const };
  }

  return {
    filters: {
      dateFrom,
      dateTo,
      platform: platformRaw || undefined,
      category: categoryRaw || undefined,
    },
  };
}

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const searchParams = new URL(request.url).searchParams;
    const parsed = parseIncomeFilters(searchParams);
    if ('error' in parsed) {
      return jsonFail(ErrorCodes.VALIDATION, parsed.error);
    }
    const { page, pageSize } = parsePagination(searchParams);
    const sort = parseIncomeSort(searchParams.get('sort'));
    const result = await listIncome(opc.id, parsed.filters, page, pageSize, sort);
    return jsonOk(result);
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
