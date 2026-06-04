import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import {
  verifyAdminPassword,
  signAdminToken,
  setAdminSession,
} from '@/lib/auth/admin';
import { prisma } from '@/lib/db';

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const { username, password } = body as { username?: string; password?: string };

    if (!username || !password) {
      return jsonFail(ErrorCodes.VALIDATION, '请输入用户名和密码');
    }

    const admin = await prisma.adminUser.findFirst({
      where: { username, deleted: false },
    });
    if (!admin || !(await verifyAdminPassword(password, admin.passwordHash))) {
      return jsonFail(ErrorCodes.AUTH_INVALID, '用户名或密码错误', 401);
    }

    const token = await signAdminToken(admin.id.toString(), admin.username, admin.role);
    await setAdminSession(token);

    return jsonOk({
      token,
      admin: {
        id: admin.id.toString(),
        username: admin.username,
        role: admin.role,
      },
    });
  } catch (error) {
    return handleApiError(error);
  }
}
