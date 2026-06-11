<!-- 社保指引 P-SG -->
<template>
  <section class="space-y-6">
    <EmploymentStatusCard />

    <div class="member-panel">
      <div class="member-panel__head">
        <div>
          <h2 class="member-panel__title">社保指引</h2>
          <p class="member-panel__desc">了解创始人参保路径与雇员社保义务</p>
        </div>
      </div>

      <div v-if="loading" class="member-empty">
        <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
      </div>

      <div v-else class="member-card-grid">
        <button
          v-for="item in guides"
          :key="item.slug"
          type="button"
          class="member-card-item"
          @click="openGuide(item.slug)"
        >
          <h3>{{ item.title }}</h3>
          <p class="line-clamp-2">{{ item.summary }}</p>
        </button>
      </div>
    </div>

    <ElDrawer v-model="drawerVisible" :title="detail?.title || '指引详情'" size="520px">
      <div v-if="detailLoading" class="py-8 text-center">
        <ArtSvgIcon icon="ri:loader-4-line" class="text-2xl animate-spin mx-auto text-clay-accent" />
      </div>
      <div v-else-if="detail" class="social-guide-md">
        <MdPreview
          :model-value="detail.contentMd"
          :editor-id="previewId"
          preview-theme="github"
          :no-mermaid="true"
          :no-katex="true"
        />
      </div>
    </ElDrawer>
  </section>
</template>

<script setup lang="ts">
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import EmploymentStatusCard from '@/components/member/EmploymentStatusCard.vue'
import { getSocialGuides, getSocialGuideDetail, type SocialGuideDetail } from '@/api/frontend/compliance/social'

defineOptions({ name: 'ComplianceSocialGuide' })

const loading = ref(true)
const guides = ref<Awaited<ReturnType<typeof getSocialGuides>>['list']>([])
const drawerVisible = ref(false)
const detailLoading = ref(false)
const detail = ref<SocialGuideDetail | null>(null)
const previewId = 'social-guide-preview'

async function loadGuides() {
  loading.value = true
  try {
    const res = await getSocialGuides()
    guides.value = res.list || []
  } finally {
    loading.value = false
  }
}

async function openGuide(slug: string) {
  drawerVisible.value = true
  detailLoading.value = true
  detail.value = null
  try {
    detail.value = await getSocialGuideDetail(slug)
  } finally {
    detailLoading.value = false
  }
}

onMounted(loadGuides)
</script>

<style scoped>
:deep(.social-guide-md) {
  font-size: 14px;
  line-height: 1.7;
  color: var(--el-text-color-primary);
}
</style>
