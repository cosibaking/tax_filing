<template>

  <ElDialog

    v-model="visible"

    title="OPC 任务详情"

    width="720px"

    destroy-on-close

    @closed="handleClosed"

  >

    <div v-if="loading" class="py-12 text-center">

      <ElIcon class="is-loading text-2xl"><Loading /></ElIcon>

    </div>



    <template v-else-if="detail">

      <ElDescriptions :column="2" border size="small" class="mb-4">

        <ElDescriptionsItem label="任务 ID">{{ detail.id }}</ElDescriptionsItem>

        <ElDescriptionsItem label="当前节点">

          <ElTag>{{ OPC_STATUS_LABELS[detail.status] || detail.status }}</ElTag>

        </ElDescriptionsItem>

        <ElDescriptionsItem label="主播">{{ detail.memberName || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="手机">{{ detail.memberPhone || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="套餐">{{ detail.planName || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="提交日">{{ detail.materialsSubmittedAt || '—' }}</ElDescriptionsItem>

      </ElDescriptions>



      <ElDivider content-position="left">拟设公司</ElDivider>

      <ElDescriptions :column="1" border size="small" class="mb-4">

        <ElDescriptionsItem label="备选名称">{{ (detail.proposedNames || []).join(' / ') || detail.proposedName || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="注册资本">{{ detail.registeredCapital ? `${detail.registeredCapital} 万元` : '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="经营范围">{{ detail.businessScope || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="注册地址">{{ detail.registerAddress || '—' }}</ElDescriptionsItem>

      </ElDescriptions>



      <ElDivider content-position="left">法人信息（脱敏）</ElDivider>

      <ElDescriptions :column="2" border size="small" class="mb-4">

        <ElDescriptionsItem label="姓名">{{ detail.legalPersonName || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="身份证">{{ detail.idCardNumberMasked || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="手机">{{ detail.phone || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="邮箱">{{ detail.email || '—' }}</ElDescriptionsItem>

      </ElDescriptions>



      <ElDivider v-if="detail.companyName || detail.creditCode" content-position="left">工商结果</ElDivider>

      <ElDescriptions v-if="detail.companyName || detail.creditCode" :column="2" border size="small" class="mb-4">

        <ElDescriptionsItem label="核准名称">{{ detail.companyName || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="信用代码">{{ detail.creditCode || '—' }}</ElDescriptionsItem>

        <ElDescriptionsItem label="成立日">{{ detail.establishedAt || '—' }}</ElDescriptionsItem>

      </ElDescriptions>



      <ElDivider v-if="detail.progressLogs?.length" content-position="left">进度时间轴</ElDivider>

      <ElTimeline v-if="detail.progressLogs?.length" class="mb-4 pl-1">

        <ElTimelineItem

          v-for="(log, idx) in detail.progressLogs"

          :key="idx"

          :timestamp="log.createdAt"

          placement="top"

        >

          <span class="font-medium">{{ log.step }}</span>

          <span class="text-g-500 ml-2">{{ log.status }}</span>

          <p v-if="log.note" class="text-sm text-g-600 mt-1">{{ log.note }}</p>

        </ElTimelineItem>

      </ElTimeline>



      <ElDivider v-if="detail.rejectNote" content-position="left">驳回备注</ElDivider>

      <p v-if="detail.rejectNote" class="text-orange-600 text-sm mb-4">{{ detail.rejectNote }}</p>



      <div v-if="showActions" class="mt-4 p-4 bg-gray-50 rounded-lg space-y-4">

        <template v-if="detail.status === 'materials_review'">

          <ElInput v-model="actionForm.note" type="textarea" :rows="2" placeholder="驳回时必填备注（≥10字）" />

          <ElSpace>

            <ElButton type="success" :loading="acting" @click="doAction('approve_materials')">通过资料审核</ElButton>

            <ElButton type="warning" :loading="acting" @click="doAction('reject_materials')">驳回资料</ElButton>

          </ElSpace>

        </template>



        <template v-else-if="detail.status === 'registering'">

          <ElForm label-width="100px" size="small">

            <ElFormItem label="核准名称"><ElInput v-model="actionForm.companyName" /></ElFormItem>

            <ElFormItem label="信用代码"><ElInput v-model="actionForm.creditCode" maxlength="18" /></ElFormItem>

            <ElFormItem label="成立日期"><ElDatePicker v-model="actionForm.establishedAt" type="date" value-format="YYYY-MM-DD" class="w-full" /></ElFormItem>

            <ElFormItem label="执照附件"><ArtFileUpload v-model="actionForm.licenseFileId" accept=".pdf,.jpg,.jpeg,.png" compact /></ElFormItem>

          </ElForm>

          <ElButton type="primary" :loading="acting" @click="doAction('issue_license')">执照已下发</ElButton>

        </template>



        <template v-else-if="detail.status === 'tax'">

          <ElFormItem label="税务激活日">

            <ElDatePicker v-model="actionForm.taxActivatedAt" type="date" value-format="YYYY-MM-DD" />

          </ElFormItem>

          <ElButton type="primary" :loading="acting" @click="doAction('complete_tax')">税务登记完成</ElButton>

        </template>



        <template v-else-if="detail.status === 'bank'">

          <ElForm label-width="100px" size="small">

            <ElFormItem label="对公账号"><ElInput v-model="actionForm.bankAccount" /></ElFormItem>

            <ElFormItem label="回执附件"><ArtFileUpload v-model="actionForm.bankReceiptFileId" accept=".pdf,.jpg,.jpeg,.png" compact /></ElFormItem>

          </ElForm>

          <ElButton type="primary" :loading="acting" @click="doAction('complete_bank')">银行开户完成</ElButton>

        </template>

      </div>

    </template>



    <template #footer>

      <ElButton @click="visible = false">关闭</ElButton>

    </template>

  </ElDialog>

</template>



<script setup lang="ts">

import { Loading } from '@element-plus/icons-vue'

import ArtFileUpload from '@/components/core/forms/art-file-upload/index.vue'

import { useAuth } from '@/hooks/core/useAuth'

import {

  getOpcTaskDetail,

  patchOpcTask,

  OPC_STATUS_LABELS,

  type OpcTaskDetail

} from '@/api/backend/compliance'



interface Props {

  taskId?: number | string

}



const props = defineProps<Props>()

const emit = defineEmits<{ success: [] }>()



const { hasAuth } = useAuth()

const visible = defineModel<boolean>('visible', { default: false })



const loading = ref(false)

const acting = ref(false)

const detail = ref<OpcTaskDetail | null>(null)



const actionForm = reactive({

  note: '',

  companyName: '',

  creditCode: '',

  establishedAt: '',

  licenseFileId: '',

  taxActivatedAt: '',

  bankAccount: '',

  bankReceiptFileId: ''

})



const showActions = computed(() =>

  hasAuth('advance') &&

  detail.value &&

  ['materials_review', 'registering', 'tax', 'bank'].includes(detail.value.status)

)



async function loadDetail() {

  if (!props.taskId) return

  loading.value = true

  try {

    detail.value = await getOpcTaskDetail(props.taskId)

    actionForm.companyName = detail.value.companyName || ''

    actionForm.creditCode = detail.value.creditCode || ''

  } catch {

    ElMessage.error('加载任务详情失败')

    detail.value = null

  } finally {

    loading.value = false

  }

}



async function doAction(action: string) {

  if (!detail.value) return

  if (action === 'reject_materials' && actionForm.note.trim().length < 10) {

    ElMessage.warning('驳回备注至少 10 字')

    return

  }

  acting.value = true

  try {

    await patchOpcTask({

      id: detail.value.id,

      action: action as any,

      note: actionForm.note || undefined,

      companyName: actionForm.companyName || undefined,

      creditCode: actionForm.creditCode || undefined,

      establishedAt: actionForm.establishedAt || undefined,

      licenseFileId: actionForm.licenseFileId || undefined,

      taxActivatedAt: actionForm.taxActivatedAt || undefined,

      bankAccount: actionForm.bankAccount || undefined,

      bankReceiptFileId: actionForm.bankReceiptFileId || undefined

    })

    ElMessage.success('操作成功')

    emit('success')

    actionForm.note = ''

    await loadDetail()

  } catch {

    ElMessage.error('操作失败')

  } finally {

    acting.value = false

  }

}



function handleClosed() {

  detail.value = null

  Object.assign(actionForm, {

    note: '', companyName: '', creditCode: '', establishedAt: '',

    licenseFileId: '', taxActivatedAt: '', bankAccount: '', bankReceiptFileId: ''

  })

}



watch(

  () => [visible.value, props.taskId] as const,

  ([v, id]) => {

    if (v && id) loadDetail()

  }

)

</script>


