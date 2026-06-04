import { prisma, serializeBigInt } from '@/lib/db';
import {
  DEFAULT_EXPENSE_SORT,
  DEFAULT_INCOME_SORT,
  EXPENSE_KEYWORD_BLACKLIST,
  type ExpenseSort,
  type IncomeSort,
  parseExpenseSort,
  parseIncomeSort,
} from '@/lib/api/constants';
import type { Prisma } from '@prisma/client';
import { writeAuditLog } from '@/lib/services/compliance/audit/audit-log';
import {
  generateVouchersForPeriod,
  recalculateProfit,
} from '@/lib/services/compliance/ledger/ledger-service';

const COST_THRESHOLD = parseFloat(process.env.COMPLIANCE_COST_RATIO_THRESHOLD ?? '0.8');

const INCOME_ORDER_BY: Record<IncomeSort, Prisma.IncomeEntryOrderByWithRelationInput> = {
  createdAt_desc: { createdAt: 'desc' },
  createdAt_asc: { createdAt: 'asc' },
  occurredAt_desc: { occurredAt: 'desc' },
  occurredAt_asc: { occurredAt: 'asc' },
};

const EXPENSE_ORDER_BY: Record<ExpenseSort, Prisma.ExpenseEntryOrderByWithRelationInput> = {
  createdAt_desc: { createdAt: 'desc' },
  createdAt_asc: { createdAt: 'asc' },
  occurredAt_desc: { occurredAt: 'desc' },
  occurredAt_asc: { occurredAt: 'asc' },
  amount_desc: { amount: 'desc' },
  amount_asc: { amount: 'asc' },
};

export async function requireActiveOpc(memberId: bigint) {
  const opc = await prisma.opcEntity.findFirst({
    where: { memberId, status: 'active', deleted: false },
  });
  if (!opc) throw new Error('OPC 未激活');
  return opc;
}

export type IncomeListFilters = {
  dateFrom?: Date;
  dateTo?: Date;
  platform?: string;
  category?: string;
};

function buildIncomeWhere(opcId: bigint, filters: IncomeListFilters): Prisma.IncomeEntryWhereInput {
  const where: Prisma.IncomeEntryWhereInput = {
    opcId,
    deleted: false,
  };
  if (filters.dateFrom != null || filters.dateTo != null) {
    where.occurredAt = {};
    if (filters.dateFrom != null) where.occurredAt.gte = filters.dateFrom;
    if (filters.dateTo != null) where.occurredAt.lte = filters.dateTo;
  }
  if (filters.platform) where.platform = filters.platform;
  if (filters.category) where.category = filters.category;
  return where;
}

export async function listIncome(
  opcId: bigint,
  filters: IncomeListFilters,
  page = 1,
  pageSize = 10,
  sort: IncomeSort = DEFAULT_INCOME_SORT,
) {
  const where = buildIncomeWhere(opcId, filters);
  const orderBy = INCOME_ORDER_BY[parseIncomeSort(sort)];

  const [total, items, agg] = await Promise.all([
    prisma.incomeEntry.count({ where }),
    prisma.incomeEntry.findMany({
      where,
      orderBy,
      skip: (page - 1) * pageSize,
      take: pageSize,
    }),
    prisma.incomeEntry.aggregate({
      where,
      _sum: { grossAmount: true, platformFee: true, netAmount: true },
    }),
  ]);

  return serializeBigInt({
    items,
    total,
    page,
    pageSize,
    totalPages: total === 0 ? 0 : Math.ceil(total / pageSize),
    sort: parseIncomeSort(sort),
    summary: {
      grossAmount: Number(agg._sum.grossAmount ?? 0),
      platformFee: Number(agg._sum.platformFee ?? 0),
      netAmount: Number(agg._sum.netAmount ?? 0),
    },
  });
}

export async function createIncome(
  opcId: bigint,
  memberId: bigint,
  data: {
    platform: string;
    category: string;
    grossAmount: number;
    platformFee: number;
    occurredAt: Date;
    note?: string;
  },
) {
  if (['loan', 'gift'].includes(data.category)) {
    throw new Error('禁止录入借款/赠与类收入');
  }
  if (data.note && EXPENSE_KEYWORD_BLACKLIST.some((k) => data.note!.includes(k))) {
    throw new Error('收入备注含违规关键词');
  }
  const netAmount = data.grossAmount - data.platformFee;
  const entry = await prisma.incomeEntry.create({
    data: { opcId, ...data, netAmount, source: 'manual' },
  });
  const period = `${data.occurredAt.getFullYear()}-${String(data.occurredAt.getMonth() + 1).padStart(2, '0')}`;
  await generateVouchersForPeriod(opcId, period);
  await recalculateProfit(opcId, data.occurredAt.getFullYear(), data.occurredAt.getMonth() + 1);
  await writeAuditLog({
    entityType: 'income_entry',
    entityId: entry.id,
    action: 'income.create',
    operatorId: memberId,
    operatorType: 'member',
    after: data,
  });
  return serializeBigInt(entry);
}

export type ExpenseListFilters = {
  dateFrom?: Date;
  dateTo?: Date;
  category?: string;
};

function buildExpenseWhere(opcId: bigint, filters: ExpenseListFilters): Prisma.ExpenseEntryWhereInput {
  const where: Prisma.ExpenseEntryWhereInput = {
    opcId,
    deleted: false,
  };
  if (filters.dateFrom != null || filters.dateTo != null) {
    where.occurredAt = {};
    if (filters.dateFrom != null) where.occurredAt.gte = filters.dateFrom;
    if (filters.dateTo != null) where.occurredAt.lte = filters.dateTo;
  }
  if (filters.category) where.category = filters.category;
  return where;
}

export async function listExpense(
  opcId: bigint,
  filters: ExpenseListFilters,
  page = 1,
  pageSize = 10,
  sort: ExpenseSort = DEFAULT_EXPENSE_SORT,
) {
  const where = buildExpenseWhere(opcId, filters);
  const orderBy = EXPENSE_ORDER_BY[parseExpenseSort(sort)];

  const [total, items, agg] = await Promise.all([
    prisma.expenseEntry.count({ where }),
    prisma.expenseEntry.findMany({
      where,
      orderBy,
      skip: (page - 1) * pageSize,
      take: pageSize,
    }),
    prisma.expenseEntry.aggregate({
      where,
      _sum: { amount: true },
    }),
  ]);

  return serializeBigInt({
    items,
    total,
    page,
    pageSize,
    totalPages: total === 0 ? 0 : Math.ceil(total / pageSize),
    sort: parseExpenseSort(sort),
    summary: {
      totalAmount: Number(agg._sum.amount ?? 0),
    },
  });
}

export async function createExpense(
  opcId: bigint,
  memberId: bigint,
  data: {
    category: string;
    amount: number;
    invoiceType: string;
    attachmentId?: bigint;
    description?: string;
    occurredAt: Date;
    forceConfirm?: boolean;
  },
) {
  if (data.description && EXPENSE_KEYWORD_BLACKLIST.some((k) => data.description!.includes(k))) {
    throw new Error('费用描述含违规关键词，禁止保存');
  }

  const year = data.occurredAt.getFullYear();
  const month = data.occurredAt.getMonth() + 1;
  const start = new Date(year, month - 1, 1);
  const end = new Date(year, month, 0, 23, 59, 59);

  const incomeSum = await prisma.incomeEntry.aggregate({
    where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
    _sum: { grossAmount: true },
  });
  const expenseSum = await prisma.expenseEntry.aggregate({
    where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
    _sum: { amount: true },
  });

  const revenue = Number(incomeSum._sum.grossAmount ?? 0);
  const totalCost = Number(expenseSum._sum.amount ?? 0) + data.amount;
  let warningFlag = false;

  if (revenue > 0 && totalCost / revenue > COST_THRESHOLD) {
    if (!data.forceConfirm) {
      throw new Error(`COST_RATIO_WARNING:成本占收入比 ${((totalCost / revenue) * 100).toFixed(1)}% 超过阈值`);
    }
    warningFlag = true;
  }

  const entry = await prisma.expenseEntry.create({
    data: {
      opcId,
      category: data.category,
      amount: data.amount,
      invoiceType: data.invoiceType,
      attachmentId: data.attachmentId,
      description: data.description,
      warningFlag,
      occurredAt: data.occurredAt,
    },
  });

  const period = `${year}-${String(month).padStart(2, '0')}`;
  await generateVouchersForPeriod(opcId, period);
  await recalculateProfit(opcId, year, month);

  await writeAuditLog({
    entityType: 'expense_entry',
    entityId: entry.id,
    action: 'expense.create',
    operatorId: memberId,
    operatorType: 'member',
    after: data,
  });

  return serializeBigInt(entry);
}

export async function importIncomeCsv(
  opcId: bigint,
  memberId: bigint,
  rows: Array<{
    platform: string;
    category: string;
    grossAmount: number;
    platformFee: number;
    occurredAt: string;
  }>,
) {
  let success = 0;
  let failed = 0;
  for (const row of rows) {
    try {
      await createIncome(opcId, memberId, {
        ...row,
        occurredAt: new Date(row.occurredAt),
      });
      success++;
    } catch {
      failed++;
    }
  }
  return { success, failed };
}

export async function importBankCsv(
  opcId: bigint,
  rows: Array<{ occurredAt: string; amount: number; description?: string }>,
) {
  const created = [];
  for (const row of rows) {
    const tx = await prisma.bankTransaction.create({
      data: {
        opcId,
        occurredAt: new Date(row.occurredAt),
        amount: row.amount,
        description: row.description,
      },
    });
    const income = await prisma.incomeEntry.findFirst({
      where: {
        opcId,
        deleted: false,
        netAmount: row.amount,
        occurredAt: {
          gte: new Date(new Date(row.occurredAt).getTime() - 86400000),
          lte: new Date(new Date(row.occurredAt).getTime() + 86400000),
        },
      },
    });
    if (income) {
      await prisma.bankTransaction.update({
        where: { id: tx.id },
        data: { matched: true, incomeId: income.id },
      });
    }
    created.push(tx);
  }
  const unmatched = await prisma.bankTransaction.count({ where: { opcId, matched: false, deleted: false } });
  return serializeBigInt({ imported: created.length, unmatched });
}

export async function getExpenseCategories() {
  return prisma.expenseCategory.findMany({ where: { active: true } });
}
