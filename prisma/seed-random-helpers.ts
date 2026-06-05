import { PrismaClient } from '@prisma/client';
import {
  calculateCitQuarterly,
  calculateSurcharge,
  calculateVat,
} from '../lib/services/compliance/diagnosis/tax-calculator';

export const PLATFORMS = ['douyin', 'kuaishou', 'bilibili', 'channels', 'xiaohongshu'] as const;
export const INCOME_CATEGORIES = ['tip', 'commission', 'ad', 'slot_fee', 'offline', 'other'] as const;
export const EXPENSE_CATEGORIES = [
  'equipment',
  'network',
  'venue',
  'promotion',
  'outsource',
  'travel',
  'other',
] as const;
export const INVOICE_TYPES = ['general', 'special', 'electronic'] as const;

export const SURNAMES = ['张', '李', '王', '刘', '陈', '杨', '赵', '黄', '周', '吴'] as const;
export const NICKNAMES = ['小鱼', '阿宁', '小鹿', '糖糖', '大白', '柚子', '橙子', '七七', '安安', '米米'] as const;

export const INCOME_NOTES = [
  '直播打赏收入',
  '带货佣金结算',
  '品牌广告合作',
  '坑位费收入',
  '线下商演',
  '平台活动奖励',
] as const;

export const EXPENSE_DESCRIPTIONS = [
  '直播设备采购',
  '宽带月费',
  '摄影棚租赁',
  'Dou+ 推广投放',
  '剪辑外包',
  '出差交通住宿',
  '办公耗材',
  '灯光支架补充',
] as const;

export const OPC_STATUS_PIPELINE = [
  'pending',
  'materials',
  'materials_review',
  'registering',
  'tax',
  'bank',
  'active',
] as const;

export type OpcPipelineStatus = (typeof OPC_STATUS_PIPELINE)[number];

export function pick<T>(arr: readonly T[]): T {
  return arr[Math.floor(Math.random() * arr.length)];
}

export function randomInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

export function randomAmount(min: number, max: number): number {
  return Math.round((Math.random() * (max - min) + min) * 100) / 100;
}

/** 基于序号生成唯一手机号，避免重复插入 */
export function testPhone(index: number): string {
  const suffix = String(10000000 + index).slice(-8);
  return `138${suffix}`;
}

export function randomCreditCode(): string {
  let code = '91';
  for (let i = 0; i < 16; i++) code += String(randomInt(0, 9));
  return code;
}

export function randomDateInPastMonths(monthsBack: number): Date {
  const now = new Date();
  const daysBack = randomInt(0, monthsBack * 30);
  const d = new Date(now);
  d.setDate(d.getDate() - daysBack);
  d.setHours(randomInt(9, 21), randomInt(0, 59), 0, 0);
  return d;
}

export function periodKey(date: Date): string {
  const y = date.getFullYear();
  const m = date.getMonth() + 1;
  return `${y}-${String(m).padStart(2, '0')}`;
}

async function aggregateMonthProfit(
  prisma: PrismaClient,
  opcId: bigint,
  year: number,
  month: number,
) {
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
  return { revenue, cost, profit };
}

export async function recalculateProfit(
  prisma: PrismaClient,
  opcId: bigint,
  year: number,
  month: number,
) {
  const { revenue, cost, profit } = await aggregateMonthProfit(prisma, opcId, year, month);
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

export async function generateVouchersForPeriod(
  prisma: PrismaClient,
  opcId: bigint,
  period: string,
) {
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

async function existingTask(
  prisma: PrismaClient,
  opcId: bigint,
  taxType: string,
  period: string,
) {
  const found = await prisma.taxFilingTask.findFirst({
    where: { opcId, taxType, period, deleted: false },
  });
  return Boolean(found);
}

export async function ensureTaxTasksForPeriod(
  prisma: PrismaClient,
  opcId: bigint,
  year: number,
  month: number,
) {
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

  if (!(await existingTask(prisma, opcId, 'vat', period))) {
    await prisma.taxFilingTask.create({
      data: { opcId, taxType: 'vat', period, dueDate, status: 'pending', calculatedAmount: vat },
    });
  }

  if (!(await existingTask(prisma, opcId, 'surcharge', period))) {
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
    const summaries = await prisma.profitSummary.findMany({
      where: { opcId, year, month: { gte: month - 2, lte: month } },
    });
    const cumulative = summaries.reduce((s, p) => s + Number(p.profit), 0);
    const cit = calculateCitQuarterly(cumulative);
    const qPeriod = `Q${Math.ceil(month / 3)}-${year}`;
    if (!(await existingTask(prisma, opcId, 'cit_quarterly', qPeriod))) {
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

export async function backfillLedgerForOpc(prisma: PrismaClient, opcId: bigint) {
  const entries = await prisma.expenseEntry.findMany({
    where: { opcId, deleted: false },
    select: { occurredAt: true },
  });
  const incomeEntries = await prisma.incomeEntry.findMany({
    where: { opcId, deleted: false },
    select: { occurredAt: true },
  });

  const periods = new Set<string>();
  for (const e of [...entries, ...incomeEntries]) {
    periods.add(periodKey(e.occurredAt));
  }

  for (const period of periods) {
    const [yearStr, monthStr] = period.split('-');
    const year = parseInt(yearStr, 10);
    const month = parseInt(monthStr, 10);
    await recalculateProfit(prisma, opcId, year, month);
    await generateVouchersForPeriod(prisma, opcId, period);
    await ensureTaxTasksForPeriod(prisma, opcId, year, month);
  }

  return periods.size;
}
