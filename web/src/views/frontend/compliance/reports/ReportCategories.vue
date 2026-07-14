<template>
  <section class="report-categories" aria-labelledby="report-categories-title">
    <h2 id="report-categories-title">分类检查结果</h2>
    <ElCollapse v-if="categories.length" class="report-categories__collapse">
      <ElCollapseItem v-for="category in categories" :key="category.code" :name="category.code">
        <template #title>
          <div class="report-categories__title">
            <strong>{{ category.name }}</strong>
            <ElTag :type="statusTagType(category.status)" size="small">
              {{ statusLabel(category.status) }}
            </ElTag>
          </div>
        </template>

        <p class="report-categories__summary">{{ category.summary }}</p>
        <ul v-if="category.checks.length" class="report-categories__checks">
          <li v-for="check in category.checks" :key="check.code">
            <div class="report-categories__check-heading">
              <strong>{{ check.name }}</strong>
              <ElTag :type="statusTagType(check.status)" effect="plain" size="small">
                {{ statusLabel(check.status) }}
              </ElTag>
            </div>
            <p>{{ check.message || checkFallback(check.status) }}</p>
          </li>
        </ul>
        <p v-else class="report-categories__empty">
          {{ category.summary || '当前分类暂无可核验数据' }}
        </p>
      </ElCollapseItem>
    </ElCollapse>
    <p v-else class="report-categories__empty">当前暂无分类检查结果</p>
  </section>
</template>

<script setup lang="ts">
  import type { ReportCategory } from '@/api/frontend/compliance/assistant'

  defineProps<{ categories: ReportCategory[] }>()

  function statusLabel(status: string) {
    const labels: Record<string, string> = {
      normal: '正常',
      attention: '需关注',
      urgent: '需尽快处理',
      high: '高风险',
      medium: '中风险',
      low: '低风险',
      insufficient: '资料不足'
    }
    return labels[status] || status || '待核验'
  }

  function statusTagType(status: string) {
    if (status === 'normal') return 'success'
    if (status === 'high' || status === 'urgent') return 'danger'
    if (status === 'attention' || status === 'medium' || status === 'insufficient') return 'warning'
    return 'info'
  }

  function checkFallback(status: string) {
    return status === 'insufficient' ? '资料不足，请补充相关资料后重新核验。' : '当前暂无补充说明。'
  }
</script>

<style scoped>
  .report-categories {
    min-width: 0;
    padding: 24px;
    color: #32325d;
    background: #fff;
    border: 1px solid #e8edf3;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgb(15 23 42 / 4%);
  }

  .report-categories h2 {
    margin: 0 0 14px;
    font-size: 18px;
  }

  .report-categories__collapse {
    min-width: 0;
    border-top: 0;
  }

  .report-categories__title,
  .report-categories__check-heading {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    align-items: center;
    min-width: 0;
  }

  .report-categories__summary,
  .report-categories__checks p,
  .report-categories__empty {
    line-height: 1.7;
    color: #64748b;
    overflow-wrap: anywhere;
  }

  .report-categories__summary {
    margin: 0 0 12px;
  }

  .report-categories__checks {
    display: grid;
    gap: 10px;
    padding: 0;
    margin: 0;
    list-style: none;
  }

  .report-categories__checks li {
    min-width: 0;
    padding: 12px 14px;
    background: #f8fafc;
    border-left: 3px solid #cbd5e1;
    border-radius: 6px;
  }

  .report-categories__checks p {
    margin: 8px 0 0;
  }

  .report-categories__empty {
    margin: 8px 0;
  }

  :deep(.el-collapse-item__header),
  :deep(.el-collapse-item__wrap),
  :deep(.el-collapse-item__content) {
    min-width: 0;
  }

  @media (width <= 640px) {
    .report-categories {
      padding: 18px;
      overflow: hidden;
    }
  }
</style>
