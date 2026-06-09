<!-- 合规客户列表 -->
<template>
  <div class="compliance-customers-page art-full-height">
    <ArtSearchBar
      v-model="searchForm"
      :items="formItems"
      :rules="{}"
      @reset="resetSearchParams"
      @search="handleSearch"
    />

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
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
import { useTable } from '@/hooks/core/useTable'
import { useAuth } from '@/hooks/core/useAuth'
import {
  exportComplianceAudit,
  getComplianceCustomerList,
  OPC_STATUS_LABELS,
  OPC_STATUS_OPTIONS,
  type ComplianceCustomerItem
} from '@/api/backend/compliance'
import type { OpcStatus } from '@/api/frontend/compliance/opc'
import { ElTag } from 'element-plus'

defineOptions({ name: 'ComplianceCustomers' })

const { hasAuth } = useAuth()

const searchForm = ref({
  q: undefined as string | undefined,
  opcStatus: undefined as OpcStatus | '' | undefined
})

const formItems = computed(() => [
  {
    label: '关键词',
    key: 'q',
    type: 'input',
    placeholder: '昵称 / 手机号',
    clearable: true
  },
  {
    label: 'OPC 状态',
    key: 'opcStatus',
    type: 'select',
    props: {
      placeholder: '全部',
      options: OPC_STATUS_OPTIONS.map(o => ({ label: o.label, value: o.value }))
    }
  }
])

async function handleExport(row: ComplianceCustomerItem) {
  try {
    const blob = await exportComplianceAudit(row.memberId)
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `compliance-audit-${row.memberId}.csv`
    a.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch {
    ElMessage.error('导出失败')
  }
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
    apiFn: getComplianceCustomerList,
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
      { prop: 'memberId', label: '会员ID', width: 90 },
      { prop: 'nickname', label: '昵称', minWidth: 120 },
      { prop: 'mobile', label: '手机号', minWidth: 130 },
      { prop: 'planName', label: '套餐', width: 100 },
      { prop: 'opcCompanyName', label: 'OPC 主体', minWidth: 160 },
      {
        prop: 'opcStatus',
        label: 'OPC 状态',
        width: 120,
        formatter: (row: ComplianceCustomerItem) =>
          row.opcStatus
            ? h(ElTag, { size: 'small' }, () => OPC_STATUS_LABELS[row.opcStatus!] || row.opcStatus)
            : h('span', { class: 'text-gray-400' }, '未设立')
      },
      { prop: 'signedAt', label: '签约时间', width: 160 },
      {
        prop: 'action',
        label: '操作',
        width: 100,
        fixed: 'right',
        formatter: (row: ComplianceCustomerItem) =>
          hasAuth('export')
            ? h(ArtButtonTable, {
                type: 'view',
                icon: 'ri:download-line',
                onClick: () => handleExport(row)
              })
            : null
      }
    ]
  }
})

const handleSearch = () => {
  Object.assign(searchParams, searchForm.value)
  getData()
}
</script>
