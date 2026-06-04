import { PrismaClient } from '@prisma/client';
const p = new PrismaClient();
const rows = await p.sysAttachment.findMany({ take: 10 });
console.log(rows.map((r) => ({ id: r.id.toString(), memberId: r.memberId?.toString(), fileName: r.fileName })));
await p.$disconnect();
