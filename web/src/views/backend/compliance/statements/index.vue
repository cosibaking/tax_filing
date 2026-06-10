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

    <ElDialog v-model="summaryVisible" title="对账单摘要" width="480px" destroy-on-close>
      <ElDescriptions v-if="summaryRow" :column="1" border size="small">
        <ElDescriptionsItem label="主播">{{ summaryRow.memberName || `会员#${summaryRow.memberId}` }}</ElDescriptionsItem>
        <ElDescriptionsItem label="OPC 主体">{{ summaryRow.opcCompanyName || '—' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="月份">{{ summaryRow.period }}</ElDescriptionsItem>
        <ElDescriptionsItem label="收入">¥{{ summaryRow.revenue?.toFixed(2) ?? '0.00' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="成本">¥{{ summaryRow.cost?.toFixed(2) ?? '0.00' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="利润">¥{{ summaryRow.profit?.toFixed(2) ?? '0.00' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="预缴税额">¥{{ summaryRow.prepaidTax?.toFixed(2) ?? '0.00' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="状态">
          <ElTag :type="statusTagType(summaryRow.status) as any" size="small">
            {{ STATEMENT_STATUS_LABELS[summaryRow.status] || summaryRow.status }}
          </ElTag>
        </ElDescriptionsItem>
        <ElDescriptionsItem label="发送时间">{{ summaryRow.sentAt || '—' }}</ElDescriptionsItem>
      </ElDescriptions>
      <template #footer>
        <ElButton @click="summaryVisible = false">关闭</ElButton>
        <ElButton type="primary" :loading="pdfLoading" @click="downloadSummaryPdf">下载 PDF</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
import { useTable } from '@/hooks/core/useTable'
import { useAuth } from '@/hooks/core/useAuth'
import {
  getAdminStatementList,
  batchGenerateStatements,
  sendStatementNotifications,
  getAdminStatementPdfUrl,
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
const summaryVisible = ref(false)
const summaryRow = ref<AdminStatementItem | null>(null)
const pdfLoading = ref(false)

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
        prop: 'cost',
        label: '成本',
        width: 100,
        formatter: (row: AdminStatementItem) => `¥${row.cost?.toFixed(2) ?? '0.00'}`
      },
      {
        prop: 'profit',
        label: '利润',
        width: 100,
        formatter: (row: AdminStatementItem) => `¥${row.profit?.toFixed(2) ?? '0.00'}`
      },
      {
        prop: 'prepaidTax',
        label: '预缴税额',
        width: 110,
        formatter: (row: AdminStatementItem) => `¥${row.prepaidTax?.toFixed(2) ?? '0.00'}`
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
          hasAuth('view') || hasAuth('send')
            ? h(ElButton, { type: 'primary', link: true, size: 'small', onClick: () => openSummary(row) }, () => '查看摘要')
            : null
      }
    ]
  }
})

function openSummary(row: AdminStatementItem) {
  summaryRow.value = row
  summaryVisible.value = true
}

async function downloadSummaryPdf() {
  if (!summaryRow.value) return
  if (summaryRow.value.pdfUrl) {
    window.open(summaryRow.value.pdfUrl, '_blank')
    return
  }
  pdfLoading.value = true
  try {
    const res = await getAdminStatementPdfUrl(summaryRow.value.id)
    if (res.url) window.open(res.url, '_blank')
    else ElMessage.warning('PDF 生成失败')
  } catch {
    ElMessage.error('PDF 下载失败')
  } finally {
    pdfLoading.value = false
  }
}

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
