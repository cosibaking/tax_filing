import { NextResponse } from 'next/server';
import { jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { getAdminFromRequest } from '@/lib/auth/admin';
import { getAttachment, readUpload } from '@/lib/storage/local';
import { assertAntiHotlink, isImageMime, verifyMediaSignature } from '@/lib/media/signed-access';

type RouteParams = { params: Promise<{ id: string }> };

async function assertAttachmentAccess(
  request: Request,
  attachmentId: bigint,
  attachmentMemberId: bigint | null,
) {
  const url = new URL(request.url);
  const uid = url.searchParams.get('uid');
  const exp = url.searchParams.get('exp');
  const sig = url.searchParams.get('sig');

  if (uid && exp && sig) {
    if (!verifyMediaSignature(attachmentId.toString(), uid, exp, sig)) {
      throw new Error('FORBIDDEN');
    }
    if (attachmentMemberId != null && attachmentMemberId.toString() !== uid) {
      throw new Error('FORBIDDEN');
    }
    return;
  }

  try {
    await getAdminFromRequest(request);
    return;
  } catch {
    const member = await getMemberFromRequest(request);
    if (attachmentMemberId != null && attachmentMemberId !== member.id) {
      throw new Error('FORBIDDEN');
    }
  }
}

export async function GET(request: Request, { params }: RouteParams) {
  try {
    const { id } = await params;
    const attachment = await getAttachment(BigInt(id));
    if (!attachment) {
      return jsonFail(ErrorCodes.NOT_FOUND, '附件不存在', 404);
    }

    assertAntiHotlink(request);
    await assertAttachmentAccess(request, attachment.id, attachment.memberId);

    const buffer = await readUpload(attachment.storageKey);
    const inline = isImageMime(attachment.mimeType);

    return new NextResponse(new Uint8Array(buffer), {
      headers: {
        'Content-Type': attachment.mimeType,
        'Content-Disposition': inline
          ? 'inline'
          : `attachment; filename="${encodeURIComponent(attachment.fileName)}"`,
        'Content-Length': String(buffer.length),
        'Cache-Control': 'private, no-store',
        'X-Content-Type-Options': 'nosniff',
      },
    });
  } catch (error) {
    return handleApiError(error);
  }
}
