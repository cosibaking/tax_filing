<!-- +----------------------------------------------------------------------
  | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
  +----------------------------------------------------------------------
  | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
  +----------------------------------------------------------------------
  | Licensed ( https://opensource.org/licenses/MIT )
  +----------------------------------------------------------------------
  | Author: 喜羊羊 <751300685@qq.com>
  +---------------------------------------------------------------------- -->
<template>
  <div class="min-h-[60vh] flex items-center justify-center">
    <div class="text-center">
      <ArtSvgIcon icon="ri:loader-4-line" class="text-4xl text-clay-accent animate-spin" />
      <p class="mt-4 text-sm font-bold text-clay-muted">{{ message }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * 前台路由加载桥接页（对齐 BuildAdmin loading.vue）
 * 
 * 流程：调菜单 API → 注册动态路由 → 导航到目标页
 */
import { getMemberMenus, type MemberMenuItem } from '@/api/frontend/member/user'
import { useMemberMenuStore } from '@/store/modules/memberMenu'

defineOptions({ name: 'FrontendLoading' })

const router = useRouter()
const route = useRoute()
const message = ref('加载中...')

const viewModules = import.meta.glob('/src/views/frontend/**/*.vue')

onMounted(async () => {
  console.log('[FrontendLoading] 开始加载, params:', route.params)

  // 解析目标路径
  let targetPath = '/'
  let targetQuery: Record<string, any> = {}
  const toParam = route.params.to as string
  if (toParam) {
    try {
      const parsed = JSON.parse(decodeURIComponent(toParam))
      targetPath = parsed.path || '/'
      targetQuery = parsed.query || {}
    } catch {
      console.warn('[FrontendLoading] 解析 to 参数失败:', toParam)
    }
  }

  console.log('[FrontendLoading] 目标路径:', targetPath)

  try {
    // 1. 拉取菜单
    const data = await getMemberMenus()
    console.log('[FrontendLoading] API 返回:', data)

    const memberMenuStore = useMemberMenuStore()
    memberMenuStore.navMenus = data.nav || []
    memberMenuStore.centerMenus = data.menus || []
    memberMenuStore.routeRules = data.rules || []
    memberMenuStore.loaded = true

    // 2. 注册所有路由
    const allItems = [...(data.nav || []), ...(data.menus || []), ...(data.rules || [])]
    let registered = 0
    for (const item of allItems) {
      if (registerRoute(item)) registered++
    }
    console.log(`[FrontendLoading] 注册了 ${registered} 条路由`)

    // 3. 把 catch-all 路由移到最后（让动态注册的路由优先匹配）
    if (router.hasRoute('FrontendCatchAll')) {
      const catchAllRoute = router.getRoutes().find(r => r.name === 'FrontendCatchAll')
      router.removeRoute('FrontendCatchAll')
      if (catchAllRoute) {
        router.addRoute('FrontendLayout', {
          path: ':path(.*)*',
          name: 'FrontendCatchAll',
          redirect: (to: any) => ({
            name: 'FrontendLoading',
            params: { to: JSON.stringify({ path: to.path, query: to.query }) },
          }),
        })
      }
    }

    // 4. 导航到目标路径（replace 触发完整导航，确保动态注册的路由被重新解析）
    console.log('[FrontendLoading] 导航到:', targetPath)
    message.value = '跳转中...'
    await nextTick()
    await router.replace({ path: targetPath, query: targetQuery })
  } catch (e) {
    console.error('[FrontendLoading] 加载失败:', e)
    message.value = '加载失败，返回首页...'
    setTimeout(() => router.replace('/'), 1500)
  }
})

function registerRoute(item: MemberMenuItem): boolean {
  if (!item.path || !item.component) return false
  if (item.type === 'button') return false
  if (item.menuType === 'link') return false

  const compPath = `/src/views/frontend/${item.component}.vue`
  const compModule = viewModules[compPath]
  if (!compModule) {
    console.warn(`[FrontendLoading] 组件未找到: ${compPath}`)
    return false
  }

  let path = item.path
  if (path.startsWith('/')) path = path.substring(1)

  const routeName = `Frontend_${item.name || item.id}`
  if (router.hasRoute(routeName)) return false

  router.addRoute('FrontendLayout', {
    path,
    name: routeName,
    component: compModule,
    meta: {
      title: item.title,
      icon: item.icon,
      menuType: item.menuType,
      type: item.type,
    }
  })
  return true
}
</script>

<style scoped>
.text-clay-accent { color: #5a8dee; }
.text-clay-muted { color: #8898aa; }
</style>
