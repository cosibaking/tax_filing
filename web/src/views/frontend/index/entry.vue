<!-- +----------------------------------------------------------------------
  | 门户首页入口：未登录展示免费诊断，已登录展示营销主页
  +---------------------------------------------------------------------- -->
<template>
  <DiagnosisIndex v-if="!isLoggedIn" />
  <MarketingHome v-else />
</template>

<script setup lang="ts">
import DiagnosisIndex from '../diagnosis/index.vue'
import MarketingHome from './index.vue'
import { useMemberStore } from '@/store/modules/member'

defineOptions({ name: 'FrontendHomeEntry' })

const memberStore = useMemberStore()
const isLoggedIn = computed(() => memberStore.isLogin)
</script>
