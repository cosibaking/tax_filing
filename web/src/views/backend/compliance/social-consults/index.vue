<!-- 社保咨询工单 -->
<template>
  <div class="social-consults-page art-full-height">
    <ElCard shadow="never" class="mb-4">
      <ElForm :inline="true" @submit.prevent="handleSearch">
        <ElFormItem label="状态">
          <ElSelect v-model="searchForm.status" clearable placeholder="全部" style="width: 120px">
            <ElOption label="待回复" value="open" />
            <ElOption label="已回复" value="replied" />
            <ElOption label="已关闭" value="closed" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="关键词">
          <ElInput v-model="searchForm.q" clearable placeholder="问题/客户/公司" style="width: 200px" />
        </ElFormItem>
        <ElFormItem>
          <ElButton type="primary" @click="handleSearch">查询</ElButton>
          <ElButton @click="resetSearch">重置</ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <ElCard shadow="never" class="art-table-card">
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
        <template #memberName="{ row }">
          <CellEllipsis :text="row.memberName || `会员#${row.memberId}`" />
        </template>
        <template #companyName="{ row }">
          <CellEllipsis :text="row.companyName || '—'" />
        </template>
        <template #category="{ row }">
          {{ categoryLabel(row.category) }}
        </template>
        <template #regionCode="{ row }">
          <CellEllipsis :text="row.regionCode || '—'" />
        </template>
        <template #question="{ row }">
          <CellEllipsis :text="row.question || '—'" :lines="2" />
        </template>
        <template #status="{ row }">
          <ElTag :type="statusTagType(row.status)" size="small">
            {{ statusLabel(row.status) }}
          </ElTag>
        </template>
        <template #createdAt="{ row }">
          {{ row.createdAt || '—' }}
        </template>
        <template #action="{ row }">
          <div class="flex gap-1">
            <ElButton v-if="row.status !== 'closed'" type="primary" link @click="openReply(row)">
              回复
            </ElButton>
            <ElButton v-if="row.status !== 'closed'" type="info" link @click="closeRow(row)">
              关闭
            </ElButton>
          </div>
        </template>
      </ArtTable>
    </ElCard>

    <ElDialog v-model="replyVisible" title="回复咨询" width="640px" destroy-on-close>
      <div v-if="currentRow" class="reply-dialog">
        <dl class="reply-dialog__meta">
          <div class="reply-dialog__meta-item">
            <dt>客户</dt>
            <dd>{{ currentRow.memberName || `会员#${currentRow.memberId}` }}</dd>
          </div>
          <div class="reply-dialog__meta-item">
            <dt>OPC 主体</dt>
            <dd>{{ currentRow.companyName || '—' }}</dd>
          </div>
          <div class="reply-dialog__meta-item">
            <dt>咨询类型</dt>
            <dd>{{ categoryLabel(currentRow.category) }}</dd>
          </div>
          <div class="reply-dialog__meta-item">
            <dt>注册地</dt>
            <dd>{{ currentRow.regionCode || '—' }}</dd>
          </div>
          <div class="reply-dialog__meta-item">
            <dt>提交时间</dt>
            <dd>{{ currentRow.createdAt || '—' }}</dd>
          </div>
          <div class="reply-dialog__meta-item">
            <dt>当前状态</dt>
            <dd>
              <ElTag :type="statusTagType(currentRow.status)" size="small">
                {{ statusLabel(currentRow.status) }}
              </ElTag>
            </dd>
          </div>
        </dl>

        <section class="reply-dialog__section">
          <h4 class="reply-dialog__section-title">咨询问题（全文）</h4>
          <div class="reply-dialog__question">{{ currentRow.question || '—' }}</div>
        </section>

        <section v-if="currentRow.reply" class="reply-dialog__section">
          <h4 class="reply-dialog__section-title">历史回复</h4>
          <div class="reply-dialog__history">
            <p class="reply-dialog__history-text">{{ currentRow.reply }}</p>
            <p v-if="currentRow.repliedAt" class="reply-dialog__history-time">{{ currentRow.repliedAt }}</p>
          </div>
        </section>

        <section class="reply-dialog__section">
          <h4 class="reply-dialog__section-title">顾问回复</h4>
          <ElInput v-model="replyText" type="textarea" :rows="5" placeholder="填写回复内容，至少 2 个字" />
        </section>
      </div>
      <template #footer>
        <ElButton @click="replyVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="replying" @click="submitReply">发送回复</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
import { useTable } from '@/hooks/core/useTable'
import {
  getSocialConsultList,
  replySocialConsult,
  ADMIN_SOCIAL_CONSULT_STATUS_LABELS,
  ADMIN_SOCIAL_CONSULT_CATEGORY_LABELS,
  type AdminSocialConsultItem,
  type AdminSocialConsultStatus
} from '@/api/backend/compliance/social'
import { useComplianceMenuBadges } from '@/hooks/compliance/useComplianceMenuBadges'
import CellEllipsis from './modules/cell-ellipsis.vue'
import { ElButton, ElMessage, ElTag } from 'element-plus'

defineOptions({ name: 'ComplianceSocialConsults' })

const searchForm = ref({
  status: '' as '' | AdminSocialConsultStatus,
  q: ''
})

const replyVisible = ref(false)
const replyText = ref('')
const replying = ref(false)
const currentRow = ref<AdminSocialConsultItem | null>(null)
const { refreshComplianceMenuBadges } = useComplianceMenuBadges()

function categoryLabel(category?: string) {
  if (!category) return '—'
  return ADMIN_SOCIAL_CONSULT_CATEGORY_LABELS[category] || category
}

function statusLabel(status?: AdminSocialConsultStatus | string) {
  if (!status) return '—'
  return ADMIN_SOCIAL_CONSULT_STATUS_LABELS[status as AdminSocialConsultStatus] || status
}

function statusTagType(status?: AdminSocialConsultStatus | string) {
  if (status === 'open') return 'warning'
  if (status === 'replied') return 'success'
  return 'info'
}

const {
  columns,
  data,
  loading,
  pagination,
  getData,
  searchParams,
  handleSizeChange,
  handleCurrentChange,
  refreshData
} = useTable({
  core: {
    apiFn: getSocialConsultList,
    apiParams: { page: 1, pageSize: 20 },
    paginationKey: {
      current: 'page',
      size: 'pageSize'
    },
    columnsFactory: () => [
      { prop: 'memberName', label: '客户', minWidth: 100, useSlot: true },
      { prop: 'companyName', label: 'OPC', minWidth: 140, useSlot: true },
      { prop: 'category', label: '类型', width: 110, useSlot: true },
      { prop: 'regionCode', label: '注册地', minWidth: 100, useSlot: true },
      { prop: 'question', label: '问题', minWidth: 220, useSlot: true },
      { prop: 'status', label: '状态', width: 100, useSlot: true },
      { prop: 'createdAt', label: '提交时间', width: 170, useSlot: true },
      { prop: 'action', label: '操作', width: 140, fixed: 'right', useSlot: true }
    ]
  }
})

function handleSearch() {
  Object.assign(searchParams, {
    status: searchForm.value.status || undefined,
    q: searchForm.value.q || undefined,
    page: 1
  })
  getData()
}

function resetSearch() {
  searchForm.value = { status: '', q: '' }
  Object.assign(searchParams, { status: undefined, q: undefined, page: 1 })
  getData()
}

function openReply(row: AdminSocialConsultItem) {
  currentRow.value = row
  replyText.value = row.reply || ''
  replyVisible.value = true
}

async function submitReply() {
  if (!currentRow.value || replyText.value.trim().length < 2) {
    ElMessage.warning('请填写回复内容')
    return
  }
  replying.value = true
  try {
    await replySocialConsult(currentRow.value.id, { action: 'reply', reply: replyText.value.trim() })
    ElMessage.success('已回复')
    replyVisible.value = false
    refreshData()
    refreshComplianceMenuBadges()
  } catch (e: any) {
    ElMessage.error(e?.message || '回复失败')
  } finally {
    replying.value = false
  }
}

async function closeRow(row: AdminSocialConsultItem) {
  try {
    await replySocialConsult(row.id, { action: 'close' })
    ElMessage.success('已关闭')
    refreshData()
    refreshComplianceMenuBadges()
  } catch (e: any) {
    ElMessage.error(e?.message || '操作失败')
  }
}
</script>

<style lang="scss" scoped>
.reply-dialog {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.reply-dialog__meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 20px;
  margin: 0;
  padding: 14px 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.reply-dialog__meta-item {
  min-width: 0;

  dt {
    margin: 0 0 4px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }

  dd {
    margin: 0;
    font-size: 14px;
    color: var(--el-text-color-primary);
    word-break: break-word;
  }
}

.reply-dialog__section-title {
  margin: 0 0 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.reply-dialog__question,
.reply-dialog__history {
  max-height: 220px;
  overflow: auto;
  padding: 12px 14px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  background: #fff;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
  font-size: 14px;
  color: var(--el-text-color-primary);
}

.reply-dialog__history {
  background: var(--el-fill-color-blank);
}

.reply-dialog__history-text {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.reply-dialog__history-time {
  margin: 8px 0 0;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
</style>
