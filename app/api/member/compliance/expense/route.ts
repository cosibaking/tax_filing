import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import {
  ErrorCodes,
  isValidExpenseCategory,
  isValidInvoiceType,
  parseExpenseSort,
} from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import {
  requireActiveOpc,
  listExpense,
  createExpense,
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

function parseExpenseFilters(searchParams: URLSearchParams) {
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

  const categoryRaw = searchParams.get('category')?.trim() ?? '';
  if (categoryRaw && !isValidExpenseCategory(categoryRaw)) {
    return { error: '无效费用类别' as const };
  }

  return {
    filters: {
      dateFrom,
      dateTo,
      category: categoryRaw || undefined,
    },
  };
}

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const opc = await requireActiveOpc(member.id);
    const searchParams = new URL(request.url).searchParams;
    const parsed = parseExpenseFilters(searchParams);
    if ('error' in parsed) {
      return jsonFail(ErrorCodes.VALIDATION, parsed.error);
    }
    const { page, pageSize } = parsePagination(searchParams);
    const sort = parseExpenseSort(searchParams.get('sort'));
    const result = await listExpense(opc.id, parsed.filters, page, pageSize, sort);
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
    if (!isValidExpenseCategory(category)) {
      return jsonFail(ErrorCodes.VALIDATION, '无效费用类别');
    }
    if (!isValidInvoiceType(invoiceType)) {
      return jsonFail(ErrorCodes.VALIDATION, '无效发票类型');
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
