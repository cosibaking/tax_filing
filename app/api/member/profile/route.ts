import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import {
  serializeMemberProfile,
  updateMemberProfile,
} from '@/lib/services/member/profile-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    return jsonOk(serializeMemberProfile(member));
  } catch (error) {
    return handleApiError(error);
  }
}

export async function PATCH(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const body = (await request.json()) as { name?: string; avatarFileId?: string | null };
    const data = await updateMemberProfile(member.id, body, request);
    return jsonOk(data, '资料已更新');
  } catch (error) {
    return handleApiError(error);
  }
}
