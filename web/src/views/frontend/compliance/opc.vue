<!-- +----------------------------------------------------------------------
  | XYGo Admin — OPC 落地进度 P-06
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-8 pb-16 px-6">
    <div class="max-w-3xl mx-auto space-y-8">
      <!-- Header -->
      <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-10">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-8">
          <div>
            <h1 class="font-heading font-black text-3xl text-clay-foreground mb-1">OPC 设立进度</h1>
            <p class="text-clay-muted font-medium">预计 {{ progress?.estimatedSlaDays ?? 14 }} 个工作日完成</p>
          </div>
          <div v-if="progress?.companyName" class="px-4 py-2 rounded-2xl bg-blue-50 text-clay-accent font-bold text-sm">
            {{ progress.companyName }}
          </div>
        </div>

        <!-- Timeline -->
        <div v-if="loading" class="text-center py-8">
          <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
        </div>
        <div v-else class="relative pl-6 space-y-0">
          <div
            v-for="(step, idx) in timelineSteps"
            :key="step.key"
            class="relative pb-8 last:pb-0"
          >
            <div
              v-if="idx < timelineSteps.length - 1"
              class="absolute left-[11px] top-6 bottom-0 w-0.5"
              :class="step.status === 'done' ? 'bg-clay-accent' : 'bg-[#d1d9e6]'"
            />
            <div class="flex items-start gap-4">
              <div
                class="w-6 h-6 rounded-full shrink-0 flex items-center justify-center z-10"
                :class="stepDotClass(step.status)"
              >
                <ArtSvgIcon v-if="step.status === 'done'" icon="ri:check-line" class="text-white text-xs" />
                <span v-else-if="step.status === 'current'" class="w-2 h-2 rounded-full bg-white" />
              </div>
              <div class="flex-1 pt-0.5">
                <div class="flex items-center gap-2">
                  <span class="font-heading font-black text-clay-foreground">{{ step.label }}</span>
                  <span class="text-xs font-bold px-2 py-0.5 rounded-full" :class="stepBadgeClass(step.status)">
                    {{ stepSubLabel(step) }}
                  </span>
                </div>
                <p v-if="step.date" class="text-xs text-clay-muted mt-1">{{ step.date }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Reject note -->
        <div v-if="progress?.rejectNote" class="mt-6 p-4 rounded-2xl bg-orange-50 border border-orange-200">
          <p class="text-sm font-bold text-orange-700">驳回原因：{{ progress.rejectNote }}</p>
        </div>
      </div>

      <!-- Materials form -->
      <div v-if="showMaterialsForm" class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-10">
        <div
          v-if="verifyMock"
          class="mb-6 px-4 py-3 rounded-2xl bg-amber-50 border border-amber-200 text-sm text-amber-800 font-medium"
        >
          开发模式：身份证 OCR、手机实名、三要素等真实性校验已 Mock，填写格式正确即可提交。
        </div>
        <h2 class="font-heading font-black text-xl text-clay-foreground mb-6">注册资料提交</h2>

        <ElForm ref="formRef" :model="formData" :rules="rules" label-position="top" class="space-y-6">
          <!-- 拟设公司 -->
          <section>
            <h3 class="text-sm font-black text-clay-muted mb-4 uppercase tracking-wide">拟设公司信息</h3>
            <div class="space-y-4">
              <ElFormItem label="备选公司名称（1~3个）" prop="proposedNames">
                <div class="space-y-2">
                  <ElInput
                    v-for="(_, idx) in formData.proposedNames"
                    :key="idx"
                    v-model.trim="formData.proposedNames[idx]"
                    :placeholder="`备选名称 ${idx + 1}`"
                    size="large"
                    class="clay-input"
                  />
                  <button
                    v-if="formData.proposedNames.length < 3"
                    type="button"
                    class="text-xs font-bold text-clay-accent hover:underline"
                    @click="formData.proposedNames.push('')"
                  >
                    + 添加备选名称
                  </button>
                </div>
              </ElFormItem>
              <div class="grid md:grid-cols-2 gap-4">
                <ElFormItem label="注册资本（万元）" prop="registeredCapital">
                  <ElInputNumber v-model="formData.registeredCapital" :min="0.01" :max="1000" :precision="2" class="w-full" />
                </ElFormItem>
                <ElFormItem label="认缴期限" prop="capitalTermYears">
                  <ElSelect v-model="formData.capitalTermYears" placeholder="请选择" size="large" class="w-full clay-select">
                    <ElOption v-for="y in [5, 10, 20, 30]" :key="y" :label="`${y} 年`" :value="y" />
                  </ElSelect>
                </ElFormItem>
              </div>
              <ElFormItem label="营业期限" prop="businessTermType">
                <ElRadioGroup v-model="formData.businessTermType" class="flex flex-col gap-2">
                  <ElRadio value="long_term" class="clay-radio">长期</ElRadio>
                  <ElRadio value="fixed" class="clay-radio">固定期限</ElRadio>
                </ElRadioGroup>
              </ElFormItem>
              <ElFormItem v-if="formData.businessTermType === 'fixed'" label="营业期限截止" prop="businessTermEnd">
                <ElDatePicker v-model="formData.businessTermEnd" type="date" value-format="YYYY-MM-DD" class="w-full" />
              </ElFormItem>
              <ElFormItem label="经营范围" prop="businessScope">
                <ElInput v-model="formData.businessScope" type="textarea" :rows="4" class="clay-textarea" />
                <button type="button" class="text-xs font-bold text-clay-accent hover:underline mt-2" @click="applyRecommendedScope">
                  使用推荐模板
                </button>
              </ElFormItem>
              <ChinaRegionSelect
                v-model:province="formData.registerProvince"
                v-model:city="formData.registerCity"
                v-model:district="formData.registerDistrict"
              />
              <ElFormItem label="详细地址" prop="registerAddress">
                <ElInput v-model.trim="formData.registerAddress" placeholder="街道门牌号" size="large" class="clay-input" />
              </ElFormItem>
              <ElFormItem label="地址证明" prop="addressProofFileId">
                <MemberFileUpload v-model="formData.addressProofFileId" accept=".pdf,.jpg,.jpeg,.png" />
              </ElFormItem>
            </div>
          </section>

          <!-- 法人信息 -->
          <section>
            <h3 class="text-sm font-black text-clay-muted mb-4 uppercase tracking-wide">法人（股东）信息</h3>
            <div class="space-y-4">
              <ElFormItem label="法人姓名" prop="legalPersonName">
                <ElInput v-model.trim="formData.legalPersonName" size="large" class="clay-input" />
              </ElFormItem>
              <ElFormItem label="身份证号" prop="idCardNumber">
                <ElInput v-model.trim="formData.idCardNumber" maxlength="18" size="large" class="clay-input" />
              </ElFormItem>
              <div class="grid md:grid-cols-2 gap-4">
                <ElFormItem label="身份证有效期起" prop="idCardValidFrom">
                  <ElDatePicker v-model="formData.idCardValidFrom" type="date" value-format="YYYY-MM-DD" class="w-full" />
                </ElFormItem>
                <ElFormItem label="身份证有效期止" prop="idCardValidTo">
                  <ElDatePicker v-model="formData.idCardValidTo" type="date" value-format="YYYY-MM-DD" class="w-full" />
                </ElFormItem>
              </div>
              <ElFormItem label="民族">
                <ElInput v-model.trim="formData.ethnicity" size="large" class="clay-input" />
              </ElFormItem>
              <ElFormItem label="户籍地址" prop="householdAddress">
                <ElInput v-model.trim="formData.householdAddress" size="large" class="clay-input" />
              </ElFormItem>
              <ElFormItem label="现居住地址" prop="residenceAddress">
                <ElInput v-model.trim="formData.residenceAddress" size="large" class="clay-input" />
              </ElFormItem>
              <div class="grid md:grid-cols-2 gap-4">
                <ElFormItem label="手机号" prop="phone">
                  <ElInput v-model.trim="formData.phone" maxlength="11" size="large" class="clay-input" />
                </ElFormItem>
                <ElFormItem label="邮箱" prop="email">
                  <ElInput v-model.trim="formData.email" size="large" class="clay-input" />
                </ElFormItem>
              </div>
              <ElFormItem label="身份证正面" prop="idCardFrontFileId">
                <MemberFileUpload v-model="formData.idCardFrontFileId" accept=".jpg,.jpeg,.png" />
              </ElFormItem>
              <ElFormItem label="身份证反面" prop="idCardBackFileId">
                <MemberFileUpload v-model="formData.idCardBackFileId" accept=".jpg,.jpeg,.png" />
              </ElFormItem>
            </div>
          </section>

          <!-- 授权确认 -->
          <section>
            <h3 class="text-sm font-black text-clay-muted mb-4 uppercase tracking-wide">授权确认</h3>
            <div class="space-y-3">
              <ElCheckbox v-model="formData.confirmations.infoTrue" class="clay-checkbox w-full">本人确认以上信息真实、完整</ElCheckbox>
              <ElCheckbox v-model="formData.confirmations.authConsent" class="clay-checkbox w-full">同意用于工商、税务、银行开户申报</ElCheckbox>
              <ElCheckbox v-model="formData.confirmations.opcLimitAck" class="clay-checkbox w-full">知晓自然人 3 年内不得再设立新 OPC</ElCheckbox>
              <ElCheckbox v-model="formData.confirmations.eSignAuth" class="clay-checkbox w-full">电子签名授权（沿用签约姓名）</ElCheckbox>
            </div>
          </section>

          <button
            type="button"
            class="w-full h-14 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black shadow-clay-btn hover:shadow-clay-btn-hover active:scale-95 transition-all flex items-center justify-center gap-2 disabled:opacity-50"
            :disabled="submitting"
            @click="handleSubmitMaterials"
          >
            <ArtSvgIcon v-if="submitting" icon="ri:loader-4-line" class="text-xl animate-spin" />
            {{ submitting ? '提交中...' : '提交注册资料' }}
          </button>
        </ElForm>
      </div>

      <!-- Readonly materials -->
      <div v-else-if="progress?.materialsReadonly" class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-10">
        <h2 class="font-heading font-black text-xl text-clay-foreground mb-4">已提交资料</h2>
        <p class="text-sm text-clay-muted font-medium">资料审核中或已进入后续流程，如需修改请等待顾问联系。</p>
      </div>

      <!-- Bank section -->
      <div v-if="showBankSection" class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-10">
        <h2 class="font-heading font-black text-xl text-clay-foreground mb-4">银行开户指引</h2>
        <ElCollapse>
          <ElCollapseItem title="开户材料清单（点击展开）" name="guide">
            <ul class="space-y-2 text-sm text-clay-foreground font-medium">
              <li>· 营业执照正副本原件</li>
              <li>· 公章、财务章、法人章（刻章完成后）</li>
              <li>· 公司章程</li>
              <li>· 法人身份证原件</li>
              <li>· 税务登记相关回执</li>
            </ul>
          </ElCollapseItem>
        </ElCollapse>

        <div v-if="progress?.opcStatus === 'bank'" class="mt-6 space-y-4">
          <ElFormItem label="开户银行（选填）">
            <ElInput v-model.trim="bankForm.bankName" placeholder="如：中国工商银行 XX 支行" size="large" class="clay-input" />
          </ElFormItem>
          <ElFormItem label="上传开户回执">
            <MemberFileUpload v-model="bankForm.bankReceiptFileId" accept=".pdf,.jpg,.jpeg,.png" />
          </ElFormItem>
          <button
            type="button"
            class="px-8 py-3 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black shadow-clay-btn hover:shadow-clay-btn-hover active:scale-95 transition-all disabled:opacity-50"
            :disabled="!bankForm.bankReceiptFileId || bankSubmitting"
            @click="handleSubmitBankReceipt"
          >
            提交开户回执
          </button>
        </div>
      </div>

      <!-- Active success -->
      <div v-if="progress?.opcStatus === 'active'" class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 text-center">
        <ArtSvgIcon icon="ri:checkbox-circle-line" class="text-5xl text-clay-success mx-auto mb-4" />
        <h2 class="font-heading font-black text-2xl text-clay-foreground mb-2">OPC 设立完成</h2>
        <p class="text-clay-muted font-medium">收入台账、费用台账等功能已解锁</p>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import ChinaRegionSelect from '@/components/frontend/ChinaRegionSelect.vue'
import MemberFileUpload from '@/components/frontend/MemberFileUpload.vue'
import {
  getOpcProgress,
  submitOpcMaterials,
  submitBankReceipt,
  RECOMMENDED_BUSINESS_SCOPE,
  type OpcProgressResult,
  type OpcProgressStep,
  type OpcStepStatus
} from '@/api/frontend/compliance/opc'
import {
  ID_CARD_FORMAT_PATTERN,
  PHONE_MOCK_PATTERN,
  isComplianceVerifyMock
} from '@/config/complianceVerify'
import { validateChineseIDCard, validateEmail, validatePhone } from '@/utils/form/validator'
import { useMemberStore } from '@/store/modules/member'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

defineOptions({ name: 'ComplianceOpc' })

const router = useRouter()
const memberStore = useMemberStore()
const verifyMock = isComplianceVerifyMock()

const loading = ref(true)
const submitting = ref(false)
const bankSubmitting = ref(false)
const progress = ref<OpcProgressResult | null>(null)
const formRef = ref<FormInstance>()

const DEFAULT_STEPS: OpcProgressStep[] = [
  { key: 'materials', label: '资料提交', status: 'current' },
  { key: 'business', label: '工商注册', status: 'pending' },
  { key: 'tax', label: '税务登记', status: 'pending' },
  { key: 'bank', label: '银行开户', status: 'pending' }
]

const timelineSteps = computed(() => progress.value?.steps?.length ? progress.value.steps : DEFAULT_STEPS)

const showMaterialsForm = computed(() =>
  progress.value?.materialsEditable !== false &&
  ['pending', 'materials'].includes(progress.value?.opcStatus ?? 'pending')
)

const showBankSection = computed(() =>
  ['tax', 'bank', 'active'].includes(progress.value?.opcStatus ?? '')
)

const formData = reactive({
  proposedNames: [''],
  registeredCapital: 10,
  capitalTermYears: 20,
  businessTermType: 'long_term' as 'long_term' | 'fixed',
  businessTermEnd: '',
  businessScope: '',
  registerProvince: '',
  registerCity: '',
  registerDistrict: '',
  registerAddress: '',
  addressProofFileId: '' as number | string,
  legalPersonName: '',
  idCardNumber: '',
  idCardValidFrom: '',
  idCardValidTo: '',
  ethnicity: '',
  householdAddress: '',
  residenceAddress: '',
  phone: '',
  email: '',
  idCardFrontFileId: '' as number | string,
  idCardBackFileId: '' as number | string,
  confirmations: {
    infoTrue: false,
    authConsent: false,
    opcLimitAck: false,
    eSignAuth: false
  }
})

const bankForm = reactive({
  bankName: '',
  bankReceiptFileId: '' as number | string
})


function validateFormPhone(_r: unknown, value: string, cb: (err?: Error) => void) {
  const phone = String(value ?? '').trim()
  if (!phone) {
    cb(new Error('请输入手机号'))
    return
  }
  const ok = verifyMock ? PHONE_MOCK_PATTERN.test(phone) : validatePhone(phone)
  cb(ok ? undefined : new Error('手机号格式不正确'))
}

function validateFormIdCard(_r: unknown, value: string, cb: (err?: Error) => void) {
  const id = String(value ?? '').trim()
  if (!id) {
    cb(new Error('请输入身份证号'))
    return
  }
  const ok = verifyMock ? ID_CARD_FORMAT_PATTERN.test(id) : validateChineseIDCard(id)
  cb(ok ? undefined : new Error('身份证号格式不正确'))
}

function validateFormEmail(_r: unknown, value: string, cb: (err?: Error) => void) {
  const email = String(value ?? '').trim()
  if (!email) {
    cb(new Error('请输入邮箱'))
    return
  }
  cb(validateEmail(email) ? undefined : new Error('邮箱格式不正确'))
}

const rules: FormRules = {
  proposedNames: [{
    validator: (_r, _v, cb) => {
      const names = formData.proposedNames.filter(n => n.trim())
      if (!names.length) cb(new Error('至少填写一个备选名称'))
      else if (names.some(n => n.length < 2 || n.length > 30)) cb(new Error('名称长度 2-30 字'))
      else cb()
    },
    trigger: 'blur'
  }],
  registeredCapital: [{ required: true, message: '请填写注册资本', trigger: 'blur' }],
  capitalTermYears: [{ required: true, message: '请选择认缴期限', trigger: 'change' }],
  businessScope: [
    { required: true, message: '请填写经营范围', trigger: 'blur' },
    { min: 10, max: 2000, message: '10-2000 字', trigger: 'blur' }
  ],
  registerProvince: [{ required: true, message: '请选择省份', trigger: 'change' }],
  registerCity: [{ required: true, message: '请选择城市', trigger: 'change' }],
  registerDistrict: [{ required: true, message: '请选择区县', trigger: 'change' }],
  registerAddress: [{ required: true, min: 5, message: '详细地址至少 5 字', trigger: 'blur' }],
  businessTermType: [{ required: true, message: '请选择营业期限', trigger: 'change' }],
  addressProofFileId: [{
    validator: (_r, v, cb) => (Number(v) > 0 ? cb() : cb(new Error('请上传地址证明'))),
    trigger: 'change',
  }],
  legalPersonName: [{ required: true, min: 2, max: 20, message: '2-20 字中文姓名', trigger: 'blur' }],
  idCardNumber: [{ validator: validateFormIdCard, trigger: 'blur' }],
  idCardValidFrom: [{ required: true, message: '请选择有效期起', trigger: 'change' }],
  idCardValidTo: [{ required: true, message: '请选择有效期止', trigger: 'change' }],
  householdAddress: [{ required: true, min: 5, message: '户籍地址至少 5 字', trigger: 'blur' }],
  residenceAddress: [{ required: true, min: 5, message: '现居地址至少 5 字', trigger: 'blur' }],
  phone: [{ validator: validateFormPhone, trigger: 'blur' }],
  email: [{ validator: validateFormEmail, trigger: 'blur' }],
  idCardFrontFileId: [{
    validator: (_r, v, cb) => (Number(v) > 0 ? cb() : cb(new Error('请上传身份证正面'))),
    trigger: 'change',
  }],
  idCardBackFileId: [{
    validator: (_r, v, cb) => (Number(v) > 0 ? cb() : cb(new Error('请上传身份证反面'))),
    trigger: 'change',
  }],
}

function stepDotClass(status: OpcStepStatus) {
  if (status === 'done') return 'bg-gradient-to-br from-blue-400 to-blue-600 shadow-clay-btn'
  if (status === 'current') return 'bg-clay-accent shadow-clay-btn ring-4 ring-blue-100'
  if (status === 'rejected') return 'bg-orange-500'
  return 'bg-[#d1d9e6]'
}

function stepBadgeClass(status: OpcStepStatus) {
  if (status === 'done') return 'bg-green-100 text-green-700'
  if (status === 'current') return 'bg-blue-100 text-clay-accent'
  if (status === 'rejected') return 'bg-orange-100 text-orange-700'
  return 'bg-gray-100 text-clay-muted'
}

function stepSubLabel(step: OpcProgressStep) {
  if (step.subLabel) return step.subLabel
  const map: Record<OpcStepStatus, string> = {
    done: '已完成',
    current: '进行中',
    pending: '待开始',
    rejected: '已退回'
  }
  return map[step.status] ?? '待开始'
}

function applyRecommendedScope() {
  formData.businessScope = RECOMMENDED_BUSINESS_SCOPE
}

async function loadProgress() {
  loading.value = true
  try {
    progress.value = await getOpcProgress()
  } catch {
    progress.value = {
      opcId: 0,
      opcStatus: 'pending',
      estimatedSlaDays: 14,
      steps: DEFAULT_STEPS,
      materialsEditable: true
    }
  } finally {
    loading.value = false
  }
}

async function handleSubmitMaterials() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  const c = formData.confirmations
  if (!c.infoTrue || !c.authConsent || !c.opcLimitAck || !c.eSignAuth) {
    ElMessage.warning('请勾选全部授权确认项')
    return
  }

  submitting.value = true
  try {
    await submitOpcMaterials({
      ...formData,
      proposedNames: formData.proposedNames.filter(n => n.trim())
    })
    ElMessage.success('资料提交成功，等待顾问审核')
    await loadProgress()
  } catch (e: any) {
    const msg = e?.message || e?.msg || '提交失败，请检查表单后重试'
    ElMessage.error(msg)
  } finally {
    submitting.value = false
  }
}

async function handleSubmitBankReceipt() {
  if (!bankForm.bankReceiptFileId) return
  bankSubmitting.value = true
  try {
    await submitBankReceipt(bankForm)
    ElMessage.success('开户回执已提交')
    await loadProgress()
  } catch {
    ElMessage.error('提交失败，请重试')
  } finally {
    bankSubmitting.value = false
  }
}

onMounted(() => {
  if (!memberStore.isLogin) {
    router.replace({ path: '/user/login', query: { redirect: '/user/compliance/opc' } })
    return
  }
  loadProgress()
})
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.bg-clay-accent { background-color: #5a8dee; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.shadow-clay-deep {
  box-shadow: 30px 30px 60px #d1d9e6, -30px -30px 60px #ffffff,
    inset 10px 10px 20px rgba(90, 141, 238, 0.05), inset -10px -10px 20px rgba(255, 255, 255, 0.8);
}
.shadow-clay-btn {
  box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-btn-hover {
  box-shadow: 16px 16px 32px rgba(90, 141, 238, 0.4), -10px -10px 20px rgba(255, 255, 255, 0.5);
}

:deep(.clay-input) .el-input__wrapper {
  height: 48px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff; border: none;
}
:deep(.clay-textarea) .el-textarea__inner {
  border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff; border: none;
}
:deep(.clay-select) .el-select__wrapper {
  height: 48px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff; border: none;
}
:deep(.clay-checkbox), :deep(.clay-radio) {
  padding: 12px 16px; border-radius: 16px; background: #f0f3f8;
  box-shadow: inset 6px 6px 12px #e0e5ec, inset -6px -6px 12px #ffffff;
  margin-right: 0 !important; height: auto;
}
</style>
