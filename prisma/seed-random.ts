/**
 * 随机测试数据生成脚本
 *
 * 用法：
 *   npm run db:seed-random
 *   SEED_RANDOM_MEMBERS=10 npm run db:seed-random
 *
 * 环境变量：
 *   SEED_RANDOM_MEMBERS  - 生成会员数量（默认 8，含各 OPC 状态）
 *   SEED_RANDOM_PASSWORD - 会员统一密码（默认 test123）
 *   SEED_RANDOM_SKIP     - 设为 true 时跳过执行（部署脚本可条件调用）
 */
import { Prisma, PrismaClient } from '@prisma/client';
import bcrypt from 'bcryptjs';
import {
  EXPENSE_CATEGORIES,
  EXPENSE_DESCRIPTIONS,
  INCOME_CATEGORIES,
  INCOME_NOTES,
  INVOICE_TYPES,
  NICKNAMES,
  OPC_STATUS_PIPELINE,
  PLATFORMS,
  SURNAMES,
  backfillLedgerForOpc,
  pick,
  randomAmount,
  randomCreditCode,
  randomDateInPastMonths,
  randomInt,
  testPhone,
  type OpcPipelineStatus,
} from './seed-random-helpers';

const prisma = new PrismaClient();

const MEMBER_COUNT = Math.max(1, parseInt(process.env.SEED_RANDOM_MEMBERS ?? '8', 10));
const DEFAULT_PASSWORD = process.env.SEED_RANDOM_PASSWORD ?? 'test123';

const MONTHLY_INCOME_RANGES = ['0-2万', '2-5万', '5-10万', '10-20万', '20万+'] as const;
const EXISTING_ENTITIES = ['none', 'individual', 'labor'] as const;

interface CreatedMemberSummary {
  phone: string;
  name: string;
  opcStatus: OpcPipelineStatus;
  opcId: string;
}

function buildDisplayName(index: number): string {
  return `${pick(SURNAMES)}${pick(NICKNAMES)}${index}`;
}

function buildCompanyName(name: string): string {
  return `${name}文化传媒有限公司`;
}

async function ensureBaseSeedData() {
  const [admin, planCount, categoryCount] = await Promise.all([
    prisma.adminUser.findFirst({ where: { username: 'admin', deleted: false } }),
    prisma.servicePlan.count({ where: { active: true } }),
    prisma.expenseCategory.count({ where: { active: true } }),
  ]);

  if (!admin || planCount === 0 || categoryCount === 0) {
    throw new Error('请先运行 npm run db:seed 初始化管理员、套餐与费用科目');
  }

  return prisma.servicePlan.findFirst({ where: { tier: 'basic', active: true } });
}

function opcStatusForIndex(index: number): OpcPipelineStatus {
  if (index < OPC_STATUS_PIPELINE.length) {
    return OPC_STATUS_PIPELINE[index];
  }
  return 'active';
}

function materialsBaseFields(name: string, phone: string) {
  const now = new Date();
  return {
    proposedNames: [buildCompanyName(name), `${name}直播服务有限公司`, `${name}科技工作室`],
    registeredCapital: randomAmount(10, 100) * 10000,
    capitalTermYears: 20,
    businessTermType: 'long_term' as const,
    businessScope:
      '一般项目：个人互联网直播服务；文艺创作；摄像及视频制作服务；广告设计、代理。（除依法须经批准的项目外，凭营业执照依法自主开展经营活动）',
    registerProvince: '浙江省',
    registerCity: '杭州市',
    registerDistrict: pick(['西湖区', '滨江区', '余杭区']),
    registerAddress: `杭州市${pick(['文三路', '网商路', '五常大道'])}${randomInt(1, 999)}号`,
    legalPersonName: name,
    phone,
    email: `test${phone.slice(-4)}@example.com`,
    ethnicity: '汉族',
    householdAddress: `浙江省杭州市${pick(['西湖区', '滨江区'])}测试路${randomInt(1, 200)}号`,
    residentialAddress: `浙江省杭州市${pick(['西湖区', '滨江区'])}测试路${randomInt(1, 200)}号`,
    esignAuthorized: true,
    materialsSubmittedAt: now,
    idCardValidFrom: new Date(2018, 0, 1),
    idCardValidTo: new Date(2038, 0, 1),
  };
}

function activeOpcFields(name: string, phone: string) {
  const established = randomDateInPastMonths(8);
  return {
    ...materialsBaseFields(name, phone),
    companyName: buildCompanyName(name),
    creditCode: randomCreditCode(),
    status: 'active' as const,
    materialsApprovedAt: established,
    establishedAt: established,
    taxActivatedAt: new Date(established.getTime() + 7 * 86400000),
    taxpayerType: 'small_scale',
    bankName: pick(['中国工商银行', '招商银行', '中国建设银行']),
    bankAccountEnc: `enc_${phone.slice(-8)}`,
  };
}

async function createProgressLogs(opcId: bigint, status: OpcPipelineStatus) {
  const statusIndex = OPC_STATUS_PIPELINE.indexOf(status);
  const steps: { step: string; status: string; note: string }[] = [
    { step: 'materials', status: 'pending', note: '等待提交资料' },
  ];

  if (statusIndex >= 2) {
    steps.push({ step: 'materials', status: 'approved', note: '资料审核通过' });
  }
  if (status === 'materials') {
    steps.push({ step: 'materials', status: 'rejected', note: '地址证明不清晰，请重新上传' });
  }
  if (statusIndex >= 3) {
    steps.push({ step: 'business', status: 'in_progress', note: '工商注册办理中' });
  }
  if (statusIndex >= 4) {
    steps.push({ step: 'business', status: 'completed', note: '营业执照已核发' });
    steps.push({ step: 'tax', status: 'in_progress', note: '税务登记办理中' });
  }
  if (statusIndex >= 5) {
    steps.push({ step: 'tax', status: 'completed', note: '税务登记完成' });
    steps.push({ step: 'bank', status: 'in_progress', note: '银行开户办理中' });
  }
  if (status === 'active') {
    steps.push({ step: 'bank', status: 'completed', note: '银行开户完成' });
    steps.push({ step: 'complete', status: 'completed', note: 'OPC 设立完成' });
  }

  for (const log of steps) {
    await prisma.opcProgressLog.create({
      data: { opcId, step: log.step, status: log.status, note: log.note },
    });
  }
}

async function seedLedgerData(opcId: bigint) {
  const incomeCount = randomInt(25, 40);
  const expenseCount = randomInt(15, 30);

  const incomes = Array.from({ length: incomeCount }, (_, i) => {
    const gross = randomAmount(500, 25000);
    const platformFee = Math.round(gross * randomAmount(0.03, 0.08) * 100) / 100;
    return {
      opcId,
      platform: pick(PLATFORMS),
      category: pick(INCOME_CATEGORIES),
      grossAmount: gross,
      platformFee,
      netAmount: Math.round((gross - platformFee) * 100) / 100,
      occurredAt: randomDateInPastMonths(6),
      source: 'manual' as const,
      note: `${pick(INCOME_NOTES)}（随机 #${i + 1}）`,
    };
  });

  const expenses = Array.from({ length: expenseCount }, (_, i) => ({
    opcId,
    category: pick(EXPENSE_CATEGORIES),
    amount: randomAmount(100, 8000),
    invoiceType: pick(INVOICE_TYPES),
    description: `${pick(EXPENSE_DESCRIPTIONS)}（随机 #${i + 1}）`,
    occurredAt: randomDateInPastMonths(6),
    warningFlag: false,
  }));

  await prisma.incomeEntry.createMany({ data: incomes });
  await prisma.expenseEntry.createMany({ data: expenses });

  return backfillLedgerForOpc(prisma, opcId);
}

async function createRandomMember(
  slot: number,
  passwordHash: string,
  planId: bigint,
): Promise<CreatedMemberSummary> {
  const phone = testPhone(slot);
  const name = buildDisplayName(slot);
  const opcStatus = opcStatusForIndex(slot - 1);

  const existing = await prisma.member.findFirst({ where: { phone, deleted: false } });
  if (existing) {
    const opc = await prisma.opcEntity.findFirst({ where: { memberId: existing.id, deleted: false } });
    return {
      phone: existing.phone,
      name: existing.name ?? name,
      opcStatus: (opc?.status as OpcPipelineStatus) ?? 'pending',
      opcId: opc?.id.toString() ?? '-',
    };
  }

  const member = await prisma.member.create({
    data: {
      phone,
      passwordHash,
      name,
      email: `test${phone.slice(-4)}@example.com`,
    },
  });

  const diagnosis = await prisma.complianceDiagnosis.create({
    data: {
      memberId: member.id,
      platforms: [pick(PLATFORMS), pick(PLATFORMS)],
      monthlyIncomeRange: pick(MONTHLY_INCOME_RANGES),
      annualCostEstimate: randomAmount(50000, 300000),
      existingEntity: pick(EXISTING_ENTITIES),
      hasFiledTax: Math.random() > 0.5,
      taxBureauContact: Math.random() > 0.7,
      recommendedPlan: 'opc',
      taxComparison: {
        labor: { taxAmount: randomInt(20000, 80000), effectiveRate: 0.12 },
        individual: { taxAmount: randomInt(15000, 60000), effectiveRate: 0.09 },
        opc: { taxAmount: randomInt(10000, 40000), effectiveRate: 0.06 },
      },
    },
  });

  const order = await prisma.serviceOrder.create({
    data: {
      memberId: member.id,
      planId,
      diagnosisId: diagnosis.id,
      status: 'active',
      amount: 299,
      signedAt: randomDateInPastMonths(3),
    },
  });

  await prisma.complianceConsent.createMany({
    data: [
      {
        memberId: member.id,
        orderId: order.id,
        type: 'risk_disclosure',
        documentVersion: '1.0',
        agreedAt: new Date(),
      },
      {
        memberId: member.id,
        orderId: order.id,
        type: 'plan_confirm',
        documentVersion: '1.0',
        agreedAt: new Date(),
      },
      {
        memberId: member.id,
        orderId: order.id,
        type: 'contract_sign',
        documentVersion: '1.0',
        agreedAt: new Date(),
      },
    ],
  });

  let opcData: Prisma.OpcEntityUncheckedCreateInput;

  if (opcStatus === 'pending') {
    opcData = { memberId: member.id, status: 'pending' };
  } else if (opcStatus === 'materials') {
    opcData = {
      memberId: member.id,
      status: 'materials',
      ...materialsBaseFields(name, phone),
      rejectNote: '地址证明不清晰，请补充加盖公章的租赁合同首页与签字页',
    };
  } else if (opcStatus === 'materials_review') {
    opcData = {
      memberId: member.id,
      status: 'materials_review',
      ...materialsBaseFields(name, phone),
    };
  } else if (opcStatus === 'registering') {
    opcData = {
      memberId: member.id,
      status: 'registering',
      ...materialsBaseFields(name, phone),
      materialsApprovedAt: randomDateInPastMonths(2),
    };
  } else if (opcStatus === 'tax') {
    opcData = {
      memberId: member.id,
      status: 'tax',
      ...materialsBaseFields(name, phone),
      companyName: buildCompanyName(name),
      creditCode: randomCreditCode(),
      materialsApprovedAt: randomDateInPastMonths(2),
      establishedAt: randomDateInPastMonths(1),
    };
  } else if (opcStatus === 'bank') {
    opcData = {
      memberId: member.id,
      status: 'bank',
      ...materialsBaseFields(name, phone),
      companyName: buildCompanyName(name),
      creditCode: randomCreditCode(),
      materialsApprovedAt: randomDateInPastMonths(3),
      establishedAt: randomDateInPastMonths(2),
      taxActivatedAt: randomDateInPastMonths(1),
    };
  } else {
    opcData = {
      memberId: member.id,
      ...activeOpcFields(name, phone),
    };
  }

  const opc = await prisma.opcEntity.create({ data: opcData });
  await createProgressLogs(opc.id, opcStatus);

  if (opcStatus === 'active') {
    const periodCount = await seedLedgerData(opc.id);
    await prisma.memberNotice.createMany({
      data: [
        {
          memberId: member.id,
          type: 'statement',
          title: '月度对账单已生成',
          content: `您 ${new Date().getMonth() || 12} 月的 OPC 合规对账单已生成，请登录会员中心查看。`,
        },
        {
          memberId: member.id,
          type: 'tax_reminder',
          title: '申报提醒',
          content: '本月增值税申报截止日期临近，请确认收入与成本台账已更新。',
        },
      ],
    });
    console.log(`  └─ OPC #${opc.id} 台账：${periodCount} 个会计期间`);
  }

  return { phone, name, opcStatus, opcId: opc.id.toString() };
}

async function main() {
  if (process.env.SEED_RANDOM_SKIP === 'true') {
    console.log('SEED_RANDOM_SKIP=true，跳过随机测试数据生成');
    return;
  }

  const plan = await ensureBaseSeedData();
  if (!plan) throw new Error('未找到基础版套餐，请先运行 db:seed');

  const passwordHash = await bcrypt.hash(DEFAULT_PASSWORD, 12);
  const summaries: CreatedMemberSummary[] = [];

  console.log(`开始生成 ${MEMBER_COUNT} 个随机测试会员（密码：${DEFAULT_PASSWORD}）…`);

  for (let i = 0; i < MEMBER_COUNT; i++) {
    const summary = await createRandomMember(i + 1, passwordHash, plan.id); // slot 1-based，状态按 i 分配
    summaries.push(summary);
    console.log(
      `  [${i + 1}/${MEMBER_COUNT}] ${summary.name} | ${summary.phone} | OPC ${summary.opcStatus} (#${summary.opcId})`,
    );
  }

  console.log('\n随机测试数据生成完成：');
  console.log('─'.repeat(60));
  console.log('手机号'.padEnd(14), '姓名'.padEnd(10), 'OPC 状态'.padEnd(16), 'OPC ID');
  console.log('─'.repeat(60));
  for (const s of summaries) {
    console.log(s.phone.padEnd(14), s.name.padEnd(10), s.opcStatus.padEnd(16), s.opcId);
  }
  console.log('─'.repeat(60));
  console.log(`会员登录密码：${DEFAULT_PASSWORD}`);
  console.log(`管理员：admin / admin123（来自 db:seed）`);
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
