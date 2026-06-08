<!-- P-16 申报工作台 -->
<template>
  <div class="compliance-filing-page art-full-height">
    <FilingSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData" />

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />

      <MarkFiledDialog
        v-model:visible="markDialogVisible"
        :task="currentTask"
        @success="refreshData"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { useTable } from '@/hooks/core/useTable'
import { useAuth } from '@/hooks/core/useAuth'
import {
  getFilingList,
  FILING_TAX_TYPE_LABELS,
  FILING_STATUS_LABELS,
  type FilingTaskItem,
  type FilingStatus
} from '@/api/backend/compliance'
import FilingSearch from './modules/filing-search.vue'
import MarkFiledDialog from './modules/mark-filed-dialog.vue'
import { ElTag, ElButton } from 'element-plus'

defineOptions({ name: 'ComplianceFiling' })

const { hasAuth } = useAuth()

const markDialogVisible = ref(false)
const currentTask = ref<FilingTaskItem | null>(null)

const searchForm = ref({
  q: undefined as string | undefined,
  taxType: undefined as string | undefined,
  period: undefined as string | undefined,
  status: undefined as string | undefined
})

const statusTagType = (status: FilingStatus) => {
  const map: Record<FilingStatus, string> = {
    pending: 'warning',
    filed: 'success',
    overdue: 'danger'
  }
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
    apiFn: getFilingList,
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
      { prop: 'id', label: 'ID', width: 80 },
      {
        prop: 'memberName',
        label: '主播',
        minWidth: 120,
        formatter: (row: FilingTaskItem) => row.memberName || `会员#${row.memberId}`
      },
      { prop: 'opcCompanyName', label: 'OPC', minWidth: 160 },
      {
        prop: 'taxType',
        label: '税种',
        width: 130,
        formatter: (row: FilingTaskItem) => FILING_TAX_TYPE_LABELS[row.taxType] || row.taxType
      },
      { prop: 'period', label: '周期', width: 100 },
      { prop: 'dueDate', label: '截止日', width: 120 },
      {
        prop: 'calculatedAmount',
        label: '预估税额',
        width: 110,
        formatter: (row: FilingTaskItem) => `¥${row.calculatedAmount?.toFixed(2) ?? '0.00'}`
      },
      {
        prop: 'status',
        label: '状态',
        width: 100,
        formatter: (row: FilingTaskItem) =>
          h(ElTag, { type: statusTagType(row.status) as any, size: 'small' }, () =>
            FILING_STATUS_LABELS[row.status] || row.status
          )
      },
      {
        prop: 'action',
        label: '操作',
        width: 120,
        fixed: 'right',
        formatter: (row: FilingTaskItem) =>
          row.status !== 'filed' && hasAuth('filed')
            ? h(ElButton, { type: 'primary', link: true, size: 'small', onClick: () => openMarkDialog(row) }, () => '标记已申报')
            : null
      }
    ]
  }
})

function openMarkDialog(row: FilingTaskItem) {
  currentTask.value = row
  markDialogVisible.value = true
}

const handleSearch = (params: Record<string, any>) => {
  Object.assign(searchParams, params)
  getData()
}
</script>
