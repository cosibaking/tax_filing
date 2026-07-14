<template>
  <section class="assistant-page">
    <header class="page-header"
      ><div><h1>合规日历</h1><p>按截止时间查看企业待办与逾期事项。</p></div></header
    >
    <ElCard shadow="never">
      <div class="filters">
        <ElDatePicker
          v-model="month"
          type="month"
          value-format="YYYY-MM"
          placeholder="选择月份"
          @change="load"
        />
        <ElSelect v-model="status" clearable placeholder="全部状态" @change="load">
          <ElOption
            v-for="item in statusOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </ElSelect>
      </div>
    </ElCard>
    <ElCard v-loading="loading" shadow="never">
      <ElEmpty v-if="!tasks.length" description="当前期间暂无合规任务" />
      <ElTimeline v-else>
        <ElTimelineItem
          v-for="item in tasks"
          :key="item.id"
          :timestamp="formatDate(item.dueAt)"
          :type="timelineType(item.status)"
        >
          <RouterLink :to="`/user/compliance/tasks/${item.id}`" class="task-link">
            <strong>{{ item.title }}</strong
            ><ElTag :type="tagType(item.status)" size="small">{{ statusLabel[item.status] }}</ElTag>
          </RouterLink>
          <p>{{ item.description }}</p>
        </ElTimelineItem>
      </ElTimeline>
    </ElCard>
  </section>
</template>

<script setup lang="ts">
  import { getComplianceTasks } from '@/api/frontend/compliance/assistant'
  import type { ComplianceTask, TaskStatus } from '@/api/frontend/compliance/assistant'

  defineOptions({ name: 'ComplianceAssistantCalendar' })
  const month = ref(new Date().toISOString().slice(0, 7))
  const status = ref<TaskStatus | ''>('')
  const tasks = ref<ComplianceTask[]>([])
  const loading = ref(false)
  const statusLabel: Record<TaskStatus, string> = {
    not_started: '未开始',
    preparing: '待准备材料',
    pending_confirmation: '待确认',
    completed: '已完成',
    overdue: '已逾期',
    cancelled: '已取消'
  }
  const statusOptions = Object.entries(statusLabel).map(([value, label]) => ({ value, label }))

  onMounted(load)
  async function load() {
    loading.value = true
    try {
      tasks.value =
        (await getComplianceTasks({ periodKey: month.value, status: status.value })).list || []
    } finally {
      loading.value = false
    }
  }
  const formatDate = (value: string) => (value ? new Date(value).toLocaleString('zh-CN') : '-')
  const tagType = (value: TaskStatus) =>
    value === 'overdue'
      ? 'danger'
      : value === 'completed'
        ? 'success'
        : value === 'pending_confirmation'
          ? 'warning'
          : 'info'
  const timelineType = (value: TaskStatus) =>
    value === 'overdue' ? 'danger' : value === 'completed' ? 'success' : 'primary'
</script>

<style scoped>
  .assistant-page {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
  .page-header h1 {
    margin: 0 0 8px;
    font-size: 26px;
  }
  .page-header p {
    margin: 0;
    color: #64748b;
  }
  .filters {
    display: flex;
    gap: 12px;
  }
  .task-link {
    display: flex;
    align-items: center;
    gap: 10px;
    color: #1e293b;
    text-decoration: none;
  }
  .task-link:hover {
    color: #2563eb;
  }
  .task-link + p {
    margin: 8px 0;
    color: #64748b;
  }
</style>
