import { PrismaClient } from '@prisma/client';
import bcrypt from 'bcryptjs';

const prisma = new PrismaClient();

async function main() {
  const adminHash = await bcrypt.hash('admin123', 12);

  await prisma.adminUser.upsert({
    where: { username: 'admin' },
    update: {},
    create: {
      username: 'admin',
      passwordHash: adminHash,
      role: 'admin',
    },
  });

  const plans = [
    {
      name: '基础版',
      tier: 'basic',
      monthlyPrice: 299,
      features: ['OPC注册', '月度记账', '季度申报', '年度汇算'],
    },
    {
      name: '进阶版',
      tier: 'advanced',
      monthlyPrice: 799,
      features: ['含基础版', '税务筹划', '发票管理', '风险预警'],
    },
    {
      name: '尊享版',
      tier: 'premium',
      monthlyPrice: 1999,
      features: ['含进阶版', '架构设计', '稽查应对', '专属顾问'],
    },
  ];

  for (const plan of plans) {
    const existing = await prisma.servicePlan.findFirst({ where: { tier: plan.tier } });
    if (!existing) {
      await prisma.servicePlan.create({
        data: {
          name: plan.name,
          tier: plan.tier,
          monthlyPrice: plan.monthlyPrice,
          features: plan.features,
        },
      });
    }
  }

  const categories = [
    { code: 'equipment', name: '设备', voucherHint: '需提供购置发票' },
    { code: 'network', name: '网费', voucherHint: '需提供运营商发票' },
    { code: 'venue', name: '场地', voucherHint: '需提供租赁合同+发票' },
    { code: 'promotion', name: '投流', voucherHint: '需提供平台推广发票' },
    { code: 'outsource', name: '外包', voucherHint: '需提供合同+发票' },
    { code: 'travel', name: '差旅', voucherHint: '需提供车票/住宿票' },
    { code: 'other', name: '其他', voucherHint: '需提供说明+凭证' },
  ];

  for (const cat of categories) {
    await prisma.expenseCategory.upsert({
      where: { code: cat.code },
      update: {},
      create: cat,
    });
  }

  const accounts = [
    { code: '1001', name: '银行存款', type: 'asset' },
    { code: '6001', name: '主营业务收入', type: 'revenue' },
    { code: '6401', name: '管理费/推广费', type: 'expense' },
  ];

  for (const acc of accounts) {
    await prisma.ledgerAccount.upsert({
      where: { code: acc.code },
      update: {},
      create: acc,
    });
  }

  console.log('Seed completed: admin/admin123, 3 plans, expense categories, ledger accounts');
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
