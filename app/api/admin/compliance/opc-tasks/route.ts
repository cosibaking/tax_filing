import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getAdminFromRequest } from '@/lib/auth/admin';
import {
  applyOpcAdminAction,
  listOpcTasks,
  type OpcAdminAction,
} from '@/lib/services/compliance/opc/opc-service';

export async function GET(request: Request) {
  try {
    await getAdminFromRequest(request);
    const params = new URL(request.url).searchParams;
    const status = params.get('status') ?? undefined;
    const q = params.get('q')?.trim() || params.get('search')?.trim() || undefined;

    const tasks = await listOpcTasks({ status, q });
    return jsonOk(tasks);
  } catch (error) {
    return handleApiError(error);
  }
}

export async function PATCH(request: Request) {
  try {
    const admin = await getAdminFromRequest(request);
    const body = await request.json();
    const { opcId, action, note, payload } = body as {
      opcId?: string;
      action?: OpcAdminAction;
      note?: string;
      payload?: Record<string, string>;
    };

    if (!opcId || !action) {
      return jsonFail(ErrorCodes.VALIDATION, '请提供 opcId 与 action');
    }

    const result = await applyOpcAdminAction(
      BigInt(opcId),
      action,
      admin.id,
      { note, ...payload },
    );
    return jsonOk(result);
  } catch (error) {
    if (error instanceof Error) {
      return jsonFail(ErrorCodes.VALIDATION, error.message);
    }
    return handleApiError(error);
  }
}
