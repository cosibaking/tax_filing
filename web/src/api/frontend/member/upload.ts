/**
 * 会员端文件上传 API（返回 attachmentId 供合规资料引用）
 */
import { memberRequest } from '@/utils/http'

export interface MemberUploadResult {
  url: string
  name: string
  size: number
  attachmentId: number
}

export function uploadMemberFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return memberRequest.post<MemberUploadResult>({
    url: '/user/upload',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}
