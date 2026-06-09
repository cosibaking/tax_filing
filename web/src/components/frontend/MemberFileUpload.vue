<!-- 会员端文件上传（存储 attachmentId，供 OPC 合规资料提交） -->
<template>
  <div class="member-file-upload">
    <div v-if="fileInfo" class="file-preview">
      <div class="file-icon">
        <ArtSvgIcon :icon="getFileIcon(fileInfo.name)" :size="24" />
      </div>
      <div class="file-info">
        <div class="file-name">{{ fileInfo.name }}</div>
        <div class="file-size">{{ formatFileSize(fileInfo.size) }}</div>
      </div>
      <div class="file-actions">
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
</template>

<script setup lang="ts">
import type { UploadProps, UploadRequestOptions } from 'element-plus'
import { ElMessage } from 'element-plus'
import { uploadMemberFile } from '@/api/frontend/member/upload'

defineOptions({ name: 'MemberFileUpload' })

const props = withDefaults(defineProps<{
  modelValue?: number | string
  accept?: string
  maxSize?: number
}>(), {
  accept: '.jpg,.jpeg,.png,.pdf',
  maxSize: 10,
})

const emit = defineEmits<{
  'update:modelValue': [value: number | string]
  'change': [value: number | string]
}>()

interface FileItem {
  id: number | string
  name: string
  size: number
}

const fileInfo = ref<FileItem | null>(null)

watch(() => props.modelValue, (val) => {
  if (!val) {
    fileInfo.value = null
    return
  }
  if (!fileInfo.value || String(fileInfo.value.id) !== String(val)) {
    fileInfo.value = { id: val, name: `附件 #${val}`, size: 0 }
  }
}, { immediate: true })

const acceptHint = computed(() => {
  const types = props.accept.split(',').map(t => t.trim().replace('.', '').toUpperCase()).filter(Boolean)
  return `支持 ${types.join('、')} 格式，不超过 ${props.maxSize}MB`
})

const getFileIcon = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase()
  if (ext === 'pdf') return 'ri:file-pdf-line'
  if (['jpg', 'jpeg', 'png', 'gif', 'webp'].includes(ext || '')) return 'ri:image-line'
  return 'ri:file-line'
}

const formatFileSize = (bytes: number) => {
  if (!bytes) return '已上传'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  if (file.size / 1024 / 1024 >= props.maxSize) {
    ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
    return false
  }
  return true
}

const customUpload = async (options: UploadRequestOptions) => {
  try {
    const file = options.file as File
    const data = await uploadMemberFile(file)
    if (!data?.attachmentId) throw new Error('上传响应缺少 attachmentId')
    fileInfo.value = { id: data.attachmentId, name: data.name || file.name, size: data.size || file.size }
    emit('update:modelValue', data.attachmentId)
    emit('change', data.attachmentId)
    ElMessage.success('上传成功')
    options.onSuccess(data)
  } catch (error) {
    options.onError(error as Parameters<NonNullable<typeof options.onError>>[0])
  }
}

const handleRemove = () => {
  fileInfo.value = null
  emit('update:modelValue', '')
  emit('change', '')
}
</script>

<style lang="scss" scoped>
.member-file-upload {
  width: 100%;
}
.file-preview {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 16px;
  background: #f0f3f8;
  box-shadow: inset 6px 6px 12px #e0e5ec, inset -6px -6px 12px #ffffff;
}
.file-info { flex: 1; min-width: 0; }
.file-name {
  font-size: 13px;
  font-weight: 600;
  color: #32325d;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.file-size { font-size: 12px; color: #8898aa; margin-top: 2px; }
.upload-tip { font-size: 12px; color: #8898aa; margin-top: 8px; }
</style>
