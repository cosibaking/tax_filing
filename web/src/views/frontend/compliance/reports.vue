<template>
  <section class="assistant-page">
    <header class="page-header">
      <div>
        <h1>月度经营体检</h1>
        <p>按经营数据、资料完整度和风险提示分类整理；结论基于已确认资料，重要事项需人工复核。</p>
      </div>
    </header>

    <ElCard shadow="never" class="toolbar-card">
      <div class="toolbar">
        <label class="period-filter">
          <span>体检期间</span>
          <ElDatePicker
            v-model="period"
            class="period-picker"
            type="month"
            value-format="YYYY-MM"
            :clearable="false"
            :disabled="loading"
            @change="load"
          />
        </label>
        <ElButton type="primary" :loading="creating" :disabled="!period" @click="create">
          生成新版本
        </ElButton>
      </div>
      <p class="toolbar-note">数据不足时会明确标注，报告用于经营合规提醒，不替代正式申报意见。</p>
    </ElCard>

    <div v-loading="loading" class="report-list">
      <ElEmpty v-if="!loading && !items.length" description="当前期间暂无体检报告">
        <ElButton type="primary" :loading="creating" @click="create">生成首份报告</ElButton>
      </ElEmpty>

      <article v-for="item in items" :key="item.id" class="report-version">
        <div class="version-header">
          <div class="version-title">
            <h2>{{ item.periodKey }} 月度体检 · v{{ item.version }}</h2>
            <div class="version-meta">
              <ElTag :type="item.status === 'published' ? 'success' : 'info'">
                {{ statusLabel(item.status) }}
              </ElTag>
              <span>生成于 {{ formatTimestamp(item.createdAt) }}</span>
              <span v-if="item.publishedAt">发布于 {{ formatTimestamp(item.publishedAt) }}</span>
              <span v-else>报告草稿，发布前请核对资料与提示内容</span>
            </div>
          </div>
          <div class="version-actions">
            <ElButton
              v-if="item.status === 'draft'"
              type="primary"
              :loading="publishingId === item.id"
              :disabled="publishingId !== null"
              @click="publish(item.id)"
            >
              确认并发布
            </ElButton>
            <ElButton
              :loading="exportingId === item.id"
              :disabled="exportingId !== null"
              @click="exportPdf(item)"
            >
              导出 PDF
            </ElButton>
          </div>
        </div>

        <div v-if="item.structuredReport" class="structured-report">
          <ReportSummary :summary="item.structuredReport.summary" />
          <ReportCategories :categories="item.structuredReport.categories" />
          <ReportAnomalyList
            :report-id="item.id"
            :anomalies="item.structuredReport.anomalies"
            @add-materials="addMaterials(item, $event)"
            @manual-review="requestManualReview(item, $event)"
          />
        </div>

        <section v-else class="legacy-report" aria-label="历史报告内容">
          <h3>历史版本内容</h3>
          <p class="legacy-note">该版本为旧格式报告，以下保留原始内容供查阅。</p>
          <pre>{{ item.content || '该历史版本暂无可展示内容。' }}</pre>
        </section>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import {
    createComplianceReport,
    exportComplianceReportPdf,
    getComplianceReports,
    publishComplianceReport
  } from '@/api/frontend/compliance/assistant'
  import type { ComplianceReport, ReportAnomaly } from '@/api/frontend/compliance/assistant'
  import ReportSummary from './reports/ReportSummary.vue'
  import ReportCategories from './reports/ReportCategories.vue'
  import ReportAnomalyList from './reports/ReportAnomalyList.vue'

  const router = useRouter()
  const period = ref(new Date().toISOString().slice(0, 7))
  const items = ref<ComplianceReport[]>([])
  const loading = ref(false)
  const creating = ref(false)
  const publishingId = ref<number | null>(null)
  const exportingId = ref<number | null>(null)

  onMounted(load)

  async function load() {
    loading.value = true
    try {
      items.value = (await getComplianceReports(period.value)).list || []
    } catch {
      ElMessage.error('体检报告加载失败，请稍后重试')
    } finally {
      loading.value = false
    }
  }

  async function create() {
    if (!period.value || creating.value) return
    creating.value = true
    try {
      await createComplianceReport(period.value)
      await load()
      ElMessage.success('报告草稿已生成，请核对后发布')
    } catch {
      ElMessage.error('报告生成失败，请确认资料后重试')
    } finally {
      creating.value = false
    }
  }

  async function publish(id: number) {
    if (publishingId.value !== null) return
    publishingId.value = id
    try {
      await publishComplianceReport(id)
      await load()
      ElMessage.success('报告已发布')
    } catch {
      ElMessage.error('报告发布失败，请稍后重试')
    } finally {
      if (publishingId.value === id) publishingId.value = null
    }
  }

  async function exportPdf(report: ComplianceReport) {
    if (exportingId.value !== null) return
    exportingId.value = report.id
    let objectUrl = ''
    try {
      const file = await exportComplianceReportPdf(report.id)
      if (!(file instanceof Blob)) throw new TypeError('PDF response is not a Blob')
      objectUrl = URL.createObjectURL(file)
      const link = document.createElement('a')
      link.href = objectUrl
      link.download = `月度经营体检-${report.periodKey}-v${report.version}.pdf`
      document.body.appendChild(link)
      link.click()
      link.remove()
      ElMessage.success('PDF 已开始下载')
    } catch {
      ElMessage.error('PDF 导出失败，请稍后重试或申请人工协助')
    } finally {
      if (objectUrl) URL.revokeObjectURL(objectUrl)
      if (exportingId.value === report.id) exportingId.value = null
    }
  }

  function addMaterials(report: ComplianceReport, anomaly: ReportAnomaly) {
    router.push({
      name: 'MemberComplianceDocuments',
      query: reportQuery(report, anomaly)
    })
  }

  function requestManualReview(report: ComplianceReport, anomaly: ReportAnomaly) {
    router.push({
      name: 'MemberComplianceTickets',
      query: {
        ...reportQuery(report, anomaly),
        ticketType: 'compliance_review'
      }
    })
  }

  function reportQuery(report: ComplianceReport, anomaly: ReportAnomaly) {
    return {
      reportId: String(report.id),
      periodKey: report.periodKey,
      anomalyId: anomaly.code,
      anomalyCode: anomaly.code,
      anomalyTitle: anomaly.title,
      riskRuleId: anomaly.code,
      ruleVersion: anomaly.ruleVersion,
      title: anomaly.title
    }
  }

  function statusLabel(status: ComplianceReport['status']) {
    return status === 'published' ? '已发布' : '草稿'
  }

  function formatTimestamp(timestamp: number) {
    const milliseconds = timestamp < 1_000_000_000_000 ? timestamp * 1000 : timestamp
    return new Date(milliseconds).toLocaleString('zh-CN', { hour12: false })
  }
</script>

<style scoped>
  .assistant-page,
  .report-list,
  .structured-report {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .assistant-page,
  .report-list {
    gap: 18px;
  }

  .structured-report {
    gap: 16px;
    padding: 0 20px 20px;
  }

  .page-header h1,
  .version-title h2,
  .legacy-report h3 {
    margin: 0;
  }

  .page-header p {
    max-width: 840px;
    margin: 8px 0 0;
    line-height: 1.7;
    color: #64748b;
    overflow-wrap: anywhere;
  }

  .toolbar {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    align-items: flex-end;
    justify-content: space-between;
  }

  .period-filter {
    display: grid;
    flex: 0 1 260px;
    gap: 8px;
    min-width: 200px;
    font-weight: 600;
    color: #475569;
  }

  .period-picker {
    --el-date-editor-width: 100%;

    width: 100%;
  }

  .toolbar-note {
    margin: 14px 0 0;
    line-height: 1.6;
    color: #8898aa;
    overflow-wrap: anywhere;
  }

  .report-list {
    min-height: 160px;
  }

  .report-version {
    min-width: 0;
    overflow: hidden;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
  }

  .version-header {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    align-items: flex-start;
    justify-content: space-between;
    padding: 20px;
    background: #fff;
    border-bottom: 1px solid #e8edf3;
  }

  .version-title {
    min-width: 0;
  }

  .version-title h2 {
    font-size: 18px;
    overflow-wrap: anywhere;
  }

  .version-meta,
  .version-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    align-items: center;
  }

  .version-meta {
    margin-top: 10px;
    line-height: 1.6;
    color: #64748b;
  }

  .legacy-report {
    min-width: 0;
    padding: 20px;
    background: #fff;
  }

  .legacy-note {
    margin: 8px 0 14px;
    line-height: 1.6;
    color: #8898aa;
  }

  .legacy-report pre {
    padding: 16px;
    margin: 0;
    font: inherit;
    line-height: 1.8;
    color: #475569;
    overflow-wrap: anywhere;
    white-space: pre-wrap;
    background: #f8fafc;
    border-radius: 8px;
  }

  @media (width <= 640px) {
    .toolbar,
    .version-header,
    .version-actions {
      align-items: stretch;
    }

    .period-filter,
    .toolbar > .el-button,
    .version-actions,
    .version-actions .el-button {
      width: 100%;
      min-width: 0;
    }

    .version-actions {
      display: grid;
      grid-template-columns: minmax(0, 1fr);
    }

    .version-actions .el-button + .el-button {
      margin-left: 0;
    }

    .structured-report {
      padding: 0 12px 12px;
    }

    .version-header,
    .legacy-report {
      padding: 16px;
    }
  }
</style>
