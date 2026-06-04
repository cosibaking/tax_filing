import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { hashPassword, signMemberToken, getClientIp, getUserAgent } from '@/lib/auth/member';
import { prisma } from '@/lib/db';

const PHONE_RE = /^1[3-9]\d{9}$/;

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { phone, password, agreeTerms, agreePrivacy } = body as {
      phone?: string;
      password?: string;
      agreeTerms?: boolean;
      agreePrivacy?: boolean;
    };

    if (!phone || !PHONE_RE.test(phone)) {
      return jsonFail(ErrorCodes.VALIDATION, '手机号格式不正确');
    }
    if (!password || password.length < 6) {
      return jsonFail(ErrorCodes.VALIDATION, '密码至少 6 位');
    }
    if (!agreeTerms || !agreePrivacy) {
      return jsonFail(ErrorCodes.VALIDATION, '请同意服务条款与隐私政策');
    }

    const existing = await prisma.member.findFirst({ where: { phone, deleted: false } });
    if (existing) {
      return jsonFail(ErrorCodes.AUTH_INVALID, '该手机号已注册');
    }

    const termsVersion = process.env.LEGAL_TERMS_VERSION ?? '1.0';
    const privacyVersion = process.env.LEGAL_PRIVACY_VERSION ?? '1.0';
    const passwordHash = await hashPassword(password);
    const ip = getClientIp(request);
    const userAgent = getUserAgent(request);
    const now = new Date();

    const member = await prisma.member.create({
      data: {
        phone,
        passwordHash,
        consents: {
          create: [
            { type: 'terms', documentVersion: termsVersion, ip, userAgent, agreedAt: now },
            { type: 'privacy', documentVersion: privacyVersion, ip, userAgent, agreedAt: now },
          ],
        },
      },
    });

    const sub = member.id.toString();
    const accessToken = await signMemberToken({ sub, phone }, 'access');
    const refreshToken = await signMemberToken({ sub, phone }, 'refresh');

    return jsonOk({ accessToken, refreshToken, memberId: sub });
  } catch (error) {
    return handleApiError(error);
  }
}
