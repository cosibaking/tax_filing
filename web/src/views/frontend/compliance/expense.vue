<!-- P-09 费用台账 -->
<template>
  <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <div>
        <h2 class="font-heading font-black text-2xl text-clay-foreground">费用台账</h2>
        <p class="text-xs text-clay-muted mt-1">成本录入与凭证上传</p>
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
      </div>
    </div>

    <!-- 虚增拦截警告 -->
    <div v-if="costRatioWarning" class="mb-4 p-4 rounded-2xl bg-yellow-50 border border-yellow-200 text-sm text-yellow-800 font-medium">
      本月成本/收入比已达 {{ (summary?.costRatio ?? 0).toFixed(0) }}%，超过 80% 阈值，录入时请确认费用真实性。
    </div>

    <div v-if="loading" class="py-16 text-center">
      <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
    </div>

    <div v-else-if="list.length === 0" class="py-16 text-center text-clay-muted font-medium">
      本月暂无费用记录
    </div>

    <div v-else class="overflow-x-auto rounded-2xl">
      <table class="w-full text-sm">
        <thead>
          <tr class="text-left text-[10px] font-black text-clay-muted uppercase tracking-widest border-b border-gray-200/50">
            <th class="py-3 px-4">日期</th>
            <th class="py-3 px-4">费用类型</th>
            <th class="py-3 px-4">金额</th>
            <th class="py-3 px-4">发票类型</th>
            <th class="py-3 px-4">凭证</th>
            <th class="py-3 px-4 w-20">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="row in list"
            :key="row.id"
            class="border-b border-gray-100 hover:bg-white/50 transition-colors"
            :class="{ 'bg-yellow-50/50': row.warningFlag }"
          >
            <td class="py-3 px-4 font-medium">{{ row.occurredAt }}</td>
            <td class="py-3 px-4">{{ row.categoryName || row.category }}</td>
            <td class="py-3 px-4 font-bold">¥{{ formatMoney(row.amount) }}</td>
            <td class="py-3 px-4">{{ invoiceLabel(row.invoiceType) }}</td>
            <td class="py-3 px-4">
              <ArtSvgIcon v-if="row.attachmentId" icon="ri:file-check-line" class="text-clay-success" />
              <span v-else class="text-clay-muted">—</span>
            </td>
            <td class="py-3 px-4">
              <button type="button" class="text-xs text-red-500 font-bold hover:underline" @click="handleDelete(row)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="summary" class="mt-6 p-4 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed text-sm font-bold text-clay-foreground">
      本月费用合计：¥{{ formatMoney(summary.totalAmount) }}
    </div>

    <!-- 记一笔 Modal -->
    <ElDialog v-model="showAddModal" title="记一笔费用" width="480px" destroy-on-close @closed="resetAddForm">
      <div v-if="blacklistHit" class="mb-4 p-3 rounded-xl bg-red-50 border border-red-200 text-sm text-red-700 font-bold">
        描述含敏感关键词（咨询费/服务费/借款/赠与），禁止提交。
      </div>
      <div v-else-if="costRatioAckRequired" class="mb-4 p-3 rounded-xl bg-yellow-50 border border-yellow-200 text-sm text-yellow-800">
        成本/收入比超过 80%，请确认后继续。
        <ElCheckbox v-model="addForm.costRatioAck" class="mt-2">我已确认费用真实合理</ElCheckbox>
      </div>

      <ElForm ref="addFormRef" :model="addForm" :rules="addRules" label-position="top">
        <ElFormItem label="日期" prop="occurredAt">
          <ElDatePicker v-model="addForm.occurredAt" type="date" value-format="YYYY-MM-DD" class="w-full" />
        </ElFormItem>
        <ElFormItem prop="category">
          <template #label>
            <span>费用类型</span>
            <ElTooltip v-if="selectedCategoryHint" :content="selectedCategoryHint" placement="top">
              <ArtSvgIcon icon="ri:question-line" class="ml-1 text-clay-accent cursor-help inline" />
            </ElTooltip>
          </template>
          <ElSelect v-model="addForm.category" placeholder="请选择" class="w-full">
            <ElOption v-for="c in categories" :key="c.code" :label="c.name" :value="c.code" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="金额（元）" prop="amount">
          <ElInputNumber v-model="addForm.amount" :min="0" :precision="2" class="w-full" />
        </ElFormItem>
        <ElFormItem label="发票类型" prop="invoiceType">
          <ElRadioGroup v-model="addForm.invoiceType">
            <ElRadio v-for="t in INVOICE_TYPES" :key="t.value" :value="t.value">{{ t.label }}</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="凭证上传">
          <ArtFileUpload v-model="addForm.attachmentId" accept=".pdf,.jpg,.jpeg,.png" />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput v-model.trim="addForm.description" type="textarea" :rows="2" placeholder="费用说明（选填）" />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="showAddModal = false">取消</ElButton>
        <ElButton type="primary" :loading="submitting" :disabled="blacklistHit || (costRatioAckRequired && !addForm.costRatioAck)" @click="handleAdd">
          确认
        </ElButton>
      </template>
    </ElDialog>
  </section>
</template>

<script setup lang="ts">
import ArtFileUpload from '@/components/core/forms/art-file-upload/index.vue'
import {
  getExpenseList,
  createExpense,
  deleteExpense,
  getExpenseCategories,
  INVOICE_TYPES,
  EXPENSE_BLACKLIST_KEYWORDS,
  formatMoney,
  type ExpenseEntry,
  type ExpenseCategory
} from '@/api/frontend/compliance/member'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

defineOptions({ name: 'ComplianceExpense' })

const loading = ref(false)
const submitting = ref(false)
const list = ref<ExpenseEntry[]>([])
const categories = ref<ExpenseCategory[]>([])
const summary = ref<{ totalAmount: number; costRatio?: number } | null>(null)

const now = new Date()
const selectedMonth = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)

const showAddModal = ref(false)
const addFormRef = ref<FormInstance>()
const addForm = reactive({
  occurredAt: '',
  category: '',
  amount: 0,
  invoiceType: 'general',
  attachmentId: '',
  description: '',
  costRatioAck: false
})

const addRules: FormRules = {
  occurredAt: [{ required: true, message: '请选择日期', trigger: 'change' }],
  category: [{ required: true, message: '请选择费用类型', trigger: 'change' }],
  amount: [{ required: true, message: '请填写金额', trigger: 'blur' }],
  invoiceType: [{ required: true, message: '请选择发票类型', trigger: 'change' }]
}

const costRatioWarning = computed(() => (summary.value?.costRatio ?? 0) > 80)
const costRatioAckRequired = computed(() => costRatioWarning.value)
const selectedCategoryHint = computed(() => categories.value.find(c => c.code === addForm.category)?.voucherHint)
const blacklistHit = computed(() =>
  EXPENSE_BLACKLIST_KEYWORDS.some(kw => (addForm.description || '').includes(kw))
)

function invoiceLabel(type: string) {
  return INVOICE_TYPES.find(t => t.value === type)?.label || type
}

async function loadCategories() {
  try {
    const res = await getExpenseCategories()
    categories.value = res.list || []
    if (categories.value.length && !addForm.category) {
      addForm.category = categories.value[0].code
    }
  } catch {
    categories.value = [
      { code: 'equipment', name: '设备', voucherHint: '需提供购置发票' },
      { code: 'network', name: '网费', voucherHint: '需提供通信费发票' },
      { code: 'venue', name: '场地', voucherHint: '合同 + 发票' },
      { code: 'promotion', name: '投流', voucherHint: '平台推广费发票' },
      { code: 'outsource', name: '外包', voucherHint: '合同 + 发票' },
      { code: 'travel', name: '差旅', voucherHint: '车票/住宿票' },
      { code: 'other', name: '其他', voucherHint: '说明 + 凭证' }
    ]
  }
}

async function loadData() {
  loading.value = true
  try {
    const res = await getExpenseList({ month: selectedMonth.value, page: 1, pageSize: 100 })
    list.value = res.list || []
    summary.value = res.summary || null
  } catch {
    list.value = []
    summary.value = null
  } finally {
    loading.value = false
  }
}

async function handleAdd() {
  if (!addFormRef.value || blacklistHit.value) return
  if (costRatioAckRequired.value && !addForm.costRatioAck) {
    ElMessage.warning('请勾选确认项')
    return
  }
  const valid = await addFormRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    await createExpense({ ...addForm })
    ElMessage.success('费用已记录')
    showAddModal.value = false
    await loadData()
  } catch {
    ElMessage.error('保存失败，请重试')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: ExpenseEntry) {
  try {
    await ElMessageBox.confirm('确定删除该条费用记录？', '提示', { type: 'warning' })
    await deleteExpense(row.id)
    ElMessage.success('已删除')
    await loadData()
  } catch { /* cancel */ }
}

function resetAddForm() {
  Object.assign(addForm, {
    occurredAt: '', category: categories.value[0]?.code || '', amount: 0,
    invoiceType: 'general', attachmentId: '', description: '', costRatioAck: false
  })
}

onMounted(async () => {
  await loadCategories()
  await loadData()
})
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }
.shadow-clay-card { box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9); }
.shadow-clay-btn { box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4); }
.shadow-clay-pressed { box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff; }
:deep(.clay-date-picker) .el-input__wrapper {
  border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 6px 6px 12px #e0e5ec, inset -6px -6px 12px #ffffff; border: none;
}
</style>
