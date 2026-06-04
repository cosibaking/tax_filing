import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import {
  generateMonthlyStatement,
  sendStatementNotification,
} from '@/lib/services/compliance/statement/statement-service';
import { prisma } from '@/lib/db';

function verifyCron(request: Request) {
  const secret = process.env.CRON_SECRET;
  const auth = request.headers.get('authorization');
  if (!secret || auth !== `Bearer ${secret}`) {
    throw new Error('UNAUTHORIZED');
  }
}

function previousMonth() {
  const now = new Date();
  if (now.getMonth() === 0) {
    return { year: now.getFullYear() - 1, month: 12 };
  }
  return { year: now.getFullYear(), month: now.getMonth() };
}

export async function POST(request: Request) {
  try {
    verifyCron(request);

    const { year, month } = previousMonth();
    const opcs = await prisma.opcEntity.findMany({
      where: { status: 'active', deleted: false },
    });

    let generated = 0;
    const errors: string[] = [];

    for (const opc of opcs) {
      try {
        await generateMonthlyStatement(opc.id, year, month);
        await sendStatementNotification(opc.id, year, month);
        generated++;
      } catch (err) {
        errors.push(`${opc.id}: ${err instanceof Error ? err.message : 'unknown'}`);
      }
    }

    return jsonOk({ year, month, generated, total: opcs.length, errors });
  } catch (error) {
    if (error instanceof Error && error.message === 'UNAUTHORIZED') {
      return jsonFail(ErrorCodes.UNAUTHORIZED, '未授权', 401);
    }
    return handleApiError(error);
  }
}
