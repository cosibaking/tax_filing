import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { signMemberToken, verifyMemberToken } from '@/lib/auth/member';
import { prisma } from '@/lib/db';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { refreshToken } = body as { refreshToken?: string };

    if (!refreshToken) {
      return jsonFail(ErrorCodes.VALIDATION, '缺少 refreshToken');
    }

    const payload = await verifyMemberToken(refreshToken);
    if (payload.type !== 'refresh') {
      return jsonFail(ErrorCodes.UNAUTHORIZED, '无效的刷新令牌', 401);
    }

    const member = await prisma.member.findFirst({
      where: { id: BigInt(payload.sub), deleted: false },
    });
    if (!member) {
      return jsonFail(ErrorCodes.UNAUTHORIZED, '未授权', 401);
    }

    const sub = member.id.toString();
    const accessToken = await signMemberToken({ sub, phone: member.phone }, 'access');

    return jsonOk({ accessToken });
  } catch (error) {
    return handleApiError(error);
  }
}
