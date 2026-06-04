import { PrismaClient } from '@prisma/client';

const prisma = new PrismaClient();

const CATEGORIES = ['equipment', 'network', 'venue', 'promotion', 'outsource', 'travel', 'other'] as const;
const INVOICE_TYPES = ['general', 'special', 'electronic'] as const;

const DESCRIPTIONS = [
  '直播设备采购',
  '宽带月费',
  '摄影棚租赁',
  'Dou+ 推广投放',
  '剪辑外包',
  '出差交通住宿',
  '办公耗材',
  '灯光支架补充',
  '5G 流量套餐',
  '季度场地续约',
];

function pick<T>(arr: readonly T[]): T {
  return arr[Math.floor(Math.random() * arr.length)];
}

function randomAmount(): number {
  return Math.round((Math.random() * 4800 + 100) * 100) / 100;
}

function randomOccurredAt(): Date {
  const now = new Date();
  const daysBack = Math.floor(Math.random() * 180);
  const d = new Date(now);
  d.setDate(d.getDate() - daysBack);
  d.setHours(12, 0, 0, 0);
  return d;
}

async function main() {
  const opc =
    (await prisma.opcEntity.findFirst({
      where: { deleted: false, status: 'active' },
      orderBy: { id: 'desc' },
    })) ??
    (await prisma.opcEntity.findFirst({
      where: { deleted: false },
      orderBy: { id: 'desc' },
    }));

  if (!opc) {
    throw new Error('未找到 OPC 主体，请先完成签约并激活 OPC 后再运行本脚本');
  }

  const rows = Array.from({ length: 30 }, (_, i) => ({
    opcId: opc.id,
    category: pick(CATEGORIES),
    amount: randomAmount(),
    invoiceType: pick(INVOICE_TYPES),
    description: `${pick(DESCRIPTIONS)}（测试 #${i + 1}）`,
    occurredAt: randomOccurredAt(),
    warningFlag: false,
  }));

  const result = await prisma.expenseEntry.createMany({ data: rows });

  console.log(`已为 OPC #${opc.id} 插入 ${result.count} 条费用测试数据`);
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
