import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

const rows = await prisma.opcEntity.findMany({
  where: { deleted: false },
  orderBy: { id: 'asc' },
});

for (const r of rows) {
  console.log('--- opc id', r.id.toString());
  console.log({
    status: r.status,
    proposedNames: r.proposedNames,
    legalPersonName: r.legalPersonName,
    businessScope: r.businessScope?.slice(0, 40),
    registerProvince: r.registerProvince,
    phone: r.phone,
    materialsSubmittedAt: r.materialsSubmittedAt,
    hasIdCard: !!r.idCardEncrypted,
  });
}

await prisma.$disconnect();
