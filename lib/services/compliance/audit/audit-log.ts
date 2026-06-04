import { prisma } from '@/lib/db';

export async function writeAuditLog(params: {
  entityType: string;
  entityId: bigint;
  action: string;
  operatorId: bigint;
  operatorType: 'member' | 'admin';
  before?: unknown;
  after?: unknown;
  ip?: string | null;
}) {
  await prisma.complianceAuditLog.create({
    data: {
      entityType: params.entityType,
      entityId: params.entityId,
      action: params.action,
      operatorId: params.operatorId,
      operatorType: params.operatorType,
      before: params.before ? (params.before as object) : undefined,
      after: params.after ? (params.after as object) : undefined,
      ip: params.ip ?? undefined,
    },
  });
}

export async function exportAuditForMember(memberId: bigint) {
  const consents = await prisma.complianceConsent.findMany({
    where: { memberId },
    orderBy: { agreedAt: 'desc' },
  });
  const opc = await prisma.opcEntity.findFirst({ where: { memberId, deleted: false } });
  const auditLogs = opc
    ? await prisma.complianceAuditLog.findMany({
        where: { entityId: opc.id },
        orderBy: { createdAt: 'desc' },
        take: 500,
      })
    : [];
  return { consents, auditLogs };
}
