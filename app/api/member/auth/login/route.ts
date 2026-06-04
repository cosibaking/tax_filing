import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { verifyPassword, signMemberToken } from '@/lib/auth/member';
import { verifyOtp } from '@/lib/auth/otp';
import { prisma } from '@/lib/db';
import { persistClaimedGuestDiagnoses } from '@/lib/services/compliance/diagnosis/claim-guest-diagnoses';
import { resolveGuestSessionId } from '@/lib/diagnosis/guest-session';

const PHONE_RE = /^1[3-9]\d{9}$/;

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { phone, password, otp, loginType, guestSessionId } = body as {
      phone?: string;
      password?: string;
      otp?: string;
      loginType?: string;
      guestSessionId?: string;
    };

    if (!phone || !PHONE_RE.test(phone)) {
      return jsonFail(ErrorCodes.VALIDATION, '手机号格式不正确');
    }

    const member = await prisma.member.findFirst({ where: { phone, deleted: false } });
    if (!member) {
      return jsonFail(ErrorCodes.AUTH_INVALID, '账号不存在');
    }

    if (loginType === 'otp') {
      if (!otp) return jsonFail(ErrorCodes.VALIDATION, '请输入验证码');
      const valid = await verifyOtp(phone, 'login', otp);
      if (!valid) return jsonFail(ErrorCodes.AUTH_OTP_EXPIRED, '验证码无效或已过期');
    } else {
      if (!password) return jsonFail(ErrorCodes.VALIDATION, '请输入密码');
      if (!member.passwordHash || !(await verifyPassword(password, member.passwordHash))) {
        return jsonFail(ErrorCodes.AUTH_INVALID, '手机号或密码错误');
      }
    }

    const sub = member.id.toString();
    const accessToken = await signMemberToken({ sub, phone }, 'access');
    const refreshToken = await signMemberToken({ sub, phone }, 'refresh');

    const guestSession = resolveGuestSessionId(guestSessionId, request.headers.get('cookie'));
    const claimedDiagnosisIds = await persistClaimedGuestDiagnoses(guestSession, member.id);

    return jsonOk({ accessToken, refreshToken, memberId: sub, claimedDiagnosisIds });
  } catch (error) {
    return handleApiError(error);
  }
}
