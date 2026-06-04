import 'server-only';
import { prisma, serializeBigInt } from '@/lib/db';
import {
  loadTaxCalendar,
  refreshOverdueTaxTasks,
} from '@/lib/services/compliance/ledger/tax-calendar-service';
import { buildMemberPendingTasks } from './pending-tasks';

async function loadMonthlySummary(opcId: bigint) {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth() + 1;
  const start = new Date(year, month - 1, 1);
  const end = new Date(year, month, 0, 23, 59, 59, 999);

  const [incomeSum, expenseSum] = await Promise.all([
    prisma.incomeEntry.aggregate({
      where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
      _sum: { grossAmount: true },
    }),
    prisma.expenseEntry.aggregate({
      where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
      _sum: { amount: true },
    }),
  ]);

  const income = Number(incomeSum._sum.grossAmount ?? 0);
  const expense = Number(expenseSum._sum.amount ?? 0);
  return { income, expense, profit: income - expense };
}

export async function getMemberOverview(memberId: bigint) {
  const [activeOrder, opc, pendingTasks] = await Promise.all([
    prisma.serviceOrder.findFirst({
      where: { memberId, deleted: false, status: 'active' },
    }),
    prisma.opcEntity.findFirst({ where: { memberId, deleted: false } }),
    buildMemberPendingTasks(memberId),
  ]);

  const opcStatus = opc?.status ?? null;
  let nextTaxDeadline: string | undefined;
  let monthlySummary: { income: number; expense: number; profit: number } | undefined;

  if (opc && opc.status === 'active') {
    const now = new Date();
    await refreshOverdueTaxTasks(opc.id);
    const calendar = await loadTaxCalendar(opc.id, now.getFullYear(), now.getMonth() + 1);
    if (calendar.nextDue) {
      const { dueDate, daysRemaining, taxTypeLabel: label } = calendar.nextDue;
      const remain =
        daysRemaining < 0
          ? `已逾期 ${Math.abs(daysRemaining)} 天`
          : daysRemaining === 0
            ? '今日截止'
            : `剩余 ${daysRemaining} 天`;
      nextTaxDeadline = `${label} · ${dueDate}（${remain}）`;
    }

    monthlySummary = await loadMonthlySummary(opc.id);
  }

  return serializeBigInt({
    opcStatus,
    hasOrder: Boolean(activeOrder),
    pendingTasks,
    nextTaxDeadline,
    monthlySummary,
  });
}
