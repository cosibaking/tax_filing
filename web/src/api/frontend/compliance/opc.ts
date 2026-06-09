/**
 * 前台 OPC 落地 API（需登录）
 * @module api/frontend/compliance/opc
 */
import { memberRequest } from '@/utils/http'

/** OPC 主状态 */
export type OpcStatus =
  | 'unsigned'
  | 'pending'
  | 'materials'
  | 'materials_review'
  | 'registering'
  | 'tax'
  | 'bank'
  | 'active'

/** 时间轴步骤状态 */
export type OpcStepStatus = 'done' | 'current' | 'pending' | 'rejected'

/** 时间轴节点 */
export interface OpcProgressStep {
  key: 'materials' | 'business' | 'tax' | 'bank'
  label: string
  status: OpcStepStatus
  date?: string | null
  subLabel?: string
}

/** OPC 进度 */
export interface OpcProgressResult {
  opcId: number | string
  opcStatus: OpcStatus
  companyName?: string
  creditCode?: string
  estimatedSlaDays?: number
  steps: OpcProgressStep[]
  rejectNote?: string | null
  materialsEditable?: boolean
  materialsReadonly?: Record<string, unknown> | null
  materialsSubmitted?: boolean
  planTier?: string
  planName?: string
  planAmount?: number
  signedAt?: string
}

/** 资料实体卡片 */
export interface MaterialsEntityOverview {
  opcId: number | string
  companyNameMasked: string
  legalPersonSummary: string
  status: 'pending' | 'reviewing' | 'approved' | 'rejected'
  statusLabel: string
}

/** 资料详情分组 */
export interface MaterialsDetailSection {
  key: string
  title: string
  fields: MaterialsSectionField[]
}

/** 单个 OPC 实体完整资料（脱敏） */
export interface MaterialsEntityDetail {
  opcId: number | string
  status: string
  statusLabel: string
  sections: MaterialsDetailSection[]
}

/** 完整资料分组（明文） */
export interface MaterialsRevealSection {
  key: string
  title: string
  fields: MaterialsSectionField[]
  attachments?: MaterialsAttachmentReveal[]
}

/** 全部完整资料 */
export interface MaterialsRevealAllResult {
  opcId: number | string
  sections: MaterialsRevealSection[]
}

/** 资料字段 */
export interface MaterialsSectionField {
  label: string
  value: string
}

/** 完整资料附件项 */
export interface MaterialsAttachmentReveal {
  label: string
  fileId: number | string
  fileName: string
  mimeType: string
  accessUrl: string
}

/** 资料分组详情 */
export interface MaterialsSectionDetail {
  key: string
  title: string
  status: string
  statusLabel: string
  fields: MaterialsSectionField[]
}

export const PLAN_TIER_LABELS: Record<string, string> = {
  basic: '基础套餐',
  advanced: '进阶套餐',
  premium: '尊享套餐'
}

/** 注册资料提交参数 */
export interface OpcMaterialsParams {
  proposedNames: string[]
  registeredCapital: number
  capitalTermYears: number
  businessTermType: 'long_term' | 'fixed'
  businessTermEnd?: string
  businessScope: string
  registerProvince: string
  registerCity: string
  registerDistrict: string
  registerAddress: string
  addressProofFileId: number | string
  legalPersonName: string
  idCardNumber: string
  idCardValidFrom: string
  idCardValidTo: string
  ethnicity?: string
  householdAddress: string
  residenceAddress: string
  phone: string
  email: string
  idCardFrontFileId: number | string
  idCardBackFileId: number | string
  confirmations: {
    infoTrue: boolean
    authConsent: boolean
    opcLimitAck: boolean
    eSignAuth: boolean
  }
}

/** 银行回执提交 */
export interface OpcBankReceiptParams {
  bankName?: string
  bankReceiptFileId: number | string
}

/** OPC 主体信息 */
export interface OpcEntityInfo {
  opcId?: number | string
  id?: number | string
  companyName?: string
  creditCode?: string
  status: OpcStatus
  bankAccountMasked?: string
  employmentStatus?: 'unknown' | 'no_employee' | 'has_employee'
  employmentConfirmedAt?: string
  planTier?: string
  isActive?: boolean
}

/** 获取 OPC 进度与时间轴 */
export function getOpcProgress() {
  return memberRequest.get<OpcProgressResult>({
    url: '/compliance/opc/progress'
  })
}

/** 获取已提交资料概览（按 OPC 实体） */
export function getMaterialsOverview() {
  return memberRequest.get<{ entities: MaterialsEntityOverview[] }>({
    url: '/compliance/opc/materials/overview'
  })
}

/** 获取单个 OPC 实体完整资料（脱敏） */
export function getMaterialsDetail(opcId: number | string) {
  return memberRequest.get<MaterialsEntityDetail>({
    url: '/compliance/opc/materials/detail',
    params: { opcId }
  })
}

/** 获取资料分组详情（脱敏） */
export function getMaterialsSection(section: string, opcId?: number | string) {
  return memberRequest.get<MaterialsSectionDetail>({
    url: '/compliance/opc/materials/section',
    params: { section, opcId }
  })
}

/** 密码验证后查看完整资料（单分组） */
export function revealMaterialsSection(section: string, password: string, opcId: number | string) {
  return memberRequest.post<{
    key: string
    fields: MaterialsSectionField[]
    attachments?: MaterialsAttachmentReveal[]
  }>({
    url: '/compliance/opc/materials/reveal',
    data: { section, password, opcId }
  })
}

/** 密码验证后查看全部完整资料 */
export function revealMaterialsAll(password: string, opcId: number | string) {
  return memberRequest.post<MaterialsRevealAllResult>({
    url: '/compliance/opc/materials/reveal-all',
    data: { password, opcId }
  })
}

/** 获取 OPC 主体信息（未签约时返回 null） */
export async function getOpcEntity() {
  const res = await memberRequest.get<OpcEntityInfo | null>({
    url: '/compliance/opc'
  })
  if (!res || res.status === 'unsigned') return null
  return res
}

/** 后端材料提交请求体 */
export interface OpcMaterialsApiPayload {
  proposedNames: string[]
  registeredCapital: number
  capitalTermYears: number
  businessTermType: 'long_term' | 'fixed'
  businessTermEnd?: string
  businessScope: string
  registerProvince: string
  registerCity: string
  registerDistrict: string
  registerAddress: string
  addressProofFileId: number
  legalPersonName: string
  idCard: string
  idCardValidFrom: string
  idCardValidTo: string
  ethnicity?: string
  householdAddress: string
  residentialAddress: string
  phone: string
  email: string
  idCardFrontFileId: number
  idCardBackFileId: number
  esignAuthorized: boolean
  confirmations: {
    truthful: boolean
    usageConsent: boolean
    opcLimitAck: boolean
  }
}

function toAttachmentId(value: number | string, label: string): number {
  const id = Number(value)
  if (!id || id < 1) {
    throw new Error(`请重新上传${label}`)
  }
  return id
}

/** 将表单数据映射为后端 API 格式 */
export function mapMaterialsToApi(data: OpcMaterialsParams): OpcMaterialsApiPayload {
  return {
    proposedNames: data.proposedNames,
    registeredCapital: data.registeredCapital,
    capitalTermYears: data.capitalTermYears,
    businessTermType: data.businessTermType,
    businessTermEnd: data.businessTermEnd || '',
    businessScope: data.businessScope,
    registerProvince: data.registerProvince,
    registerCity: data.registerCity,
    registerDistrict: data.registerDistrict,
    registerAddress: data.registerAddress,
    addressProofFileId: toAttachmentId(data.addressProofFileId, '地址证明'),
    legalPersonName: data.legalPersonName,
    idCard: data.idCardNumber,
    idCardValidFrom: data.idCardValidFrom,
    idCardValidTo: data.idCardValidTo,
    ethnicity: data.ethnicity,
    householdAddress: data.householdAddress,
    residentialAddress: data.residenceAddress,
    phone: data.phone,
    email: data.email,
    idCardFrontFileId: toAttachmentId(data.idCardFrontFileId, '身份证正面'),
    idCardBackFileId: toAttachmentId(data.idCardBackFileId, '身份证反面'),
    esignAuthorized: data.confirmations.eSignAuth,
    confirmations: {
      truthful: data.confirmations.infoTrue,
      usageConsent: data.confirmations.authConsent,
      opcLimitAck: data.confirmations.opcLimitAck,
    },
  }
}

/** 提交 OPC 注册资料 */
export function submitOpcMaterials(data: OpcMaterialsParams) {
  return memberRequest.post<{ opcId: number | string }>({
    url: '/compliance/opc/materials',
    data: mapMaterialsToApi(data),
  })
}

/** 上传银行开户回执 */
export function submitBankReceipt(data: OpcBankReceiptParams) {
  const fileId = Number(data.bankReceiptFileId)
  if (!fileId || fileId < 1) {
    return Promise.reject(new Error('请上传开户回执'))
  }
  return memberRequest.post<void>({
    url: '/compliance/opc/bank-receipt',
    data: {
      bankName: data.bankName,
      bankReceiptFileId: fileId,
    },
  })
}

/** 状态中文映射 */
export const OPC_STATUS_LABELS: Record<OpcStatus, string> = {
  unsigned: '未签约',
  pending: '待提交资料',
  materials: '资料待补正',
  materials_review: '资料审核中',
  registering: '工商注册中',
  tax: '税务登记中',
  bank: '银行开户中',
  active: '已激活'
}

/** 推荐经营范围模板 */
export const RECOMMENDED_BUSINESS_SCOPE =
  '一般项目：技术服务、技术开发、技术咨询；文化艺术交流活动组织；文艺创作；个人互联网直播服务；摄像及视频制作服务；数字内容制作服务（不含出版发行）；广告设计、代理；广告发布。（除依法须经批准的项目外，凭营业执照依法自主开展经营活动）'
