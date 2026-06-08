<!-- OPC 任务队列 -->

<template>

  <div class="opc-tasks-page art-full-height">

    <OpcTaskSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />



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



      <OpcTaskDetailDialog

        v-model:visible="detailVisible"

        :task-id="currentTaskId"

        @success="refreshData"

      />

    </ElCard>

  </div>

</template>



<script setup lang="ts">

import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'

import { useTable } from '@/hooks/core/useTable'

import { useAuth } from '@/hooks/core/useAuth'

import {

  getOpcTaskList,

  OPC_STATUS_LABELS,

  type OpcTaskItem

} from '@/api/backend/compliance'

import OpcTaskSearch from './modules/opc-task-search.vue'

import OpcTaskDetailDialog from './modules/opc-task-detail-dialog.vue'

import { ElTag } from 'element-plus'



defineOptions({ name: 'ComplianceOpcTasks' })



const route = useRoute()

const { hasAuth } = useAuth()



const detailVisible = ref(false)

const currentTaskId = ref<number | string>('')



const searchForm = ref({

  q: undefined as string | undefined,

  status: undefined as string | undefined

})



const statusTagType = (status: string) => {

  const map: Record<string, string> = {

    pending: 'info',

    materials: 'warning',

    materials_review: '',

    registering: 'primary',

    tax: 'primary',

    bank: 'primary',

    active: 'success'

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

    apiFn: getOpcTaskList,

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

        formatter: (row: OpcTaskItem) => row.memberName || `会员#${row.memberId}`

      },

      { prop: 'memberPhone', label: '手机', minWidth: 130 },

      {

        prop: 'proposedName',

        label: '拟设名称',

        minWidth: 160,

        formatter: (row: OpcTaskItem) => row.proposedName || '未命名'

      },

      {

        prop: 'status',

        label: '当前节点',

        width: 120,

        formatter: (row: OpcTaskItem) =>

          h(ElTag, { type: statusTagType(row.status) as any, size: 'small' }, () =>

            OPC_STATUS_LABELS[row.status] || row.status

          )

      },

      { prop: 'materialsSubmittedAt', label: '提交日', width: 160 },

      {

        prop: 'slaDays',

        label: 'SLA',

        width: 90,

        formatter: (row: OpcTaskItem) => (row.slaDays != null ? `第 ${row.slaDays} 天` : '—')

      },

      {

        prop: 'action',

        label: '操作',

        width: 120,

        fixed: 'right',

        formatter: (row: OpcTaskItem) =>

          hasAuth('view')

            ? h(ArtButtonTable, {

                type: 'view',

                onClick: () => openDetail(row)

              })

            : null

      }

    ]

  }

})



function openDetail(row: OpcTaskItem) {

  currentTaskId.value = row.id

  detailVisible.value = true

}



const handleSearch = (params: Record<string, any>) => {

  Object.assign(searchParams, params)

  getData()

}



onMounted(() => {

  const id = route.query.id

  if (id) {

    currentTaskId.value = String(id)

    detailVisible.value = true

  }

})

</script>


