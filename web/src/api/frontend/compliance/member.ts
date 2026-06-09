/**

 * 会员合规 API（需登录，Xy-User-Token）

 * @module api/frontend/compliance/member

 */

import { memberRequest } from '@/utils/http'

import type { RecommendedPlan, TaxComparison } from './diagnosis'

import type { CompliancePlanState, OpcStatus } from '@/config/complianceMenu'

import { getOpcEntity } from './opc'
import { getActiveOrder } from './order'



/** 诊断历史项 */

export interface DiagnosisHistoryItem {

  id: number

  platforms: string[]

  monthlyIncomeRange: string

  annualCostEstimate: number

  existingEntity: string

  hasFiledTax: string

  taxBureauContact: boolean

  recommendedPlan: RecommendedPlan

  taxComparison: TaxComparison

  createdAt: string

}



/** 诊断历史分页响应 */

export interface DiagnosisHistoryResult {

  list: DiagnosisHistoryItem[]

  page: number

  pageSize: number

  total: number

}



/** 收入渠道/平台 */

export const INCOME_PLATFORMS = [

  { value: '', label: '全部' },

  { value: 'douyin', label: '抖音' },

  { value: 'kuaishou', label: '快手' },

  { value: 'bilibili', label: 'B站' },

  { value: 'channels', label: '视频号' },

  { value: 'xiaohongshu', label: '小红书' },

  { value: 'taobao', label: '淘宝/天猫' },

  { value: 'wechat', label: '微信生态' },

  { value: 'alipay', label: '支付宝' },

  { value: 'offline', label: '线下收款' },

  { value: 'other', label: '其他' }

] as const



export type IncomePlatform = (typeof INCOME_PLATFORMS)[number]['value']



/** 收入类型（F-27） */

export const INCOME_CATEGORIES = [

  { value: 'tip', label: '打赏/酬劳' },

  { value: 'commission', label: '佣金' },

  { value: 'ad', label: '广告' },

  { value: 'service_fee', label: '服务费' },

  { value: 'product_sales', label: '商品销售' },

  { value: 'slot_fee', label: '坑位费' },

  { value: 'offline', label: '线下业务' },

  { value: 'other', label: '其他' }

] as const



export type IncomeCategory = (typeof INCOME_CATEGORIES)[number]['value']



/** 收入台账项 */

export interface IncomeEntry {

  id: number | string

  occurredAt: string

  platform: string

  category: string

  grossAmount: number

  platformFee: number

  netAmount: number

  source?: string

  settlementType?: string

  mcnName?: string

  mcnSplitRatio?: number

  mcnShareAmount?: number

  grossBeforeSplit?: number

  remark?: string

}



/** 收入列表查询 */

export interface IncomeListParams {

  month?: string

  platform?: string

  page?: number

  pageSize?: number

}



/** 收入列表响应 */

export interface IncomeListResult {

  list: IncomeEntry[]

  total: number

  page: number

  pageSize: number

  summary?: {

    grossTotal: number

    netTotal: number

  }

}



/** 新增收入参数 */

export interface IncomeCreateParams {

  occurredAt: string

  platform: string

  category: string

  grossAmount: number

  platformFee: number

  settlementType?: string

  mcnName?: string

  mcnSplitRatio?: number

  grossBeforeSplit?: number

  remark?: string

}

/** 收入一致性比对 */

export interface IncomeConsistencyResult {

  month: string

  ledgerGross: number

  ledgerNet: number

  bankInflow: number

  bankMatched: number

  bankUnmatched: number

  platformReported: number

  varianceLedgerBank: number

  varianceLedgerPlatform: number

  varianceRate: number

  status: 'ok' | 'warning' | 'danger'

  incomeMatchPassed: boolean

  hints?: string[]

  platforms?: {

    platform: string

    ledgerNet: number

    platformReported: number

    variance: number

    status: string

  }[]

}

/** OCR 预览响应 */

export interface IncomeOCRPreviewResult {

  platform: string

  ocrSource: string

  rawTextPreview?: string

  rows: IncomeImportPreviewRow[]

  validCount: number

  invalidCount: number

  totalGross: number

  totalNet: number

}



/** CSV 导入预览行 */

export interface IncomeImportPreviewRow {

  row: number

  occurredAt: string

  platform: string

  category: string

  grossAmount: number

  platformFee: number

  valid: boolean

  error?: string

}



/** CSV 导入预览响应 */

export interface IncomeImportPreviewResult {

  rows: IncomeImportPreviewRow[]

  validCount: number

  invalidCount: number

}



/** 银行流水未匹配项 */

export interface BankUnmatchedItem {

  id: number | string

  occurredAt: string

  amount: number

  description?: string

}



/** 费用发票类型 */

export const INVOICE_TYPES = [

  { value: 'general', label: '普票' },

  { value: 'special', label: '专票' },

  { value: 'none', label: '无票' }

] as const



/** 费用类型（F-37） */

export interface ExpenseCategory {

  code: string

  name: string

  voucherHint: string

}



/** 费用台账项 */

export interface ExpenseEntry {

  id: number | string

  occurredAt: string

  category: string

  categoryName?: string

  amount: number

  invoiceType: string

  attachmentId?: string

  description?: string

  warningFlag?: boolean

}



/** 费用列表查询 */

export interface ExpenseListParams {

  month?: string

  page?: number

  pageSize?: number

}



/** 费用列表响应 */

export interface ExpenseListResult {

  list: ExpenseEntry[]

  total: number

  page: number

  pageSize: number

  summary?: {

    totalAmount: number

    incomeTotal?: number

    costRatio?: number

  }

}



/** 新增费用参数 */

export interface ExpenseCreateParams {

  occurredAt: string

  category: string

  amount: number

  invoiceType: string

  attachmentId?: string

  description?: string

  costRatioAck?: boolean

}



/** 虚增拦截关键词（F-38） */

export const EXPENSE_BLACKLIST_KEYWORDS = ['咨询费', '服务费', '借款', '赠与']



/** 利润表周期 */

export type ProfitPeriodType = 'month' | 'quarter' | 'year'



/** 利润表摘要 */

export interface ProfitSummary {

  period: string

  periodType: ProfitPeriodType

  revenue: number

  cost: number

  profit: number

  cumulativeProfit: number

}



/** 会计分录 */

export interface LedgerVoucher {

  id: number | string

  occurredAt: string

  summary: string

  debitAccount: string

  creditAccount: string

  amount: number

}



/** 税种任务状态 */

export type TaxTaskStatus = 'pending' | 'filed' | 'overdue'



/** 税种类型 */

export type TaxType = 'vat' | 'cit_quarterly' | 'cit_annual' | 'surcharge'



/** 申报任务 */

export interface TaxTask {

  id: number | string

  taxType: TaxType

  taxTypeLabel: string

  period: string

  dueDate: string

  status: TaxTaskStatus

  calculatedAmount: number

  filedAmount?: number

  receiptUrl?: string

  detail?: TaxTaskDetail

}



/** 任务计算明细 */

export interface TaxTaskDetail {

  revenueExTax: number

  vat: number

  surcharge: number

  cit: number

  total: number

}



/** 申报日历 */

export interface TaxCalendarResult {

  year: number

  month: number

  nextDueDate?: string

  daysUntilDue?: number

  dueDates: string[]

  tasks: TaxTask[]

}



/** 自查清单项（F-69） */

export interface TaxChecklistItem {

  key: string

  label: string

  checked: boolean

  na?: boolean

  hint?: string

}



/** 月度对账单 */

export interface StatementItem {

  id: number | string

  period: string

  revenue: number

  cost: number

  profit: number

  prepaidTax: number

  cumulativeProfit: number

  filingStatus: string

  status: 'draft' | 'completed'

  pdfUrl?: string

}



/** 对账单列表响应 */

export interface StatementListResult {

  list: StatementItem[]

  total: number

}



/** 报税前自查清单九项 */

export const TAX_CHECKLIST_ITEMS: { key: string; label: string }[] = [

  { key: 'income_match', label: '收入流水是否与平台数据一致？' },

  { key: 'invoice_filed', label: '是否有该申报未申报的发票？' },

  { key: 'cost_booked', label: '成本费用发票是否已入账？' },

  { key: 'payroll_tax', label: '经营主体是否有员工要报工资个税？' },

  { key: 'vat_filed', label: '本季度增值税申报了吗？' },

  { key: 'cit_prepaid', label: '企业所得税预缴了吗？' },

  { key: 'social_insurance', label: '有员工的话，社保申报了吗？' },

  { key: 'other_platform', label: '是否还有其他渠道收入要合并计算？' },

  { key: 'prior_correction', label: '上一期申报是否有错误要更正？' }

]



export const TAX_TYPE_LABELS: Record<TaxType, string> = {

  vat: '增值税（小规模）',

  cit_quarterly: '企税季度预缴',

  cit_annual: '企税年度汇算',

  surcharge: '附加税费'

}



export const TAX_STATUS_LABELS: Record<TaxTaskStatus, string> = {

  pending: '待申报',

  filed: '已申报',

  overdue: '已逾期'

}



/** 获取会员诊断历史 */

export function getDiagnosisHistory(params?: { page?: number; pageSize?: number }) {

  return memberRequest.get<DiagnosisHistoryResult>({

    url: '/compliance/diagnosis',

    params

  })

}



/** 登录后同步访客诊断缓存 */

export function syncGuestDiagnoses(items: import('./diagnosis').DiagnosisSubmitParams[]) {

  return memberRequest.post<{ savedIds: number[]; count: number }>({

    url: '/compliance/diagnosis/sync',

    data: { items },

    showErrorMessage: false

  })

}



/** 绑定匿名诊断记录到当前会员 */

export function bindDiagnosisRecords(diagnosisIds: number[]) {

  return memberRequest.post<{ boundCount: number }>({

    url: '/compliance/diagnosis/bind',

    data: { diagnosisIds },

    showErrorMessage: false

  })

}



/** 获取诊断详情（已登录，从历史页查看） */

export function getDiagnosisDetail(id: number | string) {

  return memberRequest.get<import('./diagnosis').DiagnosisSubmitResult>({

    url: `/compliance/diagnosis/${id}`,

    showErrorMessage: false

  })

}



/**

 * 获取合规套餐/订单状态（驱动侧栏可见性）

 */

export async function getCompliancePlanState(): Promise<CompliancePlanState> {
  try {
    const order = await getActiveOrder()
    if (!order || order.status === 'cancelled') {
      return { hasActiveOrder: false, opcStatus: 'none', hasPendingOrder: false }
    }
    if (order.status === 'pending') {
      return { hasActiveOrder: false, opcStatus: 'none', hasPendingOrder: true }
    }

    let opcStatus: OpcStatus = 'pending'
    try {
      const entity = await getOpcEntity()
      if (entity?.status === 'active') {
        opcStatus = 'active'
      } else if (entity?.status && entity.status !== 'unsigned') {
        opcStatus = entity.status as OpcStatus
      }
    } catch {
      /* ignore */
    }
    return { hasActiveOrder: true, opcStatus, hasPendingOrder: false }
  } catch {
    return { hasActiveOrder: false, opcStatus: 'none', hasPendingOrder: false }
  }
}



/** 收入台账列表 */

export async function getIncomeList(params: IncomeListParams) {

  const raw = await memberRequest.get<Record<string, any>>({

    url: '/compliance/income',

    params

  })

  const summary = raw?.summary

  return {

    list: (raw?.list || []) as IncomeEntry[],

    total: Number(raw?.total) || 0,

    page: Number(raw?.page) || 1,

    pageSize: Number(raw?.pageSize) || 20,

    summary: summary

      ? {

          grossTotal: Number(summary.totalGross ?? summary.grossTotal) || 0,

          netTotal: Number(summary.totalNet ?? summary.netTotal) || 0

        }

      : undefined

  } satisfies IncomeListResult

}



/** 新增收入 */

export function createIncome(data: IncomeCreateParams) {

  return memberRequest.post<IncomeEntry>({

    url: '/compliance/income',

    data

  })

}



/** 删除收入 */

export function deleteIncome(id: number | string) {

  return memberRequest.request<void>({

    url: '/compliance/income',

    method: 'DELETE',

    params: { id }

  })

}



/** CSV 导入预览 */

export function previewIncomeImport(file: File) {

  const formData = new FormData()

  formData.append('file', file)

  return memberRequest.post<IncomeImportPreviewResult>({

    url: '/compliance/income/import/preview',

    data: formData

  })

}



/** 确认 CSV 导入 */

export function confirmIncomeImport(file: File) {

  const formData = new FormData()

  formData.append('file', file)

  return memberRequest.post<{ imported: number }>({

    url: '/compliance/income/import',

    data: formData

  })

}

/** 平台流水 OCR 预览 */

export function previewIncomeOCR(params: { file?: File; platform?: string; ocrText?: string; attachmentId?: number }) {

  const formData = new FormData()

  if (params.file) formData.append('file', params.file)

  if (params.platform) formData.append('platform', params.platform)

  if (params.ocrText) formData.append('ocrText', params.ocrText)

  if (params.attachmentId) formData.append('attachmentId', String(params.attachmentId))

  return memberRequest.post<IncomeOCRPreviewResult>({

    url: '/compliance/income/ocr/preview',

    data: formData

  })

}

/** 平台流水 OCR 确认导入 */

export function confirmIncomeOCR(params: { file?: File; platform?: string; ocrText?: string; attachmentId?: number }) {

  const formData = new FormData()

  if (params.file) formData.append('file', params.file)

  if (params.platform) formData.append('platform', params.platform)

  if (params.ocrText) formData.append('ocrText', params.ocrText)

  if (params.attachmentId) formData.append('attachmentId', String(params.attachmentId))

  return memberRequest.post<{ imported: number; failed?: number }>({

    url: '/compliance/income/ocr/import',

    data: formData

  })

}

/** 收入一致性比对 */

export function getIncomeConsistency(params?: { month?: string }) {

  return memberRequest.get<IncomeConsistencyResult>({

    url: '/compliance/income/consistency',

    params

  })

}

/** 银行流水 CSV 导入 */

export function importBankStatement(file: File) {

  const formData = new FormData()

  formData.append('file', file)

  return memberRequest.post<{ imported: number; matched: number; unmatched: number }>({

    url: '/compliance/bank/import',

    data: formData

  })

}



/** 银行流水未匹配列表 */

export function getBankUnmatched(params?: { month?: string }) {

  return memberRequest.get<{ list: BankUnmatchedItem[]; count: number }>({

    url: '/compliance/bank/unmatched',

    params

  })

}



/** 费用类型库 */

export function getExpenseCategories() {

  return memberRequest.get<{ list: ExpenseCategory[] }>({

    url: '/compliance/expense/categories'

  })

}



/** 费用台账列表 */

export function getExpenseList(params: ExpenseListParams) {

  return memberRequest.get<ExpenseListResult>({

    url: '/compliance/expense',

    params

  })

}



/** 新增费用 */

export function createExpense(data: ExpenseCreateParams) {

  return memberRequest.post<ExpenseEntry>({

    url: '/compliance/expense',

    data

  })

}



/** 删除费用 */

export function deleteExpense(id: number | string) {

  return memberRequest.request<void>({

    url: '/compliance/expense',

    method: 'DELETE',

    params: { id }

  })

}



/** 利润表摘要 */

export async function getProfitSummary(params: { period: string; periodType: ProfitPeriodType }) {

  const raw = await memberRequest.get<Record<string, any>>({

    url: '/compliance/ledger/profit',

    params,

    showErrorMessage: false

  })

  const items = Array.isArray(raw?.items) ? raw.items : []

  return {

    period: raw?.period ?? '',

    periodType: params.periodType,

    revenue: Number(raw?.revenue) || 0,

    cost: Number(raw?.cost) || 0,

    profit: Number(raw?.profit) || 0,

    cumulativeProfit: Number(raw?.cumulativeProfit ?? items[0]?.cumulativeProfit) || 0

  } satisfies ProfitSummary

}



/** 会计分录列表 */

export function getLedgerVouchers(params: { period: string }) {

  return memberRequest.get<{ list: LedgerVoucher[] }>({

    url: '/compliance/ledger/vouchers',

    params

  })

}



/** 申报日历与任务 */

export function getTaxCalendar(params: { year: number; month: number }) {

  return memberRequest.get<TaxCalendarResult>({

    url: '/compliance/tax/calendar',

    params

  })

}



/** 报税前自查清单 */

export function getTaxChecklist(params?: { period?: string }) {

  return memberRequest.get<{ items: TaxChecklistItem[] }>({

    url: '/compliance/tax/checklist',

    params

  })

}



/** 申报任务详情 */

export function getTaxTaskDetail(id: number | string) {

  return memberRequest.get<TaxTask>({

    url: '/compliance/tax/tasks/detail',

    params: { id }

  })

}



/** 月度对账单列表 */

export function getStatementList(params?: { page?: number; pageSize?: number }) {

  return memberRequest.get<StatementListResult>({

    url: '/compliance/statements',

    params

  })

}



/** 对账单 PDF 下载地址 */

export function getStatementPdfUrl(id: number | string) {

  return memberRequest.get<{ url: string }>({

    url: `/compliance/statements/${id}/pdf`

  })

}



/** 平台中文 */

export function platformLabel(code: string) {

  return INCOME_PLATFORMS.find(p => p.value === code)?.label || code

}



/** 收入类型中文 */

export function incomeCategoryLabel(code: string) {

  return INCOME_CATEGORIES.find(c => c.value === code)?.label || code

}



/** 金额格式化 */

export function formatMoney(amount: number | undefined | null) {

  if (amount == null) return '—'

  return amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })

}



/** 报税 Excel 申报汇总预览 */

export interface TaxFilingSummaryPreview {

  period: string

  revenue: number

  expenseTotal: number

  profit: number

  vatAmount: number

  citAmount: number

  surchargeAmount: number

  remark?: string

  valid: boolean

  error?: string

}



/** 报税 Excel 报销明细预览行 */

export interface TaxFilingExpensePreviewRow {

  row: number

  occurredAt: string

  category: string

  categoryName?: string

  amount: number

  invoiceType: string

  description?: string

  valid: boolean

  error?: string

}



export interface TaxFilingImportPreviewResult {

  summary: TaxFilingSummaryPreview

  expenseRows: TaxFilingExpensePreviewRow[]

  validExpenseCount: number

  invalidExpenseCount: number

}



export interface TaxFilingImportResult {

  submissionId: number

  period: string

  expenseImported: number

  expenseFailed: number

}



export interface TaxFilingSubmissionItem {

  id: number

  period: string

  revenue: number

  expenseTotal: number

  profit: number

  vatAmount: number

  citAmount: number

  surchargeAmount: number

  expenseImported: number

  remark?: string

  createdAt: string

}



/** 下载报税 Excel 模板 */

export async function downloadTaxFilingTemplate() {

  const { useMemberStore } = await import('@/store/modules/member')

  const token = useMemberStore().getToken()

  const base = (import.meta.env.VITE_API_URL as string) || ''

  const res = await fetch(`${base}/member/compliance/tax/filing/template`, {

    headers: { 'Xy-User-Token': token },

  })

  if (!res.ok) throw new Error('模板下载失败')

  const blob = await res.blob()

  const url = URL.createObjectURL(blob)

  const a = document.createElement('a')

  a.href = url

  a.download = 'tax-filing-import.xlsx'

  a.click()

  URL.revokeObjectURL(url)

}



/** 报税 Excel 导入预览 */

export function previewTaxFilingImport(file: File) {

  const formData = new FormData()

  formData.append('file', file)

  return memberRequest.post<TaxFilingImportPreviewResult>({

    url: '/compliance/tax/filing/import/preview',

    data: formData,

  })

}



/** 确认报税 Excel 导入 */

export function confirmTaxFilingImport(file: File, excelFileId?: number) {

  const formData = new FormData()

  formData.append('file', file)

  if (excelFileId) formData.append('excelFileId', String(excelFileId))

  return memberRequest.post<TaxFilingImportResult>({

    url: '/compliance/tax/filing/import',

    data: formData,

  })

}



/** 报税提交历史 */

export function getTaxFilingSubmissions(params?: { page?: number; pageSize?: number }) {

  return memberRequest.get<{ list: TaxFilingSubmissionItem[]; total: number; page: number; pageSize: number }>({

    url: '/compliance/tax/filing/submissions',

    params,

  })

}


