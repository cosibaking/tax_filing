import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { prisma } from '@/lib/db';
import { buildOpcProfileSummary } from '@/lib/services/compliance/opc/opc-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const order = await prisma.serviceOrder.findFirst({
      where: {
        memberId: member.id,
        deleted: false,
        status: { in: ['pending', 'active'] },
      },
    });
    const opc = await prisma.opcEntity.findFirst({
      where: { memberId: member.id, deleted: false },
    });
    return jsonOk(buildOpcProfileSummary(opc, !!order));
  } catch (error) {
    return handleApiError(error);
  }
}
