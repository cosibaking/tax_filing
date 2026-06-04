import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

async function aggregateMonthProfit(opcId: bigint, year: number, month: number) {
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

async function recalculateProfit(opcId: bigint, year: number, month: number) {
  const { revenue, cost, profit } = await aggregateMonthProfit(opcId, year, month);
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

async function generateVouchersForPeriod(opcId: bigint, period: string) {
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

async function main() {
  const opc =
    (await prisma.opcEntity.findFirst({
      where: { deleted: false, status: 'active' },
      orderBy: { id: 'desc' },
    })) ??
    (await prisma.opcEntity.findFirst({
      where: { deleted: false },
      orderBy: { id: 'desc' },
    }));

  if (!opc) {
    throw new Error('未找到 OPC 主体');
  }

  const entries = await prisma.expenseEntry.findMany({
    where: { opcId: opc.id, deleted: false },
    select: { occurredAt: true },
  });
  const incomeEntries = await prisma.incomeEntry.findMany({
    where: { opcId: opc.id, deleted: false },
    select: { occurredAt: true },
  });

  const periods = new Set<string>();
  for (const e of [...entries, ...incomeEntries]) {
    const y = e.occurredAt.getFullYear();
    const m = e.occurredAt.getMonth() + 1;
    periods.add(`${y}-${String(m).padStart(2, '0')}`);
  }

  for (const period of periods) {
    const [yearStr, monthStr] = period.split('-');
    const year = parseInt(yearStr, 10);
    const month = parseInt(monthStr, 10);
    await recalculateProfit(opc.id, year, month);
    await generateVouchersForPeriod(opc.id, period);
  }

  console.log(`已回填 ${periods.size} 个期间的利润汇总与会计分录（OPC #${opc.id}）`);
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
