<template>
  <section class="assistant-page" v-loading="loading">
    <ElButton text @click="router.back()">← 返回合规日历</ElButton>
    <ElCard v-if="task" shadow="never">
      <template #header
        ><div class="title"
          ><strong>{{ task.title }}</strong
          ><ElTag :type="tagType">{{ statusLabel[task.status] }}</ElTag></div
        ></template
      >
      <ElDescriptions :column="2" border>
        <ElDescriptionsItem label="截止时间">{{ formatDate(task.dueAt) }}</ElDescriptionsItem>
        <ElDescriptionsItem label="优先级">{{ priorityLabel[task.priority] }}</ElDescriptionsItem>
        <ElDescriptionsItem label="任务说明" :span="2">{{ task.description }}</ElDescriptionsItem>
      </ElDescriptions>
      <h3>所需材料</h3>
      <ul
        ><li v-for="material in task.materials" :key="material">{{ material }}</li></ul
      >
      <ElInput v-model="note" type="textarea" :rows="3" placeholder="补充完成说明或退回原因" />
      <div class="actions">
        <ElButton
          v-if="task.status === 'not_started' || task.status === 'overdue'"
          @click="act('start')"
          >开始准备</ElButton
        >
        <ElButton v-if="task.status === 'preparing'" type="primary" @click="act('submit')"
          >提交确认</ElButton
        >
        <ElButton
          v-if="task.status === 'pending_confirmation'"
          type="success"
          @click="act('complete')"
          >确认完成</ElButton
        >
      </div>
    </ElCard>
  </section>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { getComplianceTask, transitionComplianceTask } from '@/api/frontend/compliance/assistant'
  import type { ComplianceTask, TaskStatus } from '@/api/frontend/compliance/assistant'

  defineOptions({ name: 'ComplianceAssistantTaskDetail' })
  const route = useRoute()
  const router = useRouter()
  const task = ref<ComplianceTask>()
  const loading = ref(false)
  const note = ref('')
  const statusLabel: Record<TaskStatus, string> = {
    not_started: '未开始',
    preparing: '待准备材料',
    pending_confirmation: '待确认',
    completed: '已完成',
    overdue: '已逾期',
    cancelled: '已取消'
  }
  const priorityLabel = { low: '低', medium: '中', high: '高' }
  const tagType = computed(() =>
    task.value?.status === 'overdue'
      ? 'danger'
      : task.value?.status === 'completed'
        ? 'success'
        : 'warning'
  )
  onMounted(load)
  async function load() {
    loading.value = true
    try {
      task.value = (await getComplianceTask(String(route.params.id))).task
    } finally {
      loading.value = false
    }
  }
  async function act(action: string) {
    task.value = (await transitionComplianceTask(String(route.params.id), action, note.value)).task
    ElMessage.success('任务状态已更新')
  }
  const formatDate = (value: string) => (value ? new Date(value).toLocaleString('zh-CN') : '-')
</script>

<style scoped>
  .assistant-page {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 20px;
  }
  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 20px;
  }
  h3 {
    margin-top: 24px;
  }
  li {
    margin: 8px 0;
    color: #475569;
  }
</style>
