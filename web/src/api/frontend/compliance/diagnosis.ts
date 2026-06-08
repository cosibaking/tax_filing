/**
 * 前台合规诊断公开 API（无需登录）
 * @module api/frontend/compliance/diagnosis
 */
import { siteRequest } from '@/utils/http'

/** 月收入区间 */
export type MonthlyIncomeRange = '0-2万' | '2-5万' | '5-15万' | '15万+'

/** 现有主体类型 */
export type ExistingEntity = 'none' | 'individual' | 'company' | 'other'

/** 报税状态 */
export type HasFiledTax = 'yes' | 'no' | 'unsure'

/** 推荐方案 */
export type RecommendedPlan = 'opc' | 'individual' | 'labor' | 'transitional' | 'none'

/** 单方案税负对比项 */
export interface TaxComparisonItem {
  /** 方案标识：none / labor / individual / opc */
  plan: string
  /** 方案名称 */
  label: string
  /** 预估年税负（元） */
  annualTax: number
  /** 综合税负率 */
  taxRate: number
  /** 是否为推荐方案 */
  recommended?: boolean
  /** 是否为高危警告（不报税列） */
  warning?: boolean
  /** 补充说明 */
  note?: string
}

/** 税负对比结果 */
export interface TaxComparison {
  annualIncome: number
  annualCost: number
  items: TaxComparisonItem[]
}

/** 诊断问卷提交参数 */
export interface DiagnosisSubmitParams {
  platforms: string[]
  monthlyIncomeRange: MonthlyIncomeRange
  annualCostEstimate?: number
  existingEntity: ExistingEntity
  hasFiledTax: HasFiledTax
  taxBureauContact: boolean
  notes?: string
  costBreakdown?: Record<string, number>
}

/** 诊断提交响应 */
export interface DiagnosisSubmitResult {
  id: number | string
  recommendedPlan: RecommendedPlan
  taxComparison: TaxComparison
  reasons?: string[]
  assumptionHints?: string[]
}

/** 税负计算器参数 */
export interface TaxCalculatorParams {
  annualIncome: number
  annualCost: number
  diagnosisId?: number | string
}

/** 税负计算器响应 */
export interface TaxCalculatorResult {
  taxComparison: TaxComparison
  recommendedPlan?: RecommendedPlan
  reasons?: string[]
  assumptionHints?: string[]
}

/** 服务套餐 */
export interface ServicePlan {
  id: number | string
  name: string
  tier: 'basic' | 'advanced' | 'premium'
  monthlyPrice: number | null
  priceLabel?: string
  features: string[]
  recommended?: boolean
}

/** 套餐列表响应 */
export interface ServicePlansResult {
  list: ServicePlan[]
}

const DIAGNOSIS_RESULT_KEY = 'compliance_diagnosis_result'

/** 缓存诊断结果（供结果页刷新后读取） */
export function cacheDiagnosisResult(id: number | string, result: DiagnosisSubmitResult) {
  sessionStorage.setItem(`${DIAGNOSIS_RESULT_KEY}_${id}`, JSON.stringify(result))
}

/** 读取缓存的诊断结果 */
export function getCachedDiagnosisResult(id: number | string): DiagnosisSubmitResult | null {
  const raw = sessionStorage.getItem(`${DIAGNOSIS_RESULT_KEY}_${id}`)
  if (!raw) return null
  try {
    return JSON.parse(raw) as DiagnosisSubmitResult
  } catch {
    return null
  }
}

/** 提交合规诊断问卷 */
export function submitDiagnosis(params: DiagnosisSubmitParams) {
  return siteRequest.post<DiagnosisSubmitResult>({
    url: '/compliance/diagnosis',
    data: params
  })
}

/** 税负计算器（可调年收入/成本） */
export function calculateTax(params: TaxCalculatorParams) {
  return siteRequest.post<TaxCalculatorResult>({
    url: '/compliance/calculator',
    data: params
  })
}

/** 获取服务套餐列表 */
export async function fetchServicePlans(): Promise<ServicePlansResult> {
  const res = await siteRequest.get<ServicePlansResult | ServicePlan[]>({
    url: '/compliance/plans'
  })
  if (Array.isArray(res)) return { list: res }
  return { list: (res as ServicePlansResult)?.list ?? [] }
}
