import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
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

    const result = await prisma.taxFilingTask.updateMany({
      where: {
        deleted: false,
        status: 'pending',
        dueDate: { lt: new Date() },
      },
      data: { status: 'overdue' },
    });

    return jsonOk({ marked: result.count });
  } catch (error) {
    if (error instanceof Error && error.message === 'UNAUTHORIZED') {
      return jsonFail(ErrorCodes.UNAUTHORIZED, '未授权', 401);
    }
    return handleApiError(error);
  }
}
