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
 * 后台路由模块汇总
 *
 * 生产环境菜单由后端 xy_admin_menu 控制；此处仅保留 SaaS 基础路由供前端模式或动态组件加载。
 */
import { AppRouteRecord } from '@/types/router'
import { dashboardRoutes } from './dashboard'
import { systemRoutes } from './system'
import { safeguardRoutes } from './safeguard'

export const backendRoutes: AppRouteRecord[] = [
  dashboardRoutes,
  systemRoutes,
  safeguardRoutes
]

export const routeModules = backendRoutes

export { dashboardRoutes, systemRoutes, safeguardRoutes }
