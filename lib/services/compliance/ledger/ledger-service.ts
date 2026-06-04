import { prisma, serializeBigInt } from '@/lib/db';
import {
  calculateCitQuarterly,
  calculateSurcharge,
  calculateVat,
} from '@/lib/services/compliance/diagnosis/tax-calculator';

export async function generateVouchersForPeriod(opcId: bigint, period: string) {
  const [yearStr, monthStr] = period.split('-');
  const year = parseInt(yearStr, 10);
  const month = parseInt(monthStr, 10);
  const start = new Date(year, month - 1, 1);
  const end = new Date(year, month, 0, 23, 59, 59);

  await prisma.ledgerVoucher.deleteMany({ where: { opcId, period } });

  const incomes = await prisma.incomeEntry.findMany({
    where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
  });
  for (const income of incomes) {
    await prisma.ledgerVoucher.create({
      data: {
        opcId,
        period,
        debitAccount: '1001',
        creditAccount: '6001',
        amount: income.netAmount,
        refType: 'income',
        refId: income.id,
      },
    });
  }

  const expenses = await prisma.expenseEntry.findMany({
    where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
  });
  for (const expense of expenses) {
    await prisma.ledgerVoucher.create({
      data: {
        opcId,
        period,
        debitAccount: '6401',
        creditAccount: '1001',
        amount: expense.amount,
        refType: 'expense',
        refId: expense.id,
      },
    });
  }
}

export async function recalculateProfit(opcId: bigint, year: number, month: number) {
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

  const revenue = Number(incomeSum._sum.grossAmount ?? 0) / 1.01;
  const cost = Number(expenseSum._sum.amount ?? 0);
  const profit = revenue - cost;

  const prev = await prisma.profitSummary.findMany({
    where: { opcId, year, month: { lt: month } },
  });
  const cumulativeProfit = prev.reduce((s, p) => s + Number(p.profit), 0) + profit;

  await prisma.profitSummary.upsert({
    where: { opcId_year_month: { opcId, year, month } },
    create: { opcId, year, month, revenue, cost, profit, cumulativeProfit },
    update: { revenue, cost, profit, cumulativeProfit },
  });
}

export async function getProfitReport(opcId: bigint, year: number, month?: number) {
  if (month) {
    const summary = await prisma.profitSummary.findUnique({
      where: { opcId_year_month: { opcId, year, month } },
    });
    const vouchers = await prisma.ledgerVoucher.findMany({
      where: { opcId, period: `${year}-${String(month).padStart(2, '0')}` },
      orderBy: { createdAt: 'desc' },
    });
    return serializeBigInt({ summary, vouchers });
  }
  const summaries = await prisma.profitSummary.findMany({
    where: { opcId, year },
    orderBy: { month: 'asc' },
  });
  return serializeBigInt(summaries);
}

export async function ensureTaxTasksForOpc(opcId: bigint, year: number, month: number) {
  const period = `${year}-${String(month).padStart(2, '0')}`;
  const dueDate = new Date(year, month, 15);

  const start = new Date(year, month - 1, 1);
  const end = new Date(year, month, 0, 23, 59, 59);
  const gross = await prisma.incomeEntry.aggregate({
    where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
    _sum: { grossAmount: true },
  });
  const monthlyGross = Number(gross._sum.grossAmount ?? 0);
  const vat = calculateVat(monthlyGross);
  const surcharge = calculateSurcharge(vat);

  if (!(await existingTask(opcId, 'vat', period))) {
    await prisma.taxFilingTask.create({
      data: { opcId, taxType: 'vat', period, dueDate, status: 'pending', calculatedAmount: vat },
    });
  }

  if (!(await existingTask(opcId, 'surcharge', period))) {
    await prisma.taxFilingTask.create({
      data: {
        opcId,
        taxType: 'surcharge',
        period,
        dueDate,
        status: 'pending',
        calculatedAmount: surcharge,
      },
    });
  }

  if ([3, 6, 9, 12].includes(month)) {
    const qStart = new Date(year, month - 3, 1);
    const summaries = await prisma.profitSummary.findMany({
      where: { opcId, year, month: { gte: month - 2, lte: month } },
    });
    const cumulative = summaries.reduce((s, p) => s + Number(p.profit), 0);
    const cit = calculateCitQuarterly(cumulative);
    const qPeriod = `Q${Math.ceil(month / 3)}-${year}`;
    if (!(await existingTask(opcId, 'cit_quarterly', qPeriod))) {
      await prisma.taxFilingTask.create({
        data: {
          opcId,
          taxType: 'cit_quarterly',
          period: qPeriod,
          dueDate,
          status: 'pending',
          calculatedAmount: cit,
        },
      });
    }
  }
}

async function existingTask(opcId: bigint, taxType: string, period: string) {
  return prisma.taxFilingTask.findFirst({
    where: { opcId, taxType, period, deleted: false },
  });
}

export async function getTaxCalendar(opcId: bigint, year: number, month: number) {
  const period = `${year}-${String(month).padStart(2, '0')}`;
  const tasks = await prisma.taxFilingTask.findMany({
    where: {
      opcId,
      deleted: false,
      OR: [{ period }, { period: { contains: String(year) } }],
    },
    orderBy: { dueDate: 'asc' },
  });
  const nextDue = tasks.find((t) => t.status === 'pending');
  return serializeBigInt({ tasks, nextDueDate: nextDue?.dueDate ?? null });
}

export async function markFilingFiled(
  taskId: bigint,
  adminId: bigint,
  filedAmount: number,
  checklist: boolean[],
  receiptFileId?: bigint,
) {
  if (checklist.length !== 9 || !checklist.every(Boolean)) {
    throw new Error('必须完成全部9项自查清单');
  }

  const task = await prisma.taxFilingTask.findFirst({ where: { id: taskId, deleted: false } });
  if (!task) throw new Error('NOT_FOUND');

  const ledgerSum = await getLedgerSumForTask(task);
  if (filedAmount < ledgerSum * 0.95) {
    throw new Error('申报额低于台账汇总，禁止低报');
  }

  await prisma.taxFilingTask.update({
    where: { id: taskId },
    data: {
      status: 'filed',
      filedAmount,
      receiptFileId,
      checklist: checklist as object,
    },
  });

  const { writeAuditLog } = await import('@/lib/services/compliance/audit/audit-log');
  await writeAuditLog({
    entityType: 'tax_filing_task',
    entityId: taskId,
    action: 'filing.mark_filed',
    operatorId: adminId,
    operatorType: 'admin',
    after: { filedAmount, checklist },
  });
}

async function getLedgerSumForTask(task: {
  opcId: bigint;
  taxType: string;
  period: string;
  calculatedAmount: { toString(): string };
}) {
  return Number(task.calculatedAmount);
}
