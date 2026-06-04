import {
  differenceInCalendarDays,
  endOfMonth,
  format,
  startOfDay,
  startOfMonth,
} from 'date-fns';
import { TAX_CHECKLIST_ITEMS, TAX_TYPE_LABELS } from '@/lib/api/constants';
import type {
  TaxCalendarDto,
  TaxCalendarTaskDto,
  TaxTaskBreakdownLine,
  TaxTaskDetailDto,
} from '@/lib/compliance/tax-calendar-shared';
import { prisma, serializeBigInt } from '@/lib/db';
import {
  calculateCitQuarterly,
  calculateSurcharge,
  calculateVat,
} from '@/lib/services/compliance/diagnosis/tax-calculator';
import { ensureTaxTasksForOpc } from '@/lib/services/compliance/ledger/ledger-service';

export type { TaxCalendarDto, TaxCalendarTaskDto, TaxTaskDetailDto } from '@/lib/compliance/tax-calendar-shared';

function taxTypeLabel(taxType: string): string {
  return TAX_TYPE_LABELS[taxType] ?? taxType;
}

function mapTask(task: {
  id: bigint;
  taxType: string;
  period: string;
  dueDate: Date;
  status: string;
  calculatedAmount: { toString(): string };
  filedAmount: { toString(): string } | null;
  receiptFileId: bigint | null;
}): TaxCalendarTaskDto {
  const due = task.dueDate;
  return {
    id: task.id.toString(),
    taxType: task.taxType,
    taxTypeLabel: taxTypeLabel(task.taxType),
    period: task.period,
    dueDate: format(due, 'yyyy-MM-dd'),
    dueDay: due.getDate(),
    status: task.status,
    calculatedAmount: Number(task.calculatedAmount),
    filedAmount: task.filedAmount != null ? Number(task.filedAmount) : null,
    receiptFileId: task.receiptFileId?.toString() ?? null,
  };
}

export async function refreshOverdueTaxTasks(opcId: bigint) {
  const today = startOfDay(new Date());
  await prisma.taxFilingTask.updateMany({
    where: {
      opcId,
      deleted: false,
      status: 'pending',
      dueDate: { lt: today },
    },
    data: { status: 'overdue' },
  });
}

export async function loadTaxCalendar(
  opcId: bigint,
  year: number,
  month: number,
): Promise<TaxCalendarDto> {
  await ensureTaxTasksForOpc(opcId, year, month);
  await refreshOverdueTaxTasks(opcId);

  const period = `${year}-${String(month).padStart(2, '0')}`;
  const monthStart = startOfMonth(new Date(year, month - 1, 1));
  const monthEnd = endOfMonth(monthStart);
  const quarter = Math.ceil(month / 3);
  const qPeriod = `Q${quarter}-${year}`;

  const rawTasks = await prisma.taxFilingTask.findMany({
    where: {
      opcId,
      deleted: false,
      OR: [
        { period },
        { period: qPeriod },
        { dueDate: { gte: monthStart, lte: monthEnd } },
      ],
    },
    orderBy: [{ dueDate: 'asc' }, { taxType: 'asc' }],
  });

  const tasks = rawTasks.map(mapTask);
  const dueDays = [...new Set(tasks.map((t) => t.dueDay))].sort((a, b) => a - b);

  const allPending = await prisma.taxFilingTask.findMany({
    where: { opcId, deleted: false, status: { in: ['pending', 'overdue'] } },
    orderBy: { dueDate: 'asc' },
    take: 1,
  });
  const next = allPending[0];
  const today = startOfDay(new Date());
  const nextDue = next
    ? {
        dueDate: format(next.dueDate, 'yyyy-MM-dd'),
        daysRemaining: differenceInCalendarDays(startOfDay(next.dueDate), today),
        taxTypeLabel: taxTypeLabel(next.taxType),
      }
    : null;

  return { year, month, dueDays, nextDue, tasks };
}

async function monthlyGrossRevenue(opcId: bigint, year: number, month: number) {
  const start = new Date(year, month - 1, 1);
  const end = new Date(year, month, 0, 23, 59, 59);
  const gross = await prisma.incomeEntry.aggregate({
    where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
    _sum: { grossAmount: true },
  });
  return Number(gross._sum.grossAmount ?? 0);
}

async function quarterCumulativeProfit(opcId: bigint, year: number, month: number) {
  const qStartMonth = Math.floor((month - 1) / 3) * 3 + 1;
  const summaries = await prisma.profitSummary.findMany({
    where: { opcId, year, month: { gte: qStartMonth, lte: month } },
  });
  return summaries.reduce((s, p) => s + Number(p.profit), 0);
}

export async function buildTaxTaskBreakdown(
  opcId: bigint,
  task: { taxType: string; period: string; calculatedAmount: { toString(): string } },
): Promise<TaxTaskBreakdownLine[]> {
  const amount = Number(task.calculatedAmount);

  if (task.taxType === 'vat') {
    const [y, m] = task.period.split('-').map((v) => parseInt(v, 10));
    const gross = await monthlyGrossRevenue(opcId, y, m);
    const net = gross / 1.01;
    const vat = calculateVat(gross);
    return [
      { label: '含税收入', amount: gross },
      { label: '不含税收入（÷1.01）', amount: net },
      { label: '增值税（预估）', amount: vat },
      { label: '系统测算税额', amount: amount },
    ];
  }

  if (task.taxType === 'surcharge') {
    const [y, m] = task.period.split('-').map((v) => parseInt(v, 10));
    const gross = await monthlyGrossRevenue(opcId, y, m);
    const vat = calculateVat(gross);
    const surcharge = calculateSurcharge(vat);
    return [
      { label: '增值税', amount: vat },
      { label: '附加税（增值税×6%）', amount: surcharge },
      { label: '系统测算税额', amount: amount },
    ];
  }

  if (task.taxType === 'cit_quarterly') {
    const match = task.period.match(/Q(\d)-(\d+)/);
    const quarter = match ? parseInt(match[1], 10) : 1;
    const year = match ? parseInt(match[2], 10) : new Date().getFullYear();
    const endMonth = quarter * 3;
    const profit = await quarterCumulativeProfit(opcId, year, endMonth);
    const cit = calculateCitQuarterly(profit);
    return [
      { label: '季度累计利润', amount: profit },
      { label: '企税预缴（小微 5%）', amount: cit },
      { label: '系统测算税额', amount: amount },
    ];
  }

  return [{ label: '系统测算税额', amount: amount }];
}

export async function getTaxTaskDetail(
  opcId: bigint,
  taskId: bigint,
): Promise<TaxTaskDetailDto | null> {
  const task = await prisma.taxFilingTask.findFirst({
    where: { id: taskId, opcId, deleted: false },
  });
  if (!task) return null;

  const breakdown = await buildTaxTaskBreakdown(opcId, task);
  const checklist = task.checklist as boolean[] | null;

  return {
    ...mapTask(task),
    breakdown,
    checklist: Array.isArray(checklist) ? checklist : null,
  };
}

export async function getMemberTaxChecklist(opcId: bigint) {
  const latestWithChecklist = await prisma.taxFilingTask.findFirst({
    where: { opcId, deleted: false, status: 'filed' },
    orderBy: { id: 'desc' },
  });

  const checked = latestWithChecklist?.checklist as boolean[] | undefined;
  const items = TAX_CHECKLIST_ITEMS.map((text, index) => ({
    id: index + 1,
    text,
    checked: checked?.[index] === true,
  }));

  const allChecked = items.every((i) => i.checked);
  const source = latestWithChecklist
    ? `${taxTypeLabel(latestWithChecklist.taxType)} · ${latestWithChecklist.period}`
    : null;

  return serializeBigInt({ items, allChecked, source });
}
