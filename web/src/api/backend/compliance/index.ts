/**
 * 后台合规服务 API
 * @module api/backend/compliance
 */
import { adminRequest } from '@/utils/http'
import { FILING_CHECKLIST_ITEMS as CHECKLIST_DEFS } from './checklist'
import type { OpcStatus } from '@/api/frontend/compliance/opc'

export { FILING_CHECKLIST_ITEMS, isChecklistItemNA } from './checklist'

/** OPC 任务列表项 */
export interface OpcTaskItem {
  id: number | string
  memberId: number
  memberName?: string
  memberPhone?: string
  proposedName?: string
  status: OpcStatus
  planName?: string
  materialsSubmittedAt?: string
  slaDays?: number
  createdAt?: string
}

/** OPC 任务列表查询 */
export interface OpcTaskListParams {
  page?: number
  pageSize?: number
  status?: OpcStatus | ''
  q?: string
}

/** OPC 任务列表响应 */
export interface OpcTaskListResult {
  list: OpcTaskItem[]
  total: number
  page: number
  pageSize: number
}

/** OPC 任务详情 */
export interface OpcTaskDetail extends OpcTaskItem {
  proposedNames?: string[]
  registeredCapital?: number
  capitalTermYears?: number
  businessScope?: string
  registerAddress?: string
  addressProofUrl?: string
  legalPersonName?: string
  idCardNumberMasked?: string
  idCardValidFrom?: string
  idCardValidTo?: string
  householdAddress?: string
  residenceAddress?: string
  phone?: string
  email?: string
  idCardFrontUrl?: string
  idCardBackUrl?: string
  companyName?: string
  creditCode?: string
  establishedAt?: string
  licenseFileUrl?: string
  taxActivatedAt?: string
  taxpayerType?: string
  bankName?: string
  bankAccountMasked?: string
  bankReceiptUrl?: string
  rejectNote?: string
  progressLogs?: Array<{
    step: string
    status: string
    note?: string
    createdAt: string
  }>
}

/** 推进 OPC 任务 */
export interface OpcTaskPatchParams {
  id: number | string
  action:
    | 'approve_materials'
    | 'reject_materials'
    | 'issue_license'
    | 'complete_tax'
    | 'complete_bank'
  note?: string
  companyName?: string
  creditCode?: string
  establishedAt?: string
  licenseFileId?: string | number
  taxActivatedAt?: string
  bankAccount?: string
  bankReceiptFileId?: string | number
}

/** 合规客户列表项 */
export interface ComplianceCustomerItem {
  id: number
  memberId: number
  opcId?: number
  username?: string
  nickname?: string
  mobile?: string
  planName?: string
  opcStatus?: OpcStatus
  opcCompanyName?: string
  signedAt?: string
  createdAt?: string
}

/** 合规客户列表查询 */
export interface ComplianceCustomerListParams {
  page?: number
  pageSize?: number
  q?: string
  opcStatus?: OpcStatus | ''
}

/** 合规客户列表响应 */
export interface ComplianceCustomerListResult {
  list: ComplianceCustomerItem[]
  total: number
  page: number
  pageSize: number
}

/** OPC 状态选项（筛选用） */
export const OPC_STATUS_OPTIONS = [
  { label: '全部', value: '' },
  { label: '待提交资料', value: 'pending' },
  { label: '资料待补正', value: 'materials' },
  { label: '资料审核中', value: 'materials_review' },
  { label: '工商注册中', value: 'registering' },
  { label: '税务登记中', value: 'tax' },
  { label: '银行开户中', value: 'bank' },
  { label: '已激活', value: 'active' }
] as const

/** OPC 状态中文 */
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

const PLAN_TIER_LABELS: Record<string, string> = {
  basic: '基础套餐',
  advanced: '进阶套餐',
  premium: '尊享套餐'
}

function mapOpcListItem(raw: Record<string, any>): OpcTaskItem {
  return {
    id: raw.opcId ?? raw.id,
    memberId: raw.memberId,
    memberName: raw.memberName,
    memberPhone: raw.memberPhoneMasked ?? raw.memberPhone,
    proposedName: raw.proposedName,
    status: raw.status,
    planName: PLAN_TIER_LABELS[raw.planTier] ?? raw.planTier ?? raw.planName,
    materialsSubmittedAt: raw.materialsSubmittedAt,
    slaDays: raw.daysSinceSubmit ?? raw.slaDays
  }
}

function mapOpcDetail(raw: Record<string, any>): OpcTaskDetail {
  const rejectLog = (raw.progressLogs || []).find(
    (l: { status?: string; note?: string }) => l.status === 'rejected' && l.note
  )
  return {
    ...mapOpcListItem(raw),
    memberName: raw.memberName,
    materialsSubmittedAt: raw.materialsSubmittedAt,
    proposedNames: raw.company?.proposedNames ?? raw.proposedNames,
    registeredCapital: raw.company?.registeredCapital,
    capitalTermYears: raw.company?.capitalTermYears,
    businessScope: raw.company?.businessScope,
    registerAddress: raw.company?.registerAddress,
    legalPersonName: raw.legalPerson?.legalPersonName,
    idCardNumberMasked: raw.legalPerson?.idCardMasked,
    idCardValidFrom: raw.legalPerson?.idCardValidFrom,
    idCardValidTo: raw.legalPerson?.idCardValidTo,
    householdAddress: raw.legalPerson?.householdAddress,
    residenceAddress: raw.legalPerson?.residentialAddress,
    phone: raw.legalPerson?.phone,
    email: raw.legalPerson?.email,
    companyName: raw.business?.companyName,
    creditCode: raw.business?.creditCode,
    establishedAt: raw.business?.establishedAt,
    taxActivatedAt: raw.tax?.taxActivatedAt,
    taxpayerType: raw.tax?.taxpayerType,
    bankName: raw.bank?.bankName,
    bankAccountMasked: raw.bank?.bankAccountMasked,
    rejectNote: rejectLog?.note ?? raw.rejectNote,
    progressLogs: raw.progressLogs
  }
}

function mapCustomerItem(raw: Record<string, any>): ComplianceCustomerItem {
  return {
    id: raw.memberId,
    memberId: raw.memberId,
    opcId: raw.opcId,
    nickname: raw.memberName,
    mobile: raw.memberPhoneMasked ?? raw.mobile,
    planName: PLAN_TIER_LABELS[raw.planTier] ?? raw.planTier ?? raw.planName,
    opcStatus: raw.opcStatus,
    opcCompanyName: raw.opcCompanyName,
    signedAt: raw.signedAt
  }
}

function toChecklistItems(checklist: Record<string, boolean>) {
  return CHECKLIST_DEFS.map((item) => ({
    key: item.key,
    label: item.label,
    checked: !!checklist[item.key]
  }))
}

function toFileId(value?: string | number) {
  if (value == null || value === '') return 0
  const n = Number(value)
  return Number.isFinite(n) ? n : 0
}

/** 获取 OPC 任务列表 */
export async function getOpcTaskList(params: OpcTaskListParams) {
  const res = await adminRequest.get<OpcTaskListResult>({
    url: '/compliance/opc-tasks',
    params
  })
  return {
    ...res,
    list: (res.list || []).map(mapOpcListItem)
  }
}

/** 获取 OPC 任务详情 */
export async function getOpcTaskDetail(id: number | string) {
  const raw = await adminRequest.get<Record<string, any>>({
    url: `/compliance/opc-tasks/${id}`
  })
  return mapOpcDetail(raw)
}

/** 推进 OPC 任务进度 */
export async function patchOpcTask(data: OpcTaskPatchParams) {
  const raw = await adminRequest.request<Record<string, any>>({
    url: '/compliance/opc-tasks',
    method: 'PATCH',
    data: {
      opcId: Number(data.id),
      action: data.action,
      note: data.note,
      companyName: data.companyName,
      creditCode: data.creditCode,
      establishedAt: data.establishedAt,
      licenseFileId: toFileId(data.licenseFileId),
      taxActivatedAt: data.taxActivatedAt,
      bankAccount: data.bankAccount,
      bankReceiptFileId: toFileId(data.bankReceiptFileId)
    }
  })
  return mapOpcDetail(raw)
}

/** 获取合规客户列表 */
export async function getComplianceCustomerList(params: ComplianceCustomerListParams) {
  const res = await adminRequest.get<ComplianceCustomerListResult>({
    url: '/compliance/customers',
    params
  })
  return {
    ...res,
    list: (res.list || []).map(mapCustomerItem)
  }
}

/** 导出会员合规留痕 CSV */
export function exportComplianceAudit(memberId: number) {
  return adminRequest.get<Blob>({
    url: '/compliance/audit/export',
    params: { memberId },
    responseType: 'blob'
  })
}

/** 税种类型 */
export type FilingTaxType = 'vat' | 'cit_quarterly' | 'cit_annual' | 'surcharge'

/** 申报任务状态 */
export type FilingStatus = 'pending' | 'filed' | 'overdue'

/** 申报工作台列表项 */
export interface FilingTaskItem {
  id: number | string
  memberId: number
  memberName?: string
  memberPhone?: string
  opcCompanyName?: string
  opcId?: number
  employmentStatus?: string
  taxType: FilingTaxType
  period: string
  dueDate: string
  calculatedAmount: number
  filedAmount?: number
  status: FilingStatus
  receiptFileId?: string
  checklist?: Record<string, boolean>
}

/** 申报列表查询 */
export interface FilingListParams {
  page?: number
  pageSize?: number
  taxType?: FilingTaxType | ''
  period?: string
  status?: FilingStatus | ''
  q?: string
}

/** 申报列表响应 */
export interface FilingListResult {
  list: FilingTaskItem[]
  total: number
  page: number
  pageSize: number
}

/** 标记已申报参数 */
export interface MarkFilingFiledParams {
  id: number | string
  checklist: Record<string, boolean>
  filedAmount?: number
  receiptFileId: string | number
}

/** 对账单管理列表项 */
export interface AdminStatementItem {
  id: number | string
  memberId: number
  memberName?: string
  memberPhone?: string
  opcCompanyName?: string
  period: string
  revenue: number
  cost: number
  profit: number
  prepaidTax: number
  status: 'draft' | 'sent' | 'completed'
  sentAt?: string
  pdfUrl?: string
}

/** 对账单列表查询 */
export interface AdminStatementListParams {
  page?: number
  pageSize?: number
  period?: string
  status?: string
  q?: string
}

/** 对账单列表响应 */
export interface AdminStatementListResult {
  list: AdminStatementItem[]
  total: number
  page: number
  pageSize: number
}

/** 批量生成对账单参数 */
export interface BatchGenerateStatementsParams {
  period: string
  memberIds?: number[]
}

/** 发送对账单通知参数 */
export interface SendStatementsParams {
  ids: (number | string)[]
}

/** 税种选项 */
export const FILING_TAX_TYPE_OPTIONS = [
  { label: '全部', value: '' },
  { label: '增值税', value: 'vat' },
  { label: '企税季度预缴', value: 'cit_quarterly' },
  { label: '企税年度汇算', value: 'cit_annual' },
  { label: '附加税费', value: 'surcharge' }
] as const

/** 申报状态选项 */
export const FILING_STATUS_OPTIONS = [
  { label: '全部', value: '' },
  { label: '待申报', value: 'pending' },
  { label: '已申报', value: 'filed' },
  { label: '已逾期', value: 'overdue' }
] as const

/** 税种中文 */
export const FILING_TAX_TYPE_LABELS: Record<FilingTaxType, string> = {
  vat: '增值税',
  cit_quarterly: '企税季度预缴',
  cit_annual: '企税年度汇算',
  surcharge: '附加税费'
}

/** 申报状态中文 */
export const FILING_STATUS_LABELS: Record<FilingStatus, string> = {
  pending: '待申报',
  filed: '已申报',
  overdue: '已逾期'
}

/** 对账单状态中文 */
export const STATEMENT_STATUS_LABELS: Record<string, string> = {
  draft: '草稿',
  sent: '已发送',
  completed: '已完成'
}

function mapFilingItem(raw: Record<string, any>): FilingTaskItem {
  const taxType = raw.taxType ?? raw.tax_type
  const status = raw.status
  const calculatedAmount = Number(raw.calculatedAmount ?? raw.calculated_amount ?? 0)
  return {
    id: raw.id,
    memberId: raw.memberId ?? raw.member_id,
    memberName: raw.memberName ?? raw.member_name,
    opcCompanyName: raw.companyName ?? raw.opcCompanyName ?? raw.company_name,
    opcId: raw.opcId ?? raw.opc_id,
    employmentStatus: raw.employmentStatus ?? raw.employment_status,
    taxType,
    period: raw.period,
    dueDate: raw.dueDate ?? raw.due_date,
    calculatedAmount: Number.isFinite(calculatedAmount) ? calculatedAmount : 0,
    filedAmount: raw.filedAmount ?? raw.filed_amount,
    status
  }
}

/** 获取申报任务列表 */
export async function getFilingList(params: FilingListParams) {
  const res = await adminRequest.get<FilingListResult>({
    url: '/compliance/filing',
    params
  })
  return {
    ...res,
    list: (res.list || []).map(mapFilingItem)
  }
}

/** 标记已申报 */
export function markFilingFiled(data: MarkFilingFiledParams) {
  return adminRequest.request<FilingTaskItem>({
    url: '/compliance/filing',
    method: 'PATCH',
    data: {
      taskId: Number(data.id),
      checklist: toChecklistItems(data.checklist),
      filedAmount: data.filedAmount,
      receiptFileId: toFileId(data.receiptFileId)
    }
  })
}

/** 获取对账单管理列表 */
export function getAdminStatementList(params: AdminStatementListParams) {
  return adminRequest.get<AdminStatementListResult>({
    url: '/compliance/statements',
    params
  })
}

/** 批量生成对账单 */
export function batchGenerateStatements(data: BatchGenerateStatementsParams) {
  return adminRequest.post<{ generated: number }>({
    url: '/compliance/statements/generate',
    data
  })
}

/** 发送对账单通知 */
export function sendStatementNotifications(data: SendStatementsParams) {
  return adminRequest.post<{ sent: number }>({
    url: '/compliance/statements/send',
    data: {
      ids: data.ids.map((id) => Number(id))
    }
  })
}

/** 获取对账单 PDF 下载地址 */
export function getAdminStatementPdfUrl(id: number | string) {
  return adminRequest.get<{ url: string }>({
    url: `/compliance/statements/${id}/pdf`
  })
}

/** 工作台概览数据 */
export interface ComplianceDashboardData {
  activeCustomers: number
  pendingOpcTasks: number
  pendingFilings: number
  overdueFilings: number
  draftStatements: number
  opcTasks: Array<{
    opcId: number
    memberName?: string
    proposedName?: string
    status: string
    statusLabel?: string
    materialsSubmittedAt?: string
    daysSinceSubmit?: number
  }>
  filingTasks: Array<{
    id: number
    memberName?: string
    companyName?: string
    taxType: string
    taxTypeLabel?: string
    period: string
    dueDate: string
    status: string
    calculatedAmount: number
  }>
}

/** 获取合规工作台概览 */
export function getComplianceDashboard() {
  return adminRequest.get<ComplianceDashboardData>({
    url: '/compliance/dashboard'
  })
}
