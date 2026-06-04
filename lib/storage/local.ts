import { mkdir, writeFile, readFile } from 'fs/promises';
import path from 'path';
import { prisma } from '@/lib/db';

const UPLOAD_DIR = process.env.UPLOAD_DIR ?? './uploads';

export async function saveUpload(
  file: File,
  memberId?: bigint,
): Promise<{ id: string; storageKey: string }> {
  await mkdir(UPLOAD_DIR, { recursive: true });
  const buffer = Buffer.from(await file.arrayBuffer());
  const storageKey = `${Date.now()}-${file.name.replace(/[^a-zA-Z0-9._-]/g, '_')}`;
  const filePath = path.join(UPLOAD_DIR, storageKey);
  await writeFile(filePath, buffer);
  const attachment = await prisma.sysAttachment.create({
    data: {
      memberId,
      fileName: file.name,
      mimeType: file.type || 'application/octet-stream',
      size: buffer.length,
      storageKey,
    },
  });
  return { id: attachment.id.toString(), storageKey };
}

export async function readUpload(storageKey: string): Promise<Buffer> {
  const filePath = path.join(UPLOAD_DIR, storageKey);
  return readFile(filePath);
}

export async function getAttachment(id: bigint) {
  return prisma.sysAttachment.findFirst({ where: { id, deleted: false } });
}
