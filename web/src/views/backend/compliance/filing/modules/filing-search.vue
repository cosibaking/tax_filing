<template>
  <ArtSearchBar
    ref="searchBarRef"
    v-model="formData"
    :items="formItems"
    :rules="rules"
    @reset="handleReset"
    @search="handleSearch"
  />
</template>

<script setup lang="ts">
import { FILING_TAX_TYPE_OPTIONS, FILING_STATUS_OPTIONS } from '@/api/backend/compliance'

interface Props {
  modelValue: Record<string, any>
}
interface Emits {
  (e: 'update:modelValue', value: Record<string, any>): void
  (e: 'search', params: Record<string, any>): void
  (e: 'reset'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const searchBarRef = ref()
const formData = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const rules = {}

const formItems = computed(() => [
  {
    label: '关键词',
    key: 'q',
    type: 'input',
    placeholder: '主播名 / 手机号 / OPC',
    clearable: true
  },
  {
    label: '税种',
    key: 'taxType',
    type: 'select',
    props: {
      placeholder: '全部',
      options: FILING_TAX_TYPE_OPTIONS.map(o => ({ label: o.label, value: o.value }))
    }
  },
  {
    label: '截止月',
    key: 'period',
    type: 'date',
    props: { type: 'month', valueFormat: 'YYYY-MM', placeholder: '选择月份' }
  },
  {
    label: '状态',
    key: 'status',
    type: 'select',
    props: {
      placeholder: '全部',
      options: FILING_STATUS_OPTIONS.map(o => ({ label: o.label, value: o.value }))
    }
  }
])

function handleReset() {
  emit('reset')
}

async function handleSearch() {
  await searchBarRef.value?.validate?.()
  emit('search', formData.value)
}
</script>
