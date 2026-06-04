import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { getOpcProgress } from '@/lib/services/compliance/order/order-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const progress = await getOpcProgress(member.id);
    return jsonOk(progress);
  } catch (error) {
    return handleApiError(error);
  }
}
