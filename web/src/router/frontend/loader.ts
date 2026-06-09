// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

/**
 * 前台动态路由加载器
 *
 * 对齐 BuildAdmin 的前台路由加载方案：
 * 1. 调用菜单 API 获取菜单数据
 * 2. 用 import.meta.glob 扫描所有前台 Vue 组件
 * 3. 根据菜单的 component 字段动态映射到 Vue 组件
 * 4. 调用 router.addRoute 注册到 FrontendLayout 下
 */

import type { Router, RouteRecordRaw } from 'vue-router'
import { getMemberMenus, type MemberMenuItem } from '@/api/frontend/member/user'
import { useMemberMenuStore } from '@/store/modules/memberMenu'

// 用 import.meta.glob 扫描所有前台页面组件
const viewModules = import.meta.glob('/src/views/frontend/**/*.vue')

// 已从门户导航移除的菜单（不再注册动态路由）
const hiddenNavNames = new Set(['cases', 'community', 'changelog'])

// 标记是否已注册
let isRegistered = false
// 存储已注册的路由名称（用于卸载）
const registeredRouteNames: string[] = []

/**
 * 加载并注册前台动态路由
 * @returns 是否成功加载
 */
export async function loadFrontendRoutes(router: Router): Promise<boolean> {
  if (isRegistered) return true

  try {
    // 1. 拉取菜单（白名单接口，登录/未登录都能调）
    const data = await getMemberMenus()
    const allItems = [...(data.nav || []), ...(data.menus || []), ...(data.rules || [])]

    // 2. 同步到 Store
    const menuStore = useMemberMenuStore()
    menuStore.navMenus = data.nav || []
    menuStore.centerMenus = data.menus || []
    menuStore.routeRules = data.rules || []
    menuStore.loaded = true

    // 3. 注册动态路由
    for (const item of allItems) {
      if (hiddenNavNames.has(item.name)) continue
      if (!item.component && item.type === 'button') continue
      if (!item.path) continue

      const route = buildRoute(item)
      if (route) {
        router.addRoute('FrontendLayout', route)
        if (route.name) {
          registeredRouteNames.push(route.name as string)
        }
      }
    }

    isRegistered = true
    return true
  } catch (e) {
    console.error('[FrontendRouteLoader] 加载前台菜单失败:', e)
    return false
  }
}

/**
 * 根据菜单项构建路由记录
 */
function buildRoute(item: MemberMenuItem): RouteRecordRaw | null {
  // 路径：去掉开头的 /（因为是作为 FrontendLayout 的子路由）
  let path = item.path || ''
  if (path.startsWith('/')) {
    path = path.substring(1)
  }
  if (!path) return null

  // 组件映射
  let component: any = null
  if (item.component) {
    // component 字段格式：'docs/index' → '/src/views/frontend/docs/index.vue'
    const compPath = `/src/views/frontend/${item.component}.vue`
    if (viewModules[compPath]) {
      component = viewModules[compPath]
    } else {
      console.warn(`[FrontendRouteLoader] 组件未找到: ${compPath} (菜单: ${item.title})`)
      return null
    }
  } else if (item.menuType === 'link') {
    // 外链不注册路由
    return null
  } else {
    // 没有 component 的非外链菜单，跳过
    return null
  }

  const routeName = `Frontend_${item.name || item.id}`

  return {
    path,
    name: routeName,
    component,
    meta: {
      title: item.title,
      icon: item.icon,
      menuType: item.menuType,
      requiresAuth: Number(item.noLoginValid) !== 1,
    }
  }
}

/**
 * 卸载所有动态注册的前台路由
 */
export function unloadFrontendRoutes(router: Router): void {
  for (const name of registeredRouteNames) {
    if (router.hasRoute(name)) {
      router.removeRoute(name)
    }
  }
  registeredRouteNames.length = 0
  isRegistered = false
}

/**
 * 是否已加载
 */
export function isFrontendRoutesLoaded(): boolean {
  return isRegistered
}

/**
 * 重置状态（登出时调用）
 */
export function resetFrontendRouteState(): void {
  isRegistered = false
  registeredRouteNames.length = 0
}
