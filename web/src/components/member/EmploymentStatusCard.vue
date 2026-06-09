<!-- 用工状态采集卡片（F-91） -->
<template>
  <div
    v-if="visible"
    class="p-6 rounded-[32px] border mb-6"
    :class="isUnknown ? 'bg-amber-50 border-amber-200' : 'bg-[#f0f3f8] border-[#d1d9e6]/40 shadow-clay-pressed'"
  >
    <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
      <div>
        <h3 class="font-heading font-black text-lg text-clay-foreground">
          {{ isUnknown ? '请确认用工情况' : '用工状态' }}
        </h3>
        <p class="text-sm text-clay-muted mt-1">
          <template v-if="isUnknown">
            一人 OPC 通常无雇员。确认后将影响报税前自查清单中工资/社保项是否适用。
          </template>
          <template v-else>
            当前：{{ EMPLOYMENT_STATUS_LABELS[status] }}
            <span v-if="confirmedAt" class="ml-2">（{{ confirmedAt }} 确认）</span>
          </template>
        </p>
      </div>
      <div v-if="!isUnknown && status === 'no_employee'" class="shrink-0">
        <RouterLink
          to="/user/compliance/social-guide"
          class="text-sm font-bold text-clay-accent hover:underline"
        >
          查看社保指引 →
        </RouterLink>
      </div>
    </div>

    <div v-if="isUnknown || editing" class="mt-4 flex flex-wrap gap-3">
      <button
        type="button"
        class="px-5 py-2.5 rounded-2xl font-bold text-sm transition-all"
        :class="selected === 'no_employee' ? 'bg-clay-accent text-white shadow-clay-btn' : 'bg-white shadow-clay-pressed text-clay-foreground'"
        :disabled="submitting"
        @click="selected = 'no_employee'"
      >
        目前无雇员
      </button>
      <button
        type="button"
        class="px-5 py-2.5 rounded-2xl font-bold text-sm transition-all"
        :class="selected === 'has_employee' ? 'bg-clay-accent text-white shadow-clay-btn' : 'bg-white shadow-clay-pressed text-clay-foreground'"
        :disabled="submitting"
        @click="selected = 'has_employee'"
      >
        有雇员（助理/剪辑等）
      </button>
      <button
        v-if="selected"
        type="button"
        class="px-5 py-2.5 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-bold text-sm shadow-clay-btn"
        :disabled="submitting"
        @click="handleSubmit"
      >
        {{ submitting ? '保存中…' : '确认' }}
      </button>
    </div>

    <p v-if="status === 'has_employee'" class="mt-3 text-xs text-orange-700 font-medium">
      员工档案与工资申报功能将在后续版本开放，如有紧急需求请提交「社保咨询」。
    </p>
  </div>
</template>

<script setup lang="ts">
import {
  getEmploymentStatus,
  setEmploymentStatus,
  EMPLOYMENT_STATUS_LABELS,
  type EmploymentStatus
} from '@/api/frontend/compliance/social'
import { ElMessage } from 'element-plus'

defineOptions({ name: 'EmploymentStatusCard' })

const props = withDefaults(defineProps<{
  /** 仅 unknown 时展示（用于嵌入申报页等） */
  onlyWhenUnknown?: boolean
}>(), {
  onlyWhenUnknown: false
})

const emit = defineEmits<{ saved: [status: EmploymentStatus] }>()

const loading = ref(true)
const submitting = ref(false)
const status = ref<EmploymentStatus>('unknown')
const confirmedAt = ref('')
const selected = ref<'no_employee' | 'has_employee' | ''>('')
const editing = ref(false)

const isUnknown = computed(() => status.value === 'unknown')
const visible = computed(() => !props.onlyWhenUnknown || isUnknown.value)

async function load() {
  loading.value = true
  try {
    const res = await getEmploymentStatus()
    status.value = res.employmentStatus || 'unknown'
    confirmedAt.value = res.employmentConfirmedAt || ''
  } catch {
    status.value = 'unknown'
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  if (!selected.value) return
  submitting.value = true
  try {
    const res = await setEmploymentStatus(selected.value)
    status.value = res.employmentStatus
    confirmedAt.value = res.employmentConfirmedAt || ''
    editing.value = false
    ElMessage.success('用工状态已保存')
    emit('saved', status.value)
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>
