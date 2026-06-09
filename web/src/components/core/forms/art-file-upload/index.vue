<!-- +----------------------------------------------------------------------
  | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
  +----------------------------------------------------------------------
  | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
  +----------------------------------------------------------------------
  | Licensed ( https://opensource.org/licenses/MIT )
  +----------------------------------------------------------------------
  | Author: 喜羊羊 <751300685@qq.com>
  +---------------------------------------------------------------------- -->
<!-- 文件上传组件 -->
<template>
  <div class="art-file-upload" :class="{ 'is-compact': compact }">
    <!-- 单文件上传 -->
    <div v-if="!multiple" class="single-file-upload">
      <div v-if="fileInfo" class="file-preview">
        <div class="file-icon">
          <ArtSvgIcon :icon="getFileIcon(fileInfo.name)" :size="24" />
        </div>
        <div class="file-info">
          <div class="file-name">{{ fileInfo.name }}</div>
          <div class="file-size">{{ formatFileSize(fileInfo.size) }}</div>
        </div>
        <div class="file-actions">
          <ElButton size="small" text @click="handleDownload">
            <ArtSvgIcon icon="ri:download-line" :size="16" />
          </ElButton>
          <ElButton size="small" text type="danger" @click="handleRemove">
            <ArtSvgIcon icon="ri:delete-bin-line" :size="16" />
          </ElButton>
        </div>
      </div>
      <ElUpload
        v-else
        :http-request="customUpload"
        :show-file-list="false"
        :accept="accept"
        :before-upload="beforeUpload"
      >
        <ElButton type="primary" plain>
          <ArtSvgIcon icon="ri:upload-line" :size="16" style="margin-right: 6px;" />
          点击上传文件
        </ElButton>
        <template #tip>
          <div class="upload-tip">{{ acceptHint }}</div>
        </template>
      </ElUpload>
    </div>

    <!-- 多文件上传 -->
    <div v-else class="multiple-file-upload">
      <div v-if="fileList.length > 0" class="file-list">
        <div v-for="(file, index) in fileList" :key="file.id" class="file-item">
          <div class="file-icon">
            <ArtSvgIcon :icon="getFileIcon(file.name)" :size="20" />
          </div>
          <div class="file-info">
            <div class="file-name">{{ file.name }}</div>
            <div class="file-size">{{ formatFileSize(file.size) }}</div>
          </div>
          <div class="file-actions">
            <ElButton size="small" text @click="handleDownloadMultiple(file.url)">
              <ArtSvgIcon icon="ri:download-line" :size="14" />
            </ElButton>
            <ElButton size="small" text type="danger" @click="handleRemoveMultiple(index)">
              <ArtSvgIcon icon="ri:delete-bin-line" :size="14" />
            </ElButton>
          </div>
        </div>
      </div>

      <ElUpload
        v-if="!limit || fileList.length < limit"
        :http-request="customUploadMultiple"
        :show-file-list="false"
        :accept="accept"
        :before-upload="beforeUpload"
      >
        <ElButton type="primary" plain size="small">
          <ArtSvgIcon icon="ri:add-line" :size="16" style="margin-right: 4px;" />
          添加文件
        </ElButton>
        <template #tip>
          <div class="upload-tip">{{ acceptHint }}</div>
        </template>
      </ElUpload>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { UploadProps, UploadRequestOptions } from 'element-plus'
import { uploadFileApi } from '@/api/backend/common/upload'

defineOptions({ name: 'ArtFileUpload' })

interface Props {
  modelValue?: string | string[]
  multiple?: boolean
  limit?: number
  accept?: string
  maxSize?: number // MB
  compact?: boolean // 紧凑模式（用于数组编辑器等场景）
  /** url=存访问地址（默认）；id=存附件表 attachmentId，供合规等需 fileId 的场景 */
  valueType?: 'url' | 'id'
}

const props = withDefaults(defineProps<Props>(), {
  multiple: false,
  limit: 10,
  accept: '*',
  maxSize: 10,
  compact: false,
  valueType: 'url'
})

const emit = defineEmits<{
  'update:modelValue': [value: string | string[]]
  'change': [value: string | string[]]
}>()

interface FileItem {
  id: number
  url: string
  name: string
  size: number
}

// 单文件模式
const fileInfo = ref<FileItem | null>(null)

const emitValue = (value: string) => {
  emit('update:modelValue', value)
  emit('change', value)
}

// 解析单文件展示信息
const parseSingleFile = (val: string) => {
  if (!val) return null
  if (props.valueType === 'id') {
    return {
      id: Number(val) || Date.now(),
      url: '',
      name: `附件 #${val}`,
      size: 0
    }
  }
  const name = val.split('/').pop() || 'file'
  return {
    id: Date.now(),
    url: val,
    name,
    size: 0
  }
}

// 监听单文件值变化
watch(() => props.modelValue, (val) => {
  if (!props.multiple && typeof val === 'string') {
    fileInfo.value = val ? parseSingleFile(val) : null
  }
}, { immediate: true })

// 多文件模式
const fileList = computed({
  get: () => {
    if (!props.multiple) return []
    const val = props.modelValue
    if (Array.isArray(val)) {
      return val.map((url, index) => ({
        id: index,
        url,
        name: url.split('/').pop() || `file-${index}`,
        size: 0
      }))
    }
    if (typeof val === 'string' && val) {
      return val.split(',').map((url, index) => ({
        id: index,
        url: url.trim(),
        name: url.split('/').pop() || `file-${index}`,
        size: 0
      }))
    }
    return []
  },
  set: (val) => {
    const urls = val.map(item => item.url)
    const value = urls.join(',')
    emit('update:modelValue', value)
    emit('change', value)
  }
})

// 文件类型提示
const acceptHint = computed(() => {
  if (props.accept === '*') {
    return `支持所有文件类型，不超过 ${props.maxSize}MB`
  }
  const types = props.accept.split(',').map(t => {
    const ext = t.trim().replace('.', '').toUpperCase()
    return ext
  })
  return `支持 ${types.join('、')} 格式，不超过 ${props.maxSize}MB`
})

// 根据文件名获取图标
const getFileIcon = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()
  const iconMap: Record<string, string> = {
    // 文档
    'pdf': 'ri:file-pdf-line',
    'doc': 'ri:file-word-line',
    'docx': 'ri:file-word-line',
    'xls': 'ri:file-excel-line',
    'xlsx': 'ri:file-excel-line',
    'ppt': 'ri:file-ppt-line',
    'pptx': 'ri:file-ppt-line',
    // 压缩包
    'zip': 'ri:file-zip-line',
    'rar': 'ri:file-zip-line',
    '7z': 'ri:file-zip-line',
    'tar': 'ri:file-zip-line',
    'gz': 'ri:file-zip-line',
    // 代码
    'js': 'ri:file-code-line',
    'ts': 'ri:file-code-line',
    'json': 'ri:file-code-line',
    'html': 'ri:file-code-line',
    'css': 'ri:file-code-line',
    'vue': 'ri:file-code-line',
    'go': 'ri:file-code-line',
    'php': 'ri:file-code-line',
    'py': 'ri:file-code-line',
    // 文本
    'txt': 'ri:file-text-line',
    'md': 'ri:file-text-line',
    // 视频
    'mp4': 'ri:file-video-line',
    'avi': 'ri:file-video-line',
    'mov': 'ri:file-video-line',
    // 音频
    'mp3': 'ri:file-music-line',
    'wav': 'ri:file-music-line',
  }
  return iconMap[ext || ''] || 'ri:file-line'
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '未知大小'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

// 上传前校验
const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isValidSize = file.size / 1024 / 1024 < props.maxSize
  if (!isValidSize) {
    ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
    return false
  }
  return true
}

// 自定义上传 - 单文件
const customUpload = async (options: UploadRequestOptions) => {
  try {
    const file = options.file as File
    // request.post 成功时直接返回 data.data，失败时抛出异常
    const data = await uploadFileApi(file)
    
    if (data?.url) {
      const stored =
        props.valueType === 'id' && data.attachmentId
          ? String(data.attachmentId)
          : data.url
      fileInfo.value = {
        id: data.attachmentId || Date.now(),
        url: data.url,
        name: file.name,
        size: file.size
      }
      emitValue(stored)
      ElMessage.success('上传成功')
      options.onSuccess(data)
    } else {
      throw new Error('上传响应数据格式错误')
    }
  } catch (error: any) {
    console.error('上传失败:', error)
    // 错误已经在拦截器中显示了，这里不再重复显示
    options.onError(error as any)
  }
}

// 自定义上传 - 多文件
const customUploadMultiple = async (options: UploadRequestOptions) => {
  try {
    const file = options.file as File
    // request.post 成功时直接返回 data.data，失败时抛出异常
    const data = await uploadFileApi(file)
    
    if (data?.url) {
      const newFile: FileItem = {
        id: Date.now(),
        url: data.url,
        name: file.name,
        size: file.size
      }
      const newList = [...fileList.value, newFile]
      fileList.value = newList
      ElMessage.success('上传成功')
      options.onSuccess(data)
    } else {
      throw new Error('上传响应数据格式错误')
    }
  } catch (error: any) {
    console.error('上传失败:', error)
    // 错误已经在拦截器中显示了，这里不再重复显示
    options.onError(error as any)
  }
}

// 下载文件
const handleDownload = () => {
  if (fileInfo.value?.url) {
    window.open(fileInfo.value.url, '_blank')
  }
}

// 删除文件
const handleRemove = () => {
  fileInfo.value = null
  emitValue('')
}

// 多文件下载
const handleDownloadMultiple = (url: string) => {
  window.open(url, '_blank')
}

// 多文件删除
const handleRemoveMultiple = (index: number) => {
  const newList = fileList.value.filter((_, i) => i !== index)
  fileList.value = newList
}
</script>

<style scoped lang="scss">
.art-file-upload {
  width: 100%;
}

// 单文件上传
.single-file-upload {
  width: 100%;
  max-width: 500px;
}

.file-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border: 1px solid var(--art-gray-300);
  border-radius: 8px;
  background: var(--art-gray-50);
  transition: all 0.2s;

  &:hover {
    border-color: var(--theme-color);
    background: var(--theme-color-alpha-5);
  }
}

.file-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--default-box-color);
  border-radius: 6px;
  color: var(--theme-color);
}

.file-info {
  flex: 1;
  min-width: 0;

  .file-name {
    font-size: 14px;
    font-weight: 500;
    color: var(--art-gray-800);
    margin-bottom: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .file-size {
    font-size: 12px;
    color: var(--art-gray-500);
  }
}

.file-actions {
  flex-shrink: 0;
  display: flex;
  gap: 4px;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: var(--art-gray-500);
}

// 多文件上传
.multiple-file-upload {
  width: 100%;
}

.file-list {
  margin-bottom: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border: 1px solid var(--art-gray-300);
  border-radius: 6px;
  background: var(--default-box-color);
  transition: all 0.2s;

  &:hover {
    border-color: var(--theme-color);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  }

  .file-icon {
    width: 36px;
    height: 36px;
    font-size: 20px;
  }

  .file-info {
    .file-name {
      font-size: 13px;
    }
  }
}

// 紧凑模式：适用于数组编辑器等场景
.art-file-upload.is-compact {
  .single-file-upload {
    max-width: 100% !important;
  }

  .file-preview,
  .file-item {
    height: 32px !important;
    min-height: 32px !important;
    max-height: 32px !important;
    padding: 0 10px !important;
    gap: 8px !important;
    border-radius: 4px !important;
    margin: 0 !important;
  }

  .file-icon {
    width: 22px !important;
    height: 22px !important;
    min-width: 22px !important;
    max-width: 22px !important;
    min-height: 22px !important;
    max-height: 22px !important;

    .art-svg-icon {
      font-size: 14px !important;
    }
  }

  .file-info {
    gap: 6px !important;

    .file-name {
      font-size: 12px !important;
      line-height: 1 !important;
      margin: 0 !important;
    }

    .file-size {
      font-size: 11px !important;
      line-height: 1 !important;
      margin: 0 !important;
    }
  }

  .file-actions {
    gap: 2px !important;
    height: 32px !important;

    .el-button {
      padding: 4px !important;
      height: auto !important;

      .art-svg-icon {
        font-size: 14px !important;
      }
    }
  }

  .el-upload {
    .el-button {
      height: 32px !important;
      padding: 0 10px !important;
      font-size: 12px !important;

      .art-svg-icon {
        font-size: 14px !important;
        margin-right: 4px !important;
      }
    }
  }

  .upload-tip {
    display: none !important;
  }

  .file-list {
    gap: 4px !important;
  }
}

// 暗黑模式
.dark {
  .file-preview {
    background: var(--art-gray-800);
  }

  .file-icon {
    background: var(--art-gray-700);
  }

  .file-item {
    background: var(--art-gray-800);
  }
}
</style>
