<!-- +----------------------------------------------------------------------
  | 金税管家 — 免费合规诊断问卷
  +---------------------------------------------------------------------- -->
<template>
  <main class="diag-page">
    <div class="diag-wrap">
      <div class="diag-card">
        <!-- Header -->
        <header class="diag-header">
          <span class="diag-badge">
            <ArtSvgIcon icon="ri:shield-check-line" class="text-base" />
            约 3 分钟 · 免费
          </span>
          <h1 class="diag-title">金税管家合规诊断</h1>
          <p class="diag-subtitle">
            根据您的收入结构与合规现状，生成劳务 / 个体户 / 小微公司税负对比与个性化建议
          </p>
        </header>

        <!-- Progress -->
        <div class="diag-progress">
          <div class="diag-progress__bar">
            <div class="diag-progress__fill" :style="{ width: `${progressPercent}%` }" />
          </div>
          <div class="diag-progress__meta">
            <span>第 {{ currentStep }} / {{ totalSteps }} 步</span>
            <span>{{ currentStepMeta.title }}</span>
          </div>
        </div>

        <!-- API Error -->
        <div v-if="apiError" class="diag-error">
          <ArtSvgIcon icon="ri:error-warning-line" class="text-lg shrink-0" />
          <p>{{ apiError }}</p>
          <button type="button" @click="apiError = ''">关闭</button>
        </div>

        <ElForm ref="formRef" :model="formData" :rules="rules" class="diag-form">
          <!-- Step 1: 身份与渠道 -->
          <section v-show="currentStep === 1" class="diag-step">
            <div class="diag-step-head">
              <h2>{{ currentStepMeta.title }}</h2>
              <p>{{ currentStepMeta.subtitle }}</p>
            </div>

            <ElFormItem prop="persona" class="diag-field">
              <label class="diag-label">您的从业类型 <span class="req">*</span></label>
              <div class="diag-card-grid diag-card-grid--3">
                <button
                  v-for="item in personaOptions"
                  :key="item.value"
                  type="button"
                  class="diag-option"
                  :class="{ 'is-active': formData.persona === item.value }"
                  @click="formData.persona = item.value"
                >
                  <ArtSvgIcon :icon="item.icon" class="diag-option__icon" />
                  <span class="diag-option__title">{{ item.label }}</span>
                  <span class="diag-option__desc">{{ item.desc }}</span>
                </button>
              </div>
            </ElFormItem>

            <ElFormItem v-if="formData.persona === 'MCN签约主播'" prop="mcnSettlement" class="diag-field">
              <label class="diag-label">MCN 结算方式 <span class="req">*</span></label>
              <p class="diag-hint">影响 OPC 落地节奏与收入台账口径</p>
              <div class="diag-list">
                <button
                  v-for="item in mcnSettlementOptions"
                  :key="item.value"
                  type="button"
                  class="diag-list-item"
                  :class="{ 'is-active': formData.mcnSettlement === item.value }"
                  @click="formData.mcnSettlement = item.value"
                >
                  <span class="diag-list-item__title">{{ item.label }}</span>
                </button>
              </div>
            </ElFormItem>

            <ElFormItem prop="platforms" class="diag-field">
              <label class="diag-label">主要收入渠道 <span class="req">*</span></label>
              <p class="diag-hint">可多选，包含平台、店铺或接单渠道</p>
              <div class="diag-chip-grid">
                <button
                  v-for="item in channelOptions"
                  :key="item.value"
                  type="button"
                  class="diag-chip"
                  :class="{ 'is-active': formData.platforms.includes(item.value) }"
                  @click="toggleArrayItem(formData.platforms, item.value)"
                >
                  <ArtSvgIcon :icon="item.icon" class="text-sm" />
                  {{ item.label }}
                </button>
              </div>
            </ElFormItem>

            <ElFormItem prop="incomeTypes" class="diag-field">
              <label class="diag-label">收入构成 <span class="req">*</span></label>
              <p class="diag-hint">请选择您的主要收入类型，便于判断申报口径</p>
              <div class="diag-chip-grid">
                <button
                  v-for="item in incomeTypeOptions"
                  :key="item.value"
                  type="button"
                  class="diag-chip"
                  :class="{ 'is-active': formData.incomeTypes.includes(item.value) }"
                  @click="toggleArrayItem(formData.incomeTypes, item.value)"
                >
                  <ArtSvgIcon :icon="item.icon" class="text-sm" />
                  {{ item.label }}
                </button>
              </div>
            </ElFormItem>
          </section>

          <!-- Step 2: 收入规模 -->
          <section v-show="currentStep === 2" class="diag-step">
            <div class="diag-step-head">
              <h2>{{ currentStepMeta.title }}</h2>
              <p>{{ currentStepMeta.subtitle }}</p>
            </div>

            <ElFormItem prop="monthlyIncomeRange" class="diag-field">
              <label class="diag-label">月均收入区间 <span class="req">*</span></label>
              <div class="diag-income-grid">
                <button
                  v-for="item in incomeRangeOptions"
                  :key="item.value"
                  type="button"
                  class="diag-income"
                  :class="{ 'is-active': formData.monthlyIncomeRange === item.value }"
                  @click="formData.monthlyIncomeRange = item.value"
                >
                  <span class="diag-income__label">{{ item.label }}</span>
                  <span class="diag-income__desc">{{ item.desc }}</span>
                  <span class="diag-income__hint">{{ item.hint }}</span>
                </button>
              </div>
            </ElFormItem>

            <ElFormItem prop="annualCostEstimate" class="diag-field">
              <label class="diag-label">年度可扣除成本（估算）</label>
              <p class="diag-hint">设备、投流、场地、采购等可入账经营费用合计，可在下一步细化</p>
              <ElInput
                v-model.number="formData.annualCostEstimate"
                placeholder="如 80000"
                size="large"
                type="number"
                class="diag-input"
              >
                <template #prefix>
                  <ArtSvgIcon icon="ri:money-cny-circle-line" class="text-lg text-clay-muted" />
                </template>
                <template #suffix>元/年</template>
              </ElInput>
            </ElFormItem>
          </section>

          <!-- Step 3: 合规现状 -->
          <section v-show="currentStep === 3" class="diag-step">
            <div class="diag-step-head">
              <h2>{{ currentStepMeta.title }}</h2>
              <p>{{ currentStepMeta.subtitle }}</p>
            </div>

            <ElFormItem prop="existingEntity" class="diag-field">
              <label class="diag-label">经营主体情况 <span class="req">*</span></label>
              <div class="diag-list">
                <button
                  v-for="item in entityOptions"
                  :key="item.value"
                  type="button"
                  class="diag-list-item"
                  :class="{ 'is-active': formData.existingEntity === item.value }"
                  @click="formData.existingEntity = item.value"
                >
                  <span class="diag-list-item__title">{{ item.label }}</span>
                  <span class="diag-list-item__desc">{{ item.desc }}</span>
                </button>
              </div>
            </ElFormItem>

            <ElFormItem prop="hasFiledTax" class="diag-field">
              <label class="diag-label">纳税申报情况 <span class="req">*</span></label>
              <div class="diag-list">
                <button
                  v-for="item in taxFiledOptions"
                  :key="item.value"
                  type="button"
                  class="diag-list-item"
                  :class="{ 'is-active': formData.hasFiledTax === item.value }"
                  @click="formData.hasFiledTax = item.value"
                >
                  <span class="diag-list-item__title">{{ item.label }}</span>
                  <span class="diag-list-item__desc">{{ item.desc }}</span>
                </button>
              </div>
            </ElFormItem>
          </section>

          <!-- Step 4: 成本与诉求 -->
          <section v-show="currentStep === 4" class="diag-step">
            <div class="diag-step-head">
              <h2>{{ currentStepMeta.title }}</h2>
              <p>{{ currentStepMeta.subtitle }}</p>
            </div>

            <ElFormItem prop="riskSignals" class="diag-field">
              <label class="diag-label">合规风险信号 <span class="req">*</span></label>
              <div class="diag-risk-grid">
                <button
                  v-for="item in riskSignalOptions"
                  :key="item.value"
                  type="button"
                  class="diag-risk"
                  :class="[
                    `diag-risk--${item.severity}`,
                    { 'is-active': formData.riskSignals.includes(item.value) }
                  ]"
                  @click="toggleRiskSignal(item.value)"
                >
                  <span class="diag-risk__title">{{ item.label }}</span>
                  <span class="diag-risk__desc">{{ item.desc }}</span>
                </button>
              </div>
            </ElFormItem>

            <ElFormItem class="diag-field">
              <label class="diag-label">您最关心的问题</label>
              <div class="diag-chip-grid">
                <button
                  v-for="item in concernOptions"
                  :key="item.value"
                  type="button"
                  class="diag-chip"
                  :class="{ 'is-active': formData.concerns.includes(item.value) }"
                  @click="toggleArrayItem(formData.concerns, item.value)"
                >
                  {{ item.label }}
                </button>
              </div>
            </ElFormItem>

            <div class="diag-cost-block">
              <div class="diag-cost-block__head">
                <label class="diag-label">成本细项（选填）</label>
                <p class="diag-hint">按类目填写可提升税负测算精度</p>
              </div>
              <div class="diag-cost-grid">
                <div v-for="item in costItemOptions" :key="item.key" class="diag-cost-item">
                  <div class="diag-cost-item__top">
                    <span class="diag-cost-item__label">{{ item.label }}</span>
                    <span class="diag-cost-item__hint">{{ item.hint }}</span>
                  </div>
                  <ElInput
                    v-model.number="formData.costBreakdown[item.key]"
                    :placeholder="item.placeholder"
                    size="large"
                    type="number"
                    class="diag-input"
                  >
                    <template #suffix>元/年</template>
                  </ElInput>
                </div>
              </div>
            </div>

            <ElFormItem prop="notes" class="diag-field">
              <label class="diag-label">补充说明</label>
              <ElInput
                v-model="formData.notes"
                type="textarea"
                :rows="3"
                placeholder="如：多平台收入汇总困难、近期收到补税通知等（选填）"
                class="diag-textarea"
              />
            </ElFormItem>
          </section>
        </ElForm>

        <!-- Navigation -->
        <footer class="diag-footer">
          <button
            v-if="currentStep > 1"
            type="button"
            class="diag-btn diag-btn--ghost"
            @click="prevStep"
          >
            上一步
          </button>
          <div v-else />

          <button
            v-if="currentStep < totalSteps"
            type="button"
            class="diag-btn diag-btn--primary"
            @click="nextStep"
          >
            下一步
            <ArtSvgIcon icon="ri:arrow-right-line" />
          </button>
          <button
            v-else
            type="button"
            class="diag-btn diag-btn--primary"
            :disabled="submitting"
            @click="handleSubmit"
          >
            <ArtSvgIcon v-if="submitting" icon="ri:loader-4-line" class="animate-spin" />
            {{ submitting ? '分析中...' : '生成诊断报告' }}
          </button>
        </footer>

        <button
          v-if="currentStep === totalSteps"
          type="button"
          class="diag-skip"
          :disabled="submitting"
          @click="handleSubmit"
        >
          跳过成本细项，直接查看结果
        </button>

        <p class="diag-disclaimer">
          诊断结果基于您填写的信息与现行税收政策测算，仅供参考，不构成税务法律意见。
        </p>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import type { HistoryState } from 'vue-router'
import {
  submitDiagnosis,
  cacheDiagnosisResult,
  type DiagnosisSubmitParams,
  type MonthlyIncomeRange,
  type ExistingEntity,
  type HasFiledTax
} from '@/api/frontend/compliance/diagnosis'
import { useMemberStore } from '@/store/modules/member'
import {
  generateGuestDiagnosisId,
  saveGuestDiagnosis
} from '@/utils/compliance/guestDiagnosisCache'
import {
  diagnosisSteps,
  personaOptions,
  channelOptions,
  incomeTypeOptions,
  incomeRangeOptions,
  entityOptions,
  taxFiledOptions,
  riskSignalOptions,
  concernOptions,
  costItemOptions,
  mcnSettlementOptions
} from '@/data/frontend/diagnosisQuestionnaire'

defineOptions({ name: 'ComplianceDiagnosis' })

const router = useRouter()
const memberStore = useMemberStore()
const formRef = ref<FormInstance>()

const totalSteps = diagnosisSteps.length
const currentStep = ref(1)
const submitting = ref(false)
const apiError = ref('')

const costKeys = costItemOptions.map((item) => item.key)

const formData = reactive({
  persona: '',
  mcnSettlement: '',
  platforms: [] as string[],
  incomeTypes: [] as string[],
  monthlyIncomeRange: '' as MonthlyIncomeRange | '',
  annualCostEstimate: undefined as number | undefined,
  existingEntity: '' as ExistingEntity | '',
  hasFiledTax: '' as HasFiledTax | '',
  riskSignals: [] as string[],
  concerns: [] as string[],
  notes: '',
  costBreakdown: Object.fromEntries(costKeys.map((key) => [key, undefined])) as Record<string, number | undefined>
})

const currentStepMeta = computed(() => diagnosisSteps[currentStep.value - 1])
const progressPercent = computed(() => Math.round((currentStep.value / totalSteps) * 100))

const rules: FormRules = {
  persona: [{ required: true, message: '请选择从业类型', trigger: 'change' }],
  mcnSettlement: [{
    validator: (_r, _v, cb) => {
      if (formData.persona === 'MCN签约主播' && !formData.mcnSettlement) {
        cb(new Error('请选择 MCN 结算方式'))
        return
      }
      cb()
    },
    trigger: 'change'
  }],
  platforms: [{
    type: 'array',
    required: true,
    min: 1,
    message: '请至少选择一个收入渠道',
    trigger: 'change'
  }],
  incomeTypes: [{
    type: 'array',
    required: true,
    min: 1,
    message: '请至少选择一种收入构成',
    trigger: 'change'
  }],
  monthlyIncomeRange: [{ required: true, message: '请选择月均收入区间', trigger: 'change' }],
  existingEntity: [{ required: true, message: '请选择经营主体情况', trigger: 'change' }],
  hasFiledTax: [{ required: true, message: '请选择纳税申报情况', trigger: 'change' }],
  riskSignals: [{
    type: 'array',
    required: true,
    min: 1,
    message: '请选择合规风险信号',
    trigger: 'change'
  }]
}

const stepFields: Record<number, string[]> = {
  1: ['persona', 'mcnSettlement', 'platforms', 'incomeTypes'],
  2: ['monthlyIncomeRange'],
  3: ['existingEntity', 'hasFiledTax'],
  4: ['riskSignals']
}

function toggleArrayItem(list: string[], value: string) {
  const idx = list.indexOf(value)
  if (idx >= 0) list.splice(idx, 1)
  else list.push(value)
}

function toggleRiskSignal(value: string) {
  if (value === 'none') {
    formData.riskSignals = ['none']
    return
  }
  const filtered = formData.riskSignals.filter((item) => item !== 'none')
  toggleArrayItem(filtered, value)
  formData.riskSignals = filtered
}

function buildNotes(): string | undefined {
  const parts: string[] = []
  if (formData.persona) parts.push(`从业类型：${formData.persona}`)
  if (formData.persona === 'MCN签约主播' && formData.mcnSettlement) {
    parts.push(`MCN结算：${formData.mcnSettlement}`)
  }
  if (formData.incomeTypes.length) parts.push(`收入构成：${formData.incomeTypes.join('、')}`)
  if (formData.riskSignals.length) {
    const labels = formData.riskSignals
      .map((value) => riskSignalOptions.find((item) => item.value === value)?.label || value)
      .join('、')
    parts.push(`风险信号：${labels}`)
  }
  if (formData.concerns.length) parts.push(`核心诉求：${formData.concerns.join('、')}`)
  if (formData.notes.trim()) parts.push(`补充说明：${formData.notes.trim()}`)
  return parts.length ? parts.join('\n') : undefined
}

function resolveTaxBureauContact(): boolean {
  return formData.riskSignals.some((item) => item === 'tax_bureau' || item === 'platform_notice')
}

function buildPayload(): DiagnosisSubmitParams {
  const breakdown = Object.fromEntries(
    Object.entries(formData.costBreakdown).filter(([, value]) => value != null && value > 0)
  ) as Record<string, number>
  const breakdownSum = Object.values(breakdown).reduce((sum, value) => sum + value, 0)
  const annualCost = formData.annualCostEstimate || breakdownSum || undefined

  return {
    platforms: formData.platforms,
    monthlyIncomeRange: formData.monthlyIncomeRange as MonthlyIncomeRange,
    annualCostEstimate: annualCost,
    existingEntity: formData.existingEntity as ExistingEntity,
    hasFiledTax: formData.hasFiledTax as HasFiledTax,
    taxBureauContact: resolveTaxBureauContact(),
    notes: buildNotes(),
    costBreakdown: Object.keys(breakdown).length ? breakdown : undefined
  }
}

async function validateStep(step: number): Promise<boolean> {
  const fields = stepFields[step]
  if (!fields?.length || !formRef.value) return true
  try {
    await formRef.value.validateField(fields)
    return true
  } catch {
    return false
  }
}

function scrollPageToTop() {
  nextTick(() => {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  })
}

async function nextStep() {
  const valid = await validateStep(currentStep.value)
  if (!valid) return
  if (currentStep.value < totalSteps) {
    currentStep.value++
    scrollPageToTop()
  }
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
    scrollPageToTop()
  }
}

async function handleSubmit() {
  for (let step = 1; step <= 3; step++) {
    const valid = await validateStep(step)
    if (!valid) {
      currentStep.value = step
      return
    }
  }
  if (!(await validateStep(4))) {
    currentStep.value = 4
    return
  }

  submitting.value = true
  apiError.value = ''
  try {
    const payload = buildPayload()
    const result = await submitDiagnosis(payload)
    const isLoggedIn = memberStore.getIsLogin

    let resultId: string | number = result.id
    if (!isLoggedIn || !result.id) {
      resultId = generateGuestDiagnosisId()
      const guestResult = { ...result, id: resultId }
      saveGuestDiagnosis(resultId, payload, guestResult)
      cacheDiagnosisResult(resultId, guestResult)
      router.push({
        path: '/diagnosis/result',
        query: { id: String(resultId) },
        state: { diagnosisResult: guestResult } as unknown as HistoryState
      })
      return
    }

    cacheDiagnosisResult(result.id, result)
    router.push({
      path: '/diagnosis/result',
      query: { id: String(result.id) },
      state: { diagnosisResult: result } as unknown as HistoryState
    })
  } catch (e: unknown) {
    const msg = (e as { message?: string })?.message
    apiError.value = msg && msg !== 'Network Error' ? msg : '诊断提交失败，请检查网络后重试'
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.diag-page {
  padding: 32px 24px 64px;
}

.diag-wrap {
  max-width: 720px;
  margin: 0 auto;
}

.diag-card {
  padding: 32px 28px 28px;
  border-radius: 32px;
  background: rgba(255, 255, 255, 0.82);
  border: 1px solid rgba(209, 217, 230, 0.5);
  box-shadow:
    24px 24px 48px #d1d9e6,
    -16px -16px 40px #fff;
  backdrop-filter: blur(20px);
}

.diag-header {
  text-align: center;
  margin-bottom: 28px;
}

.diag-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 999px;
  background: #f0f6ff;
  color: #5a8dee;
  font-size: 13px;
  font-weight: 700;
}

.diag-title {
  margin: 16px 0 8px;
  font-size: clamp(24px, 4vw, 30px);
  font-weight: 900;
  color: #1a1f36;
}

.diag-subtitle {
  margin: 0;
  font-size: 14px;
  line-height: 1.7;
  color: #6b7c93;
}

.diag-progress {
  margin-bottom: 28px;
}

.diag-progress__bar {
  height: 6px;
  border-radius: 999px;
  background: #e8edf5;
  overflow: hidden;
}

.diag-progress__fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #5a8dee, #3b6fd9);
  transition: width 0.35s ease;
}

.diag-progress__meta {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 12px;
  font-weight: 700;
  color: #8898aa;
}

.diag-error {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 20px;
  padding: 14px 16px;
  border-radius: 16px;
  background: #fff7ed;
  border: 1px solid #fed7aa;
  color: #c2410c;
  font-size: 13px;

  button {
    margin-left: auto;
    font-size: 12px;
    font-weight: 700;
    color: #5a8dee;
  }
}

.diag-step-head {
  margin-bottom: 20px;

  h2 {
    margin: 0 0 6px;
    font-size: 20px;
    font-weight: 800;
    color: #1a1f36;
  }

  p {
    margin: 0;
    font-size: 13px;
    color: #6b7c93;
  }
}

.diag-field {
  margin-bottom: 22px;
}

.diag-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 800;
  color: #32325d;
}

.req {
  color: #f56565;
}

.diag-hint {
  margin: -4px 0 10px;
  font-size: 12px;
  line-height: 1.6;
  color: #8898aa;
}

.diag-card-grid {
  display: grid;
  gap: 12px;

  &--3 {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

.diag-option {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  padding: 14px 12px;
  border: 2px solid transparent;
  border-radius: 16px;
  background: #f4f7fb;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #fff;
    box-shadow: 0 8px 20px rgba(90, 141, 238, 0.12);
  }

  &.is-active {
    border-color: #5a8dee;
    background: #f0f6ff;
    box-shadow: 0 8px 20px rgba(90, 141, 238, 0.15);
  }
}

.diag-option__icon {
  font-size: 20px;
  color: #5a8dee;
}

.diag-option__title {
  font-size: 13px;
  font-weight: 800;
  color: #32325d;
}

.diag-option__desc {
  font-size: 11px;
  line-height: 1.5;
  color: #8898aa;
}

.diag-chip-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.diag-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  border: 2px solid transparent;
  border-radius: 999px;
  background: #f4f7fb;
  font-size: 13px;
  font-weight: 700;
  color: #32325d;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #fff;
  }

  &.is-active {
    border-color: #5a8dee;
    background: #f0f6ff;
    color: #3b6fd9;
  }
}

.diag-income-grid {
  display: grid;
  gap: 12px;
}

.diag-income {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  width: 100%;
  padding: 16px 18px;
  border: 2px solid transparent;
  border-radius: 18px;
  background: #f4f7fb;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover,
  &.is-active {
    background: #fff;
    border-color: #5a8dee;
    box-shadow: 0 10px 24px rgba(90, 141, 238, 0.12);
  }
}

.diag-income__label {
  font-size: 16px;
  font-weight: 800;
  color: #1a1f36;
}

.diag-income__desc,
.diag-income__hint {
  font-size: 12px;
  line-height: 1.5;
  color: #8898aa;
}

.diag-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.diag-list-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  width: 100%;
  padding: 14px 16px;
  border: 2px solid transparent;
  border-radius: 16px;
  background: #f4f7fb;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover,
  &.is-active {
    background: #fff;
    border-color: #5a8dee;
  }
}

.diag-list-item__title {
  font-size: 14px;
  font-weight: 800;
  color: #32325d;
}

.diag-list-item__desc {
  font-size: 12px;
  line-height: 1.5;
  color: #8898aa;
}

.diag-risk-grid {
  display: grid;
  gap: 10px;
}

.diag-risk {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  width: 100%;
  padding: 14px 16px;
  border: 2px solid transparent;
  border-radius: 16px;
  background: #f4f7fb;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;

  &.is-active {
    border-color: #5a8dee;
    background: #f0f6ff;
  }

  &--high.is-active {
    border-color: #f56565;
    background: #fff5f5;
  }
}

.diag-risk__title {
  font-size: 14px;
  font-weight: 800;
  color: #32325d;
}

.diag-risk__desc {
  font-size: 12px;
  color: #8898aa;
}

.diag-cost-block {
  margin-bottom: 22px;
  padding: 18px;
  border-radius: 18px;
  background: #f8fafc;
  border: 1px solid #e8edf5;
}

.diag-cost-grid {
  display: grid;
  gap: 14px;
}

.diag-cost-item__top {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 6px;
  margin-bottom: 8px;
}

.diag-cost-item__label {
  font-size: 13px;
  font-weight: 800;
  color: #32325d;
}

.diag-cost-item__hint {
  font-size: 11px;
  color: #a0aec0;
}

.diag-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-top: 28px;
}

.diag-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 14px;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
  transition: all 0.2s ease;

  &--ghost {
    background: #fff;
    color: #32325d;
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.06);
  }

  &--primary {
    background: linear-gradient(135deg, #5a8dee, #3b6fd9);
    color: #fff;
    box-shadow: 0 10px 24px rgba(90, 141, 238, 0.35);

    &:hover:not(:disabled) {
      transform: translateY(-1px);
    }

    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }
  }
}

.diag-skip {
  display: block;
  width: 100%;
  margin-top: 12px;
  padding: 0;
  border: none;
  background: none;
  font-size: 13px;
  font-weight: 700;
  color: #8898aa;
  cursor: pointer;

  &:hover:not(:disabled) {
    color: #5a8dee;
  }
}

.diag-disclaimer {
  margin: 20px 0 0;
  font-size: 11px;
  line-height: 1.6;
  text-align: center;
  color: #a0aec0;
}

.text-clay-muted {
  color: #8898aa;
}

:deep(.diag-input) .el-input__wrapper,
:deep(.diag-textarea) .el-textarea__inner {
  border: none;
  border-radius: 14px;
  background: #fff;
  box-shadow: inset 0 0 0 1px #e2e8f0;
}

:deep(.diag-input) .el-input__inner,
:deep(.diag-textarea) .el-textarea__inner {
  font-weight: 600;
  color: #32325d;
}

:deep(.el-form-item__error) {
  padding-top: 4px;
  font-size: 12px;
}

@media (max-width: 640px) {
  .diag-card {
    padding: 24px 18px 20px;
    border-radius: 24px;
  }

  .diag-card-grid--3 {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .diag-footer {
    flex-direction: column-reverse;
    align-items: stretch;
  }

  .diag-btn {
    justify-content: center;
    width: 100%;
  }
}
</style>
