import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { saveUpload } from '@/lib/storage/local';

export async function POST(request: Request) {
  try {
    let memberId: bigint | undefined;
    try {
      const member = await getMemberFromRequest(request);
      memberId = member.id;
    } catch {
      /* allow anonymous upload for pre-login flows if needed */
    }

    const formData = await request.formData();
    const file = formData.get('file');
    if (!(file instanceof File)) {
      return jsonFail(ErrorCodes.VALIDATION, '请上传文件');
    }

    const result = await saveUpload(file, memberId);
    return jsonOk(result);
  } catch (error) {
    return handleApiError(error);
  }
}
