/**
 * 社保与用工合规 API（Phase 2）
 * @module api/frontend/compliance/social
 */
import { memberRequest } from '@/utils/http'

export type EmploymentStatus = 'unknown' | 'no_employee' | 'has_employee'

export const EMPLOYMENT_STATUS_LABELS: Record<EmploymentStatus, string> = {
  unknown: '未填写',
  no_employee: '无雇员',
  has_employee: '有雇员'
}

export interface EmploymentStatusResult {
  employmentStatus: EmploymentStatus
  employmentConfirmedAt?: string
  canChangeToNoEmployee?: boolean
}

export interface SocialGuideItem {
  slug: string
  title: string
  summary: string
}

export interface SocialGuideDetail extends SocialGuideItem {
  contentMd: string
  audience: string
}

export type SocialConsultStatus = 'open' | 'replied' | 'closed'

export interface SocialConsultItem {
  id: number | string
  category: string
  question: string
  regionCode?: string
  status: SocialConsultStatus
  reply?: string
  repliedAt?: string
  createdAt: string
}

export const SOCIAL_CONSULT_STATUS_LABELS: Record<SocialConsultStatus, string> = {
  open: '待回复',
  replied: '已回复',
  closed: '已关闭'
}

export function getEmploymentStatus() {
  return memberRequest.get<EmploymentStatusResult>({ url: '/compliance/employment' })
}

export function setEmploymentStatus(employmentStatus: 'no_employee' | 'has_employee') {
  return memberRequest.put<EmploymentStatusResult>({
    url: '/compliance/employment',
    data: { employmentStatus }
  })
}

export function getSocialGuides() {
  return memberRequest.get<{ list: SocialGuideItem[] }>({ url: '/compliance/social/guides' })
}

export function getSocialGuideDetail(slug: string) {
  return memberRequest.get<SocialGuideDetail>({ url: `/compliance/social/guides/${slug}` })
}

export function createSocialConsult(data: {
  category?: string
  question: string
  regionCode?: string
}) {
  return memberRequest.post<{ id: number | string }>({
    url: '/compliance/social/consults',
    data
  })
}

export function getSocialConsults(params?: { page?: number; pageSize?: number }) {
  return memberRequest.get<{ list: SocialConsultItem[]; total: number }>({
    url: '/compliance/social/consults',
    params
  })
}
