import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { listStatements } from '@/lib/services/compliance/statement/statement-service';

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const statements = await listStatements(member.id);
    return jsonOk(statements);
  } catch (error) {
    return handleApiError(error);
  }
}
