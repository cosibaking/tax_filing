import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { submitBankReceipt } from '@/lib/services/compliance/opc/opc-service';

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const body = await request.json();
    const { bankName, bankReceiptFileId } = body as {
      bankName?: string;
      bankReceiptFileId?: string;
    };

    if (!bankReceiptFileId) {
      return jsonFail(ErrorCodes.VALIDATION, '请上传开户回执');
    }

    const result = await submitBankReceipt(
      member.id,
      { bankName, bankReceiptFileId },
      request,
    );
    return jsonOk(result);
  } catch (error) {
    if (error instanceof Error) {
      return jsonFail(ErrorCodes.VALIDATION, error.message);
    }
    return handleApiError(error);
  }
}
