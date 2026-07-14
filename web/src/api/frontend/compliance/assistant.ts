import { memberRequest } from '@/utils/http'

export type TaskStatus =
  | 'not_started'
  | 'preparing'
  | 'pending_confirmation'
  | 'completed'
  | 'overdue'
  | 'cancelled'

export interface EnterpriseProfileData {
  region: 'CN-BJ'
  entityType: 'one_person_limited_company'
  taxpayerType: 'small_scale'
  vatPeriod: 'monthly' | 'quarterly'
  employeeCount: number
  invoiceEnabled: boolean
  hasRevenue: boolean
  hasPublicBankAccount: boolean
  complexity: 'low'
}

export interface EnterpriseProfile {
  id: number
  opcEntityId: number
  version: number
  data: EnterpriseProfileData
  labels: Record<string, string>
  source: string
  status: string
  confirmedAt: number
}

export interface ComplianceTask {
  id: number
  memberId: number
  taskType: string
  title: string
  description: string
  materials: string[]
  status: TaskStatus
  priority: 'low' | 'medium' | 'high'
  dueAt: string
  completedAt?: string
  completionNote?: string
  key: { periodKey: string }
}

export interface TaskListQuery {
  status?: TaskStatus | ''
  taskType?: string
  periodKey?: string
  page?: number
  pageSize?: number
}

export interface BusinessDocument {
  id: number
  attachmentId: number
  periodKey: string
  documentType: string
  processStatus: string
  confidence: number
  duplicateOfId?: number
}
export const documentTypes = [
  { value: 'bank_statement', label: '银行流水' },
  { value: 'sales_invoice', label: '销项发票' },
  { value: 'expense_invoice', label: '费用凭证' },
  { value: 'contract', label: '业务合同' },
  { value: 'tax_receipt', label: '完税证明' }
]
export function createBusinessDocument(data: { attachmentId: number | string; periodKey: string }) {
  return memberRequest.post<{ document: BusinessDocument }>({ url: '/compliance/documents', data })
}
export function getBusinessDocuments(periodKey?: string) {
  return memberRequest.get<{ list: BusinessDocument[] }>({
    url: '/compliance/documents',
    params: { periodKey }
  })
}
export function confirmBusinessDocument(
  id: number,
  documentType: string,
  fields: Record<string, unknown> = {}
) {
  return memberRequest.put<{ document: BusinessDocument }>({
    url: `/compliance/documents/${id}/confirm`,
    data: { documentType, fields }
  })
}

export function getEnterpriseProfile() {
  return memberRequest.get<{ profile?: EnterpriseProfile }>({ url: '/compliance/profile' })
}

export function saveEnterpriseProfile(data: EnterpriseProfileData) {
  return memberRequest.put<{ profile: EnterpriseProfile; changed: boolean }>({
    url: '/compliance/profile',
    data: { data }
  })
}

export function getComplianceTasks(params?: TaskListQuery) {
  return memberRequest.get<{
    list: ComplianceTask[]
    total: number
    page: number
    pageSize: number
  }>({
    url: '/compliance/tasks',
    params
  })
}

export function getComplianceTask(id: number | string) {
  return memberRequest.get<{ task: ComplianceTask }>({ url: `/compliance/tasks/${id}` })
}

export function transitionComplianceTask(id: number | string, action: string, note = '') {
  return memberRequest.post<{ task: ComplianceTask }>({
    url: `/compliance/tasks/${id}/action`,
    data: { action, note }
  })
}
