import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getMemberFromRequest } from '@/lib/auth/member';
import { prisma, serializeBigInt } from '@/lib/db';
import { decrypt, maskIdCard, maskBankAccount } from '@/lib/crypto/encrypt';

function maskOpcEntity(opc: {
  idCardEncrypted: string | null;
  bankAccountEnc: string | null;
  [key: string]: unknown;
}) {
  const data = serializeBigInt(opc) as Record<string, unknown>;
  delete data.idCardEncrypted;
  delete data.bankAccountEnc;

  if (opc.idCardEncrypted) {
    try {
      data.idCardMasked = maskIdCard(decrypt(opc.idCardEncrypted));
    } catch {
      data.idCardMasked = '****';
    }
  }
  if (opc.bankAccountEnc) {
    try {
      data.bankAccountMasked = maskBankAccount(decrypt(opc.bankAccountEnc));
    } catch {
      data.bankAccountMasked = '****';
    }
  }
  return data;
}

export async function GET(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const order = await prisma.serviceOrder.findFirst({
      where: {
        memberId: member.id,
        deleted: false,
        status: { in: ['pending', 'active'] },
      },
    });
    const opc = await prisma.opcEntity.findFirst({
      where: { memberId: member.id, deleted: false },
    });
    if (!opc) {
      return jsonOk({
        hasOrder: !!order,
        opcStatus: null,
      });
    }
    return jsonOk({
      ...maskOpcEntity(opc),
      hasOrder: !!order,
      opcStatus: opc.status,
    });
  } catch (error) {
    return handleApiError(error);
  }
}
