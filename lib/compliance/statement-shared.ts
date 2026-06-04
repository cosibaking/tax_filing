/** 会员端对账单 DTO（与 GET /api/member/compliance/statements 一致） */
export interface StatementSummaryDto {
  revenue: number;
  cost: number;
  profit: number;
  cumulativeProfit: number;
  taxTotal: number;
  filingStatus: string;
}

export interface MemberStatementDto {
  id: string;
  year: number;
  month: number;
  period: string;
  status: 'sent' | 'ready';
  statusLabel: string;
  serviceComplete: boolean;
  sentAt: string | null;
  hasPdf: boolean;
  pdfFileId: string | null;
  summary: StatementSummaryDto;
}

export function formatStatementMoney(n: number): string {
  return `¥${n.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

export function statementStatusBadgeVariant(
  status: MemberStatementDto['status'],
  serviceComplete: boolean,
): 'success' | 'warning' | 'secondary' {
  if (serviceComplete) return 'success';
  if (status === 'sent') return 'success';
  return 'warning';
}
