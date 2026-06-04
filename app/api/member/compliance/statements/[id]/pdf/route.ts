import { NextResponse } from 'next/server';
import { jsonFail, handleApiError } from '@/lib/api/envelope';
import { ErrorCodes } from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { getStatementPdfFileId } from '@/lib/services/compliance/statement/statement-service';
import { getAttachment, readUpload } from '@/lib/storage/local';

type RouteParams = { params: Promise<{ id: string }> };

export async function GET(request: Request, { params }: RouteParams) {
  try {
    const member = await getMemberFromRequest(request);
    const { id } = await params;
    const fileId = await getStatementPdfFileId(member.id, BigInt(id));
    const attachment = await getAttachment(fileId);
    if (!attachment) {
      return jsonFail(ErrorCodes.NOT_FOUND, 'PDF 不存在', 404);
    }
    if (attachment.memberId && attachment.memberId !== member.id) {
      throw new Error('FORBIDDEN');
    }

    const buffer = await readUpload(attachment.storageKey);
    return new NextResponse(new Uint8Array(buffer), {
      headers: {
        'Content-Type': attachment.mimeType || 'application/pdf',
        'Content-Disposition': `inline; filename="${encodeURIComponent(attachment.fileName)}"`,
        'Content-Length': String(buffer.length),
      },
    });
  } catch (error) {
    return handleApiError(error);
  }
}
