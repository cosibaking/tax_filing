import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { createDiagnosis } from '@/lib/services/compliance/diagnosis/diagnosis-service';
import { verifyMemberToken } from '@/lib/auth/member';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const {
      platforms,
      monthlyIncomeRange,
      annualCostEstimate,
      existingEntity,
      hasFiledTax,
      taxBureauContact,
    } = body as {
      platforms?: string[];
      monthlyIncomeRange?: string;
      annualCostEstimate?: number;
      existingEntity?: string;
      hasFiledTax?: boolean;
      taxBureauContact?: boolean;
    };

    if (!platforms?.length || !monthlyIncomeRange || annualCostEstimate == null || !existingEntity) {
      return jsonFail(ErrorCodes.DIAGNOSIS_INVALID, '诊断参数不完整');
    }

    let memberId: bigint | undefined;
    const auth = request.headers.get('authorization');
    if (auth?.startsWith('Bearer ')) {
      try {
        const payload = await verifyMemberToken(auth.slice(7));
        if (payload.type === 'access') memberId = BigInt(payload.sub);
      } catch {
        /* anonymous diagnosis */
      }
    }

    const result = await createDiagnosis({
      memberId,
      platforms,
      monthlyIncomeRange,
      annualCostEstimate,
      existingEntity,
      hasFiledTax: !!hasFiledTax,
      taxBureauContact: !!taxBureauContact,
    });

    return jsonOk(result);
  } catch (error) {
    return handleApiError(error);
  }
}
