<!-- +----------------------------------------------------------------------
  | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
  +----------------------------------------------------------------------
  | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
  +----------------------------------------------------------------------
  | Licensed ( https://opensource.org/licenses/MIT )
  +----------------------------------------------------------------------
  | Author: 喜羊羊 <751300685@qq.com>
  +---------------------------------------------------------------------- -->
<!-- 系统logo -->
<template>
  <div class="flex-cc">
    <img 
      :style="logoStyle" 
      :src="logoUrl" 
      alt="logo" 
      class="w-full h-full object-contain" 
    />
  </div>
</template>

<script setup lang="ts">
  import { useSiteStore } from '@/store/modules/site'
  import defaultLogo from '@imgs/common/jinshui-logo.png'

  defineOptions({ name: 'ArtLogo' })

  interface Props {
    /** logo 大小 */
    size?: number | string
  }

  const props = withDefaults(defineProps<Props>(), {
    size: 36
  })

  const siteStore = useSiteStore()
  const logoUrl = computed(() => siteStore.getLogo() || defaultLogo)
  const logoStyle = computed(() => ({ width: `${props.size}px` }))
</script>
