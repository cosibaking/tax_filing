import 'server-only';
import { prisma, serializeBigInt } from '@/lib/db';
import { signMediaUrl } from '@/lib/media/signed-access';
import { writeAuditLog } from '@/lib/services/compliance/audit/audit-log';
import { getClientIp } from '@/lib/auth/member';

const NAME_RE = /^[\u4e00-\u9fa5a-zA-Z·\s]{2,64}$/;

type MemberRow = {
  id: bigint;
  phone: string;
  name: string | null;
  email: string | null;
  avatarFileId: bigint | null;
  createdAt: Date;
};

export function validateMemberName(name: string): string | null {
  const trimmed = name.trim();
  if (!NAME_RE.test(trimmed)) return '姓名需为 2–64 个汉字或字母';
  return null;
}

async function assertAvatarOwned(fileId: bigint, memberId: bigint) {
  const file = await prisma.sysAttachment.findFirst({
    where: { id: fileId, deleted: false },
  });
  if (!file) throw new Error('头像文件不存在');
  if (file.memberId != null && file.memberId !== memberId) {
    throw new Error('无权使用该头像文件');
  }
  if (!file.mimeType.startsWith('image/')) {
    throw new Error('头像必须为图片格式');
  }
}

export function serializeMemberProfile(member: MemberRow) {
  const data = serializeBigInt({
    id: member.id,
    phone: member.phone,
    name: member.name,
    email: member.email,
    avatarFileId: member.avatarFileId,
    createdAt: member.createdAt,
  }) as Record<string, unknown>;

  if (member.avatarFileId) {
    data.avatarUrl = signMediaUrl(member.avatarFileId.toString(), member.id.toString());
  }

  return data;
}

export async function updateMemberProfile(
  memberId: bigint,
  input: { name?: string; avatarFileId?: string | null },
  request: Request,
) {
  const before = await prisma.member.findFirst({ where: { id: memberId, deleted: false } });
  if (!before) throw new Error('NOT_FOUND');

  const data: { name?: string; avatarFileId?: bigint | null } = {};

  if (input.name !== undefined) {
    const err = validateMemberName(input.name);
    if (err) throw new Error(err);
    data.name = input.name.trim();
  }

  if (input.avatarFileId !== undefined) {
    if (input.avatarFileId === null || input.avatarFileId === '') {
      data.avatarFileId = null;
    } else {
      const fileId = BigInt(input.avatarFileId);
      await assertAvatarOwned(fileId, memberId);
      data.avatarFileId = fileId;
    }
  }

  if (Object.keys(data).length === 0) {
    throw new Error('没有可更新的字段');
  }

  const updated = await prisma.member.update({
    where: { id: memberId },
    data,
  });

  await writeAuditLog({
    entityType: 'member',
    entityId: memberId,
    action: 'profile.update',
    operatorId: memberId,
    operatorType: 'member',
    before: serializeBigInt({ name: before.name, avatarFileId: before.avatarFileId }),
    after: serializeBigInt({ name: updated.name, avatarFileId: updated.avatarFileId }),
    ip: getClientIp(request),
  });

  return serializeMemberProfile(updated);
}
