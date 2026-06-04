import 'server-only';
import { format } from 'date-fns';
import { TAX_TYPE_LABELS } from '@/lib/api/constants';
import { prisma } from '@/lib/db';
import { refreshOverdueTaxTasks } from '@/lib/services/compliance/ledger/tax-calendar-service';

export type PendingTaskDto = {
  key: string;
  label: string;
  href: string;
  urgent?: boolean;
};

function taxTypeLabel(taxType: string): string {
  return TAX_TYPE_LABELS[taxType] ?? taxType;
}

/** 每月 1–10 日提醒上传上月银行流水（已有流水则不再展示） */
async function maybeMonthlyMaterialTask(opcId: bigint): Promise<PendingTaskDto | null> {
  const now = new Date();
  const day = now.getDate();
  if (day > 10) return null;

  const prevMonth = now.getMonth() === 0 ? 12 : now.getMonth();
  const year = now.getMonth() === 0 ? now.getFullYear() - 1 : now.getFullYear();
  const start = new Date(year, prevMonth - 1, 1);
  const end = new Date(year, prevMonth, 0, 23, 59, 59, 999);

  const bankCount = await prisma.bankTransaction.count({
    where: { opcId, deleted: false, occurredAt: { gte: start, lte: end } },
  });
  if (bankCount > 0) return null;

  const deadline = format(new Date(now.getFullYear(), now.getMonth(), 5), 'M月d日');
  return {
    key: `monthly_materials_${year}_${prevMonth}`,
    label: `上传${year}年${prevMonth}月银行流水（建议 ${deadline} 前）`,
    href: '/user/compliance/income',
    urgent: day <= 5,
  };
}

/**
 * 根据订单、OPC、申报等业务状态推导待办（不落独立待办表）。
 * 已完成/终态事项不返回。
 */
export async function buildMemberPendingTasks(memberId: bigint): Promise<PendingTaskDto[]> {
  const tasks: PendingTaskDto[] = [];

  const activeOrder = await prisma.serviceOrder.findFirst({
    where: { memberId, deleted: false, status: 'active' },
  });

  if (!activeOrder) {
    tasks.push({
      key: 'sign_contract',
      label: '完成方案签约',
      href: '/user/compliance/plan',
      urgent: true,
    });
    return tasks;
  }

  const opc = await prisma.opcEntity.findFirst({
    where: { memberId, deleted: false },
  });
  if (!opc) return tasks;

  if (opc.status === 'pending' && !opc.materialsSubmittedAt) {
    tasks.push({
      key: 'submit_opc_materials',
      label: '提交 OPC 设立资料',
      href: '/user/compliance/opc',
      urgent: true,
    });
  }

  if (opc.status === 'materials') {
    tasks.push({
      key: 'resubmit_opc_materials',
      label: '补正 OPC 设立资料',
      href: '/user/compliance/opc',
      urgent: true,
    });
  }

  if (opc.status === 'bank' && !opc.bankReceiptFileId) {
    tasks.push({
      key: 'upload_bank_receipt',
      label: '上传银行开户回执',
      href: '/user/compliance/opc',
      urgent: false,
    });
  }

  if (opc.status === 'active') {
    const materialTask = await maybeMonthlyMaterialTask(opc.id);
    if (materialTask) tasks.push(materialTask);

    await refreshOverdueTaxTasks(opc.id);
    const pendingTax = await prisma.taxFilingTask.findMany({
      where: {
        opcId: opc.id,
        deleted: false,
        status: { in: ['pending', 'overdue'] },
      },
      orderBy: { dueDate: 'asc' },
      take: 5,
    });

    for (const t of pendingTax) {
      tasks.push({
        key: `tax_${t.id}`,
        label: `${taxTypeLabel(t.taxType)}申报（截止 ${format(t.dueDate, 'M/d')}）`,
        href: '/user/compliance/tax',
        urgent: t.status === 'overdue',
      });
    }
  }

  return tasks;
}
