<template>
  <section class="assistant-page">
    <header
      ><h1>月度经营体检</h1
      ><p>报告基于已确认资料生成；数据不足会明确标注，重要结论需人工复核。</p></header
    >
    <ElCard shadow="never"
      ><div class="actions"
        ><ElDatePicker
          v-model="period"
          type="month"
          value-format="YYYY-MM"
          @change="load"
        /><ElButton type="primary" @click="create">生成新版本</ElButton></div
      ></ElCard
    >
    <ElCard v-for="item in items" :key="item.id" shadow="never">
      <template #header
        ><div class="title"
          ><b>{{ item.periodKey }} · v{{ item.version }}</b
          ><ElTag>{{ item.status }}</ElTag></div
        ></template
      >
      <pre>{{ item.content }}</pre>
      <ElButton v-if="item.status === 'draft'" type="primary" @click="publish(item.id)"
        >确认并发布</ElButton
      >
    </ElCard>
  </section>
</template>
<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import {
    createComplianceReport,
    getComplianceReports,
    publishComplianceReport
  } from '@/api/frontend/compliance/assistant'
  import type { ComplianceReport } from '@/api/frontend/compliance/assistant'
  const period = ref(new Date().toISOString().slice(0, 7))
  const items = ref<ComplianceReport[]>([])
  onMounted(load)
  async function load() {
    items.value = (await getComplianceReports(period.value)).list || []
  }
  async function create() {
    await createComplianceReport(period.value)
    await load()
    ElMessage.success('报告草稿已生成')
  }
  async function publish(id: number) {
    await publishComplianceReport(id)
    await load()
    ElMessage.success('报告已发布')
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
  .actions,
  .title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }
  pre {
    white-space: pre-wrap;
    font: inherit;
    line-height: 1.8;
  }
</style>
