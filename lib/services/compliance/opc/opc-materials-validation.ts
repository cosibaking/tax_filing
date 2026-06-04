import type { OpcMaterialsInput } from './opc-types';

const ID_CARD_RE = /^[1-9]\d{16}[\dXx]$/;
const PHONE_RE = /^1\d{10}$/;
const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export function parseDateOnly(value: string): Date {
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) throw new Error('日期格式无效');
  return d;
}

/** 将 API/表单 JSON 规范为强类型，避免字符串数字导致校验失败 */
export function parseMaterialsBody(raw: unknown): OpcMaterialsInput {
  const b = (raw ?? {}) as Record<string, unknown>;
  const conf = (b.confirmations ?? {}) as Record<string, unknown>;
  const namesRaw = b.proposedNames;
  const proposedNames = Array.isArray(namesRaw)
    ? namesRaw.map((n) => String(n).trim()).filter(Boolean)
    : [];

  return {
    proposedNames,
    registeredCapital: Number(b.registeredCapital),
    capitalTermYears: Number(b.capitalTermYears),
    businessTermType: b.businessTermType === 'fixed' ? 'fixed' : 'long_term',
    businessTermEnd: b.businessTermEnd ? String(b.businessTermEnd) : undefined,
    businessScope: String(b.businessScope ?? '').trim(),
    registerProvince: String(b.registerProvince ?? '').trim(),
    registerCity: String(b.registerCity ?? '').trim(),
    registerDistrict: String(b.registerDistrict ?? '').trim(),
    registerAddress: String(b.registerAddress ?? '').trim(),
    addressProofFileId: String(b.addressProofFileId ?? '').trim(),
    legalPersonName: String(b.legalPersonName ?? '').trim(),
    idCard: String(b.idCard ?? '').trim(),
    idCardValidFrom: String(b.idCardValidFrom ?? '').trim(),
    idCardValidTo: String(b.idCardValidTo ?? '').trim(),
    householdAddress: String(b.householdAddress ?? '').trim(),
    residentialAddress: String(b.residentialAddress ?? '').trim(),
    phone: String(b.phone ?? '').trim(),
    email: String(b.email ?? '').trim(),
    idCardFrontFileId: String(b.idCardFrontFileId ?? '').trim(),
    idCardBackFileId: String(b.idCardBackFileId ?? '').trim(),
    ethnicity: b.ethnicity ? String(b.ethnicity).trim() : undefined,
    esignAuthorized: b.esignAuthorized === true || b.esignAuthorized === 'true',
    confirmations: {
      truthful: conf.truthful === true || conf.truthful === 'true',
      usageConsent: conf.usageConsent === true || conf.usageConsent === 'true',
      opcLimitAck: conf.opcLimitAck === true || conf.opcLimitAck === 'true',
    },
  };
}

export function validateOpcMaterials(data: OpcMaterialsInput): string | null {
  if (!data.proposedNames?.length || data.proposedNames.length > 3) {
    return '请填写 1～3 个备选公司名称';
  }
  for (const name of data.proposedNames) {
    if (name.length < 2 || name.length > 30) return '公司名称长度应为 2～30 字';
  }
  if (
    !Number.isFinite(data.registeredCapital) ||
    data.registeredCapital <= 0 ||
    data.registeredCapital > 1000
  ) {
    return '注册资本应在 0～1000 万元之间';
  }
  if (![5, 10, 20, 30].includes(data.capitalTermYears)) return '请选择认缴期限';
  if (!data.businessScope || data.businessScope.length < 10) return '请填写经营范围';
  if (!data.registerProvince || !data.registerCity || !data.registerDistrict) {
    return '请选择注册地址省市区';
  }
  if (!data.registerAddress || data.registerAddress.length < 5) return '请填写详细注册地址';
  if (!data.legalPersonName || data.legalPersonName.length < 2) return '请填写法人姓名';
  if (!ID_CARD_RE.test(data.idCard)) return '身份证号格式不正确';
  if (!PHONE_RE.test(data.phone)) return '手机号格式不正确';
  if (!EMAIL_RE.test(data.email)) return '邮箱格式不正确';
  if (!data.idCardValidFrom || !data.idCardValidTo) return '请填写身份证有效期';
  if (!data.idCardFrontFileId || !data.idCardBackFileId) return '请上传身份证正反面';
  if (!data.addressProofFileId) return '请上传地址证明';
  if (!/^\d+$/.test(data.addressProofFileId)) return '地址证明未上传成功，请重新上传';
  if (!/^\d+$/.test(data.idCardFrontFileId) || !/^\d+$/.test(data.idCardBackFileId)) {
    return '身份证照片未上传成功，请重新上传';
  }
  if (!data.householdAddress || !data.residentialAddress) return '请填写户籍与现居住地址';
  if (!data.esignAuthorized) return '请确认电子签名授权';
  if (!data.confirmations?.truthful || !data.confirmations?.usageConsent || !data.confirmations?.opcLimitAck) {
    return '请勾选全部确认项';
  }
  const validTo = parseDateOnly(data.idCardValidTo);
  if (validTo < new Date()) return '身份证已过期';
  if (data.businessTermType === 'fixed' && !data.businessTermEnd) return '请选择营业期限截止日';
  return null;
}
