import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { submitOpcMaterials } from '@/lib/services/compliance/opc/opc-service';

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const body = await request.json();

    const result = await submitOpcMaterials(member.id, body, request);
    return jsonOk(result);
  } catch (error) {
    if (error instanceof Error) {
      if (error.message.includes('Unknown argument')) {
        return jsonFail(
          ErrorCodes.VALIDATION,
          '服务数据模型未同步，请执行 npm run db:generate 后重启开发服务器',
        );
      }
      if (!error.message.includes('UNAUTHORIZED')) {
        return jsonFail(ErrorCodes.VALIDATION, error.message);
      }
    }
    return handleApiError(error);
  }
}
