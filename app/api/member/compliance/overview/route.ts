import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { getMemberOverview } from '@/lib/services/compliance/overview/member-overview-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const data = await getMemberOverview(member.id);
    return jsonOk(data);
  } catch (error) {
    return handleApiError(error);
  }
}
