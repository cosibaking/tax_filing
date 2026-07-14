<template>
  <div class="art-full-height"
    ><ElCard shadow="never"
      ><ElTabs v-model="section" @tab-change="load"
        ><ElTabPane
          v-for="item in sections"
          :key="item.value"
          :label="item.label"
          :name="item.value" /></ElTabs
      ><ElTable :data="items" v-loading="loading"
        ><ElTableColumn prop="id" label="ID" width="90" /><ElTableColumn
          prop="opc_entity_id"
          label="企业ID"
          width="110"
        /><ElTableColumn prop="name" label="名称" /><ElTableColumn
          prop="title"
          label="标题"
        /><ElTableColumn prop="summary" label="风险/摘要" /><ElTableColumn
          prop="period_key"
          label="期间"
          width="110"
        /><ElTableColumn prop="status" label="状态" width="120" /><ElTableColumn
          v-if="section === 'rules'"
          label="规则操作"
          width="190"
          ><template #default="{ row }"
            ><ElButton link type="primary" @click="transition(row.id, 'submit_review')"
              >提交审核</ElButton
            ><ElButton link type="success" @click="transition(row.id, 'publish')"
              >发布</ElButton
            ></template
          ></ElTableColumn
        ></ElTable
      ><ElPagination
        v-model:current-page="page"
        :total="total"
        :page-size="20"
        @current-change="load" /></ElCard
  ></div>
</template>
<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { getAssistantWorkspace, transitionRuleVersion } from '@/api/backend/compliance/assistant'
  import type { WorkspaceSection } from '@/api/backend/compliance/assistant'
  const sections = [
    { value: 'rules', label: '规则审核' },
    { value: 'tasks', label: '任务' },
    { value: 'documents', label: '资料审核' },
    { value: 'risks', label: '风险' },
    { value: 'reports', label: '报告' },
    { value: 'tickets', label: '工单' }
  ] as const
  const section = ref<WorkspaceSection>('rules'),
    items = ref<Record<string, unknown>[]>([]),
    loading = ref(false),
    page = ref(1),
    total = ref(0)
  onMounted(load)
  async function load() {
    loading.value = true
    try {
      const res = await getAssistantWorkspace({
        section: section.value,
        page: page.value,
        pageSize: 20
      })
      items.value = res.list || []
      total.value = res.total
    } finally {
      loading.value = false
    }
  }
  async function transition(id: number, action: 'submit_review' | 'publish') {
    await transitionRuleVersion(id, action)
    await load()
    ElMessage.success('规则状态已更新')
  }
</script>
<style scoped>
  .el-pagination {
    margin-top: 16px;
    justify-content: flex-end;
  }
</style>
