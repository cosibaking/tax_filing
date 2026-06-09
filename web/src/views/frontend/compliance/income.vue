<!-- P-08 收入台账 -->
<template>
  <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <div>
        <h2 class="font-heading font-black text-2xl text-clay-foreground">收入台账</h2>
        <p class="text-xs text-clay-muted mt-1">多平台收入、OCR 流水导入、MCN 分成与一致性比对</p>
      </div>
      <div class="flex flex-wrap items-center gap-3">
        <ElDatePicker
          v-model="selectedMonth"
          type="month"
          value-format="YYYY-MM"
          placeholder="选择月份"
          class="clay-date-picker"
          @change="loadData"
        />
        <button
          type="button"
          class="px-5 py-2.5 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-bold text-sm shadow-clay-btn active:scale-95 transition-all flex items-center gap-2"
          @click="showAddModal = true"
        >
          <ArtSvgIcon icon="ri:add-line" />记一笔
        </button>
        <button
          type="button"
          class="px-5 py-2.5 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed text-clay-foreground font-bold text-sm active:scale-95 transition-all flex items-center gap-2"
          @click="showImportModal = true"
        >
          <ArtSvgIcon icon="ri:file-upload-line" />CSV
        </button>
        <button
          type="button"
          class="px-5 py-2.5 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed text-clay-foreground font-bold text-sm active:scale-95 transition-all flex items-center gap-2"
          @click="showOcrModal = true"
        >
          <ArtSvgIcon icon="ri:scan-line" />OCR 流水
        </button>
        <button
          type="button"
          class="px-5 py-2.5 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed text-clay-foreground font-bold text-sm active:scale-95 transition-all flex items-center gap-2"
          @click="showBankModal = true"
        >
          <ArtSvgIcon icon="ri:bank-line" />银行流水
        </button>
      </div>
    </div>

    <!-- 收入一致性比对 -->
    <div
      v-if="consistency"
      class="mb-6 p-4 rounded-2xl border text-sm"
      :class="consistencyClass"
    >
      <div class="flex flex-wrap items-center justify-between gap-2 mb-2">
        <span class="font-bold">收入一致性比对（{{ consistency.month }}）</span>
        <span class="text-xs font-bold uppercase">{{ consistencyStatusLabel }}</span>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
        <div><span class="text-clay-muted">台账实收</span><p class="font-bold">¥{{ formatMoney(consistency.ledgerNet) }}</p></div>
        <div><span class="text-clay-muted">银行已匹配</span><p class="font-bold">¥{{ formatMoney(consistency.bankMatched) }}</p></div>
        <div><span class="text-clay-muted">平台 OCR 报送</span><p class="font-bold">¥{{ formatMoney(consistency.platformReported) }}</p></div>
        <div><span class="text-clay-muted">偏差率</span><p class="font-bold">{{ (consistency.varianceRate * 100).toFixed(1) }}%</p></div>
      </div>
      <ul v-if="consistency.hints?.length" class="mt-2 space-y-1 text-xs text-clay-muted">
        <li v-for="(h, i) in consistency.hints" :key="i">· {{ h }}</li>
      </ul>
    </div>

    <!-- 平台 Tab -->
    <div class="flex flex-wrap gap-2 mb-6">
      <button
        v-for="p in INCOME_PLATFORMS"
        :key="p.value"
        type="button"
        class="px-4 py-2 rounded-xl text-sm font-bold transition-all"
        :class="activePlatform === p.value
          ? 'bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn'
          : 'bg-[#f0f3f8] text-clay-muted shadow-clay-pressed hover:text-clay-accent'"
        @click="activePlatform = p.value; loadData()"
      >
        {{ p.label }}
      </button>
    </div>

    <div v-if="loading" class="py-16 text-center">
      <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
    </div>

    <div v-else-if="list.length === 0" class="py-16 text-center text-clay-muted font-medium">
      本月暂无收入记录，点击「记一笔」开始录入
    </div>

    <div v-else class="overflow-x-auto rounded-2xl">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-left text-[10px] font-black text-clay-muted uppercase tracking-widest border-b border-gray-200/50">
            <th class="py-3 px-4">日期</th>
            <th class="py-3 px-4">类型</th>
            <th class="py-3 px-4">含税收入</th>
            <th class="py-3 px-4">平台费</th>
            <th class="py-3 px-4">实收</th>
            <th class="py-3 px-4">结算</th>
            <th class="py-3 px-4 w-20">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="row in list"
            :key="row.id"
            class="border-b border-gray-100 hover:bg-white/50 transition-colors"
          >
            <td class="py-3 px-4 font-medium text-clay-foreground">{{ row.occurredAt }}</td>
            <td class="py-3 px-4">{{ incomeCategoryLabel(row.category) }}</td>
            <td class="py-3 px-4 font-bold">¥{{ formatMoney(row.grossAmount) }}</td>
            <td class="py-3 px-4 text-clay-muted">¥{{ formatMoney(row.platformFee) }}</td>
            <td class="py-3 px-4 font-bold text-clay-success">¥{{ formatMoney(row.netAmount) }}</td>
            <td class="py-3 px-4 text-xs text-clay-muted">
              <span v-if="row.settlementType === 'mcn_public'">MCN {{ row.mcnName || '' }}</span>
              <span v-else>个人</span>
              <span v-if="row.mcnShareAmount" class="block">分成 ¥{{ formatMoney(row.mcnShareAmount) }}</span>
            </td>
            <td class="py-3 px-4">
              <button type="button" class="text-xs text-red-500 font-bold hover:underline" @click="handleDelete(row)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="summary" class="mt-6 p-4 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed flex flex-wrap items-center gap-x-6 gap-y-2 text-sm">
      <span class="font-bold text-clay-foreground">本月合计：含税 ¥{{ formatMoney(summary.grossTotal) }}</span>
      <span class="font-bold text-clay-success">实收 ¥{{ formatMoney(summary.netTotal) }}</span>
      <RouterLink to="/user/compliance/ledger" class="text-clay-accent font-bold hover:underline">
        查看利润报表 →
      </RouterLink>
    </div>

    <!-- 银行流水对账 -->
    <div v-if="bankUnmatched.length > 0" class="mt-6">
      <button
        type="button"
        class="w-full flex items-center justify-between p-4 rounded-2xl bg-orange-50 border border-orange-200 text-sm font-bold text-orange-700"
        @click="showBankPanel = !showBankPanel"
      >
        <span>银行流水对账 — 未匹配 {{ bankUnmatched.length }} 笔</span>
        <ArtSvgIcon :icon="showBankPanel ? 'ri:arrow-up-s-line' : 'ri:arrow-down-s-line'" />
      </button>
      <div v-if="showBankPanel" class="mt-3 space-y-2">
        <div
          v-for="item in bankUnmatched"
          :key="item.id"
          class="p-3 rounded-xl bg-white shadow-clay-pressed text-sm flex justify-between"
        >
          <span>{{ item.occurredAt }} {{ item.description || '银行入账' }}</span>
          <span class="font-bold">¥{{ formatMoney(item.amount) }}</span>
        </div>
      </div>
    </div>

    <!-- 记一笔 Modal -->
    <ElDialog v-model="showAddModal" title="记一笔收入" width="480px" destroy-on-close>
      <ElForm ref="addFormRef" :model="addForm" :rules="addRules" label-position="top">
        <ElFormItem label="日期" prop="occurredAt">
          <ElDatePicker v-model="addForm.occurredAt" type="date" value-format="YYYY-MM-DD" class="w-full" />
        </ElFormItem>
        <ElFormItem label="平台" prop="platform">
          <ElSelect v-model="addForm.platform" placeholder="请选择" class="w-full">
            <ElOption v-for="p in platformOptions" :key="p.value" :label="p.label" :value="p.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="收入类型" prop="category">
          <ElSelect v-model="addForm.category" placeholder="请选择" class="w-full">
            <ElOption v-for="c in INCOME_CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="含税收入（元）" prop="grossAmount">
          <ElInputNumber v-model="addForm.grossAmount" :min="0" :precision="2" class="w-full" />
        </ElFormItem>
        <ElFormItem label="结算方式">
          <ElSelect v-model="addForm.settlementType" class="w-full">
            <ElOption label="个人提现" value="personal" />
            <ElOption label="MCN 对公结算" value="mcn_public" />
          </ElSelect>
        </ElFormItem>
        <template v-if="addForm.settlementType === 'mcn_public'">
          <ElFormItem label="MCN 机构名称">
            <ElInput v-model.trim="addForm.mcnName" placeholder="签约 MCN 名称" />
          </ElFormItem>
          <ElFormItem label="分成前平台结算额（元）">
            <ElInputNumber v-model="addForm.grossBeforeSplit" :min="0" :precision="2" class="w-full" />
          </ElFormItem>
          <ElFormItem label="MCN 分成比例（%）">
            <ElInputNumber v-model="addForm.mcnSplitPercent" :min="0" :max="99" :precision="1" class="w-full" />
          </ElFormItem>
          <p v-if="mcnPreview.share > 0" class="text-xs text-clay-muted">
            MCN 分成 ¥{{ formatMoney(mcnPreview.share) }}，入账含税 ¥{{ formatMoney(mcnPreview.gross) }}
          </p>
        </template>
        <ElFormItem v-else label="含税收入（元）" prop="grossAmount">
          <ElInputNumber v-model="addForm.grossAmount" :min="0" :precision="2" class="w-full" />
        </ElFormItem>
        <ElFormItem label="平台服务费（元）" prop="platformFee">
          <ElInputNumber v-model="addForm.platformFee" :min="0" :precision="2" class="w-full" />
        </ElFormItem>
        <div class="text-sm text-clay-muted font-medium">
          实收：¥{{ formatMoney(computedNet) }}
        </div>
      </ElForm>
      <template #footer>
        <ElButton @click="showAddModal = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" @click="handleAdd">确认</ElButton>
      </template>
    </ElDialog>

    <!-- CSV 导入 Modal -->
    <ElDialog v-model="showImportModal" title="导入 CSV" width="640px" destroy-on-close @closed="resetImport">
      <div class="space-y-4">
        <a href="/templates/income-import.csv" class="text-sm text-clay-accent font-bold hover:underline" download>
          下载 CSV 模板
        </a>
        <ElUpload
          :auto-upload="false"
          :limit="1"
          accept=".csv"
          :on-change="handleFileChange"
        >
          <ElButton type="primary" plain>选择 CSV 文件</ElButton>
        </ElUpload>
        <div v-if="importPreview.length > 0" class="overflow-x-auto max-h-64 rounded-xl border">
          <table class="w-full text-xs">
            <thead class="bg-gray-50 sticky top-0">
              <tr>
                <th class="p-2">行</th>
                <th class="p-2">日期</th>
                <th class="p-2">平台</th>
                <th class="p-2">类型</th>
                <th class="p-2">含税</th>
                <th class="p-2">状态</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in importPreview" :key="row.row" :class="row.valid ? '' : 'bg-red-50'">
                <td class="p-2">{{ row.row }}</td>
                <td class="p-2">{{ row.occurredAt }}</td>
                <td class="p-2">{{ platformLabel(row.platform) }}</td>
                <td class="p-2">{{ incomeCategoryLabel(row.category) }}</td>
                <td class="p-2">{{ formatMoney(row.grossAmount) }}</td>
                <td class="p-2">{{ row.valid ? '✓' : row.error }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-if="importPreview.length > 0" class="text-sm text-clay-muted">
          有效 {{ importValidCount }} 条，无效 {{ importInvalidCount }} 条
        </p>
      </div>
      <template #footer>
        <ElButton @click="showImportModal = false">取消</ElButton>
        <ElButton
          type="primary"
          :loading="importing"
          :disabled="!importFile || importInvalidCount > 0"
          @click="handleImportConfirm"
        >
          确认导入
        </ElButton>
      </template>
    </ElDialog>

    <!-- OCR 导入 Modal -->
    <ElDialog v-model="showOcrModal" title="平台流水 OCR 导入" width="680px" destroy-on-close @closed="resetOcr">
      <div class="space-y-4">
        <ElSelect v-model="ocrPlatform" placeholder="选择平台" class="w-full">
          <ElOption v-for="p in platformOptions" :key="p.value" :label="p.label" :value="p.value" />
        </ElSelect>
        <ElUpload :auto-upload="false" :limit="1" accept="image/*" :on-change="handleOcrFileChange">
          <ElButton type="primary" plain>上传创作者中心截图</ElButton>
        </ElUpload>
        <ElInput v-model="ocrText" type="textarea" :rows="4" placeholder="若自动识别失败，请粘贴 OCR 识别文本（日期、金额行）" />
        <div v-if="ocrPreview.length" class="overflow-x-auto max-h-48 rounded-xl border">
          <table class="w-full text-xs">
            <thead class="bg-gray-50 sticky top-0">
              <tr><th class="p-2">日期</th><th class="p-2">类型</th><th class="p-2">含税</th><th class="p-2">状态</th></tr>
            </thead>
            <tbody>
              <tr v-for="row in ocrPreview" :key="row.row" :class="row.valid ? '' : 'bg-red-50'">
                <td class="p-2">{{ row.occurredAt }}</td>
                <td class="p-2">{{ incomeCategoryLabel(row.category) }}</td>
                <td class="p-2">{{ formatMoney(row.grossAmount) }}</td>
                <td class="p-2">{{ row.valid ? '✓' : row.error }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <template #footer>
        <ElButton @click="showOcrModal = false">取消</ElButton>
        <ElButton :loading="ocrPreviewing" @click="handleOcrPreview">预览解析</ElButton>
        <ElButton type="primary" :loading="ocrImporting" :disabled="ocrValidCount === 0" @click="handleOcrImport">确认导入</ElButton>
      </template>
    </ElDialog>

    <!-- 银行流水 Modal -->
    <ElDialog v-model="showBankModal" title="银行流水导入" width="480px" destroy-on-close @closed="bankFile = null">
      <p class="text-sm text-clay-muted mb-4">CSV 表头：occurred_at, amount（description 可选），导入后自动与收入台账对账。</p>
      <ElUpload :auto-upload="false" :limit="1" accept=".csv" :on-change="handleBankFileChange">
        <ElButton type="primary" plain>选择银行流水 CSV</ElButton>
      </ElUpload>
      <template #footer>
        <ElButton @click="showBankModal = false">取消</ElButton>
        <ElButton type="primary" :loading="bankImporting" :disabled="!bankFile" @click="handleBankImport">导入并对账</ElButton>
      </template>
    </ElDialog>
  </section>
</template>

<script setup lang="ts">
import {
  getIncomeList,
  createIncome,
  deleteIncome,
  previewIncomeImport,
  confirmIncomeImport,
  previewIncomeOCR,
  confirmIncomeOCR,
  getIncomeConsistency,
  importBankStatement,
  getBankUnmatched,
  INCOME_PLATFORMS,
  INCOME_CATEGORIES,
  formatMoney,
  platformLabel,
  incomeCategoryLabel,
  type IncomeEntry,
  type IncomeImportPreviewRow,
  type IncomeConsistencyResult
} from '@/api/frontend/compliance/member'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadFile } from 'element-plus'

defineOptions({ name: 'ComplianceIncome' })

const loading = ref(false)
const submitting = ref(false)
const importing = ref(false)
const list = ref<IncomeEntry[]>([])
const summary = ref<{ grossTotal: number; netTotal: number } | null>(null)
const bankUnmatched = ref<{ id: number | string; occurredAt: string; amount: number; description?: string }[]>([])
const showBankPanel = ref(false)
const consistency = ref<IncomeConsistencyResult | null>(null)

const showOcrModal = ref(false)
const showBankModal = ref(false)
const ocrPlatform = ref('douyin')
const ocrText = ref('')
const ocrFile = ref<File | null>(null)
const ocrPreview = ref<IncomeImportPreviewRow[]>([])
const ocrPreviewing = ref(false)
const ocrImporting = ref(false)
const bankFile = ref<File | null>(null)
const bankImporting = ref(false)

const now = new Date()
const selectedMonth = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)
const activePlatform = ref('')

const showAddModal = ref(false)
const showImportModal = ref(false)
const addFormRef = ref<FormInstance>()
const addForm = reactive({
  occurredAt: '',
  platform: 'douyin',
  category: 'tip',
  grossAmount: 0,
  platformFee: 0,
  settlementType: 'personal',
  mcnName: '',
  grossBeforeSplit: 0,
  mcnSplitPercent: 30
})

const addRules: FormRules = {
  occurredAt: [{ required: true, message: '请选择日期', trigger: 'change' }],
  platform: [{ required: true, message: '请选择平台', trigger: 'change' }],
  category: [{ required: true, message: '请选择类型', trigger: 'change' }],
  grossAmount: []
}

const platformOptions = INCOME_PLATFORMS.filter(p => p.value)

const importFile = ref<File | null>(null)
const importPreview = ref<IncomeImportPreviewRow[]>([])
const importValidCount = computed(() => importPreview.value.filter(r => r.valid).length)
const importInvalidCount = computed(() => importPreview.value.filter(r => !r.valid).length)
const ocrValidCount = computed(() => ocrPreview.value.filter(r => r.valid).length)

const mcnPreview = computed(() => {
  const before = addForm.grossBeforeSplit || 0
  const ratio = (addForm.mcnSplitPercent || 0) / 100
  const share = before * ratio
  return { share, gross: before - share }
})

const computedGross = computed(() => {
  if (addForm.settlementType === 'mcn_public' && addForm.grossBeforeSplit > 0) {
    return mcnPreview.value.gross
  }
  return addForm.grossAmount || 0
})

const computedNet = computed(() => computedGross.value - (addForm.platformFee || 0))

const consistencyClass = computed(() => {
  if (!consistency.value) return ''
  if (consistency.value.status === 'ok') return 'bg-green-50 border-green-200 text-green-800'
  if (consistency.value.status === 'warning') return 'bg-orange-50 border-orange-200 text-orange-800'
  return 'bg-red-50 border-red-200 text-red-800'
})

const consistencyStatusLabel = computed(() => {
  const s = consistency.value?.status
  if (s === 'ok') return '一致'
  if (s === 'warning') return '轻微偏差'
  return '需核对'
})

async function loadData() {
  loading.value = true
  try {
    const res = await getIncomeList({
      month: selectedMonth.value,
      platform: activePlatform.value || undefined,
      page: 1,
      pageSize: 100
    })
    list.value = res.list || []
    summary.value = res.summary || null
    const bank = await getBankUnmatched({ month: selectedMonth.value })
    bankUnmatched.value = bank.list || []
    try {
      consistency.value = await getIncomeConsistency({ month: selectedMonth.value })
    } catch {
      consistency.value = null
    }
  } catch {
    list.value = []
    summary.value = null
    bankUnmatched.value = []
    consistency.value = null
  } finally {
    loading.value = false
  }
}

async function handleAdd() {
  if (!addFormRef.value) return
  const valid = await addFormRef.value.validate().catch(() => false)
  if (!valid) return
  if (computedGross.value <= 0) {
    ElMessage.warning('请填写有效含税收入')
    return
  }
  submitting.value = true
  try {
    await createIncome({
      occurredAt: addForm.occurredAt,
      platform: addForm.platform,
      category: addForm.category,
      grossAmount: computedGross.value,
      platformFee: addForm.platformFee,
      settlementType: addForm.settlementType,
      mcnName: addForm.mcnName,
      mcnSplitRatio: (addForm.mcnSplitPercent || 0) / 100,
      grossBeforeSplit: addForm.grossBeforeSplit
    })
    ElMessage.success('收入已记录')
    showAddModal.value = false
    Object.assign(addForm, {
      occurredAt: '', platform: 'douyin', category: 'tip', grossAmount: 0, platformFee: 0,
      settlementType: 'personal', mcnName: '', grossBeforeSplit: 0, mcnSplitPercent: 30
    })
    await loadData()
  } catch {
    ElMessage.error('保存失败，请重试')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: IncomeEntry) {
  try {
    await ElMessageBox.confirm('确定删除该条收入记录？', '提示', { type: 'warning' })
    await deleteIncome(row.id)
    ElMessage.success('已删除')
    await loadData()
  } catch { /* cancel */ }
}

async function handleFileChange(uploadFile: UploadFile) {
  if (!uploadFile.raw) return
  importFile.value = uploadFile.raw
  try {
    const res = await previewIncomeImport(uploadFile.raw)
    importPreview.value = res.rows || []
  } catch {
    importPreview.value = []
    ElMessage.error('预览失败，请检查 CSV 格式')
  }
}

async function handleImportConfirm() {
  if (!importFile.value) return
  importing.value = true
  try {
    const res = await confirmIncomeImport(importFile.value)
    ElMessage.success(`成功导入 ${res.imported} 条`)
    showImportModal.value = false
    await loadData()
  } catch {
    ElMessage.error('导入失败，请重试')
  } finally {
    importing.value = false
  }
}

function resetImport() {
  importFile.value = null
  importPreview.value = []
}

async function handleOcrPreview() {
  if (!ocrFile.value && !ocrText.value.trim()) {
    ElMessage.warning('请上传截图或粘贴 OCR 文本')
    return
  }
  ocrPreviewing.value = true
  try {
    const res = await previewIncomeOCR({
      file: ocrFile.value || undefined,
      platform: ocrPlatform.value,
      ocrText: ocrText.value
    })
    ocrPreview.value = res.rows || []
    if (res.platform) ocrPlatform.value = res.platform
  } catch (e: unknown) {
    ocrPreview.value = []
    ElMessage.error(e instanceof Error ? e.message : 'OCR 预览失败')
  } finally {
    ocrPreviewing.value = false
  }
}

async function handleOcrImport() {
  ocrImporting.value = true
  try {
    const res = await confirmIncomeOCR({
      file: ocrFile.value || undefined,
      platform: ocrPlatform.value,
      ocrText: ocrText.value
    })
    ElMessage.success(`OCR 导入 ${res.imported} 条`)
    showOcrModal.value = false
    await loadData()
  } catch {
    ElMessage.error('OCR 导入失败')
  } finally {
    ocrImporting.value = false
  }
}

function handleOcrFileChange(uploadFile: UploadFile) {
  ocrFile.value = uploadFile.raw || null
}

function resetOcr() {
  ocrFile.value = null
  ocrText.value = ''
  ocrPreview.value = []
}

function handleBankFileChange(uploadFile: UploadFile) {
  bankFile.value = uploadFile.raw || null
}

async function handleBankImport() {
  if (!bankFile.value) return
  bankImporting.value = true
  try {
    const res = await importBankStatement(bankFile.value)
    ElMessage.success(`导入 ${res.imported} 笔，匹配 ${res.matched} 笔`)
    showBankModal.value = false
    await loadData()
  } catch {
    ElMessage.error('银行流水导入失败')
  } finally {
    bankImporting.value = false
  }
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }
.shadow-clay-card {
  box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9);
}
.shadow-clay-btn {
  box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4);
}
.shadow-clay-pressed {
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
}
:deep(.clay-date-picker) .el-input__wrapper {
  border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 6px 6px 12px #e0e5ec, inset -6px -6px 12px #ffffff; border: none;
}
</style>
