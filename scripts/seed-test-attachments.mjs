import { PrismaClient } from '@prisma/client';
const p = new PrismaClient();
const memberId = BigInt(2);
const ids = [];
for (const name of ['addr.pdf', 'id-front.jpg', 'id-back.jpg']) {
  const a = await p.sysAttachment.create({
    data: {
      memberId,
      fileName: name,
      mimeType: 'application/pdf',
      size: 100,
      storageKey: `test-${name}`,
    },
  });
  ids.push(a.id.toString());
}
console.log('attachment ids', ids);
await p.$disconnect();
