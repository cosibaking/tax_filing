import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { getDiagnosisHistory } from '@/lib/services/compliance/diagnosis/diagnosis-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const history = await getDiagnosisHistory(member.id);
    return jsonOk(history);
  } catch (error) {
    return handleApiError(error);
  }
}
