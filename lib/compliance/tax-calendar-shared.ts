import { addMonths, format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

export type TaxCalendarTaskDto = {
  id: string;
  taxType: string;
  taxTypeLabel: string;
  period: string;
  dueDate: string;
  dueDay: number;
  status: string;
  calculatedAmount: number;
  filedAmount: number | null;
  receiptFileId: string | null;
};

export type TaxCalendarDto = {
  year: number;
  month: number;
  dueDays: number[];
  nextDue: {
    dueDate: string;
    daysRemaining: number;
    taxTypeLabel: string;
  } | null;
  tasks: TaxCalendarTaskDto[];
};

export type TaxTaskBreakdownLine = { label: string; amount: number };

export type TaxTaskDetailDto = TaxCalendarTaskDto & {
  breakdown: TaxTaskBreakdownLine[];
  checklist: boolean[] | null;
};

export function formatNextDueSummary(nextDue: TaxCalendarDto['nextDue']): string | null {
  if (!nextDue) return null;
  const dateLabel = format(new Date(nextDue.dueDate), 'M月d日', { locale: zhCN });
  const days =
    nextDue.daysRemaining < 0
      ? `已逾期 ${Math.abs(nextDue.daysRemaining)} 天`
      : nextDue.daysRemaining === 0
        ? '今天截止'
        : `${nextDue.daysRemaining} 天`;
  return `${dateLabel}（${days}）· ${nextDue.taxTypeLabel}`;
}

export function shiftCalendarMonth(year: number, month: number, delta: number) {
  const d = addMonths(new Date(year, month - 1, 1), delta);
  return { year: d.getFullYear(), month: d.getMonth() + 1 };
}
