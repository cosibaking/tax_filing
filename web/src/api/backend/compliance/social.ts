/**
 * 后台社保咨询 API
 */
import { adminRequest } from '@/utils/http'

export type AdminSocialConsultStatus = 'open' | 'replied' | 'closed'

export interface AdminSocialConsultItem {
  id: number | string
  memberId: number
  memberName?: string
  opcId: number
  companyName?: string
  category: string
  question: string
  regionCode?: string
  status: AdminSocialConsultStatus
  reply?: string
  repliedAt?: string
  createdAt: string
}

export const ADMIN_SOCIAL_CONSULT_STATUS_LABELS: Record<AdminSocialConsultStatus, string> = {
  open: '待回复',
  replied: '已回复',
  closed: '已关闭'
}

export const ADMIN_SOCIAL_CONSULT_CATEGORY_LABELS: Record<string, string> = {
  founder: '创始人参保',
  employee: '雇员社保',
  other: '其他'
}

export function getSocialConsultList(params?: {
  page?: number
  pageSize?: number
  status?: AdminSocialConsultStatus | ''
  q?: string
}) {
  return adminRequest.get<{
    list: AdminSocialConsultItem[]
    total: number
    page: number
    pageSize: number
  }>({
    url: '/compliance/social-consults',
    params
  })
}

export function replySocialConsult(
  id: number | string,
  data: { reply?: string; action: 'reply' | 'close' }
) {
  return adminRequest.request<{ id: number | string; status: string }>({
    url: `/compliance/social-consults/${id}`,
    method: 'PATCH',
    data
  })
}
