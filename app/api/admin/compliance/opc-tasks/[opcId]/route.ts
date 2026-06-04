import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getAdminFromRequest } from '@/lib/auth/admin';
import { getOpcTaskDetail } from '@/lib/services/compliance/opc/opc-service';

export async function GET(
  request: Request,
  { params }: { params: Promise<{ opcId: string }> },
) {
  try {
    await getAdminFromRequest(request);
    const { opcId } = await params;
    const detail = await getOpcTaskDetail(BigInt(opcId));
    return jsonOk(detail);
  } catch (error) {
    return handleApiError(error);
  }
}
