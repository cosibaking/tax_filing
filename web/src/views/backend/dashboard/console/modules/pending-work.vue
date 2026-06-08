<!-- 合规待办列表 -->
<template>
  <ElRow :gutter="20">
    <ElCol :sm="24" :md="12">
      <ElCard shadow="never" class="mb-5">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-medium">OPC 落地待办</span>
            <ElButton type="primary" link @click="router.push('/compliance/opc-tasks')">查看全部</ElButton>
          </div>
        </template>
        <ElTable v-loading="loading" :data="opcTasks" size="small" empty-text="暂无待办 OPC 任务">
          <ElTableColumn prop="memberName" label="主播" min-width="90" />
          <ElTableColumn prop="proposedName" label="拟设名称" min-width="120" show-overflow-tooltip />
          <ElTableColumn prop="statusLabel" label="节点" width="110" />
          <ElTableColumn label="SLA" width="80">
            <template #default="{ row }">
              {{ row.daysSinceSubmit ? `第 ${row.daysSinceSubmit} 天` : '—' }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="操作" width="70" fixed="right">
            <template #default="{ row }">
              <ElButton type="primary" link size="small" @click="openOpc(row.opcId)">处理</ElButton>
            </template>
          </ElTableColumn>
        </ElTable>
      </ElCard>
    </ElCol>

    <ElCol :sm="24" :md="12">
      <ElCard shadow="never" class="mb-5">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-medium">申报待办</span>
            <ElButton type="primary" link @click="router.push('/compliance/filing')">查看全部</ElButton>
          </div>
        </template>
        <ElTable v-loading="loading" :data="filingTasks" size="small" empty-text="暂无待申报任务">
          <ElTableColumn prop="memberName" label="主播" min-width="90" />
          <ElTableColumn prop="taxTypeLabel" label="税种" width="110" />
          <ElTableColumn prop="period" label="周期" width="90" />
          <ElTableColumn prop="dueDate" label="截止" width="100" />
          <ElTableColumn label="状态" width="80">
            <template #default="{ row }">
              <ElTag :type="row.status === 'overdue' ? 'danger' : 'warning'" size="small">
                {{ row.status === 'overdue' ? '逾期' : '待办' }}
              </ElTag>
            </template>
          </ElTableColumn>
        </ElTable>
      </ElCard>
    </ElCol>
  </ElRow>
</template>

<script setup lang="ts">
import { getComplianceDashboard, type ComplianceDashboardData } from '@/api/backend/compliance'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const opcTasks = ref<ComplianceDashboardData['opcTasks']>([])
const filingTasks = ref<ComplianceDashboardData['filingTasks']>([])

function openOpc(opcId: number) {
  router.push({ path: '/compliance/opc-tasks', query: { id: String(opcId) } })
}

onMounted(async () => {
  loading.value = true
  try {
    const data = await getComplianceDashboard()
    opcTasks.value = data.opcTasks || []
    filingTasks.value = data.filingTasks || []
  } catch {
    ElMessage.error('加载待办列表失败')
  } finally {
    loading.value = false
  }
})
</script>
