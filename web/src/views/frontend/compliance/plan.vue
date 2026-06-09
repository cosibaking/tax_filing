<!-- P-05 方案确认与签约 -->
<template>
  <section class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-12">        <!-- Header -->
        <div class="text-center mb-10">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/50 border border-white/50 shadow-sm mb-6">
            <ArtSvgIcon icon="ri:file-text-line" class="text-lg text-clay-accent" />
            <span class="text-sm font-bold text-clay-muted">方案确认与签约</span>
          </div>
          <h1 class="font-heading font-black text-3xl text-clay-foreground mb-2">完成服务签约</h1>
          <p class="text-clay-muted font-medium">选择套餐 → 风险告知 → 电子签约</p>
        </div>

        <!-- Step indicator -->
        <div class="flex items-center justify-between mb-10 px-2">
          <div v-for="(label, idx) in stepLabels" :key="idx" class="flex flex-col items-center flex-1">
            <div
              class="w-10 h-10 rounded-2xl flex items-center justify-center font-black text-sm transition-all duration-300"
              :class="currentStep > idx + 1
                ? 'bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn'
                : currentStep === idx + 1
                  ? 'bg-white text-clay-accent shadow-clay-card border-2 border-clay-accent'
                  : 'bg-[#f0f3f8] text-clay-muted shadow-clay-pressed'"
            >
              <ArtSvgIcon v-if="currentStep > idx + 1" icon="ri:check-line" class="text-lg" />
              <span v-else>{{ idx + 1 }}</span>
            </div>
            <span class="text-[10px] font-bold mt-2 text-center" :class="currentStep === idx + 1 ? 'text-clay-accent' : 'text-clay-muted'">{{ label }}</span>
          </div>
        </div>

        <div v-if="apiError" class="mb-6 p-4 rounded-2xl bg-orange-50 border border-orange-200 flex items-start gap-3">
          <ArtSvgIcon icon="ri:error-warning-line" class="text-xl text-orange-500 shrink-0 mt-0.5" />
          <p class="text-sm font-bold text-orange-700">{{ apiError }}</p>
        </div>

        <div v-if="pendingOrderHint" class="mb-6 p-4 rounded-2xl bg-blue-50 border border-blue-100 flex items-start gap-3">
          <ArtSvgIcon icon="ri:information-line" class="text-xl text-clay-accent shrink-0 mt-0.5" />
          <p class="text-sm font-medium text-clay-foreground">{{ pendingOrderHint }}</p>
        </div>
        <!-- Step 1: 选择套餐 -->
        <div v-show="currentStep === 1" class="space-y-4">
          <div v-if="plansLoading" class="text-center py-12">
            <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
          </div>
          <div v-else class="space-y-4">
            <div
              v-for="plan in displayPlans"
              :key="plan.tier"
              class="p-5 rounded-2xl border-2 cursor-pointer transition-all"
              :class="selectedPlanTier === plan.tier
                ? 'border-clay-accent bg-blue-50/50 shadow-clay-card'
                : 'border-transparent bg-[#f0f3f8] shadow-clay-pressed hover:shadow-clay-card'"
              @click="selectedPlanTier = plan.tier"
            >
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="font-heading font-black text-lg text-clay-foreground">{{ plan.name }}</h3>
                  <p class="text-clay-accent font-black text-xl mt-1">{{ plan.priceLabel }}<span v-if="plan.monthlyPrice" class="text-sm text-clay-muted ml-1">/月起</span></p>
                </div>
                <div
                  class="w-6 h-6 rounded-full border-2 flex items-center justify-center"
                  :class="selectedPlanTier === plan.tier ? 'border-clay-accent bg-clay-accent' : 'border-clay-muted'"
                >
                  <ArtSvgIcon v-if="selectedPlanTier === plan.tier" icon="ri:check-line" class="text-white text-sm" />
                </div>
              </div>
              <ul class="mt-3 space-y-1">
                <li v-for="(feat, i) in plan.features.slice(0, 3)" :key="i" class="text-xs text-clay-muted font-medium">· {{ feat }}</li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Step 2: 风险告知 -->
        <div v-show="currentStep === 2" class="space-y-4">
          <p class="text-sm font-bold text-clay-foreground mb-2">请仔细阅读以下风险告知，全部勾选后方可继续：</p>
          <ElCheckbox v-for="(item, idx) in riskItems" :key="idx" v-model="item.checked" class="!mr-0 clay-checkbox w-full">
            <span class="text-sm font-medium text-clay-foreground leading-relaxed">{{ item.text }}</span>
          </ElCheckbox>
        </div>

        <!-- Step 3: 签约 -->
        <div v-show="currentStep === 3" class="space-y-6">
          <!-- 方案摘要 -->
          <div class="p-5 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed">
            <h3 class="font-heading font-black text-sm text-clay-muted mb-3">方案摘要</h3>
            <div class="space-y-2 text-sm font-medium text-clay-foreground">
              <p>服务套餐：<strong>{{ selectedPlan?.name }}</strong>（{{ selectedPlan?.priceLabel }}）</p>
              <p v-if="diagnosisId">关联诊断：#{{ diagnosisId }}</p>
              <p>推荐方案：OPC 一人有限责任公司</p>
            </div>
          </div>

          <!-- PDF 预览占位 -->
          <div>
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">服务协议预览</label>
            <div
              ref="pdfContainerRef"
              class="h-48 overflow-y-auto rounded-2xl bg-[#f0f3f8] shadow-clay-pressed p-5 text-xs text-clay-muted leading-relaxed"
              @scroll="onPdfScroll"
            >
              <h4 class="font-black text-clay-foreground text-sm mb-3">OPC 合规服务协议（摘要）</h4>
              <p class="mb-2">第一条 服务内容：本协议约定乙方为甲方提供 OPC 设立代办、记账、申报等合规服务。</p>
              <p class="mb-2">第二条 服务费用：按所选套餐标准收取，具体以订单确认金额为准。</p>
              <p class="mb-2">第三条 甲方义务：如实提供注册资料，配合工商、税务、银行开户流程。</p>
              <p class="mb-2">第四条 合规声明：本服务为合法合规方案，不提供逃税、虚开发票等违法服务。</p>
              <p class="mb-2">第五条 责任限制：乙方不承诺「包不被查」，税负取决于甲方真实收入与成本。</p>
              <p class="mb-2">第六条 隐私保护：甲方敏感信息加密存储，仅用于合规服务目的。</p>
              <p class="mb-2">第七条 协议生效：甲方完成电子签名确认后本协议生效。</p>
              <p>（完整协议 PDF 将在签约后存档）</p>
            </div>
            <p v-if="!pdfScrolledToBottom" class="text-xs text-orange-500 font-bold mt-2 ml-1">请滚动阅读完整协议</p>
          </div>

          <ElCheckbox v-model="agreementAccepted" class="clay-checkbox">
            <span class="text-sm font-bold text-clay-foreground">我已阅读并同意服务协议</span>
          </ElCheckbox>

          <ElFormItem class="!mb-0">
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">电子签名（输入姓名确认） <span class="text-red-400">*</span></label>
            <ElInput v-model.trim="signerName" placeholder="请输入您的真实姓名" size="large" class="clay-input" />
          </ElFormItem>
        </div>

        <!-- Actions -->
        <div class="flex gap-4 mt-10">
          <button
            v-if="currentStep > 1"
            type="button"
            class="flex-1 h-14 rounded-2xl bg-white text-clay-foreground font-black shadow-clay-btn hover:shadow-clay-btn-hover active:scale-95 transition-all"
            :disabled="submitting"
            @click="prevStep"
          >
            上一步
          </button>
          <button
            v-if="currentStep < 3"
            type="button"
            class="flex-1 h-14 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black shadow-clay-btn hover:shadow-clay-btn-hover active:scale-95 transition-all disabled:opacity-50"
            :disabled="!canNext || submitting"
            @click="nextStep"
          >
            下一步
          </button>
          <button
            v-else
            type="button"
            class="flex-1 h-14 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black shadow-clay-btn hover:shadow-clay-btn-hover active:scale-95 transition-all disabled:opacity-50 flex items-center justify-center gap-2"
            :disabled="!canSign || submitting"
            @click="handleSign"
          >
            <ArtSvgIcon v-if="submitting" icon="ri:loader-4-line" class="text-xl animate-spin" />
            {{ submitting ? '签约中...' : '完成签约' }}
          </button>
        </div>
  </section>
</template>
<script setup lang="ts">
import { fetchServicePlans, type ServicePlan } from '@/api/frontend/compliance/diagnosis'
import {
  createComplianceOrder,
  getActiveOrder,
  submitConsent,
  signAgreement,
  type RiskAcknowledgments
} from '@/api/frontend/compliance/order'
import { getCompliancePlanState } from '@/api/frontend/compliance/member'
import { useMemberStore } from '@/store/modules/member'
import { requireLogin } from '@/utils/auth/requireLogin'
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

const FALLBACK_PLANS: ServicePlan[] = [
  { id: 1, name: '基础版', tier: 'basic', monthlyPrice: 299, priceLabel: '¥299', features: ['OPC 注册代办', '月度记账', '季度申报'] },
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
    const msg = e instanceof Error ? e.message : '操作失败，请稍后重试'
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
    const res = await signAgreement({
      orderId: orderId.value,
      legalName: signerName.value,
    })
    sessionStorage.removeItem(PENDING_ORDER_KEY)
    ElMessage.success('签约成功，即将进入 OPC 落地流程')
    if (res.opcId) {
      router.push('/user/compliance/opc')
    } else {
      router.push('/user/compliance/opc')
    }
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : '签约失败，请检查信息后重试'
    apiError.value = msg
    ElMessage.error(msg)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  if (!requireLogin({ redirect: route.fullPath })) return
  restoreOrderId()
  await loadPlans()

  try {
    const state = await getCompliancePlanState()
    if (state.hasActiveOrder) {
      router.replace('/user/compliance/opc')
      return
    }
  } catch { /* ignore */ }

  await restorePendingOrder()
  nextTick(() => onPdfScroll())
})</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.border-clay-accent { border-color: #5a8dee; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.shadow-clay-deep {
  box-shadow: 30px 30px 60px #d1d9e6, -30px -30px 60px #ffffff,
    inset 10px 10px 20px rgba(90, 141, 238, 0.05), inset -10px -10px 20px rgba(255, 255, 255, 0.8);
}
.shadow-clay-card {
  box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9),
    inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
}
.shadow-clay-btn {
  box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-btn-hover {
  box-shadow: 16px 16px 32px rgba(90, 141, 238, 0.4), -10px -10px 20px rgba(255, 255, 255, 0.5),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-pressed {
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
}

:deep(.clay-input) {
  .el-input__wrapper {
    height: 48px; padding: 0 16px; border-radius: 16px; background: #f0f3f8;
    box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
    border: none;
  }
  .el-input__inner { font-weight: 500; color: #32325d; }
}

:deep(.clay-checkbox) {
  padding: 12px 16px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 6px 6px 12px #e0e5ec, inset -6px -6px 12px #ffffff;
  margin-right: 0 !important; height: auto; width: 100%;
}
</style>
