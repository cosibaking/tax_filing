import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { submitOpcMaterials } from '@/lib/services/compliance/order/order-service';

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const body = await request.json();
    const { idCard, phone, email, businessScope, registerAddress } = body as {
      idCard?: string;
      phone?: string;
      email?: string;
      businessScope?: string;
      registerAddress?: string;
    };

    if (!idCard || !phone || !email || !businessScope || !registerAddress) {
      return jsonFail(ErrorCodes.VALIDATION, '请填写完整资料');
    }

    const result = await submitOpcMaterials(
      member.id,
      { idCard, phone, email, businessScope, registerAddress },
      request,
    );
    return jsonOk(result);
  } catch (error) {
    return handleApiError(error);
  }
}
