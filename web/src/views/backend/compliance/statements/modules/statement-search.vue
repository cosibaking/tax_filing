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

const statusOptions = [
  { label: '全部', value: '' },
  { label: '草稿', value: 'draft' },
  { label: '已发送', value: 'sent' },
  { label: '已完成', value: 'completed' }
]

const formItems = computed(() => [
  {
    label: '关键词',
    key: 'q',
    type: 'input',
    placeholder: '主播名 / 手机号',
    clearable: true
  },
  {
    label: '月份',
    key: 'period',
    type: 'date',
    props: { type: 'month', valueFormat: 'YYYY-MM', placeholder: '选择月份' }
  },
  {
    label: '状态',
    key: 'status',
    type: 'select',
    props: { placeholder: '全部', options: statusOptions }
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
