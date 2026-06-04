import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { sendMaterialReminder } from '@/lib/services/compliance/statement/statement-service';
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

    const opcs = await prisma.opcEntity.findMany({
      where: { status: 'active', deleted: false },
      select: { memberId: true },
    });

    for (const opc of opcs) {
      await sendMaterialReminder(opc.memberId);
    }

    return jsonOk({ sent: opcs.length });
  } catch (error) {
    if (error instanceof Error && error.message === 'UNAUTHORIZED') {
      return jsonFail(ErrorCodes.UNAUTHORIZED, '未授权', 401);
    }
    return handleApiError(error);
  }
}
