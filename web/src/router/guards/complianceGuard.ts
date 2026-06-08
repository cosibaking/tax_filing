/**
 * 合规路由守卫：台账类路由需 OPC status=active
 * @module router/guards/complianceGuard
 */
import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { opcActiveRoutePrefixes } from '@/config/complianceMenu'
import { getCompliancePlanState } from '@/api/frontend/compliance/member'

/** 判断目标路由是否需要 OPC 已 active */
export function requiresOpcActive(path: string): boolean {
  return opcActiveRoutePrefixes.some((prefix) => path === prefix || path.startsWith(prefix + '/'))
}

/**
 * 检查 OPC 是否已 active
 * TODO: 对接后端 OPC 状态 API 后实现真实校验
 */
export async function isOpcActive(): Promise<boolean> {
  const plan = await getCompliancePlanState()
  return plan.opcStatus === 'active'
}

/**
 * 合规路由 beforeEnter 守卫
 * 台账路由未 active 时重定向到 OPC 进度页
 */
export async function complianceRouteGuard(
  to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
): Promise<void> {
  if (!requiresOpcActive(to.path)) {
    next()
    return
  }

  const active = await isOpcActive()
  if (active) {
    next()
    return
  }

  next({ path: '/user/compliance/opc', replace: true })
}
