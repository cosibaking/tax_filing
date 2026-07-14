<template>
  <section class="assistant-page">
    <header
      ><h1>经营资料库</h1
      ><p>原文件上传后永久保留，识别结果经确认后才参与台账和风险分析。</p></header
    >
    <ElCard shadow="never">
      <div class="upload-row">
        <ElDatePicker v-model="period" class="period-field" type="month" value-format="YYYY-MM" />
        <div class="upload-field">
          <MemberFileUpload v-model="attachmentId" accept=".pdf,.jpg,.jpeg,.png,.csv,.xlsx" />
        </div>
        <ElButton
          class="register-button"
          type="primary"
          :disabled="!attachmentId"
          @click="register"
        >
          登记资料
        </ElButton>
      </div>
    </ElCard>
    <ElTable :data="items" v-loading="loading">
      <ElTableColumn prop="periodKey" label="期间" width="120" /><ElTableColumn
        prop="attachmentId"
        label="附件"
        width="120"
      />
      <ElTableColumn label="类型"
        ><template #default="{ row }"
          ><ElSelect v-model="row.documentType"
            ><ElOption
              v-for="type in documentTypes"
              :key="type.value"
              :label="type.label"
              :value="type.value" /></ElSelect></template
      ></ElTableColumn>
      <ElTableColumn prop="processStatus" label="状态" /><ElTableColumn label="操作" width="110"
        ><template #default="{ row }"
          ><ElButton type="primary" link @click="confirm(row)">确认</ElButton></template
        ></ElTableColumn
      >
    </ElTable>
  </section>
</template>
<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import MemberFileUpload from '@/components/frontend/MemberFileUpload.vue'
  import {
    createBusinessDocument,
    getBusinessDocuments,
    confirmBusinessDocument,
    documentTypes
  } from '@/api/frontend/compliance/assistant'
  import type { BusinessDocument } from '@/api/frontend/compliance/assistant'
  const period = ref(new Date().toISOString().slice(0, 7))
  const attachmentId = ref<number | string>('')
  const items = ref<BusinessDocument[]>([])
  const loading = ref(false)
  onMounted(load)
  async function load() {
    loading.value = true
    try {
      items.value = (await getBusinessDocuments(period.value)).list || []
    } finally {
      loading.value = false
    }
  }
  async function register() {
    await createBusinessDocument({ attachmentId: attachmentId.value, periodKey: period.value })
    attachmentId.value = ''
    await load()
    ElMessage.success('资料已登记，等待确认')
  }
  async function confirm(row: BusinessDocument) {
    await confirmBusinessDocument(row.id, row.documentType)
    await load()
    ElMessage.success('资料分类已确认')
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
  .upload-row {
    display: grid;
    grid-template-columns: 180px minmax(0, 1fr) auto;
    gap: 16px;
    align-items: start;
  }
  .period-field {
    width: 100%;
  }
  .upload-field {
    min-width: 0;
  }
  .register-button {
    align-self: start;
  }
  @media (max-width: 760px) {
    .upload-row {
      grid-template-columns: minmax(0, 1fr);
    }
    .period-field,
    .upload-field,
    .register-button {
      width: 100%;
    }
  }
</style>
