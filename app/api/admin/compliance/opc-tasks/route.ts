import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getAdminFromRequest } from '@/lib/auth/admin';
import { updateOpcProgress } from '@/lib/services/compliance/order/order-service';
import { prisma, serializeBigInt } from '@/lib/db';

export async function GET(request: Request) {
  try {
    await getAdminFromRequest(request);
    const params = new URL(request.url).searchParams;
    const status = params.get('status');
    const q = params.get('q')?.trim();

    const tasks = await prisma.opcEntity.findMany({
      where: {
        deleted: false,
        ...(status ? { status } : {}),
        ...(q
          ? {
              OR: [
                { companyName: { contains: q } },
                { phone: { contains: q } },
                { member: { phone: { contains: q } } },
              ],
            }
          : {}),
      },
      include: {
        member: { select: { id: true, phone: true, name: true } },
        progressLogs: { orderBy: { createdAt: 'desc' }, take: 5 },
      },
      orderBy: { updatedAt: 'desc' },
      take: 200,
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
    const { opcId, step, status, note } = body as {
      opcId?: string;
      step?: string;
      status?: string;
      note?: string;
    };

    if (!opcId || !step || !status) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供 opcId、step、status');
    }

    await updateOpcProgress(
      BigInt(opcId),
      step,
      status,
      note ?? '',
      admin.id,
    );

    return jsonOk({ updated: true });
  } catch (error) {
    return handleApiError(error);
  }
}
