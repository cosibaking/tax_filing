<!-- 主体设立进度 P-06 -->
<template>
  <div class="opc-page">
    <div class="opc-page__inner">
      <!-- Signed plan -->
      <section v-if="signedPlanLabel" class="opc-panel">
        <div class="opc-panel__head">
          <div>
            <p class="opc-panel__label">当前签约套餐</p>
            <h2 class="opc-panel__title">{{ signedPlanLabel }}</h2>
            <p v-if="progress?.signedAt" class="opc-panel__desc">签约时间：{{ progress.signedAt }}</p>
          </div>
          <div v-if="progress?.planAmount" class="opc-tag opc-tag--price">
            ¥{{ progress.planAmount.toFixed(2) }}/月
          </div>
        </div>
      </section>

      <!-- Header + Timeline -->
      <section class="opc-panel">
        <div class="opc-panel__head">
          <div>
            <h1 class="opc-panel__title opc-panel__title--lg">主体设立进度</h1>
            <p class="opc-panel__desc">预计 {{ progress?.estimatedSlaDays ?? 14 }} 个工作日完成</p>
          </div>
          <div v-if="progress?.companyName" class="opc-tag">{{ progress.companyName }}</div>
        </div>

        <div v-if="loading" class="opc-loading">
          <ArtSvgIcon icon="ri:loader-4-line" class="opc-loading__icon" />
        </div>
        <div v-else class="opc-timeline">
          <div
            v-for="(step, idx) in timelineSteps"
            :key="step.key"
            class="opc-timeline__item"
            :class="{
              'is-done': step.status === 'done',
              'is-current': step.status === 'current'
            }"
          >
            <span class="opc-timeline__index">{{ idx + 1 }}</span>
            <div class="opc-timeline__body">
              <div class="opc-timeline__row">
                <h3>{{ step.label }}</h3>
                <span class="opc-badge" :class="`opc-badge--${step.status}`">{{ stepSubLabel(step) }}</span>
              </div>
              <p v-if="step.date">{{ step.date }}</p>
            </div>
          </div>
        </div>

        <div v-if="progress?.rejectNote" class="opc-alert opc-alert--warn">
          驳回原因：{{ progress.rejectNote }}
        </div>
      </section>

      <!-- Materials form -->
      <section v-if="showMaterialsForm" class="opc-panel">
        <div v-if="verifyMock" class="opc-alert opc-alert--info">
          开发模式：身份证 OCR、手机实名、三要素等真实性校验已 Mock，填写格式正确即可提交。
        </div>
        <h2 class="opc-panel__section-title">注册资料提交</h2>

        <ElForm ref="formRef" :model="formData" :rules="rules" label-position="top" class="space-y-6">
          <!-- 拟设公司 -->
          <section>
            <h3 class="opc-form-section__title">拟设公司信息</h3>
            <div class="space-y-4">
              <ElFormItem label="备选公司名称（1~3个）" prop="proposedNames">
                <div class="space-y-2">
                  <ElInput
                    v-for="(_, idx) in formData.proposedNames"
                    :key="idx"
                    v-model.trim="formData.proposedNames[idx]"
                    :placeholder="`备选名称 ${idx + 1}`"
                    size="large"
                    class="opc-field"
                  />
                  <button
                    v-if="formData.proposedNames.length < 3"
                    type="button"
                    class="opc-link-btn"
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
                  <ElSelect v-model="formData.capitalTermYears" placeholder="请选择" size="large" class="w-full opc-field">
                    <ElOption v-for="y in [5, 10, 20, 30]" :key="y" :label="`${y} 年`" :value="y" />
                  </ElSelect>
                </ElFormItem>
              </div>
              <ElFormItem label="营业期限" prop="businessTermType">
                <ElRadioGroup v-model="formData.businessTermType" class="flex flex-col gap-2">
                  <ElRadio value="long_term" class="opc-check">长期</ElRadio>
                  <ElRadio value="fixed" class="opc-check">固定期限</ElRadio>
                </ElRadioGroup>
              </ElFormItem>
              <ElFormItem v-if="formData.businessTermType === 'fixed'" label="营业期限截止" prop="businessTermEnd">
                <ElDatePicker v-model="formData.businessTermEnd" type="date" value-format="YYYY-MM-DD" class="w-full" />
              </ElFormItem>
              <ElFormItem label="经营范围" prop="businessScope">
                <ElInput v-model="formData.businessScope" type="textarea" :rows="4" class="opc-field" />
                <button type="button" class="opc-link-btn mt-2" @click="applyRecommendedScope">
                  使用推荐模板
                </button>
              </ElFormItem>
              <ChinaRegionSelect
                v-model:province="formData.registerProvince"
                v-model:city="formData.registerCity"
                v-model:district="formData.registerDistrict"
              />
              <ElFormItem label="详细地址" prop="registerAddress">
                <ElInput v-model.trim="formData.registerAddress" placeholder="街道门牌号" size="large" class="opc-field" />
              </ElFormItem>
              <ElFormItem label="地址证明" prop="addressProofFileId">
                <MemberFileUpload v-model="formData.addressProofFileId" accept=".pdf,.jpg,.jpeg,.png" />
              </ElFormItem>
            </div>
          </section>

          <!-- 法人信息 -->
          <section>
            <h3 class="opc-form-section__title">法人（股东）信息</h3>
            <div class="space-y-4">
              <ElFormItem label="法人姓名" prop="legalPersonName">
                <ElInput v-model.trim="formData.legalPersonName" size="large" class="opc-field" />
              </ElFormItem>
              <ElFormItem label="身份证号" prop="idCardNumber">
                <ElInput v-model.trim="formData.idCardNumber" maxlength="18" size="large" class="opc-field" />
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
                <ElInput v-model.trim="formData.ethnicity" size="large" class="opc-field" />
              </ElFormItem>
              <ElFormItem label="户籍地址" prop="householdAddress">
                <ElInput v-model.trim="formData.householdAddress" size="large" class="opc-field" />
              </ElFormItem>
              <ElFormItem label="现居住地址" prop="residenceAddress">
                <ElInput v-model.trim="formData.residenceAddress" size="large" class="opc-field" />
              </ElFormItem>
              <div class="grid md:grid-cols-2 gap-4">
                <ElFormItem label="手机号" prop="phone">
                  <ElInput v-model.trim="formData.phone" maxlength="11" size="large" class="opc-field" />
                </ElFormItem>
                <ElFormItem label="邮箱" prop="email">
                  <ElInput v-model.trim="formData.email" size="large" class="opc-field" />
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
            <h3 class="opc-form-section__title">授权确认</h3>
            <div class="space-y-3">
              <ElCheckbox v-model="formData.confirmations.infoTrue" class="opc-check w-full">本人确认以上信息真实、完整</ElCheckbox>
              <ElCheckbox v-model="formData.confirmations.authConsent" class="opc-check w-full">同意用于工商、税务、银行开户申报</ElCheckbox>
              <ElCheckbox v-model="formData.confirmations.opcLimitAck" class="opc-check w-full">知晓自然人 3 年内不得再设立新一人有限责任公司</ElCheckbox>
              <ElCheckbox v-model="formData.confirmations.eSignAuth" class="opc-check w-full">电子签名授权（沿用签约姓名）</ElCheckbox>
            </div>
          </section>

          <button
            type="button"
            class="opc-btn opc-btn--primary opc-btn--block"
            :disabled="submitting"
            @click="handleSubmitMaterials"
          >
            <ArtSvgIcon v-if="submitting" icon="ri:loader-4-line" class="text-lg animate-spin" />
            {{ submitting ? '提交中...' : '提交注册资料' }}
          </button>
        </ElForm>
      </section>

      <!-- Readonly materials -->
      <section v-else-if="showMaterialsOverview" class="opc-panel">
        <h2 class="opc-panel__section-title">已提交资料</h2>
        <p class="opc-panel__section-desc">以下为已提交申请资料的脱敏概览，点击可查看详情。</p>

        <div v-if="materialsLoading" class="opc-loading">
          <ArtSvgIcon icon="ri:loader-4-line" class="opc-loading__icon" />
        </div>
        <div v-else-if="materialsEntities.length === 0" class="opc-empty-text">
          资料审核中或已进入后续流程，如需修改请等待顾问联系。
        </div>
        <ul v-else class="opc-list">
          <li
            v-for="entity in materialsEntities"
            :key="entity.opcId"
            class="opc-list__item"
            @click="openMaterialsDetail(entity.opcId)"
          >
            <div class="opc-list__icon">
              <ArtSvgIcon icon="ri:building-2-line" />
            </div>
            <div class="opc-list__body">
              <div class="opc-list__row">
                <span class="opc-list__title">{{ entity.companyNameMasked }}</span>
                <span class="opc-badge" :class="materialsStatusClass(entity.status)">
                  {{ entity.statusLabel }}
                </span>
              </div>
              <p class="opc-list__desc">{{ entity.legalPersonSummary }}</p>
            </div>
            <ArtSvgIcon icon="ri:arrow-right-s-line" class="opc-list__arrow" />
          </li>
        </ul>
      </section>

      <!-- Materials detail dialog -->
      <ElDialog
        v-model="detailVisible"
        title="已提交资料详情"
        width="560px"
        class="opc-dialog"
        destroy-on-close
        @closed="resetDetailState"
      >
        <div v-if="detailLoading" class="opc-loading">
          <ArtSvgIcon icon="ri:loader-4-line" class="opc-loading__icon" />
        </div>
        <template v-else-if="detailData">
          <div class="mb-4">
            <span class="opc-badge" :class="materialsStatusClass(detailData.status)">
              {{ detailData.statusLabel }}
            </span>
          </div>

          <div v-for="section in displaySections" :key="section.key" class="opc-detail-section">
            <h3 class="opc-detail-section__title">{{ section.title }}</h3>
            <dl v-if="sectionDisplayFields(section).length" class="opc-detail-fields">
              <div v-for="field in sectionDisplayFields(section)" :key="field.label" class="opc-detail-field">
                <dt>{{ field.label }}</dt>
                <dd>{{ field.value || '—' }}</dd>
              </div>
            </dl>
            <div v-if="sectionDisplayAttachments(section).length" class="opc-detail-attachments">
              <div v-for="item in sectionDisplayAttachments(section)" :key="item.fileId" class="opc-detail-attachment">
                <p class="opc-detail-attachment__label">{{ item.label }}</p>
                <img
                  v-if="isImageMime(item.mimeType)"
                  :src="getAttachmentSrc(item)"
                  :alt="item.fileName"
                  class="opc-detail-attachment__img"
                  referrerpolicy="same-origin"
                  @error="refreshAttachmentUrl(item.fileId)"
                />
                <a
                  v-else
                  :href="getAttachmentSrc(item)"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="opc-link-btn"
                >
                  <ArtSvgIcon icon="ri:file-pdf-line" />
                  查看 {{ item.fileName || '附件' }}
                </a>
              </div>
            </div>
          </div>

          <button
            v-if="!revealed"
            type="button"
            class="opc-btn opc-btn--ghost opc-btn--block mt-6"
            @click="passwordVisible = true"
          >
            查看完整信息
          </button>
          <p v-else class="opc-hint mt-4">已展示完整信息，关闭弹窗后将恢复脱敏显示。</p>
        </template>
      </ElDialog>

      <!-- Password verify dialog -->
      <ElDialog
        v-model="passwordVisible"
        title="验证身份"
        width="400px"
        destroy-on-close
        @closed="passwordInput = ''"
      >
        <p class="opc-panel__desc mb-4">请输入登录密码以查看完整资料信息</p>
        <ElInput
          v-model="passwordInput"
          type="password"
          placeholder="请输入密码"
          size="large"
          class="opc-field"
          show-password
          @keyup.enter="handleReveal"
        />
        <template #footer>
          <button type="button" class="opc-btn opc-btn--ghost" @click="passwordVisible = false">
            取消
          </button>
          <button
            type="button"
            class="opc-btn opc-btn--primary"
            :disabled="revealSubmitting || !passwordInput"
            @click="handleReveal"
          >
            {{ revealSubmitting ? '验证中...' : '确认' }}
          </button>
        </template>
      </ElDialog>

      <!-- Bank section -->
      <section v-if="showBankSection" class="opc-panel">
        <h2 class="opc-panel__section-title">银行开户指引</h2>
        <ElCollapse>
          <ElCollapseItem title="开户材料清单（点击展开）" name="guide">
            <ul class="opc-guide-list">
              <li>营业执照正副本原件</li>
              <li>公章、财务章、法人章（刻章完成后）</li>
              <li>公司章程</li>
              <li>法人身份证原件</li>
              <li>税务登记相关回执</li>
            </ul>
          </ElCollapseItem>
        </ElCollapse>

        <div v-if="progress?.opcStatus === 'bank'" class="mt-6 space-y-4">
          <ElFormItem label="开户银行（选填）">
            <ElInput v-model.trim="bankForm.bankName" placeholder="如：中国工商银行 XX 支行" size="large" class="opc-field" />
          </ElFormItem>
          <ElFormItem label="上传开户回执">
            <MemberFileUpload v-model="bankForm.bankReceiptFileId" accept=".pdf,.jpg,.jpeg,.png" />
          </ElFormItem>
          <button
            type="button"
            class="opc-btn opc-btn--primary"
            :disabled="!bankForm.bankReceiptFileId || bankSubmitting"
            @click="handleSubmitBankReceipt"
          >
            提交开户回执
          </button>
        </div>
      </section>

      <!-- Active success -->
      <section v-if="progress?.opcStatus === 'active'" class="opc-panel opc-panel--success">
        <ArtSvgIcon icon="ri:checkbox-circle-line" class="opc-success__icon" />
        <h2 class="opc-panel__title">主体设立完成</h2>
        <p class="opc-panel__desc">收入台账、费用台账等功能已解锁</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import ChinaRegionSelect from '@/components/frontend/ChinaRegionSelect.vue'
import MemberFileUpload from '@/components/frontend/MemberFileUpload.vue'
import { getMemberAttachmentUrl } from '@/api/frontend/member/upload'
import {
  getOpcProgress,
  getMaterialsOverview,
  getMaterialsDetail,
  revealMaterialsAll,
  submitOpcMaterials,
  submitBankReceipt,
  RECOMMENDED_BUSINESS_SCOPE,
  PLAN_TIER_LABELS,
  type MaterialsEntityOverview,
  type MaterialsEntityDetail,
  type MaterialsDetailSection,
  type MaterialsSectionField,
  type MaterialsAttachmentReveal,
  type MaterialsRevealSection,
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
import { requireLogin } from '@/utils/auth/requireLogin'
import { sanitizeErrorMessage } from '@/utils/http/error'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

defineOptions({ name: 'ComplianceOpc' })

const memberStore = useMemberStore()
const verifyMock = isComplianceVerifyMock()

const loading = ref(true)
const submitting = ref(false)
const bankSubmitting = ref(false)
const progress = ref<OpcProgressResult | null>(null)
const formRef = ref<FormInstance>()

const materialsLoading = ref(false)
const materialsEntities = ref<MaterialsEntityOverview[]>([])
const detailVisible = ref(false)
const detailLoading = ref(false)
const detailData = ref<MaterialsEntityDetail | null>(null)
const activeOpcId = ref<number | string>(0)
const revealed = ref(false)
const revealedSections = ref<MaterialsRevealSection[]>([])
const passwordVisible = ref(false)
const passwordInput = ref('')
const revealSubmitting = ref(false)
const attachmentUrlOverrides = ref<Record<string, string>>({})
const attachmentUrlRefreshing = ref<Record<string, boolean>>({})

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

const signedPlanLabel = computed(() => {
  if (!progress.value?.planName && !progress.value?.planTier) return ''
  return progress.value.planName
    || PLAN_TIER_LABELS[progress.value.planTier ?? '']
    || progress.value.planTier
    || ''
})

const showMaterialsOverview = computed(() =>
  !showMaterialsForm.value && !!progress.value?.materialsSubmitted
)

const displaySections = computed(() => {
  if (revealed.value) {
    return revealedSections.value
  }
  return detailData.value?.sections ?? []
})

function sectionDisplayFields(section: MaterialsDetailSection | MaterialsRevealSection): MaterialsSectionField[] {
  return section.fields ?? []
}

function sectionDisplayAttachments(section: MaterialsDetailSection | MaterialsRevealSection): MaterialsAttachmentReveal[] {
  if (!revealed.value) return []
  return (section as MaterialsRevealSection).attachments ?? []
}

function isImageMime(mime?: string) {
  return !!mime && mime.startsWith('image/')
}

function getAttachmentSrc(item: MaterialsAttachmentReveal) {
  const key = String(item.fileId)
  return attachmentUrlOverrides.value[key] ?? item.accessUrl
}

async function refreshAttachmentUrl(fileId: number | string) {
  const key = String(fileId)
  if (attachmentUrlRefreshing.value[key]) return
  attachmentUrlRefreshing.value[key] = true
  try {
    const data = await getMemberAttachmentUrl(fileId)
    if (data?.url) {
      attachmentUrlOverrides.value[key] = data.url
    }
  } catch {
    // 忽略刷新失败，避免重复触发 error-handle
  } finally {
    attachmentUrlRefreshing.value[key] = false
  }
}

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

function materialsStatusClass(status: string) {
  const map: Record<string, string> = {
    approved: 'opc-badge--done',
    reviewing: 'opc-badge--current',
    rejected: 'opc-badge--rejected',
    pending: 'opc-badge--pending'
  }
  return map[status] ?? 'opc-badge--pending'
}

async function loadMaterialsOverview() {
  if (!showMaterialsOverview.value) {
    materialsEntities.value = []
    return
  }
  materialsLoading.value = true
  try {
    const res = await getMaterialsOverview()
    materialsEntities.value = res.entities ?? []
  } catch {
    materialsEntities.value = []
  } finally {
    materialsLoading.value = false
  }
}

async function openMaterialsDetail(opcId: number | string) {
  activeOpcId.value = opcId
  detailVisible.value = true
  detailLoading.value = true
  revealed.value = false
  revealedSections.value = []
  try {
    detailData.value = await getMaterialsDetail(opcId)
  } catch (e: unknown) {
    const msg = sanitizeErrorMessage(e instanceof Error ? e.message : '加载资料详情失败')
    ElMessage.error(msg)
    detailVisible.value = false
  } finally {
    detailLoading.value = false
  }
}

function resetDetailState() {
  detailData.value = null
  activeOpcId.value = 0
  revealed.value = false
  revealedSections.value = []
  attachmentUrlOverrides.value = {}
  attachmentUrlRefreshing.value = {}
}

async function handleReveal() {
  if (!passwordInput.value || !activeOpcId.value) return
  revealSubmitting.value = true
  try {
    const res = await revealMaterialsAll(passwordInput.value, activeOpcId.value)
    revealedSections.value = res.sections ?? []
    revealed.value = true
    passwordVisible.value = false
    ElMessage.success('验证成功')
  } catch (e: unknown) {
    const msg = sanitizeErrorMessage(e instanceof Error ? e.message : '密码验证失败')
    ElMessage.error(msg)
  } finally {
    revealSubmitting.value = false
  }
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
    await loadMaterialsOverview()
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
    const msg = sanitizeErrorMessage(e?.message || e?.msg || '提交失败，请检查表单后重试')
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
  if (!requireLogin({ redirect: '/user/compliance/opc' })) return
  loadProgress()
})
</script>

<style lang="scss" scoped>
.opc-page__inner {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.opc-panel {
  padding: 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.opc-panel__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 20px;
}

.opc-panel__label {
  margin: 0 0 6px;
  font-size: 12px;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.opc-panel__title {
  margin: 0 0 6px;
  font-size: 22px;
  font-weight: 800;
  color: #1a1f36;

  &--lg {
    font-size: 24px;
  }
}

.opc-panel__desc {
  margin: 0;
  font-size: 14px;
  color: #6b7c93;
}

.opc-panel__section-title {
  margin: 0 0 6px;
  font-size: 18px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-panel__section-desc {
  margin: 0 0 20px;
  font-size: 13px;
  color: #94a3b8;
}

.opc-panel--success {
  text-align: center;
  padding: 40px 24px;
}

.opc-tag {
  padding: 6px 12px;
  border-radius: 8px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 13px;
  font-weight: 700;

  &--price {
    font-size: 16px;
  }
}

.opc-timeline {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.opc-timeline__item {
  display: flex;
  gap: 12px;
  padding: 16px;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  background: #f8fafc;

  &.is-done {
    border-color: #bfdbfe;
    background: #eff6ff;

    .opc-timeline__index {
      background: #2563eb;
      color: #fff;
    }
  }

  &.is-current {
    border-color: #2563eb;
    box-shadow: 0 0 0 1px #2563eb;
  }
}

.opc-timeline__index {
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

.opc-timeline__body {
  flex: 1;
  min-width: 0;

  h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 700;
    color: #1a1f36;
  }

  p {
    margin: 4px 0 0;
    font-size: 12px;
    color: #94a3b8;
  }
}

.opc-timeline__row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.opc-badge {
  display: inline-flex;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;

  &--done {
    background: #dcfce7;
    color: #15803d;
  }

  &--current {
    background: #dbeafe;
    color: #2563eb;
  }

  &--pending {
    background: #f1f5f9;
    color: #64748b;
  }

  &--rejected {
    background: #ffedd5;
    color: #c2410c;
  }
}

.opc-alert {
  margin-top: 16px;
  padding: 12px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;

  &--warn {
    background: #fff7ed;
    border: 1px solid #fed7aa;
    color: #c2410c;
  }

  &--info {
    margin-bottom: 16px;
    background: #fffbeb;
    border: 1px solid #fde68a;
    color: #b45309;
  }
}

.opc-form-section__title {
  margin: 0 0 16px;
  font-size: 13px;
  font-weight: 700;
  color: #475569;
}

.opc-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &--primary {
    color: #fff;
    background: #2563eb;

    &:hover:not(:disabled) {
      background: #1d4ed8;
    }
  }

  &--ghost {
    color: #334155;
    background: #fff;
    border: 1px solid #d8dee9;

    &:hover:not(:disabled) {
      background: #f8fafc;
    }
  }

  &--block {
    width: 100%;
    padding: 12px 20px;
  }
}

.opc-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: none;
  background: none;
  padding: 0;
  font-size: 12px;
  font-weight: 700;
  color: #2563eb;
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }
}

.opc-loading {
  padding: 32px 0;
  text-align: center;
}

.opc-loading__icon {
  font-size: 28px;
  color: #2563eb;
  animation: spin 1s linear infinite;
}

.opc-empty-text {
  font-size: 14px;
  color: #64748b;
}

.opc-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.opc-list__item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  background: #f8fafc;
  cursor: pointer;
  transition: border-color 0.15s ease;

  &:hover {
    border-color: #2563eb;
  }
}

.opc-list__icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.opc-list__body {
  flex: 1;
  min-width: 0;
}

.opc-list__row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.opc-list__title {
  font-size: 14px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-list__desc {
  margin: 0;
  font-size: 13px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.opc-list__arrow {
  flex-shrink: 0;
  color: #94a3b8;
}

.opc-guide-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 14px;
  color: #334155;

  li::before {
    content: '· ';
    color: #94a3b8;
  }
}

.opc-success__icon {
  font-size: 48px;
  color: #16a34a;
  margin-bottom: 16px;
}

.opc-hint {
  font-size: 12px;
  color: #94a3b8;
}

.opc-detail-section {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
}

.opc-detail-section__title {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-detail-fields {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.opc-detail-field {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 12px;
  font-size: 13px;

  dt {
    font-weight: 600;
    color: #64748b;
  }

  dd {
    margin: 0;
    color: #1a1f36;
    word-break: break-all;
  }
}

.opc-detail-attachments {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 12px;
}

.opc-detail-attachment {
  padding: 14px;
  border: 1px solid #e8edf3;
  border-radius: 8px;
  background: #f8fafc;
}

.opc-detail-attachment__label {
  margin: 0 0 10px;
  font-size: 13px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-detail-attachment__img {
  width: 100%;
  max-height: 288px;
  object-fit: contain;
  border-radius: 8px;
  background: #fff;
}

:deep(.opc-field) .el-input__wrapper,
:deep(.opc-field) .el-select__wrapper {
  min-height: 40px;
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 0 0 1px #d8dee9 inset;
}

:deep(.opc-field) .el-textarea__inner {
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 0 0 1px #d8dee9 inset;
}

:deep(.opc-check) {
  margin-right: 0 !important;
  height: auto;
  padding: 10px 12px;
  border: 1px solid #e8edf3;
  border-radius: 8px;
  background: #f8fafc;
}

:deep(.opc-dialog) .el-dialog {
  border-radius: 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
