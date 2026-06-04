import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { clearAdminSession } from '@/lib/auth/admin';

export async function POST() {
  try {
    await clearAdminSession();
    return jsonOk({ loggedOut: true });
  } catch (error) {
    return handleApiError(error);
  }
}
