<!-- P-05 方案确认与签约 -->
<template>
  <div class="plan-page">
    <section class="overview-panel">
      <div class="overview-panel__head">
        <div>
          <h2 class="overview-panel__section-title">完成服务签约</h2>
          <p class="overview-panel__section-desc">选择套餐 → 风险告知 → 电子签约</p>
        </div>
      </div>

      <!-- 步骤指示 -->
      <div class="plan-steps">
        <div
          v-for="(label, idx) in stepLabels"
          :key="idx"
          class="plan-step"
          :class="{
            'is-done': currentStep > idx + 1,
            'is-current': currentStep === idx + 1,
          }"
        >
          <div class="plan-step__index">
            <ArtSvgIcon v-if="currentStep > idx + 1" icon="ri:check-line" />
            <span v-else>{{ idx + 1 }}</span>
          </div>
          <span class="plan-step__label">{{ label }}</span>
        </div>
      </div>

      <div v-if="apiError" class="plan-alert plan-alert--error">
        <ArtSvgIcon icon="ri:error-warning-line" />
        <p>{{ apiError }}</p>
      </div>

      <div v-if="pendingOrderHint" class="plan-alert plan-alert--info">
        <ArtSvgIcon icon="ri:information-line" />
        <p>{{ pendingOrderHint }}</p>
      </div>

      <!-- Step 1: 选择套餐 -->
      <div v-show="currentStep === 1" class="plan-body">
        <div v-if="plansLoading" class="overview-empty">
          <ArtSvgIcon icon="ri:loader-4-line" class="overview-empty__icon plan-loading" />
          <p>加载套餐中...</p>
        </div>
        <ul v-else class="plan-list">
          <li
            v-for="plan in displayPlans"
            :key="plan.tier"
            class="plan-card"
            :class="{ 'is-selected': selectedPlanTier === plan.tier }"
            @click="selectedPlanTier = plan.tier"
          >
            <div class="plan-card__head">
              <div>
                <h3 class="plan-card__title">{{ plan.name }}</h3>
                <p class="plan-card__price">
                  {{ plan.priceLabel }}
                  <span v-if="plan.monthlyPrice" class="plan-card__unit">/月起</span>
                </p>
              </div>
              <div class="plan-card__radio" :class="{ 'is-checked': selectedPlanTier === plan.tier }">
                <ArtSvgIcon v-if="selectedPlanTier === plan.tier" icon="ri:check-line" />
              </div>
            </div>
            <ul class="plan-card__features">
              <li v-for="(feat, i) in plan.features.slice(0, 3)" :key="i">{{ feat }}</li>
            </ul>
          </li>
        </ul>
      </div>

      <!-- Step 2: 风险告知 -->
      <div v-show="currentStep === 2" class="plan-body">
        <p class="plan-body__hint">请仔细阅读以下风险告知，全部勾选后方可继续：</p>
        <ul class="plan-risk-list">
          <li v-for="(item, idx) in riskItems" :key="idx" class="plan-risk-item">
            <ElCheckbox v-model="item.checked" class="plan-checkbox">
              <span>{{ item.text }}</span>
            </ElCheckbox>
          </li>
        </ul>
      </div>

      <!-- Step 3: 签约 -->
      <div v-show="currentStep === 3" class="plan-body">
        <h3 class="plan-subtitle">方案摘要</h3>
        <div class="plan-summary">
          <div class="plan-field">
            <span class="plan-field__label">服务套餐</span>
            <span class="plan-field__value">{{ selectedPlan?.name }}（{{ selectedPlan?.priceLabel }}）</span>
          </div>
          <div v-if="diagnosisId" class="plan-field">
            <span class="plan-field__label">关联诊断</span>
            <span class="plan-field__value">#{{ diagnosisId }}</span>
          </div>
          <div class="plan-field">
            <span class="plan-field__label">推荐方案</span>
            <span class="plan-field__value">{{ recommendedPlanLabel }}</span>
          </div>
          <div v-if="diagnosisSummary.monthlyIncome" class="plan-field">
            <span class="plan-field__label">月收入区间</span>
            <span class="plan-field__value">{{ diagnosisSummary.monthlyIncome }}</span>
          </div>
          <div v-if="diagnosisSummary.annualCost" class="plan-field">
            <span class="plan-field__label">年成本估算</span>
            <span class="plan-field__value">{{ diagnosisSummary.annualCost }} 元</span>
          </div>
        </div>

        <h3 class="plan-subtitle">服务协议预览</h3>
        <div
          ref="pdfContainerRef"
          class="plan-agreement"
          @scroll="onPdfScroll"
        >
          <h4>税务合规服务协议（摘要）</h4>
          <p>第一条 服务内容：本协议约定乙方为甲方提供主体设立代办、记账、申报等合规服务。</p>
          <p>第二条 服务费用：按所选套餐标准收取，具体以订单确认金额为准。</p>
          <p>第三条 甲方义务：如实提供注册资料，配合工商、税务、银行开户流程。</p>
          <p>第四条 合规声明：本服务为合法合规方案，不提供逃税、虚开发票等违法服务。</p>
          <p>第五条 责任限制：乙方不承诺「包不被查」，税负取决于甲方真实收入与成本。</p>
          <p>第六条 隐私保护：甲方敏感信息加密存储，仅用于合规服务目的。</p>
          <p>第七条 协议生效：甲方完成电子签名确认后本协议生效。</p>
          <p class="plan-agreement__note">（完整协议 PDF 将在签约后存档）</p>
        </div>
        <p v-if="!pdfScrolledToBottom" class="plan-scroll-hint">请滚动阅读完整协议</p>

        <div class="plan-sign-form">
          <div class="plan-sign-row">
            <ElCheckbox v-model="agreementAccepted" class="plan-checkbox">
              <span>我已阅读并同意服务协议</span>
            </ElCheckbox>
          </div>
          <div class="plan-sign-row">
            <label class="plan-label">
              电子签名（输入姓名确认） <span class="plan-required">*</span>
            </label>
            <ElInput
              v-model.trim="signerName"
              placeholder="请输入您的真实姓名"
              size="large"
              class="plan-input"
            />
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="plan-actions">
        <button
          v-if="currentStep > 1"
          type="button"
          class="overview-btn overview-btn--ghost"
          :disabled="submitting"
          @click="prevStep"
        >
          上一步
        </button>
        <button
          v-if="currentStep < 3"
          type="button"
          class="overview-btn overview-btn--primary plan-actions__main"
          :disabled="!canNext || submitting"
          @click="nextStep"
        >
          下一步
        </button>
        <button
          v-else
          type="button"
          class="overview-btn overview-btn--primary plan-actions__main"
          :disabled="!canSign || submitting"
          @click="handleSign"
        >
          <ArtSvgIcon v-if="submitting" icon="ri:loader-4-line" class="plan-loading" />
          {{ submitting ? '签约中...' : '完成签约' }}
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { fetchServicePlans, type ServicePlan, type RecommendedPlan } from '@/api/frontend/compliance/diagnosis'
import {
  createComplianceOrder,
  getActiveOrder,
  submitConsent,
  signAgreement,
  type RiskAcknowledgments
} from '@/api/frontend/compliance/order'
import { getCompliancePlanState, getDiagnosisDetail } from '@/api/frontend/compliance/member'
import { useMemberStore } from '@/store/modules/member'
import { requireLogin } from '@/utils/auth/requireLogin'
import { sanitizeErrorMessage } from '@/utils/http/error'
import { ElMessage } from 'element-plus'

const PENDING_ORDER_KEY = 'compliance_pending_order_id'

defineOptions({ name: 'CompliancePlan' })

const router = useRouter()
const route = useRoute()
const memberStore = useMemberStore()

const stepLabels = ['选择套餐', '风险告知', '签约']
const currentStep = ref(1)
const submitting = ref(false)
const apiError = ref('')
const plansLoading = ref(true)
const plans = ref<ServicePlan[]>([])
const selectedPlanTier = ref('basic')
const orderId = ref<number | string>('')
const pendingOrderHint = ref('')
const recommendedPlanLabel = ref('一人有限责任公司（小微公司）')
const diagnosisSummary = ref<{ monthlyIncome?: string; annualCost?: string }>({})

const diagnosisId = computed(() => route.query.diagnosisId as string | undefined)
const riskItems = ref([
  { text: '本服务为合规方案，非逃税方案', checked: false },
  { text: '不承诺「包不被查」', checked: false },
  { text: '不提供虚开发票、隐瞒收入服务', checked: false },
  { text: '税负取决于真实成本与收入', checked: false }
])

const pdfContainerRef = ref<HTMLElement>()
const pdfScrolledToBottom = ref(false)
const agreementAccepted = ref(false)
const signerName = ref('')

const PLAN_LABEL_MAP: Record<RecommendedPlan, string> = {
  opc: '一人有限责任公司（小微公司）',
  individual: '个体工商户',
  labor: '劳务报酬',
  transitional: '过渡期方案',
  none: '待评估',
}

const FALLBACK_PLANS: ServicePlan[] = [
  { id: 1, name: '基础版', tier: 'basic', monthlyPrice: 299, priceLabel: '¥299', features: ['主体注册代办', '月度记账', '季度申报'] },
  { id: 2, name: '进阶版', tier: 'advanced', monthlyPrice: 799, priceLabel: '¥799', recommended: true, features: ['基础版全部', '税务筹划', '发票管理'] },
  { id: 3, name: '尊享版', tier: 'premium', monthlyPrice: null, priceLabel: '面议', features: ['进阶版全部', '架构设计', '稽查应对'] }
]

const displayPlans = computed(() => (plans.value.length ? plans.value : FALLBACK_PLANS))
const selectedPlan = computed(() => displayPlans.value.find(p => p.tier === selectedPlanTier.value))

const allRiskChecked = computed(() => riskItems.value.every(i => i.checked))
const canNext = computed(() => {
  if (currentStep.value === 1) return !!selectedPlanTier.value
  if (currentStep.value === 2) return allRiskChecked.value
  return true
})
const canSign = computed(() =>
  pdfScrolledToBottom.value &&
  agreementAccepted.value &&
  signerName.value.length >= 2
)

function onPdfScroll() {
  const el = pdfContainerRef.value
  if (!el) return
  if (el.scrollTop + el.clientHeight >= el.scrollHeight - 8) {
    pdfScrolledToBottom.value = true
  }
}

async function loadDiagnosis() {
  if (!diagnosisId.value) return
  try {
    const detail = await getDiagnosisDetail(diagnosisId.value)
    recommendedPlanLabel.value = PLAN_LABEL_MAP[detail.recommendedPlan] || detail.recommendedPlan
    diagnosisSummary.value = {
      monthlyIncome: (detail as { monthlyIncomeRange?: string }).monthlyIncomeRange,
      annualCost: detail.taxComparison?.annualCost
        ? String(Math.round(detail.taxComparison.annualCost))
        : undefined,
    }
  } catch {
    /* ignore */
  }
}

async function loadPlans() {
  plansLoading.value = true
  try {
    const res = await fetchServicePlans()
    plans.value = res?.list?.length ? res.list : FALLBACK_PLANS
  } catch {
    plans.value = FALLBACK_PLANS
  } finally {
    plansLoading.value = false
  }
  const queryPlan = route.query.plan as string
  if (queryPlan && displayPlans.value.some(p => p.tier === queryPlan)) {
    selectedPlanTier.value = queryPlan
  }
}

function persistOrderId(id: number | string) {
  orderId.value = id
  sessionStorage.setItem(PENDING_ORDER_KEY, String(id))
}

function restoreOrderId() {
  const cached = sessionStorage.getItem(PENDING_ORDER_KEY)
  if (cached) orderId.value = cached
}

function resolveOrderId(res: { orderId?: number | string; id?: number | string }) {
  const id = res.orderId ?? res.id
  if (!id) throw new Error('创建订单失败：未返回订单 ID')
  persistOrderId(id)
}

function buildAcknowledgments(): RiskAcknowledgments {
  const items = riskItems.value
  return {
    complianceNotEvasion: !!items[0]?.checked,
    noAuditGuarantee: !!items[1]?.checked,
    noFakeInvoice: !!items[2]?.checked,
    taxDependsOnReality: !!items[3]?.checked,
  }
}

async function ensureOrder() {
  if (orderId.value) return
  const plan = selectedPlan.value
  if (!plan?.id) {
    throw new Error('请选择有效套餐')
  }
  const res = await createComplianceOrder({
    planId: plan.id,
    diagnosisId: diagnosisId.value ? Number(diagnosisId.value) : undefined,
  })
  resolveOrderId(res)
}

async function restorePendingOrder() {
  try {
    const order = await getActiveOrder()
    if (!order || order.status !== 'pending') return

    if (order.orderId) persistOrderId(order.orderId)
    if (order.planTier && displayPlans.value.some(p => p.tier === order.planTier)) {
      selectedPlanTier.value = order.planTier
    } else if (order.planId) {
      const matched = displayPlans.value.find(p => String(p.id) === String(order.planId))
      if (matched) selectedPlanTier.value = matched.tier
    }
    currentStep.value = 2
    pendingOrderHint.value = '检测到未完成的签约订单，请继续完成风险告知与电子签约。'
  } catch {
    /* ignore */
  }
}

async function nextStep() {
  apiError.value = ''
  submitting.value = true
  try {
    if (currentStep.value === 1) {
      await ensureOrder()
    } else if (currentStep.value === 2) {
      await ensureOrder()
      await submitConsent({
        orderId: orderId.value,
        type: 'risk_disclosure',
        acknowledgments: buildAcknowledgments(),
      })
      await submitConsent({
        orderId: orderId.value,
        type: 'plan_confirm',
        planConfirmed: true,
      })
    }
    currentStep.value++
    pendingOrderHint.value = ''
  } catch (e: unknown) {
    const msg = sanitizeErrorMessage(e instanceof Error ? e.message : '操作失败，请稍后重试')
    apiError.value = msg
    ElMessage.error(msg)
  } finally {
    submitting.value = false
  }
}

function prevStep() {
  if (currentStep.value > 1) currentStep.value--
}

async function handleSign() {
  if (!canSign.value) return
  apiError.value = ''
  submitting.value = true
  try {
    await ensureOrder()
    await signAgreement({
      orderId: orderId.value,
      legalName: signerName.value,
    })
    sessionStorage.removeItem(PENDING_ORDER_KEY)
    ElMessage.success('签约成功，即将进入主体设立流程')
    router.push('/user/compliance/opc')
  } catch (e: unknown) {
    const msg = sanitizeErrorMessage(e instanceof Error ? e.message : '签约失败，请检查信息后重试')
    apiError.value = msg
    ElMessage.error(msg)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  if (!requireLogin({ redirect: route.fullPath })) return
  restoreOrderId()
  await Promise.all([loadPlans(), loadDiagnosis()])

  try {
    const state = await getCompliancePlanState()
    if (state.hasActiveOrder) {
      router.replace('/user/compliance/opc')
      return
    }
  } catch { /* ignore */ }

  await restorePendingOrder()
  nextTick(() => onPdfScroll())
})
</script>

<style lang="scss" scoped>
.plan-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.overview-panel {
  padding: 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.overview-panel__head {
  margin-bottom: 20px;
}

.overview-panel__section-title {
  margin: 0 0 6px;
  font-size: 18px;
  font-weight: 700;
  color: #1a1f36;
}

.overview-panel__section-desc {
  margin: 0;
  font-size: 13px;
  color: #94a3b8;
}

.plan-steps {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.plan-step {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  background: #f8fafc;

  &.is-done {
    border-color: #bfdbfe;
    background: #eff6ff;

    .plan-step__index {
      background: #2563eb;
      color: #fff;
    }
  }

  &.is-current {
    border-color: #2563eb;
    box-shadow: 0 0 0 1px #2563eb;
  }
}

.plan-step__index {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: #e2e8f0;
  color: #64748b;
  font-size: 13px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
}

.plan-step__label {
  font-size: 13px;
  font-weight: 700;
  color: #334155;
}

.plan-alert {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 13px;
  font-weight: 600;

  p { margin: 0; }

  &--error {
    background: #fff7ed;
    border: 1px solid #fed7aa;
    color: #c2410c;
  }

  &--info {
    background: #eff6ff;
    border: 1px solid #bfdbfe;
    color: #1e40af;
  }
}

.plan-body {
  margin-bottom: 24px;
}

.plan-body__hint {
  margin: 0 0 14px;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
}

.overview-empty {
  padding: 32px 16px;
  text-align: center;

  p {
    margin: 0;
    font-size: 14px;
    color: #64748b;
    font-weight: 600;
  }
}

.overview-empty__icon {
  font-size: 40px;
  color: #2563eb;
  margin-bottom: 16px;
}

.plan-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.plan-card {
  padding: 16px;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  background: #f8fafc;
  cursor: pointer;
  transition: border-color 0.15s ease, background 0.15s ease;

  &:hover {
    border-color: #bfdbfe;
  }

  &.is-selected {
    border-color: #2563eb;
    background: #eff6ff;
    box-shadow: 0 0 0 1px #2563eb;
  }
}

.plan-card__head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.plan-card__title {
  margin: 0 0 4px;
  font-size: 15px;
  font-weight: 700;
  color: #1a1f36;
}

.plan-card__price {
  margin: 0;
  font-size: 18px;
  font-weight: 800;
  color: #2563eb;
}

.plan-card__unit {
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;
}

.plan-card__radio {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid #cbd5e1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 14px;
  color: #fff;

  &.is-checked {
    border-color: #2563eb;
    background: #2563eb;
  }
}

.plan-card__features {
  margin: 12px 0 0;
  padding: 0;
  list-style: none;

  li {
    font-size: 12px;
    color: #64748b;
    line-height: 1.6;

    &::before {
      content: '· ';
    }
  }
}

.plan-risk-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.plan-risk-item {
  padding: 12px 14px;
  border: 1px solid #e8edf3;
  border-radius: 8px;
  background: #f8fafc;
}

.plan-subtitle {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 700;
  color: #475569;
}

.plan-summary {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 24px;
  padding: 16px;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  background: #f8fafc;
}

.plan-field__label {
  display: block;
  margin-bottom: 4px;
  font-size: 11px;
  font-weight: 600;
  color: #94a3b8;
}

.plan-field__value {
  font-size: 14px;
  font-weight: 600;
  color: #1a1f36;
}

.plan-agreement {
  height: 200px;
  overflow-y: auto;
  padding: 16px;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  background: #f8fafc;
  font-size: 13px;
  line-height: 1.7;
  color: #64748b;

  h4 {
    margin: 0 0 12px;
    font-size: 14px;
    font-weight: 700;
    color: #1a1f36;
  }

  p {
    margin: 0 0 8px;
  }

  &__note {
    color: #94a3b8;
    font-size: 12px;
  }
}

.plan-scroll-hint {
  margin: 8px 0 0;
  font-size: 12px;
  font-weight: 600;
  color: #ea580c;
}

.plan-sign-form {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.plan-sign-row {
  padding: 12px 14px;
  border: 1px solid #e8edf3;
  border-radius: 8px;
  background: #fff;
}

.plan-label {
  display: block;
  margin-bottom: 10px;
  font-size: 13px;
  font-weight: 700;
  color: #334155;
}

.plan-required {
  color: #ef4444;
}

:deep(.plan-checkbox) {
  height: auto;
  align-items: flex-start;

  .el-checkbox__label {
    font-size: 13px;
    font-weight: 600;
    color: #334155;
    line-height: 1.6;
    white-space: normal;
  }
}

:deep(.plan-input) {
  .el-input__wrapper {
    border-radius: 8px;
    box-shadow: 0 0 0 1px #d8dee9 inset;
    background: #fff;
  }

  .el-input__inner {
    font-weight: 600;
    color: #1a1f36;
  }
}

.plan-actions {
  display: flex;
  gap: 12px;
  padding-top: 4px;
  border-top: 1px solid #f1f5f9;
}

.plan-actions__main {
  flex: 1;
}

.overview-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  border: none;
  cursor: pointer;
  transition: all 0.15s ease;

  &--primary {
    color: #fff;
    background: #2563eb;

    &:hover:not(:disabled) {
      background: #1d4ed8;
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }

  &--ghost {
    color: #334155;
    background: #fff;
    border: 1px solid #d8dee9;

    &:hover:not(:disabled) {
      background: #f8fafc;
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
}

.plan-loading {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 640px) {
  .plan-steps {
    grid-template-columns: 1fr;
  }

  .plan-summary {
    grid-template-columns: 1fr;
  }
}
</style>
