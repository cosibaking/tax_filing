import { NextResponse } from 'next/server';
import { jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { getAdminFromRequest } from '@/lib/auth/admin';
import { getAttachment, readUpload } from '@/lib/storage/local';

type RouteParams = { params: Promise<{ id: string }> };

async function assertDownloadAccess(request: Request, attachmentMemberId: bigint | null) {
  try {
    const admin = await getAdminFromRequest(request);
    return admin;
  } catch {
    const member = await getMemberFromRequest(request);
    if (attachmentMemberId && attachmentMemberId !== member.id) {
      throw new Error('FORBIDDEN');
    }
    return member;
  }
}

export async function GET(request: Request, { params }: RouteParams) {
  try {
    const { id } = await params;
    const attachment = await getAttachment(BigInt(id));
    if (!attachment) {
      return jsonFail(ErrorCodes.NOT_FOUND, '附件不存在', 404);
    }

    await assertDownloadAccess(request, attachment.memberId);

    const buffer = await readUpload(attachment.storageKey);
    return new NextResponse(new Uint8Array(buffer), {
      headers: {
        'Content-Type': attachment.mimeType,
        'Content-Disposition': `attachment; filename="${encodeURIComponent(attachment.fileName)}"`,
        'Content-Length': String(buffer.length),
      },
    });
  } catch (error) {
    return handleApiError(error);
  }
}
