import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getAdminFromRequest } from '@/lib/auth/admin';
import { markFilingFiled } from '@/lib/services/compliance/ledger/ledger-service';
import { prisma, serializeBigInt } from '@/lib/db';

export async function GET(request: Request) {
  try {
    await getAdminFromRequest(request);
    const params = new URL(request.url).searchParams;
    const status = params.get('status');
    const taxType = params.get('taxType');
    const period = params.get('period');

    const tasks = await prisma.taxFilingTask.findMany({
      where: {
        deleted: false,
        ...(status ? { status } : {}),
        ...(taxType ? { taxType } : {}),
        ...(period ? { period } : {}),
      },
      include: {
        opc: {
          include: { member: { select: { id: true, phone: true, name: true } } },
        },
      },
      orderBy: { dueDate: 'asc' },
      take: 500,
    });

    return jsonOk(serializeBigInt(tasks));
  } catch (error) {
    return handleApiError(error);
  }
}

export async function PATCH(request: Request) {
  try {
    const admin = await getAdminFromRequest(request);
    const body = await request.json();
    const { taskId, filedAmount, checklist, receiptFileId } = body as {
      taskId?: string;
      filedAmount?: number;
      checklist?: boolean[];
      receiptFileId?: string;
    };

    if (!taskId || filedAmount == null || !checklist) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供 taskId、filedAmount、checklist');
    }

    await markFilingFiled(
      BigInt(taskId),
      admin.id,
      filedAmount,
      checklist,
      receiptFileId ? BigInt(receiptFileId) : undefined,
    );

    return jsonOk({ filed: true });
  } catch (error) {
    return handleApiError(error);
  }
}
