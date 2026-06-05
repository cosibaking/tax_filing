import 'server-only';
import { prisma, serializeBigInt } from '@/lib/db';
import { encrypt, decrypt, maskIdCard, maskBankAccount } from '@/lib/crypto/encrypt';
import { withBasePath } from '@/lib/base-path';
import { signMediaUrl, isImageMime } from '@/lib/media/signed-access';
import { writeAuditLog } from '@/lib/services/compliance/audit/audit-log';
import { getClientIp } from '@/lib/auth/member';
import { BANK_OPENING_CHECKLIST, OPC_STATUS_LABELS } from './opc-constants';
import {
  parseMaterialsBody,
  parseDateOnly,
  validateOpcMaterials,
} from './opc-materials-validation';
import type { OpcAdminAction } from './opc-types';

export type { OpcAdminAction, OpcMaterialsInput } from './opc-types';
export { parseMaterialsBody, validateOpcMaterials } from './opc-materials-validation';

const CREDIT_CODE_RE = /^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$/;

async function assertAttachmentOwned(fileId: bigint, memberId: bigint) {
  const file = await prisma.sysAttachment.findFirst({
    where: { id: fileId, deleted: false },
  });
  if (!file) throw new Error('附件不存在');
  if (file.memberId != null && file.memberId !== memberId) {
    throw new Error('附件无权使用');
  }
}

type OpcRow = {
  id: bigint;
  memberId: bigint;
  status: string;
  companyName: string | null;
  creditCode: string | null;
  proposedNames: unknown;
  registeredCapital: { toString(): string } | null;
  capitalTermYears: number | null;
  businessTermType: string | null;
  businessTermEnd: Date | null;
  businessScope: string | null;
  registerProvince: string | null;
  registerCity: string | null;
  registerDistrict: string | null;
  registerAddress: string | null;
  addressProofFileId: bigint | null;
  legalPersonName: string | null;
  idCardEncrypted: string | null;
  idCardValidFrom: Date | null;
  idCardValidTo: Date | null;
  idCardFrontFileId: bigint | null;
  idCardBackFileId: bigint | null;
  ethnicity: string | null;
  householdAddress: string | null;
  residentialAddress: string | null;
  phone: string | null;
  email: string | null;
  esignAuthorized: boolean;
  establishedAt: Date | null;
  taxpayerType: string;
  taxActivatedAt: Date | null;
  bankName: string | null;
  bankAccountEnc: string | null;
  bankReceiptFileId: bigint | null;
  licenseFileId: bigint | null;
  materialsSubmittedAt: Date | null;
  materialsApprovedAt: Date | null;
  rejectNote: string | null;
  createdAt: Date;
  updatedAt: Date;
};

export function maskOpcForResponse(opc: OpcRow, role: 'member' | 'admin') {
  const data = serializeBigInt(opc) as Record<string, unknown>;
  delete data.idCardEncrypted;
  delete data.bankAccountEnc;

  if (opc.idCardEncrypted) {
    try {
      data.idCardMasked = maskIdCard(decrypt(opc.idCardEncrypted));
    } catch {
      data.idCardMasked = '****';
    }
  }
  if (opc.bankAccountEnc) {
    try {
      data.bankAccountMasked = maskBankAccount(decrypt(opc.bankAccountEnc));
    } catch {
      data.bankAccountMasked = '****';
    }
  }

  data.opcStatus = opc.status;
  data.statusLabel = OPC_STATUS_LABELS[opc.status] ?? opc.status;

  if (role === 'member') {
    delete data.rejectNote;
    if (opc.status !== 'materials' && opc.status !== 'pending') {
      /* 审核通过后不回显完整身份证 */
    }
  }

  const names = opc.proposedNames as string[] | null;
  data.proposedNamePrimary = names?.[0] ?? opc.companyName ?? null;
  if (opc.registeredCapital != null) {
    data.registeredCapital = opc.registeredCapital.toString();
  }

  return data;
}

/** 会员资料页 OPC 主体摘要（F-86 / US-E1-03） */
export function buildOpcProfileSummary(opc: OpcRow | null, hasOrder: boolean) {
  if (!opc) {
    return { hasOrder, opcStatus: null as string | null };
  }
  const masked = maskOpcForResponse(opc, 'member');
  return {
    hasOrder,
    companyName: (masked.companyName as string | null) ?? null,
    proposedNamePrimary: (masked.proposedNamePrimary as string | null) ?? null,
    creditCode: (masked.creditCode as string | null) ?? null,
    opcStatus: masked.opcStatus as string,
    statusLabel: masked.statusLabel as string,
    bankAccountMasked: masked.bankAccountMasked as string | undefined,
  };
}

export function buildOpcTimeline(status: string, logs: { step: string; status: string; createdAt: Date }[]) {
  const findLogDate = (step: string) =>
    logs.find((l) => l.step === step && ['approved', 'license_issued', 'completed', 'reviewing'].includes(l.status))
      ?.createdAt;

  const materialsDone = ['materials_review', 'registering', 'tax', 'bank', 'active'].includes(status);
  const businessDone = ['tax', 'bank', 'active'].includes(status);
  const taxDone = ['bank', 'active'].includes(status);
  const bankDone = status === 'active';

  const stepState = (done: boolean, current: boolean): 'done' | 'current' | 'pending' => {
    if (done) return 'done';
    if (current) return 'current';
    return 'pending';
  };

  const materialsCurrent = ['pending', 'materials', 'materials_review'].includes(status);
  const businessCurrent = status === 'registering';
  const taxCurrent = status === 'tax';
  const bankCurrent = status === 'bank';

  return [
    {
      key: 'materials',
      label: '资料提交',
      status: stepState(materialsDone, materialsCurrent),
      date: findLogDate('materials')?.toISOString().slice(0, 10),
    },
    {
      key: 'business',
      label: '工商注册',
      status: stepState(businessDone, businessCurrent),
      date: findLogDate('business')?.toISOString().slice(0, 10),
    },
    {
      key: 'tax',
      label: '税务登记',
      status: stepState(taxDone, taxCurrent),
      date: findLogDate('tax')?.toISOString().slice(0, 10),
    },
    {
      key: 'bank',
      label: '银行开户',
      status: stepState(bankDone, bankCurrent),
      date: findLogDate('bank')?.toISOString().slice(0, 10),
    },
  ];
}

export async function submitOpcMaterials(
  memberId: bigint,
  raw: unknown,
  request: Request,
) {
  const data = parseMaterialsBody(raw);
  const err = validateOpcMaterials(data);
  if (err) throw new Error(err);

  const opc = await prisma.opcEntity.findFirst({ where: { memberId, deleted: false } });
  if (!opc) throw new Error('请先完成签约');
  if (!['pending', 'materials'].includes(opc.status)) {
    throw new Error('当前状态不可修改资料');
  }

  const fileIds = [
    BigInt(data.addressProofFileId),
    BigInt(data.idCardFrontFileId),
    BigInt(data.idCardBackFileId),
  ];
  for (const fid of fileIds) {
    await assertAttachmentOwned(fid, memberId);
  }

  const now = new Date();
  const isResubmit = opc.materialsSubmittedAt != null;

  await prisma.opcEntity.update({
    where: { id: opc.id },
    data: {
      proposedNames: data.proposedNames,
      registeredCapital: data.registeredCapital,
      capitalTermYears: data.capitalTermYears,
      businessTermType: data.businessTermType,
      businessTermEnd:
        data.businessTermType === 'fixed' && data.businessTermEnd
          ? parseDateOnly(data.businessTermEnd)
          : null,
      businessScope: data.businessScope,
      registerProvince: data.registerProvince,
      registerCity: data.registerCity,
      registerDistrict: data.registerDistrict,
      registerAddress: data.registerAddress,
      addressProofFileId: BigInt(data.addressProofFileId),
      legalPersonName: data.legalPersonName,
      idCardEncrypted: encrypt(data.idCard),
      idCardValidFrom: parseDateOnly(data.idCardValidFrom),
      idCardValidTo: parseDateOnly(data.idCardValidTo),
      idCardFrontFileId: BigInt(data.idCardFrontFileId),
      idCardBackFileId: BigInt(data.idCardBackFileId),
      ethnicity: data.ethnicity ?? null,
      householdAddress: data.householdAddress,
      residentialAddress: data.residentialAddress,
      phone: data.phone,
      email: data.email,
      esignAuthorized: data.esignAuthorized,
      status: 'materials_review',
      rejectNote: null,
      materialsSubmittedAt: isResubmit ? opc.materialsSubmittedAt : now,
    },
  });

  await prisma.opcProgressLog.create({
    data: {
      opcId: opc.id,
      step: 'materials',
      status: 'reviewing',
      note: isResubmit ? '资料已重新提交' : '资料已提交，审核中',
    },
  });

  await writeAuditLog({
    entityType: 'opc_entity',
    entityId: opc.id,
    action: 'materials.submit',
    operatorId: memberId,
    operatorType: 'member',
    ip: getClientIp(request),
  });

  return getOpcProgress(memberId);
}

export async function submitBankReceipt(
  memberId: bigint,
  data: { bankName?: string; bankReceiptFileId: string },
  request: Request,
) {
  const opc = await prisma.opcEntity.findFirst({ where: { memberId, deleted: false } });
  if (!opc) throw new Error('OPC 不存在');
  if (!['tax', 'bank'].includes(opc.status)) {
    throw new Error('当前阶段不可上传开户回执');
  }

  await assertAttachmentOwned(BigInt(data.bankReceiptFileId), memberId);

  await prisma.opcEntity.update({
    where: { id: opc.id },
    data: {
      bankName: data.bankName ?? opc.bankName,
      bankReceiptFileId: BigInt(data.bankReceiptFileId),
      status: opc.status === 'tax' ? 'bank' : opc.status,
    },
  });

  await writeAuditLog({
    entityType: 'opc_entity',
    entityId: opc.id,
    action: 'bank.receipt_upload',
    operatorId: memberId,
    operatorType: 'member',
    ip: getClientIp(request),
  });

  return getOpcProgress(memberId);
}

export async function getOpcProgress(memberId: bigint) {
  const opc = await prisma.opcEntity.findFirst({
    where: { memberId, deleted: false },
    include: { progressLogs: { orderBy: { createdAt: 'asc' } } },
  });
  if (!opc) return null;

  const masked = maskOpcForResponse(opc, 'member');
  return {
    ...masked,
    steps: buildOpcTimeline(opc.status, opc.progressLogs),
    rejectNote: opc.status === 'materials' ? opc.rejectNote : null,
    bankOpeningChecklist: BANK_OPENING_CHECKLIST,
    canEditMaterials: ['pending', 'materials'].includes(opc.status),
    estimatedSlaDays: 14,
  };
}

async function loadAttachmentBrief(fileId: bigint | null, memberId?: bigint) {
  if (!fileId) return null;
  const file = await prisma.sysAttachment.findFirst({
    where: { id: fileId, deleted: false },
  });
  if (!file) return null;
  const base = {
    id: file.id.toString(),
    fileName: file.fileName,
    mimeType: file.mimeType,
  };
  if (memberId && isImageMime(file.mimeType)) {
    return {
      ...base,
      mediaUrl: signMediaUrl(file.id.toString(), memberId.toString()),
    };
  }
  return {
    ...base,
    downloadUrl: withBasePath(`/api/upload/${file.id}`),
  };
}

/** 会员查看自身 OPC 完整信息（不脱敏，F-86） */
export async function getOpcMemberFullDetail(memberId: bigint, request: Request) {
  const opc = await prisma.opcEntity.findFirst({
    where: { memberId, deleted: false },
  });
  if (!opc) throw new Error('NOT_FOUND');

  const fmtDate = (d: Date | null) => (d ? d.toISOString().slice(0, 10) : null);

  let idCard: string | null = null;
  if (opc.idCardEncrypted) {
    try {
      idCard = decrypt(opc.idCardEncrypted);
    } catch {
      idCard = null;
    }
  }

  let bankAccount: string | null = null;
  if (opc.bankAccountEnc) {
    try {
      bankAccount = decrypt(opc.bankAccountEnc);
    } catch {
      bankAccount = null;
    }
  }

  const names = opc.proposedNames as string[] | null;

  const [addressProof, idFront, idBack, license, bankReceipt] = await Promise.all([
    loadAttachmentBrief(opc.addressProofFileId, memberId),
    loadAttachmentBrief(opc.idCardFrontFileId, memberId),
    loadAttachmentBrief(opc.idCardBackFileId, memberId),
    loadAttachmentBrief(opc.licenseFileId, memberId),
    loadAttachmentBrief(opc.bankReceiptFileId, memberId),
  ]);

  await writeAuditLog({
    entityType: 'opc_entity',
    entityId: opc.id,
    action: 'opc.view_full_detail',
    operatorId: memberId,
    operatorType: 'member',
    ip: getClientIp(request),
  });

  return serializeBigInt({
    companyName: opc.companyName,
    proposedNames: names,
    creditCode: opc.creditCode,
    opcStatus: opc.status,
    statusLabel: OPC_STATUS_LABELS[opc.status] ?? opc.status,
    registeredCapital: opc.registeredCapital?.toString() ?? null,
    capitalTermYears: opc.capitalTermYears,
    businessTermType: opc.businessTermType,
    businessTermEnd: fmtDate(opc.businessTermEnd),
    businessScope: opc.businessScope,
    registerProvince: opc.registerProvince,
    registerCity: opc.registerCity,
    registerDistrict: opc.registerDistrict,
    registerAddress: opc.registerAddress,
    legalPersonName: opc.legalPersonName,
    idCard,
    idCardValidFrom: fmtDate(opc.idCardValidFrom),
    idCardValidTo: fmtDate(opc.idCardValidTo),
    ethnicity: opc.ethnicity,
    householdAddress: opc.householdAddress,
    residentialAddress: opc.residentialAddress,
    phone: opc.phone,
    email: opc.email,
    esignAuthorized: opc.esignAuthorized,
    establishedAt: fmtDate(opc.establishedAt),
    taxpayerType: opc.taxpayerType,
    taxActivatedAt: opc.taxActivatedAt?.toISOString().slice(0, 10) ?? null,
    bankName: opc.bankName,
    bankAccount,
    materialsSubmittedAt: opc.materialsSubmittedAt?.toISOString() ?? null,
    materialsApprovedAt: opc.materialsApprovedAt?.toISOString() ?? null,
    rejectNote: opc.status === 'materials' ? opc.rejectNote : null,
    attachments: {
      addressProof,
      idCardFront: idFront,
      idCardBack: idBack,
      license,
      bankReceipt,
    },
  });
}

export async function getOpcTaskDetail(opcId: bigint) {
  const opc = await prisma.opcEntity.findFirst({
    where: { id: opcId, deleted: false },
    include: {
      member: { select: { id: true, phone: true, name: true, email: true } },
      progressLogs: { orderBy: { createdAt: 'desc' }, take: 20 },
    },
  });
  if (!opc) throw new Error('任务不存在');

  const order = await prisma.serviceOrder.findFirst({
    where: { memberId: opc.memberId, deleted: false, status: 'active' },
    include: { plan: { select: { name: true, tier: true } } },
  });

  const [addressProof, idFront, idBack, license, bankReceipt] = await Promise.all([
    loadAttachmentBrief(opc.addressProofFileId),
    loadAttachmentBrief(opc.idCardFrontFileId),
    loadAttachmentBrief(opc.idCardBackFileId),
    loadAttachmentBrief(opc.licenseFileId),
    loadAttachmentBrief(opc.bankReceiptFileId),
  ]);

  const fmtDate = (d: Date | null) => (d ? d.toISOString().slice(0, 10) : null);

  return {
    ...maskOpcForResponse(opc, 'admin'),
    member: serializeBigInt(opc.member),
    plan: order?.plan ? serializeBigInt(order.plan) : null,
    progressLogs: serializeBigInt(opc.progressLogs),
    steps: buildOpcTimeline(opc.status, [...opc.progressLogs].reverse()),
    allowedActions: getAllowedAdminActions(opc.status),
    hasMaterials: Boolean(opc.materialsSubmittedAt && opc.legalPersonName),
    materialsSubmittedAt: opc.materialsSubmittedAt?.toISOString() ?? null,
    idCardValidFrom: fmtDate(opc.idCardValidFrom),
    idCardValidTo: fmtDate(opc.idCardValidTo),
    attachments: {
      addressProof,
      idCardFront: idFront,
      idCardBack: idBack,
      license,
      bankReceipt,
    },
  };
}

export function getAllowedAdminActions(status: string): OpcAdminAction[] {
  switch (status) {
    case 'materials_review':
      return ['approve_materials', 'reject_materials'];
    case 'registering':
      return ['issue_license'];
    case 'tax':
      return ['complete_tax'];
    case 'bank':
      return ['complete_bank'];
    default:
      return [];
  }
}

export async function listOpcTasks(filters: { status?: string; q?: string }) {
  const tasks = await prisma.opcEntity.findMany({
    where: {
      deleted: false,
      ...(filters.status ? { status: filters.status } : {}),
      ...(filters.q
        ? {
            OR: [
              { companyName: { contains: filters.q } },
              { legalPersonName: { contains: filters.q } },
              { phone: { contains: filters.q } },
              { member: { phone: { contains: filters.q } } },
            ],
          }
        : {}),
    },
    include: {
      member: { select: { id: true, phone: true, name: true } },
      progressLogs: { orderBy: { createdAt: 'desc' }, take: 1 },
    },
    orderBy: { updatedAt: 'desc' },
    take: 200,
  });

  return tasks.map((t) => {
    const names = t.proposedNames as string[] | null;
    const submitted = t.materialsSubmittedAt;
    const slaDay = submitted
      ? Math.floor((Date.now() - submitted.getTime()) / 86400000) + 1
      : null;
    return {
      id: t.id.toString(),
      status: t.status,
      statusLabel: OPC_STATUS_LABELS[t.status] ?? t.status,
      companyName: t.companyName,
      proposedName: names?.[0] ?? null,
      legalPersonName: t.legalPersonName,
      phone: t.phone,
      memberPhone: t.member.phone,
      memberName: t.member.name,
      materialsSubmittedAt: submitted?.toISOString() ?? null,
      slaDay,
      updatedAt: t.updatedAt.toISOString(),
    };
  });
}

export async function applyOpcAdminAction(
  opcId: bigint,
  action: OpcAdminAction,
  adminId: bigint,
  payload: {
    note?: string;
    companyName?: string;
    creditCode?: string;
    establishedAt?: string;
    licenseFileId?: string;
    taxActivatedAt?: string;
    bankAccount?: string;
    bankReceiptFileId?: string;
  },
) {
  const opc = await prisma.opcEntity.findFirst({ where: { id: opcId, deleted: false } });
  if (!opc) throw new Error('任务不存在');

  const allowed = getAllowedAdminActions(opc.status);
  if (!allowed.includes(action)) {
    throw new Error(`当前状态「${OPC_STATUS_LABELS[opc.status] ?? opc.status}」不可执行该操作`);
  }

  const now = new Date();

  if (action === 'reject_materials') {
    if (!payload.note || payload.note.trim().length < 10) {
      throw new Error('驳回原因至少 10 个字');
    }
    await prisma.opcEntity.update({
      where: { id: opcId },
      data: { status: 'materials', rejectNote: payload.note.trim() },
    });
    await prisma.opcProgressLog.create({
      data: { opcId, step: 'materials', status: 'rejected', note: payload.note, operatedBy: adminId },
    });
  } else if (action === 'approve_materials') {
    await prisma.opcEntity.update({
      where: { id: opcId },
      data: { status: 'registering', materialsApprovedAt: now, rejectNote: null },
    });
    await prisma.opcProgressLog.create({
      data: { opcId, step: 'materials', status: 'approved', note: payload.note ?? '资料审核通过', operatedBy: adminId },
    });
  } else if (action === 'issue_license') {
    if (!payload.companyName || !payload.creditCode || !payload.licenseFileId) {
      throw new Error('请填写核准名称、统一社会信用代码并上传执照');
    }
    if (!CREDIT_CODE_RE.test(payload.creditCode)) {
      throw new Error('统一社会信用代码格式不正确');
    }
    const licenseExists = await prisma.sysAttachment.findFirst({
      where: { id: BigInt(payload.licenseFileId), deleted: false },
    });
    if (!licenseExists) throw new Error('执照附件不存在');
    await prisma.opcEntity.update({
      where: { id: opcId },
      data: {
        status: 'tax',
        companyName: payload.companyName,
        creditCode: payload.creditCode,
        establishedAt: payload.establishedAt ? parseDateOnly(payload.establishedAt) : now,
        licenseFileId: BigInt(payload.licenseFileId),
      },
    });
    await prisma.opcProgressLog.create({
      data: {
        opcId,
        step: 'business',
        status: 'license_issued',
        note: payload.note ?? '执照已下发',
        operatedBy: adminId,
      },
    });
  } else if (action === 'complete_tax') {
    if (!payload.taxActivatedAt) throw new Error('请填写电子税务局激活日期');
    await prisma.opcEntity.update({
      where: { id: opcId },
      data: { status: 'bank', taxActivatedAt: new Date(payload.taxActivatedAt) },
    });
    await prisma.opcProgressLog.create({
      data: {
        opcId,
        step: 'tax',
        status: 'completed',
        note: payload.note ?? '税务登记完成',
        operatedBy: adminId,
      },
    });
  } else if (action === 'complete_bank') {
    if (!payload.bankAccount) throw new Error('请填写对公账号');
    await prisma.opcEntity.update({
      where: { id: opcId },
      data: {
        status: 'active',
        bankAccountEnc: encrypt(payload.bankAccount.replace(/\s/g, '')),
        ...(payload.bankReceiptFileId
          ? { bankReceiptFileId: BigInt(payload.bankReceiptFileId) }
          : {}),
      },
    });
    await prisma.opcProgressLog.create({
      data: {
        opcId,
        step: 'bank',
        status: 'completed',
        note: payload.note ?? '银行开户完成',
        operatedBy: adminId,
      },
    });
    await prisma.opcProgressLog.create({
      data: { opcId, step: 'complete', status: 'completed', note: 'OPC 设立完成', operatedBy: adminId },
    });
  }

  await writeAuditLog({
    entityType: 'opc_entity',
    entityId: opcId,
    action: `opc.${action}`,
    operatorId: adminId,
    operatorType: 'admin',
    after: { action, ...payload },
  });

  return getOpcTaskDetail(opcId);
}

/** 签约后创建 OPC 记录 */
export async function ensureOpcEntityOnSign(memberId: bigint) {
  let opc = await prisma.opcEntity.findFirst({ where: { memberId, deleted: false } });
  if (!opc) {
    opc = await prisma.opcEntity.create({
      data: { memberId, status: 'pending' },
    });
    await prisma.opcProgressLog.create({
      data: { opcId: opc.id, step: 'materials', status: 'pending', note: '等待提交资料' },
    });
  }
  return opc;
}
