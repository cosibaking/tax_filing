<template>
  <section class="report-summary" aria-labelledby="report-summary-title">
    <div class="report-summary__heading">
      <div>
        <h2 id="report-summary-title">体检结论</h2>
        <p>基于当前已提供资料形成的审慎提示</p>
      </div>
      <ElTag :type="conclusionTagType" size="large">{{ conclusionLabel }}</ElTag>
    </div>

    <div class="report-summary__metrics">
      <div class="report-summary__metric report-summary__metric--completeness">
        <span>资料完整度</span>
        <strong>{{ summary.completenessRate }}%</strong>
      </div>
      <div class="report-summary__metric report-summary__metric--high">
        <span>高风险提示</span>
        <strong>{{ summary.highCount }}</strong>
      </div>
      <div class="report-summary__metric report-summary__metric--medium">
        <span>中风险提示</span>
        <strong>{{ summary.mediumCount }}</strong>
      </div>
      <div class="report-summary__metric report-summary__metric--low">
        <span>低风险提示</span>
        <strong>{{ summary.lowCount }}</strong>
      </div>
    </div>

    <p v-if="summary.dataNotice" class="report-summary__notice" role="note">
      <strong>数据说明：</strong>{{ summary.dataNotice }}
    </p>
  </section>
</template>

<script setup lang="ts">
  import type { ReportSummary } from '@/api/frontend/compliance/assistant'

  const props = defineProps<{ summary: ReportSummary }>()

  const conclusionLabel = computed(
    () => ({ normal: '正常', attention: '需关注', urgent: '需尽快处理' })[props.summary.conclusion]
  )
  const conclusionTagType = computed(() => {
    return { normal: 'success', attention: 'warning', urgent: 'danger' }[
      props.summary.conclusion
    ] as 'success' | 'warning' | 'danger'
  })
</script>

<style scoped>
  .report-summary {
    min-width: 0;
    padding: 24px;
    color: #32325d;
    background: #fff;
    border: 1px solid #e8edf3;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgb(15 23 42 / 4%);
  }

  .report-summary__heading {
    display: flex;
    gap: 16px;
    align-items: flex-start;
    justify-content: space-between;
  }

  .report-summary__heading h2 {
    margin: 0 0 6px;
    font-size: 18px;
  }

  .report-summary__heading p {
    margin: 0;
    color: #8898aa;
  }

  .report-summary__metrics {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
    margin-top: 20px;
  }

  .report-summary__metric {
    min-width: 0;
    padding: 16px;
    background: #f8fafc;
    border: 1px solid #e8edf3;
    border-radius: 10px;
  }

  .report-summary__metric--completeness {
    grid-column: span 3;
  }

  .report-summary__metric span {
    display: block;
    margin-bottom: 8px;
    color: #64748b;
  }

  .report-summary__metric strong {
    font-size: 24px;
  }

  .report-summary__metric--high strong {
    color: #dc2626;
  }

  .report-summary__metric--medium strong {
    color: #d97706;
  }

  .report-summary__metric--low strong {
    color: #2563eb;
  }

  .report-summary__notice {
    padding: 12px 14px;
    margin: 16px 0 0;
    line-height: 1.7;
    color: #92400e;
    overflow-wrap: anywhere;
    background: #fffbeb;
    border: 1px solid #fde68a;
    border-radius: 8px;
  }

  @media (width <= 640px) {
    .report-summary {
      padding: 18px;
    }

    .report-summary__heading {
      flex-direction: column;
    }

    .report-summary__metrics {
      grid-template-columns: minmax(0, 1fr);
    }

    .report-summary__metric--completeness {
      grid-column: auto;
    }
  }
</style>
