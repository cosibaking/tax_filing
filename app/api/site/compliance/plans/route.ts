import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getActivePlans } from '@/lib/services/compliance/order/order-service';

export async function GET() {
  try {
    const plans = await getActivePlans();
    return jsonOk(plans);
  } catch (error) {
    return handleApiError(error);
  }
}
