<template>
  <section class="assistant-page"
    ><header
      ><h1>人工服务</h1
      ><p>软件负责资料整理与风险提示；代理记账和申报由明确标注的有资质机构承接。</p></header
    ><ElCard shadow="never"
      ><ElForm :model="form" label-width="90px"
        ><ElFormItem label="服务类型"
          ><ElSelect v-model="form.ticketType"
            ><ElOption label="合规人工复核" value="compliance_review" /><ElOption
              label="税务申报协助"
              value="tax_filing" /></ElSelect></ElFormItem
        ><ElFormItem label="标题"><ElInput v-model="form.title" /></ElFormItem
        ><ElFormItem label="问题描述"
          ><ElInput v-model="form.description" type="textarea" /></ElFormItem
        ><ElFormItem v-if="form.ticketType === 'tax_filing'" label="服务机构"
          ><ElInput
            v-model="form.providerName"
            placeholder="请填写有资质服务机构名称" /></ElFormItem
        ><ElButton type="primary" @click="submit">提交工单</ElButton></ElForm
      ></ElCard
    ><ElTable :data="items"
      ><ElTableColumn prop="title" label="工单" /><ElTableColumn
        prop="ticketType"
        label="类型" /><ElTableColumn prop="providerName" label="服务机构" /><ElTableColumn
        prop="status"
        label="状态" /></ElTable
  ></section>
</template>
<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { createComplianceTicket, getComplianceTickets } from '@/api/frontend/compliance/assistant'
  import type { ComplianceTicket } from '@/api/frontend/compliance/assistant'
  const items = ref<ComplianceTicket[]>([])
  const form = reactive({
    ticketType: 'compliance_review',
    title: '企业合规人工复核',
    description: '',
    providerName: ''
  })
  onMounted(load)
  async function load() {
    items.value = (await getComplianceTickets()).list || []
  }
  async function submit() {
    await createComplianceTicket(form)
    await load()
    ElMessage.success('工单已提交')
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
  .el-form {
    max-width: 640px;
  }
</style>
