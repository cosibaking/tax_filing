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
      />
    </ElCard>

    <ElDialog v-model="replyVisible" title="回复咨询" width="520px" destroy-on-close>
      <p class="text-sm text-gray-600 mb-4 whitespace-pre-wrap">{{ currentRow?.question }}</p>
      <ElInput v-model="replyText" type="textarea" :rows="5" placeholder="填写回复内容" />
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
  type AdminSocialConsultItem
} from '@/api/backend/compliance/social'
import { ElButton, ElMessage, ElTag } from 'element-plus'

defineOptions({ name: 'ComplianceSocialConsults' })

const searchForm = ref({
  status: '' as '' | 'open' | 'replied' | 'closed',
  q: ''
})

const replyVisible = ref(false)
const replyText = ref('')
const replying = ref(false)
const currentRow = ref<AdminSocialConsultItem | null>(null)

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
    columnsFactory: () => [
      { prop: 'memberName', label: '客户', minWidth: 100 },
      { prop: 'companyName', label: 'OPC', minWidth: 140 },
      {
        prop: 'question',
        label: '问题',
        minWidth: 200,
        showOverflowTooltip: true
      },
      {
        prop: 'status',
        label: '状态',
        width: 100,
        formatter: (row: AdminSocialConsultItem) =>
          h(ElTag, { type: row.status === 'open' ? 'warning' : row.status === 'replied' ? 'success' : 'info' }, () =>
            ADMIN_SOCIAL_CONSULT_STATUS_LABELS[row.status]
          )
      },
      { prop: 'createdAt', label: '提交时间', width: 170 },
      {
        prop: 'action',
        label: '操作',
        width: 160,
        fixed: 'right',
        formatter: (row: AdminSocialConsultItem) => {
          const buttons = []
          if (row.status !== 'closed') {
            buttons.push(
              h(ElButton, { type: 'primary', link: true, onClick: () => openReply(row) }, () => '回复')
            )
            buttons.push(
              h(ElButton, { type: 'info', link: true, onClick: () => closeRow(row) }, () => '关闭')
            )
          }
          return h('div', { class: 'flex gap-1' }, buttons)
        }
      }
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
  } catch (e: any) {
    ElMessage.error(e?.message || '操作失败')
  }
}
</script>
