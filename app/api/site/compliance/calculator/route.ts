import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { calculateTaxComparison } from '@/lib/services/compliance/diagnosis/tax-calculator';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { annualRevenue, annualCost } = body as {
      annualRevenue?: number;
      annualCost?: number;
    };

    if (annualRevenue == null || annualCost == null) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供年收入与年成本');
    }

    const comparison = calculateTaxComparison({ annualRevenue, annualCost });
    return jsonOk({ comparison });
  } catch (error) {
    return handleApiError(error);
  }
}
