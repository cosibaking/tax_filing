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
