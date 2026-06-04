import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getAdminFromRequest } from '@/lib/auth/admin';
import { prisma, serializeBigInt } from '@/lib/db';

export async function GET(request: Request) {
  try {
    await getAdminFromRequest(request);
    const params = new URL(request.url).searchParams;
    const q = params.get('q')?.trim();

    const members = await prisma.member.findMany({
      where: {
        deleted: false,
        ...(q
          ? {
              OR: [
                { phone: { contains: q } },
                { name: { contains: q } },
                { opcEntity: { companyName: { contains: q } } },
              ],
            }
          : {}),
      },
      include: {
        opcEntity: { where: { deleted: false } },
        orders: {
          where: { deleted: false },
          orderBy: { createdAt: 'desc' },
          take: 1,
          include: { plan: true },
        },
      },
      orderBy: { createdAt: 'desc' },
      take: 200,
    });

    return jsonOk(serializeBigInt(members));
  } catch (error) {
    return handleApiError(error);
  }
}
