<template>
  <ElDialog
    :model-value="visible"
    title="标记已申报"
    width="520px"
    destroy-on-close
    @update:model-value="emit('update:visible', $event)"
  >
    <p class="text-sm text-gray-500 mb-4">
      申报前须完成全部 9 项自查清单，并上传回执 PDF。
    </p>

    <div class="space-y-3 mb-6">
      <label
        v-for="item in FILING_CHECKLIST_ITEMS"
        :key="item.key"
        class="flex items-start gap-2 text-sm"
      >
        <ElCheckbox v-model="checklist[item.key]" />
        <span>{{ item.label }}</span>
      </label>
    </div>

    <ElForm label-position="top">
      <ElFormItem label="实际申报税额（元，选填）">
        <ElInputNumber v-model="filedAmount" :min="0" :precision="2" class="w-full" />
      </ElFormItem>
      <ElFormItem label="回执 PDF" required>
        <ArtFileUpload v-model="receiptFileId" accept=".pdf" />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <ElButton @click="emit('update:visible', false)">取消</ElButton>
      <ElButton type="primary" :loading="submitting" :disabled="!canSubmit" @click="handleSubmit">
        确认标记
      </ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
import ArtFileUpload from '@/components/core/forms/art-file-upload/index.vue'
import { markFilingFiled, FILING_CHECKLIST_ITEMS, type FilingTaskItem } from '@/api/backend/compliance'
import { ElMessage } from 'element-plus'

interface Props {
  visible: boolean
  task: FilingTaskItem | null
}
interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const submitting = ref(false)
const receiptFileId = ref('')
const filedAmount = ref<number | undefined>()
const checklist = reactive<Record<string, boolean>>(
  Object.fromEntries(FILING_CHECKLIST_ITEMS.map(i => [i.key, false]))
)

const allChecked = computed(() => FILING_CHECKLIST_ITEMS.every(i => checklist[i.key]))
const canSubmit = computed(() => allChecked.value && !!receiptFileId.value)

watch(() => props.visible, (val) => {
  if (val) {
    FILING_CHECKLIST_ITEMS.forEach(i => { checklist[i.key] = false })
    receiptFileId.value = ''
    filedAmount.value = props.task?.calculatedAmount
  }
})

async function handleSubmit() {
  if (!props.task || !canSubmit.value) return
  submitting.value = true
  try {
    await markFilingFiled({
      id: props.task.id,
      checklist: { ...checklist },
      filedAmount: filedAmount.value,
      receiptFileId: receiptFileId.value
    })
    ElMessage.success('已标记为已申报')
    emit('update:visible', false)
    emit('success')
  } catch {
    ElMessage.error('操作失败，请重试')
  } finally {
    submitting.value = false
  }
}
</script>
