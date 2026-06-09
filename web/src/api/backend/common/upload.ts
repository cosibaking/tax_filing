// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

/**
 * 文件上传 API
 * @module api/backend/common/upload
 */
import { adminRequest } from '@/utils/http'

/**
 * 上传响应数据类型
 */
export interface UploadResponse {
  url: string
  path?: string
  name?: string
  size?: number
  attachmentId?: number
}

/**
 * 上传图片
 * @param file 图片文件
 */
export const uploadImageApi = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  
  return adminRequest.post<UploadResponse>({
    url: '/upload/file',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    showErrorMessage: true
  })
}

/**
 * 上传文件
 * @param file 文件
 */
export const uploadFileApi = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  
  return adminRequest.post<UploadResponse>({
    url: '/upload/file',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    showErrorMessage: true
  })
}

/**
 * 上传文件（带进度回调）
 * @param file 文件
 * @param onProgress 进度回调
 */
export const uploadFileWithProgressApi = (
  file: File, 
  onProgress?: (progress: number) => void
) => {
  const formData = new FormData()
  formData.append('file', file)
  
  return adminRequest.post<any>({
    url: '/upload/file',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: (progressEvent: any) => {
      if (progressEvent.lengthComputable && onProgress) {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / progressEvent.total
        )
        onProgress(percentCompleted)
      }
    }
  })
}
