import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { ensureTaxTasksForOpc } from '@/lib/services/compliance/ledger/ledger-service';
import { prisma } from '@/lib/db';

function verifyCron(request: Request) {
  const secret = process.env.CRON_SECRET;
  const auth = request.headers.get('authorization');
  if (!secret || auth !== `Bearer ${secret}`) {
    throw new Error('UNAUTHORIZED');
  }
}

export async function POST(request: Request) {
  try {
    verifyCron(request);

    const now = new Date();
    const year = now.getFullYear();
    const month = now.getMonth() + 1;

    const opcs = await prisma.opcEntity.findMany({
      where: { status: 'active', deleted: false },
    });

    for (const opc of opcs) {
      await ensureTaxTasksForOpc(opc.id, year, month);
    }

    return jsonOk({ processed: opcs.length, year, month });
  } catch (error) {
    if (error instanceof Error && error.message === 'UNAUTHORIZED') {
      return jsonFail(ErrorCodes.UNAUTHORIZED, '未授权', 401);
    }
    return handleApiError(error);
  }
}
