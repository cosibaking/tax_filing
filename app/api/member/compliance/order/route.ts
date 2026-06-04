import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { createOrder } from '@/lib/services/compliance/order/order-service';

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const body = await request.json();
    const { planId, diagnosisId } = body as { planId?: string; diagnosisId?: string };

    if (!planId) {
      return jsonFail(ErrorCodes.ORDER_INVALID, '请选择套餐');
    }

    const order = await createOrder(
      member.id,
      BigInt(planId),
      diagnosisId ? BigInt(diagnosisId) : undefined,
    );
    return jsonOk(order);
  } catch (error) {
    return handleApiError(error);
  }
}
