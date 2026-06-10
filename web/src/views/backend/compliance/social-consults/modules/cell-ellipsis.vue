<!-- 表格单元格：多行截断 + hover 展示全文 -->
<template>
  <ElTooltip
    :content="text"
    placement="top-start"
    :disabled="!text || text === '—'"
    :show-after="300"
  >
    <span
      class="cell-ellipsis"
      :class="lines > 1 ? 'cell-ellipsis--multi' : 'cell-ellipsis--single'"
    >{{ text || '—' }}</span>
  </ElTooltip>
</template>

<script setup lang="ts">
defineOptions({ name: 'SocialConsultCellEllipsis' })

withDefaults(
  defineProps<{
    text?: string
    lines?: number
  }>(),
  {
    text: '',
    lines: 1
  }
)
</script>

<style scoped>
.cell-ellipsis {
  display: -webkit-box;
  overflow: hidden;
  -webkit-box-orient: vertical;
  word-break: break-word;
  line-height: 1.5;
}

.cell-ellipsis--single {
  -webkit-line-clamp: 1;
}

.cell-ellipsis--multi {
  -webkit-line-clamp: 2;
}
</style>
