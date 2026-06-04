import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes, isDev } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { signContract } from '@/lib/services/compliance/order/order-service';

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const body = await request.json();
    const { orderId, signerName } = body as { orderId?: string; signerName?: string };

    let resolvedOrderId = orderId;
    let resolvedSignerName = signerName?.trim() ?? '';

    if (isDev) {
      resolvedSignerName = resolvedSignerName || member.name || '开发测试签署人';
    } else {
      if (!orderId) {
        return jsonFail(ErrorCodes.VALIDATION, '请提供订单号');
      }
      if (!resolvedSignerName) {
        return jsonFail(ErrorCodes.VALIDATION, '请提供签署人姓名');
      }
    }

    if (!resolvedOrderId) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供订单号');
    }

    const result = await signContract(
      member.id,
      BigInt(resolvedOrderId),
      resolvedSignerName,
      request,
    );
    return jsonOk(result);
  } catch (error) {
    return handleApiError(error);
  }
}
