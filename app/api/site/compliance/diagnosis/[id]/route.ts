import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getDiagnosisByIdOrGuest } from '@/lib/services/compliance/diagnosis/diagnosis-service';

export async function GET(
  _request: Request,
  { params }: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await params;
    const diagnosis = await getDiagnosisByIdOrGuest(id);
    if (!diagnosis) {
      return jsonFail(ErrorCodes.NOT_FOUND, '诊断记录不存在或已过期');
    }
    return jsonOk(diagnosis);
  } catch (error) {
    return handleApiError(error);
  }
}
