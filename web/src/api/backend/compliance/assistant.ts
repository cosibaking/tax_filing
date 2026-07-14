import { adminRequest } from '@/utils/http'
export type WorkspaceSection = 'rules' | 'tasks' | 'documents' | 'risks' | 'reports' | 'tickets'
export function getAssistantWorkspace(params: {
  section: WorkspaceSection
  status?: string
  page?: number
  pageSize?: number
}) {
  return adminRequest.get<{ list: Record<string, unknown>[]; total: number }>({
    url: '/admin/compliance/assistant/workspace',
    params
  })
}
export function transitionRuleVersion(id: number, action: 'submit_review' | 'publish' | 'retire') {
  return adminRequest.post({
    url: `/admin/compliance/rule-versions/${id}/transition`,
    data: { action }
  })
}
