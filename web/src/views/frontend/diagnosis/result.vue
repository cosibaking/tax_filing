<!-- +----------------------------------------------------------------------
  | XYGo Admin — 合规诊断结果 P-03
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-8 pb-16 px-6">
    <div class="max-w-6xl mx-auto space-y-8">
      <!-- Header -->
      <div class="text-center">
        <h1 class="font-heading font-black text-3xl md:text-4xl text-clay-foreground mb-3">您的合规诊断结果</h1>
        <p class="text-clay-muted font-medium">基于您提供的信息，以下为各方案税负估算对比</p>
      </div>

      <!-- Loading / Error -->
      <div v-if="loading" class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-16 text-center">
        <ArtSvgIcon icon="ri:loader-4-line" class="text-4xl text-clay-accent animate-spin mx-auto mb-4" />
        <p class="text-clay-muted font-bold">加载诊断结果...</p>
      </div>

      <div v-else-if="!diagnosisResult" class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-16 text-center">
        <ArtSvgIcon icon="ri:file-search-line" class="text-5xl text-clay-muted mx-auto mb-4" />
        <p class="text-clay-muted font-bold mb-6">未找到诊断结果，请重新填写问卷</p>
        <RouterLink
          to="/diagnosis"
          class="inline-block px-8 py-3 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black shadow-clay-btn hover:shadow-clay-btn-hover transition-all"
        >
          开始诊断
        </RouterLink>
      </div>

      <template v-else>
        <!-- Four-column comparison -->
        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-6">
          <div
            v-for="col in comparisonColumns"
            :key="col.plan"
            class="relative rounded-[32px] p-6 transition-all duration-300 border"
            :class="col.cardClass"
          >
            <div v-if="col.recommended" class="absolute -top-3 left-1/2 -translate-x-1/2 px-4 py-1 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 text-white text-xs font-black shadow-clay-btn whitespace-nowrap">
              ⭐ 推荐方案
            </div>
            <div v-if="col.warning" class="absolute -top-3 left-1/2 -translate-x-1/2 px-4 py-1 rounded-full bg-orange-500 text-white text-xs font-black shadow-clay-btn whitespace-nowrap">
              ⚠️ 高危
            </div>

            <h3 class="font-heading font-black text-lg text-clay-foreground mb-1 mt-2">{{ col.label }}</h3>
            <p v-if="col.warning" class="text-xs font-bold text-orange-600 mb-4">潜在补税 + 罚款风险</p>
            <p v-else class="text-xs text-clay-muted mb-4">{{ col.note || '年度估算' }}</p>

            <div class="mb-4">
              <span class="font-heading font-black text-3xl" :class="col.warning ? 'text-orange-500' : 'text-clay-foreground'">
                {{ col.warning ? '—' : formatMoney(col.annualTax) }}
              </span>
              <span v-if="!col.warning" class="text-sm text-clay-muted ml-1">元/年</span>
            </div>
            <div v-if="!col.warning" class="inline-flex px-3 py-1 rounded-full bg-white/60 text-sm font-bold" :class="col.recommended ? 'text-clay-accent' : 'text-clay-muted'">
              综合税负 {{ formatRate(col.taxRate) }}
            </div>
          </div>
        </div>

        <!-- Recommendation summary -->
        <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-10">
          <h2 class="font-heading font-black text-xl text-clay-foreground mb-6 flex items-center gap-2">
            <ArtSvgIcon icon="ri:lightbulb-flash-line" class="text-2xl text-clay-accent" />
            方案推荐摘要
          </h2>
          <ul class="space-y-3">
            <li v-for="(reason, idx) in recommendationReasons" :key="idx" class="flex items-start gap-3">
              <div class="w-6 h-6 rounded-full bg-clay-success/20 text-clay-success flex items-center justify-center shrink-0 mt-0.5">
                <ArtSvgIcon icon="ri:check-line" class="text-sm" />
              </div>
              <span class="font-medium text-clay-foreground leading-relaxed">{{ reason }}</span>
            </li>
          </ul>
        </div>

        <!-- Tax calculator -->
        <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-10">
          <h2 class="font-heading font-black text-xl text-clay-foreground mb-6 flex items-center gap-2">
            <ArtSvgIcon icon="ri:calculator-line" class="text-2xl text-clay-accent" />
            税负计算器
          </h2>
          <div class="grid md:grid-cols-2 gap-6 mb-6">
            <div>
              <label class="block text-sm font-bold text-clay-foreground mb-3">年收入（元）</label>
              <ElInput v-model.number="calcIncome" type="number" size="large" class="clay-input" />
            </div>
            <div>
              <label class="block text-sm font-bold text-clay-foreground mb-3">年可扣除成本（元）</label>
              <ElInput v-model.number="calcCost" type="number" size="large" class="clay-input" />
            </div>
          </div>
          <button
            type="button"
            class="px-8 py-3 rounded-2xl bg-white shadow-clay-btn hover:shadow-clay-btn-hover font-bold text-clay-foreground active:scale-95 transition-all flex items-center gap-2"
            :disabled="calculating"
            @click="handleRecalculate"
          >
            <ArtSvgIcon v-if="calculating" icon="ri:loader-4-line" class="text-lg animate-spin" />
            {{ calculating ? '计算中...' : '重新计算' }}
          </button>

          <div
            v-if="assumptionHints.length"
            class="mt-6 p-5 md:p-6 rounded-2xl bg-blue-50/70 border border-blue-100/80"
          >
            <h3 class="text-sm font-black text-clay-foreground mb-3 flex items-center gap-2">
              <ArtSvgIcon icon="ri:information-line" class="text-lg text-clay-accent" />
              测算想定提示
            </h3>
            <ul class="space-y-2">
              <li
                v-for="(hint, idx) in assumptionHints"
                :key="idx"
                class="text-sm text-clay-muted leading-relaxed flex items-start gap-2"
              >
                <span class="text-clay-accent font-bold shrink-0">·</span>
                <span>{{ hint }}</span>
              </li>
            </ul>
          </div>
        </div>

        <!-- CTA -->
        <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
          <button
            v-if="!isGuestResult"
            type="button"
            class="w-full sm:w-auto px-8 py-4 rounded-2xl bg-white shadow-clay-btn hover:shadow-clay-btn-hover font-bold text-clay-foreground transition-all flex items-center justify-center gap-2"
            @click="handleDownloadPdf"
          >
            <ArtSvgIcon icon="ri:download-line" class="text-lg" />
            下载方案摘要 PDF
          </button>
          <button
            type="button"
            class="w-full sm:w-auto px-10 py-4 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black text-lg shadow-clay-btn hover:shadow-clay-btn-hover hover:-translate-y-1 active:scale-95 transition-all flex items-center justify-center gap-2"
            @click="handleSignUp"
          >
            确认方案并签约
            <ArtSvgIcon icon="ri:arrow-right-line" class="text-lg" />
          </button>
        </div>

        <!-- Disclaimer -->
        <p class="text-center text-xs text-clay-muted font-medium leading-relaxed max-w-2xl mx-auto">
          以上税负数据均为估算值，实际税负取决于真实收入、成本及税收政策变化。本服务为合规方案，非逃税方案。
        </p>
      </template>
    </div>
  </main>
</template>

<script setup lang="ts">
import {
  calculateTax,
  getCachedDiagnosisResult,
  cacheDiagnosisResult,
  type DiagnosisSubmitResult
} from '@/api/frontend/compliance/diagnosis'
import { requireLogin } from '@/utils/auth/requireLogin'
import { getDiagnosisDetail } from '@/api/frontend/compliance/member'
import { useMemberStore } from '@/store/modules/member'
import {
  getGuestDiagnosis,
  isGuestDiagnosisId,
  saveGuestDiagnosis
} from '@/utils/compliance/guestDiagnosisCache'

defineOptions({ name: 'ComplianceDiagnosisResult' })

const router = useRouter()
const route = useRoute()
const memberStore = useMemberStore()
const loading = ref(true)
const calculating = ref(false)
const diagnosisResult = ref<DiagnosisSubmitResult | null>(null)
const assumptionHints = ref<string[]>([])
const calcIncome = ref(0)
const calcCost = ref(0)
const isGuestResult = computed(() => isGuestDiagnosisId(diagnosisResult.value?.id))

interface ComparisonColumnDef {
  plan: string
  label: string
  warning?: boolean
}

const COLUMN_DEFS: ComparisonColumnDef[] = [
  { plan: 'none', label: '不报税', warning: true },
  { plan: 'labor', label: '纯劳务' },
  { plan: 'individual', label: '个体户' },
  { plan: 'opc', label: 'OPC' }
]

function isPlanRecommended(plan: string, recommended?: string, itemRecommended?: boolean) {
  if (itemRecommended) return true
  if (!recommended || recommended === 'none') return false
  if (recommended === plan) return true
  if (recommended === 'transitional' && plan === 'individual') return true
  return false
}

const comparisonColumns = computed(() => {
  const items = diagnosisResult.value?.taxComparison?.items || []
  const recommended = diagnosisResult.value?.recommendedPlan

  return COLUMN_DEFS.map((def) => {
    const item = items.find((i) => i.plan === def.plan)
    const isRecommended = isPlanRecommended(def.plan, recommended, item?.recommended)
    return {
      ...def,
      annualTax: item?.annualTax ?? 0,
      taxRate: item?.taxRate ?? 0,
      note: item?.note,
      recommended: isRecommended,
      cardClass: def.warning
        ? 'bg-orange-50/80 border-orange-200 shadow-clay-card'
        : isRecommended
          ? 'bg-white/90 border-blue-200 shadow-clay-deep ring-2 ring-blue-200/50'
          : 'bg-white/70 border-[#d1d9e6]/40 shadow-clay-card'
    }
  })
})

const recommendationReasons = computed(() => {
  if (diagnosisResult.value?.reasons?.length) {
    return diagnosisResult.value.reasons
  }
  const plan = diagnosisResult.value?.recommendedPlan
  const defaults: Record<string, string[]> = {
    opc: [
      '年收入规模适合通过 OPC 小微公司合规经营，综合税负更优',
      '可合法抵扣成本费用，降低应纳税所得额',
      '便于与平台/MCN 对公结算，降低个人账户大额流水风险'
    ],
    individual: [
      '当前收入规模适合个体户查账征收，设立成本较低',
      '经营所得可扣除真实成本，税负低于劳务报酬',
      '后续收入增长可升级 OPC 方案'
    ],
    labor: [
      '当前收入较低，劳务报酬计税相对简单',
      '建议保留完整收入凭证，为后续升级做准备',
      '收入增长后建议重新评估个体户或 OPC 方案'
    ],
    transitional: [
      '建议先完成历史申报补正，再选择合适经营主体',
      '过渡期内可借助顾问梳理合规路径',
      '稳定经营后推荐升级 OPC 方案'
    ]
  }
  return defaults[plan || 'opc'] || defaults.opc
})

function formatMoney(val: number) {
  if (val >= 10000) return `${(val / 10000).toFixed(1)}万`
  return val.toLocaleString('zh-CN', { maximumFractionDigits: 0 })
}

function formatRate(rate: number) {
  return `${(rate * 100).toFixed(1)}%`
}

function loadResult() {
  const id = route.query.id as string
  if (!id) {
    loading.value = false
    return
  }

  const state = history.state?.diagnosisResult as DiagnosisSubmitResult | undefined
  if (state?.id) {
    diagnosisResult.value = state
    cacheDiagnosisResult(id, state)
    finishLoadResult()
    return
  }

  if (isGuestDiagnosisId(id)) {
    const guest = getGuestDiagnosis(id)
    diagnosisResult.value = guest?.result ?? null
    finishLoadResult()
    return
  }

  diagnosisResult.value = getCachedDiagnosisResult(id)
  if (diagnosisResult.value) {
    finishLoadResult()
    return
  }

  if (memberStore.getIsLogin && /^\d+$/.test(id)) {
    loadResultFromServer(Number(id))
    return
  }

  finishLoadResult()
}

async function loadResultFromServer(id: number) {
  try {
    const res = await getDiagnosisDetail(id)
    if (res?.id) {
      diagnosisResult.value = res
      cacheDiagnosisResult(id, res)
    }
  } catch {
    // 拦截器已处理
  } finally {
    finishLoadResult()
  }
}

function finishLoadResult() {
  if (diagnosisResult.value?.taxComparison) {
    calcIncome.value = diagnosisResult.value.taxComparison.annualIncome || 0
    calcCost.value = diagnosisResult.value.taxComparison.annualCost || 0
  }
  assumptionHints.value = diagnosisResult.value?.assumptionHints || []
  loading.value = false
}

async function handleRecalculate() {
  if (!diagnosisResult.value) return
  calculating.value = true
  try {
    const res = await calculateTax({
      annualIncome: calcIncome.value,
      annualCost: calcCost.value,
      diagnosisId: isGuestResult.value ? undefined : diagnosisResult.value.id
    })
    if (res.taxComparison) {
      diagnosisResult.value = {
        ...diagnosisResult.value,
        taxComparison: res.taxComparison,
        recommendedPlan: res.recommendedPlan || diagnosisResult.value.recommendedPlan,
        reasons: res.reasons?.length ? res.reasons : diagnosisResult.value.reasons,
        assumptionHints: res.assumptionHints
      }
      assumptionHints.value = res.assumptionHints || []
      const id = String(diagnosisResult.value.id)
      cacheDiagnosisResult(id, diagnosisResult.value)
      if (isGuestResult.value) {
        const guest = getGuestDiagnosis(id)
        if (guest) {
          saveGuestDiagnosis(id, guest.payload, diagnosisResult.value)
        }
      }
    }
  } catch {
    // 错误由拦截器处理
  } finally {
    calculating.value = false
  }
}

function handleDownloadPdf() {
  if (isGuestResult.value) return
  const id = diagnosisResult.value?.id
  if (!id) return
  const base = import.meta.env.VITE_API_URL || ''
  window.open(`${base}/site/compliance/diagnosis/${id}/pdf`, '_blank')
}

function handleSignUp() {
  const diagnosisId = diagnosisResult.value?.id
  const query: Record<string, string> = {}
  if (diagnosisId) query.diagnosisId = String(diagnosisId)
  const target = `/user/compliance/plan${diagnosisId ? `?diagnosisId=${diagnosisId}` : ''}`

  if (!requireLogin({ redirect: target })) return
  router.push({ path: '/user/compliance/plan', query })
}

onMounted(loadResult)
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
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

:deep(.clay-input) {
  .el-input__wrapper {
    height: 48px; padding: 0 16px; border-radius: 16px; background: #f0f3f8;
    box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
    border: none;
    &.is-focus {
      background: #fff;
      box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9);
    }
  }
  .el-input__inner { font-weight: 500; color: #32325d; }
}
</style>
