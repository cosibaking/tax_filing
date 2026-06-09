<!-- 报税中心：报销凭证上传 + Excel 申报导入 -->
<template>
  <div class="filing-page">
    <section class="filing-panel">
      <div class="filing-panel__head">
        <div>
          <h1 class="filing-panel__title">报税中心</h1>
          <p class="filing-panel__desc">上传报销凭证、填写申报 Excel 并导入报税数据</p>
        </div>
        <RouterLink to="/user/compliance/tax" class="filing-link">查看申报日历 →</RouterLink>
      </div>

      <ElTabs v-model="activeTab" class="filing-tabs">
        <!-- 报销凭证 -->
        <ElTabPane label="报销凭证" name="reimburse">
          <div class="filing-toolbar">
            <ElDatePicker
              v-model="selectedMonth"
              type="month"
              value-format="YYYY-MM"
              placeholder="选择月份"
              @change="loadExpenses"
            />
            <button type="button" class="filing-btn filing-btn--primary" @click="showReimburseModal = true">
              <ArtSvgIcon icon="ri:add-line" />上传报销单
            </button>
          </div>

          <div v-if="expenseLoading" class="filing-empty">加载中...</div>
          <div v-else-if="expenseList.length === 0" class="filing-empty">
            本月暂无报销记录，点击「上传报销单」添加凭证
          </div>
          <div v-else class="filing-table-wrap">
            <table class="filing-table">
              <thead>
                <tr>
                  <th>日期</th>
                  <th>费用类型</th>
                  <th>金额</th>
                  <th>发票类型</th>
                  <th>凭证</th>
                  <th>说明</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in expenseList" :key="row.id">
                  <td>{{ row.occurredAt }}</td>
                  <td>{{ row.categoryName || row.category }}</td>
                  <td class="is-amount">¥{{ formatMoney(row.amount) }}</td>
                  <td>{{ invoiceLabel(row.invoiceType) }}</td>
                  <td>
                    <ArtSvgIcon v-if="row.attachmentId" icon="ri:file-check-line" class="text-green-600" />
                    <span v-else class="text-gray-400">—</span>
                  </td>
                  <td>{{ row.description || '—' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </ElTabPane>

        <!-- Excel 申报导入 -->
        <ElTabPane label="Excel 申报导入" name="excel">
          <div class="filing-excel-guide">
            <p>1. 下载报税 Excel 模板（含「申报汇总」与「报销明细」两个 Sheet）</p>
            <p>2. 按模板填写本期收入、税额与报销明细</p>
            <p>3. 上传文件预览校验后，确认导入系统</p>
          </div>
          <div class="filing-toolbar">
            <button type="button" class="filing-btn filing-btn--ghost" :disabled="templateLoading" @click="handleDownloadTemplate">
              <ArtSvgIcon icon="ri:download-2-line" />
              {{ templateLoading ? '下载中...' : '下载 Excel 模板' }}
            </button>
            <ElUpload
              :auto-upload="false"
              :limit="1"
              accept=".xlsx,.xls"
              :on-change="handleExcelChange"
              :on-remove="resetExcelImport"
            >
              <button type="button" class="filing-btn filing-btn--primary">
                <ArtSvgIcon icon="ri:file-upload-line" />选择 Excel 文件
              </button>
            </ElUpload>
          </div>

          <div v-if="importPreview" class="filing-preview">
            <h3 class="filing-preview__title">申报汇总</h3>
            <div class="filing-summary-grid" :class="{ 'is-error': !importPreview.summary.valid }">
              <div><span>申报期间</span><strong>{{ importPreview.summary.period || '—' }}</strong></div>
              <div><span>营业收入</span><strong>¥{{ formatMoney(importPreview.summary.revenue) }}</strong></div>
              <div><span>成本费用</span><strong>¥{{ formatMoney(importPreview.summary.expenseTotal) }}</strong></div>
              <div><span>利润总额</span><strong>¥{{ formatMoney(importPreview.summary.profit) }}</strong></div>
              <div><span>增值税</span><strong>¥{{ formatMoney(importPreview.summary.vatAmount) }}</strong></div>
              <div><span>企业所得税</span><strong>¥{{ formatMoney(importPreview.summary.citAmount) }}</strong></div>
              <div><span>附加税</span><strong>¥{{ formatMoney(importPreview.summary.surchargeAmount) }}</strong></div>
              <div v-if="importPreview.summary.remark" class="col-span-2"><span>备注</span><strong>{{ importPreview.summary.remark }}</strong></div>
            </div>
            <p v-if="!importPreview.summary.valid" class="filing-error">{{ importPreview.summary.error }}</p>

            <h3 class="filing-preview__title">报销明细（{{ importPreview.validExpenseCount }} 条有效 / {{ importPreview.invalidExpenseCount }} 条无效）</h3>
            <div v-if="importPreview.expenseRows.length === 0" class="filing-empty">无报销明细行</div>
            <div v-else class="filing-table-wrap filing-table-wrap--scroll">
              <table class="filing-table">
                <thead>
                  <tr>
                    <th>行</th>
                    <th>日期</th>
                    <th>类型</th>
                    <th>金额</th>
                    <th>发票</th>
                    <th>状态</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="row in importPreview.expenseRows" :key="row.row" :class="{ 'is-invalid': !row.valid }">
                    <td>{{ row.row }}</td>
                    <td>{{ row.occurredAt }}</td>
                    <td>{{ row.categoryName || row.category }}</td>
                    <td>¥{{ formatMoney(row.amount) }}</td>
                    <td>{{ row.invoiceType || '—' }}</td>
                    <td>{{ row.valid ? '✓' : row.error }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div class="filing-preview__actions">
              <button
                type="button"
                class="filing-btn filing-btn--primary"
                :disabled="!canConfirmImport"
                :class="{ 'is-loading': importing }"
                @click="handleConfirmImport"
              >
                {{ importing ? '导入中...' : '确认导入' }}
              </button>
            </div>
          </div>
        </ElTabPane>

        <!-- 提交记录 -->
        <ElTabPane label="提交记录" name="history">
          <div v-if="historyLoading" class="filing-empty">加载中...</div>
          <div v-else-if="submissionList.length === 0" class="filing-empty">暂无 Excel 申报提交记录</div>
          <div v-else class="filing-history">
            <div v-for="item in submissionList" :key="item.id" class="filing-history__card">
              <div class="filing-history__head">
                <strong>{{ item.period }} 申报</strong>
                <span class="filing-history__time">{{ item.createdAt }}</span>
              </div>
              <div class="filing-history__stats">
                <span>收入 ¥{{ formatMoney(item.revenue) }}</span>
                <span>费用 ¥{{ formatMoney(item.expenseTotal) }}</span>
                <span>利润 ¥{{ formatMoney(item.profit) }}</span>
                <span>导入报销 {{ item.expenseImported }} 条</span>
              </div>
              <p v-if="item.remark" class="filing-history__remark">{{ item.remark }}</p>
            </div>
          </div>
        </ElTabPane>
      </ElTabs>
    </section>

    <!-- 上传报销单 -->
    <ElDialog v-model="showReimburseModal" title="上传报销单" width="520px" destroy-on-close @closed="resetReimburseForm">
      <ElForm ref="reimburseFormRef" :model="reimburseForm" :rules="reimburseRules" label-position="top">
        <ElFormItem label="发生日期" prop="occurredAt">
          <ElDatePicker v-model="reimburseForm.occurredAt" type="date" value-format="YYYY-MM-DD" class="w-full" />
        </ElFormItem>
        <ElFormItem label="费用类型" prop="category">
          <ElSelect v-model="reimburseForm.category" placeholder="请选择" class="w-full">
            <ElOption v-for="c in categories" :key="c.code" :label="c.name" :value="c.code" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="金额（元）" prop="amount">
          <ElInputNumber v-model="reimburseForm.amount" :min="0" :precision="2" class="w-full" />
        </ElFormItem>
        <ElFormItem label="发票类型" prop="invoiceType">
          <ElRadioGroup v-model="reimburseForm.invoiceType">
            <ElRadio v-for="t in INVOICE_TYPES" :key="t.value" :value="t.value">{{ t.label }}</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="报销凭证">
          <MemberFileUpload v-model="reimburseForm.attachmentId" accept=".pdf,.jpg,.jpeg,.png" />
        </ElFormItem>
        <ElFormItem label="费用说明">
          <ElInput v-model.trim="reimburseForm.description" type="textarea" :rows="2" placeholder="选填" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="showReimburseModal = false">取消</ElButton>
        <ElButton type="primary" :loading="reimburseSubmitting" @click="handleReimburseSubmit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
import type { UploadFile } from 'element-plus'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import MemberFileUpload from '@/components/frontend/MemberFileUpload.vue'
import {
  getExpenseList,
  createExpense,
  getExpenseCategories,
  downloadTaxFilingTemplate,
  previewTaxFilingImport,
  confirmTaxFilingImport,
  getTaxFilingSubmissions,
  INVOICE_TYPES,
  formatMoney,
  type ExpenseEntry,
  type ExpenseCategory,
  type TaxFilingImportPreviewResult,
  type TaxFilingSubmissionItem,
} from '@/api/frontend/compliance/member'

defineOptions({ name: 'ComplianceFiling' })

const activeTab = ref('reimburse')
const selectedMonth = ref(new Date().toISOString().slice(0, 7))

const expenseLoading = ref(false)
const expenseList = ref<ExpenseEntry[]>([])
const categories = ref<ExpenseCategory[]>([])

const showReimburseModal = ref(false)
const reimburseSubmitting = ref(false)
const reimburseFormRef = ref<FormInstance>()
const reimburseForm = reactive({
  occurredAt: '',
  category: '',
  amount: 0,
  invoiceType: 'general',
  attachmentId: 0 as number | string,
  description: '',
})
const reimburseRules: FormRules = {
  occurredAt: [{ required: true, message: '请选择日期', trigger: 'change' }],
  category: [{ required: true, message: '请选择费用类型', trigger: 'change' }],
  amount: [{ required: true, message: '请填写金额', trigger: 'blur' }],
}

const templateLoading = ref(false)
const importFile = ref<File | null>(null)
const importPreview = ref<TaxFilingImportPreviewResult | null>(null)
const previewLoading = ref(false)
const importing = ref(false)
const excelAttachmentId = ref(0)

const historyLoading = ref(false)
const submissionList = ref<TaxFilingSubmissionItem[]>([])

const canConfirmImport = computed(() =>
  importFile.value &&
  importPreview.value?.summary.valid &&
  (importPreview.value?.invalidExpenseCount ?? 0) === 0 &&
  !importing.value
)

function invoiceLabel(type: string) {
  return INVOICE_TYPES.find(t => t.value === type)?.label || type || '—'
}

async function loadExpenses() {
  expenseLoading.value = true
  try {
    const res = await getExpenseList({ month: selectedMonth.value, pageSize: 50 })
    expenseList.value = res?.list || []
  } catch {
    expenseList.value = []
  } finally {
    expenseLoading.value = false
  }
}

async function loadCategories() {
  try {
    const res = await getExpenseCategories()
    categories.value = res?.list || []
  } catch {
    categories.value = []
  }
}

async function loadHistory() {
  historyLoading.value = true
  try {
    const res = await getTaxFilingSubmissions({ pageSize: 20 })
    submissionList.value = res?.list || []
  } catch {
    submissionList.value = []
  } finally {
    historyLoading.value = false
  }
}

function resetReimburseForm() {
  reimburseForm.occurredAt = ''
  reimburseForm.category = ''
  reimburseForm.amount = 0
  reimburseForm.invoiceType = 'general'
  reimburseForm.attachmentId = 0
  reimburseForm.description = ''
}

async function handleReimburseSubmit() {
  if (!reimburseFormRef.value) return
  const valid = await reimburseFormRef.value.validate().catch(() => false)
  if (!valid) return
  reimburseSubmitting.value = true
  try {
    await createExpense({
      occurredAt: reimburseForm.occurredAt,
      category: reimburseForm.category,
      amount: reimburseForm.amount,
      invoiceType: reimburseForm.invoiceType,
      attachmentId: reimburseForm.attachmentId ? String(reimburseForm.attachmentId) : undefined,
      description: reimburseForm.description,
      costRatioAck: true,
    })
    ElMessage.success('报销单已提交')
    showReimburseModal.value = false
    await loadExpenses()
  } catch { /* handled */ } finally {
    reimburseSubmitting.value = false
  }
}

async function handleDownloadTemplate() {
  templateLoading.value = true
  try {
    await downloadTaxFilingTemplate()
    ElMessage.success('模板已开始下载')
  } catch {
    ElMessage.error('模板下载失败')
  } finally {
    templateLoading.value = false
  }
}

async function handleExcelChange(uploadFile: UploadFile) {
  const file = uploadFile.raw
  if (!file) return
  importFile.value = file
  previewLoading.value = true
  try {
    importPreview.value = await previewTaxFilingImport(file)
  } catch {
    importPreview.value = null
    importFile.value = null
  } finally {
    previewLoading.value = false
  }
}

function resetExcelImport() {
  importFile.value = null
  importPreview.value = null
  excelAttachmentId.value = 0
}

async function handleConfirmImport() {
  if (!importFile.value || !canConfirmImport.value) return
  importing.value = true
  try {
    const res = await confirmTaxFilingImport(importFile.value, excelAttachmentId.value || undefined)
    ElMessage.success(`导入成功：${res.period}，报销 ${res.expenseImported} 条`)
    resetExcelImport()
    activeTab.value = 'history'
    await loadHistory()
    await loadExpenses()
  } catch { /* handled */ } finally {
    importing.value = false
  }
}

watch(activeTab, (tab) => {
  if (tab === 'reimburse') loadExpenses()
  if (tab === 'history') loadHistory()
})

onMounted(async () => {
  await loadCategories()
  await loadExpenses()
})
</script>

<style lang="scss" scoped>
.filing-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.filing-panel {
  padding: 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.filing-panel__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 20px;
}

.filing-panel__title {
  margin: 0 0 6px;
  font-size: 22px;
  font-weight: 800;
  color: #1a1f36;
}

.filing-panel__desc {
  margin: 0;
  font-size: 14px;
  color: #6b7c93;
}

.filing-link {
  font-size: 13px;
  font-weight: 700;
  color: #2563eb;
  text-decoration: none;
  &:hover { text-decoration: underline; }
}

.filing-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
}

.filing-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.15s ease;

  &--primary {
    color: #fff;
    background: #2563eb;
    &:hover:not(:disabled) { background: #1d4ed8; }
  }

  &--ghost {
    color: #334155;
    background: #fff;
    border: 1px solid #d8dee9;
    &:hover:not(:disabled) { background: #f8fafc; }
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.filing-empty {
  padding: 40px 16px;
  text-align: center;
  color: #94a3b8;
  font-size: 14px;
  font-weight: 600;
}

.filing-table-wrap {
  overflow-x: auto;
  border: 1px solid #e8edf3;
  border-radius: 8px;

  &--scroll {
    max-height: 280px;
    overflow-y: auto;
  }
}

.filing-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;

  th, td {
    padding: 10px 14px;
    text-align: left;
    border-bottom: 1px solid #f1f5f9;
  }

  th {
    font-size: 11px;
    font-weight: 800;
    color: #94a3b8;
    text-transform: uppercase;
    background: #f8fafc;
    position: sticky;
    top: 0;
  }

  tr.is-invalid {
    background: #fef2f2;
  }

  .is-amount {
    font-weight: 700;
    color: #1a1f36;
  }
}

.filing-excel-guide {
  margin-bottom: 16px;
  padding: 14px 16px;
  background: #f0f7ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  font-size: 13px;
  color: #1e40af;
  line-height: 1.8;

  p { margin: 0; }
}

.filing-preview {
  margin-top: 8px;
}

.filing-preview__title {
  margin: 20px 0 10px;
  font-size: 14px;
  font-weight: 700;
  color: #475569;
}

.filing-summary-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;

  &.is-error {
    border: 1px solid #fecaca;
    border-radius: 8px;
    padding: 12px;
    background: #fef2f2;
  }

  > div {
    padding: 12px;
    background: #f8fafc;
    border: 1px solid #e8edf3;
    border-radius: 8px;

    span {
      display: block;
      font-size: 11px;
      color: #94a3b8;
      margin-bottom: 4px;
    }

    strong {
      font-size: 15px;
      color: #1a1f36;
    }
  }

  .col-span-2 {
    grid-column: span 2;
  }
}

.filing-error {
  margin: 8px 0 0;
  font-size: 13px;
  color: #b91c1c;
  font-weight: 600;
}

.filing-preview__actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.filing-history {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.filing-history__card {
  padding: 16px;
  border: 1px solid #e8edf3;
  border-radius: 8px;
  background: #f8fafc;
}

.filing-history__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;

  strong { font-size: 15px; color: #1a1f36; }
}

.filing-history__time {
  font-size: 12px;
  color: #94a3b8;
}

.filing-history__stats {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 13px;
  color: #475569;
  font-weight: 600;
}

.filing-history__remark {
  margin: 8px 0 0;
  font-size: 12px;
  color: #64748b;
}

@media (max-width: 768px) {
  .filing-summary-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
