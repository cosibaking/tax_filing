import PDFDocument from 'pdfkit';
import { prisma, serializeBigInt } from '@/lib/db';
import { saveUpload } from '@/lib/storage/local';

export async function generateMonthlyStatement(opcId: bigint, year: number, month: number) {
  const summary = await prisma.profitSummary.findUnique({
    where: { opcId_year_month: { opcId, year, month } },
  });
  if (!summary) throw new Error('该月暂无利润数据');

  const tasks = await prisma.taxFilingTask.findMany({
    where: {
      opcId,
      period: `${year}-${String(month).padStart(2, '0')}`,
      deleted: false,
    },
  });
  const taxTotal = tasks.reduce((s, t) => s + Number(t.calculatedAmount), 0);

  const summaryData = {
    revenue: Number(summary.revenue),
    cost: Number(summary.cost),
    profit: Number(summary.profit),
    cumulativeProfit: Number(summary.cumulativeProfit),
    taxTotal,
    filingStatus: tasks.every((t) => t.status === 'filed') ? '已申报完成' : '待申报',
  };

  const pdfBuffer = await buildStatementPdf(year, month, summaryData);
  const file = new File([new Uint8Array(pdfBuffer)], `statement-${year}-${month}.pdf`, {
    type: 'application/pdf',
  });
  const opc = await prisma.opcEntity.findUnique({ where: { id: opcId } });
  const { id: pdfFileId } = await saveUpload(file, opc?.memberId);

  const statement = await prisma.monthlyStatement.upsert({
    where: { opcId_year_month: { opcId, year, month } },
    create: { opcId, year, month, summary: summaryData, pdfFileId: BigInt(pdfFileId) },
    update: { summary: summaryData, pdfFileId: BigInt(pdfFileId) },
  });

  return serializeBigInt(statement);
}

async function buildStatementPdf(
  year: number,
  month: number,
  summary: Record<string, unknown>,
): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    const doc = new PDFDocument();
    const chunks: Buffer[] = [];
    doc.on('data', (c) => chunks.push(c));
    doc.on('end', () => resolve(Buffer.concat(chunks)));
    doc.on('error', reject);
    doc.fontSize(16).text(`月度对账单 ${year}年${month}月`, { align: 'center' });
    doc.moveDown();
    Object.entries(summary).forEach(([k, v]) => doc.text(`${k}: ${v}`));
    doc.end();
  });
}

export async function sendStatementNotification(opcId: bigint, year: number, month: number) {
  const opc = await prisma.opcEntity.findUnique({ where: { id: opcId } });
  if (!opc) return;
  const statement = await prisma.monthlyStatement.findUnique({
    where: { opcId_year_month: { opcId, year, month } },
  });
  if (!statement) return;

  await prisma.memberNotice.create({
    data: {
      memberId: opc.memberId,
      type: 'statement',
      title: `${year}年${month}月税务服务完成报告`,
      content: `本月对账单已生成，请前往会员中心查看。申报状态：${(statement.summary as { filingStatus?: string }).filingStatus ?? '未知'}`,
    },
  });

  await prisma.monthlyStatement.update({
    where: { id: statement.id },
    data: { sentAt: new Date() },
  });
}

export async function sendMaterialReminder(memberId: bigint) {
  const now = new Date();
  const prevMonth = now.getMonth() === 0 ? 12 : now.getMonth();
  const year = now.getMonth() === 0 ? now.getFullYear() - 1 : now.getFullYear();
  await prisma.memberNotice.create({
    data: {
      memberId,
      type: 'material_reminder',
      title: '请上传上月材料',
      content: `请在每月5日前上传${year}年${prevMonth}月银行流水、平台截图及成本发票。`,
    },
  });
}

export async function listStatements(memberId: bigint) {
  const opc = await prisma.opcEntity.findFirst({ where: { memberId, deleted: false } });
  if (!opc) return [];
  const list = await prisma.monthlyStatement.findMany({
    where: { opcId: opc.id, deleted: false },
    orderBy: [{ year: 'desc' }, { month: 'desc' }],
  });
  return serializeBigInt(list);
}
