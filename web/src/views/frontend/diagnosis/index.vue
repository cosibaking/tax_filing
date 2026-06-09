<!-- +----------------------------------------------------------------------
  | XYGo Admin — OPC 合规诊断问卷 P-02
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-8 pb-16 px-6">
    <div class="max-w-2xl mx-auto relative">
      <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-12 relative z-10">
        <!-- Header -->
        <div class="text-center mb-10">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/50 border border-white/50 shadow-sm mb-6">
            <ArtSvgIcon icon="ri:shield-check-line" class="text-lg text-clay-accent" />
            <span class="text-sm font-bold text-clay-muted">约 10 分钟完成</span>
          </div>
          <h1 class="font-heading font-black text-3xl text-clay-foreground mb-2">合规诊断问卷</h1>
          <p class="text-clay-muted font-medium">了解您的收入结构与合规现状，获取个性化税负对比方案</p>
        </div>

        <!-- Step indicator -->
        <div class="flex items-center justify-between mb-10 px-2">
          <div
            v-for="(label, idx) in stepLabels"
            :key="idx"
            class="flex flex-col items-center flex-1"
          >
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
            <span
              class="text-[10px] font-bold mt-2 text-center"
              :class="currentStep === idx + 1 ? 'text-clay-accent' : 'text-clay-muted'"
            >{{ label }}</span>
          </div>
        </div>

        <!-- API Error -->
        <div
          v-if="apiError"
          class="mb-6 p-4 rounded-2xl bg-orange-50 border border-orange-200 flex items-start gap-3"
        >
          <ArtSvgIcon icon="ri:error-warning-line" class="text-xl text-orange-500 shrink-0 mt-0.5" />
          <div class="flex-1">
            <p class="text-sm font-bold text-orange-700">{{ apiError }}</p>
            <button type="button" class="text-xs font-bold text-clay-accent hover:underline mt-1" @click="apiError = ''">关闭</button>
          </div>
        </div>

        <ElForm ref="formRef" :model="formData" :rules="rules" class="space-y-6">
          <!-- Step 1: 平台与收入 -->
          <div v-show="currentStep === 1" class="space-y-6">
            <ElFormItem prop="platforms" class="!mb-0">
              <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">主要平台 <span class="text-red-400">*</span></label>
              <ElCheckboxGroup v-model="formData.platforms" class="grid grid-cols-2 gap-3">
                <ElCheckbox
                  v-for="p in platformOptions"
                  :key="p"
                  :value="p"
                  class="!mr-0 clay-checkbox"
                >
                  <span class="text-sm font-bold text-clay-foreground">{{ p }}</span>
                </ElCheckbox>
              </ElCheckboxGroup>
            </ElFormItem>

            <ElFormItem prop="monthlyIncomeRange" class="!mb-0">
              <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">月收入区间 <span class="text-red-400">*</span></label>
              <ElSelect v-model="formData.monthlyIncomeRange" placeholder="请选择月收入区间" size="large" class="w-full clay-select">
                <ElOption v-for="r in incomeRangeOptions" :key="r.value" :label="r.label" :value="r.value" />
              </ElSelect>
            </ElFormItem>

            <ElFormItem prop="annualCostEstimate" class="!mb-0">
              <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">年可扣除成本（估，元）</label>
              <ElInput
                v-model.number="formData.annualCostEstimate"
                placeholder="如设备、投流、场地等合计"
                size="large"
                type="number"
                class="clay-input"
              >
                <template #prefix>
                  <ArtSvgIcon icon="ri:money-cny-circle-line" class="text-lg text-clay-muted" />
                </template>
              </ElInput>
            </ElFormItem>
          </div>

          <!-- Step 2: 现有主体 -->
          <div v-show="currentStep === 2" class="space-y-6">
            <ElFormItem prop="existingEntity" class="!mb-0">
              <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">是否已有经营主体 <span class="text-red-400">*</span></label>
              <ElRadioGroup v-model="formData.existingEntity" class="flex flex-col gap-3 w-full">
                <ElRadio
                  v-for="e in entityOptions"
                  :key="e.value"
                  :value="e.value"
                  class="!mr-0 clay-radio w-full"
                >
                  <span class="text-sm font-bold text-clay-foreground">{{ e.label }}</span>
                </ElRadio>
              </ElRadioGroup>
            </ElFormItem>

            <ElFormItem prop="hasFiledTax" class="!mb-0">
              <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">是否按时报税 <span class="text-red-400">*</span></label>
              <ElRadioGroup v-model="formData.hasFiledTax" class="flex flex-col gap-3 w-full">
                <ElRadio v-for="t in taxFiledOptions" :key="t.value" :value="t.value" class="!mr-0 clay-radio w-full">
                  <span class="text-sm font-bold text-clay-foreground">{{ t.label }}</span>
                </ElRadio>
              </ElRadioGroup>
            </ElFormItem>
          </div>

          <!-- Step 3: 风险信号 -->
          <div v-show="currentStep === 3" class="space-y-6">
            <ElFormItem prop="taxBureauContact" class="!mb-0">
              <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">是否被税务局联系过 <span class="text-red-400">*</span></label>
              <ElRadioGroup v-model="formData.taxBureauContact" class="flex gap-4">
                <ElRadio :value="true" class="clay-radio">
                  <span class="text-sm font-bold text-clay-foreground">是</span>
                </ElRadio>
                <ElRadio :value="false" class="clay-radio">
                  <span class="text-sm font-bold text-clay-foreground">否</span>
                </ElRadio>
              </ElRadioGroup>
            </ElFormItem>

            <ElFormItem prop="notes" class="!mb-0">
              <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">补充说明</label>
              <ElInput
                v-model="formData.notes"
                type="textarea"
                :rows="4"
                placeholder="如有特殊情况可在此说明（选填）"
                class="clay-textarea"
              />
            </ElFormItem>
          </div>

          <!-- Step 4: 成本细项（可选） -->
          <div v-show="currentStep === 4" class="space-y-6">
            <p class="text-sm text-clay-muted font-medium -mt-2">填写成本细项可获得更精准的税负估算，也可直接跳过。</p>
            <div v-for="item in costItems" :key="item.key" class="space-y-2">
              <label class="block text-sm font-bold text-clay-foreground ml-1">{{ item.label }}</label>
              <ElInput
                v-model.number="formData.costBreakdown[item.key]"
                :placeholder="item.placeholder"
                size="large"
                type="number"
                class="clay-input"
              />
            </div>
          </div>
        </ElForm>

        <!-- Navigation -->
        <div class="flex items-center justify-between mt-10 gap-4">
          <button
            v-if="currentStep > 1"
            type="button"
            class="px-8 py-3 rounded-2xl bg-white shadow-clay-btn hover:shadow-clay-btn-hover font-bold text-clay-foreground active:scale-95 transition-all"
            @click="prevStep"
          >
            上一步
          </button>
          <div v-else />

          <button
            v-if="currentStep < 4"
            type="button"
            class="px-10 py-3 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black shadow-clay-btn hover:shadow-clay-btn-hover hover:-translate-y-1 active:scale-95 transition-all flex items-center gap-2"
            @click="nextStep"
          >
            下一步
            <ArtSvgIcon icon="ri:arrow-right-line" class="text-lg" />
          </button>
          <button
            v-else
            type="button"
            class="px-10 py-3 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black shadow-clay-btn hover:shadow-clay-btn-hover hover:-translate-y-1 active:scale-95 transition-all flex items-center gap-2"
            :disabled="submitting"
            @click="handleSubmit"
          >
            <ArtSvgIcon v-if="submitting" icon="ri:loader-4-line" class="text-lg animate-spin" />
            {{ submitting ? '分析中...' : '查看诊断结果' }}
          </button>
        </div>

        <button
          v-if="currentStep === 4"
          type="button"
          class="w-full mt-4 text-sm font-bold text-clay-muted hover:text-clay-accent transition-colors"
          :disabled="submitting"
          @click="handleSubmit"
        >
          跳过细项，直接查看结果
        </button>
      </div>

      <div class="absolute -bottom-10 -left-10 w-32 h-32 rounded-full bg-gradient-to-br from-cyan-300 to-cyan-500 opacity-20 blur-2xl animate-float z-0" />
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

defineOptions({ name: 'ComplianceDiagnosis' })

const router = useRouter()
const memberStore = useMemberStore()
const formRef = ref<FormInstance>()

const stepLabels = ['平台收入', '现有主体', '风险信号', '成本细项']
const currentStep = ref(1)
const submitting = ref(false)
const apiError = ref('')

const platformOptions = ['抖音', '快手', 'B站', '小红书', '视频号', '其他']
const incomeRangeOptions: { label: string; value: MonthlyIncomeRange }[] = [
  { label: '0 - 2 万', value: '0-2万' },
  { label: '2 - 5 万', value: '2-5万' },
  { label: '5 - 15 万', value: '5-15万' },
  { label: '15 万以上', value: '15万+' }
]
const entityOptions: { label: string; value: ExistingEntity }[] = [
  { label: '无经营主体', value: 'none' },
  { label: '个体户', value: 'individual' },
  { label: '公司', value: 'company' },
  { label: '其他', value: 'other' }
]
const taxFiledOptions: { label: string; value: HasFiledTax }[] = [
  { label: '是，按时申报', value: 'yes' },
  { label: '否，存在漏报', value: 'no' },
  { label: '不确定', value: 'unsure' }
]
type CostBreakdownKey = 'device' | 'marketing' | 'venue' | 'other'

const costItems: { key: CostBreakdownKey; label: string; placeholder: string }[] = [
  { key: 'device', label: '设备/器材', placeholder: '如相机、电脑等' },
  { key: 'marketing', label: '投流/推广', placeholder: '年度投流费用' },
  { key: 'venue', label: '场地/租赁', placeholder: '工作室、场地租金' },
  { key: 'other', label: '其他成本', placeholder: '助理、差旅等' }
]

const formData = reactive({
  platforms: [] as string[],
  monthlyIncomeRange: '' as MonthlyIncomeRange | '',
  annualCostEstimate: undefined as number | undefined,
  existingEntity: '' as ExistingEntity | '',
  hasFiledTax: '' as HasFiledTax | '',
  taxBureauContact: false as boolean,
  notes: '',
  costBreakdown: {
    device: undefined as number | undefined,
    marketing: undefined as number | undefined,
    venue: undefined as number | undefined,
    other: undefined as number | undefined
  }
})

const rules: FormRules = {
  platforms: [{
    type: 'array',
    required: true,
    min: 1,
    message: '请至少选择一个平台',
    trigger: 'change'
  }],
  monthlyIncomeRange: [{ required: true, message: '请选择月收入区间', trigger: 'change' }],
  existingEntity: [{ required: true, message: '请选择经营主体情况', trigger: 'change' }],
  hasFiledTax: [{ required: true, message: '请选择报税情况', trigger: 'change' }],
  taxBureauContact: [{ required: true, message: '请选择是否被税局联系', trigger: 'change' }]
}

const stepFields: Record<number, string[]> = {
  1: ['platforms', 'monthlyIncomeRange'],
  2: ['existingEntity', 'hasFiledTax'],
  3: ['taxBureauContact']
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

async function nextStep() {
  const valid = await validateStep(currentStep.value)
  if (!valid) return
  if (currentStep.value < 4) currentStep.value++
}

function prevStep() {
  if (currentStep.value > 1) currentStep.value--
}

function buildPayload(): DiagnosisSubmitParams {
  const breakdown = Object.fromEntries(
    (Object.entries(formData.costBreakdown) as [CostBreakdownKey, number | undefined][])
      .filter(([, v]) => v != null && v > 0)
  ) as Record<string, number>
  const breakdownSum = Object.values(breakdown).reduce((sum, v) => sum + v, 0)
  const annualCost = formData.annualCostEstimate || breakdownSum || undefined

  return {
    platforms: formData.platforms,
    monthlyIncomeRange: formData.monthlyIncomeRange as MonthlyIncomeRange,
    annualCostEstimate: annualCost,
    existingEntity: formData.existingEntity as ExistingEntity,
    hasFiledTax: formData.hasFiledTax as HasFiledTax,
    taxBureauContact: formData.taxBureauContact,
    notes: formData.notes || undefined,
    costBreakdown: Object.keys(breakdown).length ? breakdown : undefined
  }
}

async function handleSubmit() {
  for (let s = 1; s <= 3; s++) {
    const valid = await validateStep(s)
    if (!valid) {
      currentStep.value = s
      return
    }
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

@keyframes float { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-20px); } }
.animate-float { animation: float 8s ease-in-out infinite; }

:deep(.clay-input) {
  .el-input__wrapper {
    height: 48px; padding: 0 16px; border-radius: 16px; background: #f0f3f8;
    box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
    border: none; transition: all 0.3s;
    &.is-focus {
      background: #fff;
      box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9),
        inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
    }
  }
  .el-input__inner { font-weight: 500; color: #32325d; }
}

:deep(.clay-textarea) .el-textarea__inner {
  padding: 16px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
  border: none; font-weight: 500; color: #32325d;
  &:focus {
    background: #fff;
    box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9);
  }
}

:deep(.clay-select) .el-select__wrapper {
  height: 48px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
  border: none; font-weight: 500;
}

:deep(.clay-checkbox) {
  padding: 12px 16px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 6px 6px 12px #e0e5ec, inset -6px -6px 12px #ffffff;
  margin-right: 0 !important;
  height: auto;
}

:deep(.clay-radio) {
  padding: 14px 20px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 6px 6px 12px #e0e5ec, inset -6px -6px 12px #ffffff;
  margin-right: 0 !important;
  height: auto;
}
</style>
