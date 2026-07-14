<template>
  <section class="report-anomalies" aria-labelledby="report-anomalies-title">
    <h2 id="report-anomalies-title">待处理异常</h2>
    <p class="report-anomalies__intro">
      以下内容是基于现有资料的风险提示，不替代专业判断，请结合原始凭证核实。
    </p>

    <div v-if="anomalies.length" class="report-anomalies__list">
      <article v-for="anomaly in anomalies" :key="anomaly.code" class="report-anomalies__item">
        <header class="report-anomalies__heading">
          <div>
            <div class="report-anomalies__tags">
              <ElTag :type="severityTagType(anomaly.severity)" size="small">
                {{ severityLabel(anomaly.severity) }}风险
              </ElTag>
              <ElTag type="info" effect="plain" size="small">
                {{ categoryLabel(anomaly.categoryCode) }}
              </ElTag>
            </div>
            <h3>{{ anomaly.title }}</h3>
          </div>
          <span class="report-anomalies__rule"
            >规则版本：{{ anomaly.ruleVersion || '未标注' }}</span
          >
        </header>

        <dl class="report-anomalies__details">
          <div
            ><dt>已知事实</dt><dd>{{ anomaly.facts }}</dd></div
          >
          <div
            ><dt>判断依据</dt><dd>{{ anomaly.basis }}</dd></div
          >
          <div
            ><dt>可能影响</dt><dd>{{ anomaly.impact }}</dd></div
          >
          <div
            ><dt>处理建议</dt><dd>{{ anomaly.recommendation }}</dd></div
          >
          <div>
            <dt>建议补充资料</dt>
            <dd>{{
              anomaly.requiredMaterials.length
                ? anomaly.requiredMaterials.join('、')
                : '暂无明确资料清单'
            }}</dd>
          </div>
          <div
            ><dt>建议处理时间</dt><dd>{{ anomaly.dueDate || '建议尽快处理' }}</dd></div
          >
        </dl>

        <p v-if="anomaly.requiresManualReview" class="report-anomalies__manual" role="note">
          此项建议由专业人员人工复核后再作处理决定。
        </p>
        <footer class="report-anomalies__actions">
          <ElButton @click="emit('addMaterials', anomaly)">补充资料</ElButton>
          <ElButton type="primary" @click="emit('manualReview', anomaly)">申请人工复核</ElButton>
        </footer>
      </article>
    </div>
    <p v-else class="report-anomalies__empty">当前未发现需要处理的异常</p>
  </section>
</template>

<script setup lang="ts">
  import type { ReportAnomaly } from '@/api/frontend/compliance/assistant'

  defineProps<{ reportId: number; anomalies: ReportAnomaly[] }>()
  const emit = defineEmits<{
    addMaterials: [anomaly: ReportAnomaly]
    manualReview: [anomaly: ReportAnomaly]
  }>()

  function severityLabel(severity: ReportAnomaly['severity']) {
    return { high: '高', medium: '中', low: '低' }[severity]
  }

  function severityTagType(severity: ReportAnomaly['severity']) {
    return { high: 'danger', medium: 'warning', low: 'info' }[severity] as
      | 'danger'
      | 'warning'
      | 'info'
  }

  function categoryLabel(code: string) {
    const labels: Record<string, string> = {
      business: '经营数据',
      documents: '资料完整度',
      tax: '税务与申报',
      funds: '资金与股东往来',
      employment: '用工与社保',
      annual: '工商与年度事项'
    }
    return labels[code] || '其他合规事项'
  }
</script>

<style scoped>
  .report-anomalies {
    min-width: 0;
    color: #32325d;
  }

  .report-anomalies h2 {
    margin: 0 0 6px;
    font-size: 18px;
  }

  .report-anomalies__intro {
    margin: 0 0 14px;
    line-height: 1.7;
    color: #8898aa;
  }

  .report-anomalies__list {
    display: grid;
    gap: 14px;
  }

  .report-anomalies__item {
    min-width: 0;
    padding: 20px;
    background: #fff;
    border: 1px solid #e8edf3;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgb(15 23 42 / 4%);
  }

  .report-anomalies__heading {
    display: flex;
    gap: 16px;
    align-items: flex-start;
    justify-content: space-between;
  }

  .report-anomalies__tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .report-anomalies__heading h3 {
    margin: 10px 0 0;
    font-size: 17px;
    overflow-wrap: anywhere;
  }

  .report-anomalies__rule {
    flex: 0 0 auto;
    font-size: 12px;
    color: #8898aa;
  }

  .report-anomalies__details {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px 20px;
    margin: 18px 0 0;
  }

  .report-anomalies__details div {
    min-width: 0;
  }

  .report-anomalies__details dt {
    margin-bottom: 5px;
    font-weight: 700;
    color: #475569;
  }

  .report-anomalies__details dd {
    margin: 0;
    line-height: 1.7;
    color: #64748b;
    overflow-wrap: anywhere;
    white-space: pre-wrap;
  }

  .report-anomalies__manual {
    padding: 10px 12px;
    margin: 16px 0 0;
    line-height: 1.6;
    color: #92400e;
    overflow-wrap: anywhere;
    background: #fff7ed;
    border-left: 3px solid #f59e0b;
  }

  .report-anomalies__actions {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: flex-end;
    margin-top: 18px;
  }

  .report-anomalies__empty {
    padding: 28px;
    margin: 0;
    color: #64748b;
    text-align: center;
    background: #f8fafc;
    border: 1px dashed #cbd5e1;
    border-radius: 10px;
  }

  @media (width <= 640px) {
    .report-anomalies__item {
      padding: 18px;
      overflow: hidden;
    }

    .report-anomalies__heading {
      flex-direction: column;
    }

    .report-anomalies__rule {
      flex: auto;
    }

    .report-anomalies__details {
      grid-template-columns: minmax(0, 1fr);
    }

    .report-anomalies__actions {
      justify-content: flex-start;
    }
  }
</style>
