import { prisma, serializeBigInt } from '@/lib/db';
import { writeAuditLog } from '@/lib/services/compliance/audit/audit-log';
import { getClientIp, getUserAgent } from '@/lib/auth/member';
import { isDev } from '@/lib/api/constants';
import { ensureOpcEntityOnSign } from '@/lib/services/compliance/opc/opc-service';

export {
  submitOpcMaterials,
  getOpcProgress,
  submitBankReceipt,
} from '@/lib/services/compliance/opc/opc-service';
export { applyOpcAdminAction, listOpcTasks, getOpcTaskDetail } from '@/lib/services/compliance/opc/opc-service';

async function loadPdfDocument() {
  const { default: PDFDocument } = await import('pdfkit');
  return PDFDocument;
}

export async function getActivePlans() {
  const plans = await prisma.servicePlan.findMany({ where: { active: true } });
  return serializeBigInt(plans);
}

export async function createOrder(memberId: bigint, planId: bigint, diagnosisId?: bigint) {
  const plan = await prisma.servicePlan.findFirst({ where: { id: planId, active: true } });
  if (!plan) throw new Error('套餐不存在');
  const existing = await prisma.serviceOrder.findFirst({
    where: { memberId, status: 'active', deleted: false },
  });
  if (existing) throw new Error('已有生效订单');
  const order = await prisma.serviceOrder.create({
    data: {
      memberId,
      planId,
      diagnosisId,
      status: 'pending',
      amount: plan.monthlyPrice,
    },
    include: { plan: true },
  });
  return serializeBigInt(order);
}

export async function recordConsent(
  memberId: bigint,
  type: string,
  documentVersion: string,
  request: Request,
  orderId?: bigint,
) {
  await prisma.complianceConsent.create({
    data: {
      memberId,
      orderId,
      type,
      documentVersion,
      ip: getClientIp(request),
      userAgent: getUserAgent(request),
      agreedAt: new Date(),
    },
  });
}

export async function signContract(
  memberId: bigint,
  orderId: bigint,
  signerName: string,
  request: Request,
) {
  const order = await prisma.serviceOrder.findFirst({
    where: { id: orderId, memberId, deleted: false },
  });
  if (!order) throw new Error('订单不存在');

  let contractFileId: bigint | undefined;

  if (isDev) {
    // 开发模式：mock 电子签（跳过 PDF 生成与文件存储，待接入正式签章服务）
  } else {
    const pdfBuffer = await generateContractPdf(signerName, orderId.toString());
    const { saveUpload } = await import('@/lib/storage/local');
    const file = new File([new Uint8Array(pdfBuffer)], `contract-${orderId}.pdf`, {
      type: 'application/pdf',
    });
    const { id: fileId } = await saveUpload(file, memberId);
    contractFileId = BigInt(fileId);
  }

  await prisma.serviceOrder.update({
    where: { id: orderId },
    data: {
      status: 'active',
      signedAt: new Date(),
      ...(contractFileId != null ? { contractFileId } : {}),
    },
  });

  const opc = await ensureOpcEntityOnSign(memberId);

  await writeAuditLog({
    entityType: 'service_order',
    entityId: orderId,
    action: 'contract.sign',
    operatorId: memberId,
    operatorType: 'member',
    after: { signerName, ...(isDev ? { mockSign: true } : {}) },
    ip: getClientIp(request),
  });

  return serializeBigInt({
    orderId: orderId.toString(),
    opcId: opc.id.toString(),
    ...(isDev ? { mockSign: true } : {}),
  });
}

async function generateContractPdf(signerName: string, orderId: string): Promise<Buffer> {
  const PDFDocument = await loadPdfDocument();
  return new Promise((resolve, reject) => {
    const doc = new PDFDocument();
    const chunks: Buffer[] = [];
    doc.on('data', (c) => chunks.push(c));
    doc.on('end', () => resolve(Buffer.concat(chunks)));
    doc.on('error', reject);
    doc.fontSize(18).text('主播 OPC 合规服务协议', { align: 'center' });
    doc.moveDown();
    doc.fontSize(12).text(`订单号：${orderId}`);
    doc.text('本协议为合规服务协议，非逃税方案。');
    doc.text('签署人：' + signerName);
    doc.text('签署时间：' + new Date().toISOString());
    doc.end();
  });
}

