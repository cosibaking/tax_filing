<!-- 合规业务指标卡片 -->
<template>
  <ElRow :gutter="20" class="flex">
    <ElCol v-for="item in cards" :key="item.key" :sm="12" :md="6" :lg="6">
      <div
        class="art-card relative flex flex-col justify-center h-35 px-5 mb-5 max-sm:mb-4 cursor-pointer transition-shadow hover:shadow-md"
        @click="go(item.path)"
      >
        <span class="text-g-700 text-sm">{{ item.label }}</span>
        <ArtCountTo class="text-[26px] font-medium mt-2" :target="item.value" :duration="900" />
        <span v-if="item.hint" class="text-xs text-g-500 mt-1">{{ item.hint }}</span>
        <div class="absolute top-0 bottom-0 right-5 m-auto size-12.5 rounded-xl flex-cc bg-theme/10">
          <ArtSvgIcon :icon="item.icon" class="text-xl text-theme" />
        </div>
      </div>
    </ElCol>
  </ElRow>
</template>

<script setup lang="ts">
import { getComplianceDashboard, type ComplianceDashboardData } from '@/api/backend/compliance'
import { useRouter } from 'vue-router'

const router = useRouter()
const stats = ref<ComplianceDashboardData | null>(null)

const cards = computed(() => {
  const s = stats.value
  return [
    {
      key: 'customers',
      label: '签约客户',
      value: s?.activeCustomers ?? 0,
      icon: 'ri:team-line',
      hint: '有效服务订单',
      path: '/compliance/customers'
    },
    {
      key: 'opc',
      label: 'OPC 待办',
      value: s?.pendingOpcTasks ?? 0,
      icon: 'ri:file-list-3-line',
      hint: '未完成落地',
      path: '/compliance/opc-tasks'
    },
    {
      key: 'filing',
      label: '待申报',
      value: s?.pendingFilings ?? 0,
      icon: 'ri:calendar-check-line',
      hint: s?.overdueFilings ? `含 ${s.overdueFilings} 项逾期` : '本月申报任务',
      path: '/compliance/filing'
    },
    {
      key: 'statements',
      label: '对账单草稿',
      value: s?.draftStatements ?? 0,
      icon: 'ri:file-chart-line',
      hint: '待发送通知',
      path: '/compliance/statements'
    },
    {
      key: 'social-consults',
      label: '社保咨询待办',
      value: s?.pendingSocialConsults ?? 0,
      icon: 'ri:question-answer-line',
      hint: '待首次回复',
      path: '/compliance/social-consults'
    }
  ]
})

function go(path: string) {
  router.push(path)
}

onMounted(async () => {
  try {
    stats.value = await getComplianceDashboard()
  } catch {
    ElMessage.error('加载工作台数据失败')
  }
})
</script>
