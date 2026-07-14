<template>
  <section class="assistant-page">
    <header class="page-header">
      <div>
        <h1>企业画像</h1>
        <p>确认企业工商、税务、用工和经营信息，系统据此生成合规任务。</p>
      </div>
      <ElTag v-if="profile" type="success">版本 {{ profile.version }}</ElTag>
    </header>

    <ElAlert
      title="首期仅支持北京一人有限公司、小规模纳税人、0～5 名员工。"
      type="info"
      :closable="false"
    />

    <ElForm ref="formRef" :model="form" label-position="top" class="profile-form">
      <ElCard shadow="never">
        <template #header><strong>工商与税务</strong></template>
        <div class="form-grid">
          <ElFormItem label="注册地区"><ElInput model-value="北京市" disabled /></ElFormItem>
          <ElFormItem label="主体类型"
            ><ElInput model-value="一人有限责任公司" disabled
          /></ElFormItem>
          <ElFormItem label="纳税人类型"
            ><ElInput model-value="小规模纳税人" disabled
          /></ElFormItem>
          <ElFormItem label="增值税申报周期">
            <ElRadioGroup v-model="form.vatPeriod"
              ><ElRadio value="quarterly">按季</ElRadio
              ><ElRadio value="monthly">按月</ElRadio></ElRadioGroup
            >
          </ElFormItem>
        </div>
      </ElCard>

      <ElCard shadow="never">
        <template #header><strong>用工与经营</strong></template>
        <div class="form-grid">
          <ElFormItem label="员工人数"
            ><ElInputNumber v-model="form.employeeCount" :min="0" :max="5"
          /></ElFormItem>
          <ElFormItem label="已开票"><ElSwitch v-model="form.invoiceEnabled" /></ElFormItem>
          <ElFormItem label="已有经营收入"><ElSwitch v-model="form.hasRevenue" /></ElFormItem>
          <ElFormItem label="已有对公账户"
            ><ElSwitch v-model="form.hasPublicBankAccount"
          /></ElFormItem>
        </div>
      </ElCard>

      <div class="actions"
        ><ElButton type="primary" :loading="saving" @click="save">确认并生成画像</ElButton></div
      >
    </ElForm>
  </section>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { getEnterpriseProfile, saveEnterpriseProfile } from '@/api/frontend/compliance/assistant'
  import type {
    EnterpriseProfile,
    EnterpriseProfileData
  } from '@/api/frontend/compliance/assistant'

  defineOptions({ name: 'ComplianceAssistantProfile' })

  const profile = ref<EnterpriseProfile>()
  const saving = ref(false)
  const formRef = ref()
  const form = reactive<EnterpriseProfileData>({
    region: 'CN-BJ',
    entityType: 'one_person_limited_company',
    taxpayerType: 'small_scale',
    vatPeriod: 'quarterly',
    employeeCount: 0,
    invoiceEnabled: false,
    hasRevenue: false,
    hasPublicBankAccount: false,
    complexity: 'low'
  })

  onMounted(async () => {
    const result = await getEnterpriseProfile()
    if (result.profile) {
      profile.value = result.profile
      Object.assign(form, result.profile.data)
    }
  })

  async function save() {
    saving.value = true
    try {
      const result = await saveEnterpriseProfile({ ...form })
      profile.value = result.profile
      ElMessage.success(result.changed ? '企业画像已更新' : '企业画像没有变化')
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
  .assistant-page {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }
  .page-header h1 {
    margin: 0 0 8px;
    font-size: 26px;
  }
  .page-header p {
    margin: 0;
    color: #64748b;
  }
  .profile-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .form-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 8px 24px;
  }
  .actions {
    display: flex;
    justify-content: flex-end;
  }
  @media (max-width: 720px) {
    .form-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
