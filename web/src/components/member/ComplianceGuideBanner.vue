<!-- +----------------------------------------------------------------------
  | 合规全局引导条（M8 §5.3）
  +---------------------------------------------------------------------- -->
<template>
  <div
    class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 px-5 py-3 rounded-lg border font-bold text-sm"
    :class="bannerClass"
  >
    <div class="flex items-center gap-3">
      <ArtSvgIcon :icon="icon" class="text-xl shrink-0" />
      <span>{{ banner.message }}</span>
    </div>
    <RouterLink
      v-if="banner.actionPath"
      :to="banner.actionPath"
      class="inline-flex items-center gap-2 px-4 py-2 rounded-md bg-white/80 border border-current/10 hover:bg-white transition-all shrink-0"
    >
      {{ banner.actionLabel || '查看' }}
      <ArtSvgIcon icon="ri:arrow-right-s-line" class="text-base" />
    </RouterLink>
  </div>
</template>

<script setup lang="ts">
export interface GuideBanner {
  type: 'warning' | 'info' | 'danger'
  message: string
  actionLabel?: string
  actionPath?: string
}

defineOptions({ name: 'ComplianceGuideBanner' })

const props = defineProps<{ banner: GuideBanner }>()

const bannerClass = computed(() => {
  switch (props.banner.type) {
    case 'warning':
      return 'bg-amber-50 border-amber-200 text-amber-800'
    case 'danger':
      return 'bg-red-50 border-red-200 text-red-700'
    default:
      return 'bg-blue-50 border-blue-200 text-blue-700'
  }
})

const icon = computed(() => {
  switch (props.banner.type) {
    case 'warning':
      return 'ri:alert-line'
    case 'danger':
      return 'ri:error-warning-line'
    default:
      return 'ri:information-line'
  }
})
</script>
