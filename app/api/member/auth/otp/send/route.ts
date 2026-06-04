import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { sendOtp } from '@/lib/auth/otp';

const PHONE_RE = /^1[3-9]\d{9}$/;
const VALID_PURPOSES = ['login', 'register'];

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { phone, purpose } = body as { phone?: string; purpose?: string };

    if (!phone || !PHONE_RE.test(phone)) {
      return jsonFail(ErrorCodes.VALIDATION, '手机号格式不正确');
    }
    if (!purpose || !VALID_PURPOSES.includes(purpose)) {
      return jsonFail(ErrorCodes.VALIDATION, 'purpose 无效');
    }

    await sendOtp(phone, purpose);
    return jsonOk({ sent: true });
  } catch (error) {
    return handleApiError(error);
  }
}
