<template>
  <section class="assistant-page">
    <header
      ><h1>风险中心</h1><p>系统只提供疑点提示；请确认事实，重要事项交由专业人员审核。</p></header
    >
    <ElCard shadow="never">
      <div class="filters">
        <ElDatePicker v-model="period" type="month" value-format="YYYY-MM" @change="load" />
        <ElSelect v-model="status" clearable placeholder="全部状态" @change="load">
          <ElOption label="待处理" value="open" />
          <ElOption label="已确认" value="confirmed" />
          <ElOption label="已解决" value="resolved" />
          <ElOption label="已排除" value="dismissed" />
        </ElSelect>
      </div>
    </ElCard>
    <ElTable :data="items" v-loading="loading" empty-text="当前期间暂无风险提示">
      <ElTableColumn prop="periodKey" label="期间" width="110" />
      <ElTableColumn label="等级" width="90">
        <template #default="{ row }"
          ><ElTag :type="row.severity === 'high' ? 'danger' : 'warning'">{{
            formatSeverity(row.severity)
          }}</ElTag></template
        >
      </ElTableColumn>
      <ElTableColumn prop="summary" label="风险提示" min-width="280" />
      <ElTableColumn prop="status" label="状态" width="100" />
      <ElTableColumn label="操作" width="210">
        <template #default="{ row }">
          <template v-if="row.status === 'open' || row.status === 'confirmed'">
            <ElButton
              v-if="row.status === 'open'"
              link
              type="primary"
              @click="decide(row.id, 'confirm')"
              >确认疑点</ElButton
            >
            <ElButton link type="success" @click="decide(row.id, 'resolve')">标记解决</ElButton>
            <ElButton v-if="row.status === 'open'" link @click="decide(row.id, 'dismiss')"
              >排除</ElButton
            >
          </template>
        </template>
      </ElTableColumn>
    </ElTable>
  </section>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { getComplianceRisks, transitionComplianceRisk } from '@/api/frontend/compliance/assistant'
  import type { ComplianceRisk } from '@/api/frontend/compliance/assistant'

  const period = ref(new Date().toISOString().slice(0, 7))
  const status = ref('')
  const loading = ref(false)
  const items = ref<ComplianceRisk[]>([])
  function formatSeverity(value: ComplianceRisk['severity']) {
    return { low: '低', medium: '中', high: '高' }[value]
  }

  onMounted(load)
  async function load() {
    loading.value = true
    try {
      items.value =
        (await getComplianceRisks({ periodKey: period.value, status: status.value })).list || []
    } finally {
      loading.value = false
    }
  }
  async function decide(id: number, action: 'confirm' | 'dismiss' | 'resolve') {
    await transitionComplianceRisk(id, action)
    await load()
    ElMessage.success('风险状态已更新')
  }
</script>

<style scoped>
  .assistant-page {
    display: flex;
    flex-direction: column;
    gap: 18px;
  }
  header h1 {
    margin: 0 0 8px;
  }
  header p {
    margin: 0;
    color: #64748b;
  }
  .filters {
    display: flex;
    gap: 16px;
  }
</style>
