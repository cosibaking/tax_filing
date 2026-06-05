import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { getOpcMemberFullDetail } from '@/lib/services/compliance/opc/opc-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const data = await getOpcMemberFullDetail(member.id, request);
    return jsonOk(data);
  } catch (error) {
    return handleApiError(error);
  }
}
