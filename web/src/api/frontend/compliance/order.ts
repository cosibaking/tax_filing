/**
 * 前台合规签约 API（需登录）
 * @module api/frontend/compliance/order
 */
import { memberRequest } from '@/utils/http'
import type { ServicePlan } from './diagnosis'

/** 订单状态 */
export type OrderStatus = 'pending' | 'active' | 'cancelled'

/** 合规订单 */
export interface ComplianceOrder {
  orderId: number | string
  /** @deprecated 兼容旧字段 */
  id?: number | string
  planId: number | string
  planTier?: 'basic' | 'advanced' | 'premium'
  planName?: string
  diagnosisId?: number | string
  status: OrderStatus
  amount?: number
  signedAt?: string
  createdAt?: string
}

/** 创建订单参数 */
export interface CreateOrderParams {
  planId?: number | string
  planTier?: string
  diagnosisId?: number | string
}

/** 风险告知四条确认（对齐后端 RiskAcknowledgments） */
export interface RiskAcknowledgments {
  complianceNotEvasion: boolean
  noAuditGuarantee: boolean
  noFakeInvoice: boolean
  taxDependsOnReality: boolean
}

/** 风险告知 / 方案确认参数 */
export interface ConsentParams {
  orderId: number | string
  type: 'risk_disclosure' | 'plan_confirm'
  acknowledgments?: RiskAcknowledgments
  planConfirmed?: boolean
  documentVersion?: string
}

/** 电子签约参数 */
export interface SignParams {
  orderId: number | string
  legalName: string
}

/** 签约结果 */
export interface SignResult {
  orderId: number | string
  opcId?: number | string
  signedAt: string
}

/** 创建合规订单 */
export function createComplianceOrder(params: CreateOrderParams) {
  return memberRequest.post<ComplianceOrder>({
    url: '/compliance/order',
    data: params
  })
}

/** 提交风险告知 / 方案确认 */
export function submitConsent(params: ConsentParams) {
  return memberRequest.post<void>({
    url: '/compliance/consent',
    data: params
  })
}

/** 完成电子签约 */
export function signAgreement(params: SignParams) {
  return memberRequest.post<SignResult>({
    url: '/compliance/sign',
    data: params
  })
}

/** 获取当前活跃订单（若无有效订单则返回 null） */
export async function getActiveOrder(): Promise<ComplianceOrder | null> {
  const res = await memberRequest.get<ComplianceOrder | null>({
    url: '/compliance/order/active'
  })
  const orderId = Number(res?.orderId ?? res?.id ?? 0)
  if (!orderId || !res?.status || res.status === 'cancelled') {
    return null
  }
  return res
}

export type { ServicePlan }
