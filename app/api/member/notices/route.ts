import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { prisma, serializeBigInt } from '@/lib/db';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const notices = await prisma.memberNotice.findMany({
      where: { memberId: member.id, deleted: false },
      orderBy: { createdAt: 'desc' },
      take: 100,
    });
    return jsonOk(serializeBigInt(notices));
  } catch (error) {
    return handleApiError(error);
  }
}
