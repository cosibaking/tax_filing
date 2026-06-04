import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { serializeBigInt } from '@/lib/db';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    return jsonOk(
      serializeBigInt({
        id: member.id,
        phone: member.phone,
        name: member.name,
        email: member.email,
        createdAt: member.createdAt,
      }),
    );
  } catch (error) {
    return handleApiError(error);
  }
}
