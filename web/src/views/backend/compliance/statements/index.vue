<!-- P-17 对账单管理 -->
<template>
  <div class="compliance-statements-page art-full-height">
    <StatementSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElSpace wrap>
            <ElDatePicker
              v-model="batchPeriod"
              type="month"
              value-format="YYYY-MM"
              placeholder="选择生成月份"
              style="width: 160px"
            />
            <ElButton v-auth="'generate'" :loading="generating" @click="handleBatchGenerate" v-ripple>
              批量生成
            </ElButton>
            <ElButton
              v-auth="'send'"
              type="primary"
              :loading="sending"
              :disabled="selectedRows.length === 0"
              @click="handleSendNotifications"
              v-ripple
            >
              发送通知 ({{ selectedRows.length }})
            </ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @selection-change="handleSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { useTable } from '@/hooks/core/useTable'
import { useAuth } from '@/hooks/core/useAuth'
import {
  getAdminStatementList,
  batchGenerateStatements,
  sendStatementNotifications,
  STATEMENT_STATUS_LABELS,
  type AdminStatementItem
} from '@/api/backend/compliance'
import StatementSearch from './modules/statement-search.vue'
import { ElTag, ElButton, ElMessage, ElMessageBox } from 'element-plus'

defineOptions({ name: 'ComplianceStatements' })

const { hasAuth } = useAuth()

const selectedRows = ref<AdminStatementItem[]>([])
const batchPeriod = ref('')
const generating = ref(false)
const sending = ref(false)

const now = new Date()
batchPeriod.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`

const searchForm = ref({
  q: undefined as string | undefined,
  period: undefined as string | undefined,
  status: undefined as string | undefined
})

const statusTagType = (status: string) => {
  const map: Record<string, string> = { draft: 'info', sent: 'warning', completed: 'success' }
  return map[status] || 'info'
}

const {
  columns,
  columnChecks,
  data,
  loading,
  pagination,
  getData,
  searchParams,
  resetSearchParams,
  handleSizeChange,
  handleCurrentChange,
  refreshData
} = useTable({
  core: {
    apiFn: getAdminStatementList,
    apiParams: {
      page: 1,
      pageSize: 20,
      ...searchForm.value
    },
    paginationKey: {
      current: 'page',
      size: 'pageSize'
    },
    columnsFactory: () => [
      { type: 'selection', width: 50 },
      { prop: 'id', label: 'ID', width: 80 },
      {
        prop: 'memberName',
        label: '主播',
        minWidth: 120,
        formatter: (row: AdminStatementItem) => row.memberName || `会员#${row.memberId}`
      },
      { prop: 'opcCompanyName', label: 'OPC', minWidth: 160 },
      { prop: 'period', label: '月份', width: 100 },
      {
        prop: 'revenue',
        label: '收入',
        width: 100,
        formatter: (row: AdminStatementItem) => `¥${row.revenue?.toFixed(2) ?? '0.00'}`
      },
      {
        prop: 'profit',
        label: '利润',
        width: 100,
        formatter: (row: AdminStatementItem) => `¥${row.profit?.toFixed(2) ?? '0.00'}`
      },
      {
        prop: 'status',
        label: '状态',
        width: 100,
        formatter: (row: AdminStatementItem) =>
          h(ElTag, { type: statusTagType(row.status) as any, size: 'small' }, () =>
            STATEMENT_STATUS_LABELS[row.status] || row.status
          )
      },
      { prop: 'sentAt', label: '发送时间', width: 160 },
      {
        prop: 'action',
        label: '操作',
        width: 100,
        fixed: 'right',
        formatter: (row: AdminStatementItem) =>
          row.pdfUrl && hasAuth('view')
            ? h(ElButton, { type: 'primary', link: true, size: 'small', onClick: () => window.open(row.pdfUrl, '_blank') }, () => '预览 PDF')
            : null
      }
    ]
  }
})

const handleSearch = (params: Record<string, any>) => {
  Object.assign(searchParams, params)
  getData()
}

const handleSelectionChange = (rows: AdminStatementItem[]) => {
  selectedRows.value = rows
}

async function handleBatchGenerate() {
  if (!batchPeriod.value) {
    ElMessage.warning('请选择生成月份')
    return
  }
  try {
    await ElMessageBox.confirm(`确定为 ${batchPeriod.value} 批量生成对账单？`, '提示', { type: 'info' })
    generating.value = true
    const res = await batchGenerateStatements({ period: batchPeriod.value })
    ElMessage.success(`已生成 ${res.generated} 份对账单`)
    refreshData()
  } catch { /* cancel */ }
  finally { generating.value = false }
}

async function handleSendNotifications() {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定向 ${selectedRows.value.length} 位主播发送对账单通知？`, '提示', { type: 'info' })
    sending.value = true
    const res = await sendStatementNotifications({ ids: selectedRows.value.map(r => r.id) })
    ElMessage.success(`已发送 ${res.sent} 条通知`)
    refreshData()
  } catch { /* cancel */ }
  finally { sending.value = false }
}
</script>
