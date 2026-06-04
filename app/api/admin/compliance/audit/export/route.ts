import { NextResponse } from 'next/server';
import { jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getAdminFromRequest } from '@/lib/auth/admin';
import { exportAuditForMember } from '@/lib/services/compliance/audit/audit-log';

function escapeCsv(value: unknown): string {
  const s = value == null ? '' : String(value);
  if (s.includes(',') || s.includes('"') || s.includes('\n')) {
    return `"${s.replace(/"/g, '""')}"`;
  }
  return s;
}

export async function GET(request: Request) {
  try {
    await getAdminFromRequest(request);
    const memberId = new URL(request.url).searchParams.get('memberId');
    if (!memberId) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供 memberId');
    }

    const { consents, auditLogs } = await exportAuditForMember(BigInt(memberId));

    const lines: string[] = ['section,id,type,action,operator,createdAt,detail'];
    for (const c of consents) {
      lines.push(
        [
          'consent',
          c.id.toString(),
          c.type,
          '',
          c.memberId.toString(),
          c.agreedAt.toISOString(),
          c.documentVersion,
        ]
          .map(escapeCsv)
          .join(','),
      );
    }
    for (const log of auditLogs) {
      lines.push(
        [
          'audit',
          log.id.toString(),
          log.entityType,
          log.action,
          `${log.operatorType}:${log.operatorId}`,
          log.createdAt.toISOString(),
          JSON.stringify(log.after ?? log.before ?? {}),
        ]
          .map(escapeCsv)
          .join(','),
      );
    }

    return new NextResponse(lines.join('\n'), {
      status: 200,
      headers: {
        'Content-Type': 'text/csv; charset=utf-8',
        'Content-Disposition': `attachment; filename="audit-${memberId}.csv"`,
      },
    });
  } catch (error) {
    return handleApiError(error);
  }
}
