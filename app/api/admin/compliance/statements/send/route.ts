import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getAdminFromRequest } from '@/lib/auth/admin';
import {
  generateMonthlyStatement,
  sendStatementNotification,
} from '@/lib/services/compliance/statement/statement-service';
import { prisma } from '@/lib/db';

export async function POST(request: Request) {
  try {
    await getAdminFromRequest(request);
    const body = await request.json();
    const { year, month, opcIds } = body as {
      year?: number;
      month?: number;
      opcIds?: string[];
    };

    if (year == null || month == null) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供 year 与 month');
    }

    const opcs = opcIds?.length
      ? await prisma.opcEntity.findMany({
          where: { id: { in: opcIds.map((id) => BigInt(id)) }, deleted: false },
        })
      : await prisma.opcEntity.findMany({
          where: { status: 'active', deleted: false },
        });

    let generated = 0;
    let notified = 0;
    const errors: string[] = [];

    for (const opc of opcs) {
      try {
        await generateMonthlyStatement(opc.id, year, month);
        generated++;
        await sendStatementNotification(opc.id, year, month);
        notified++;
      } catch (err) {
        errors.push(`${opc.id}: ${err instanceof Error ? err.message : 'unknown'}`);
      }
    }

    return jsonOk({ generated, notified, total: opcs.length, errors });
  } catch (error) {
    return handleApiError(error);
  }
}
