/**
 * 登录态校验：未登录时提示「请先登录」并跳转对应登录页
 */
import { ElMessage } from 'element-plus'
import { router } from '@/router'
import { ADMIN_BASE_PATH, ADMIN_LOGIN_PATH } from '@/router/routesAlias'
import { useMemberStore } from '@/store/modules/member'
import { useUserStore } from '@/store/modules/user'

export const LOGIN_REQUIRED_MESSAGE = '请先登录'

export type AuthContext = 'member' | 'admin'

const MEMBER_LOGIN_PATH = '/user/login'
const MEMBER_REGISTER_PATH = '/user/register'
const PROMPT_DEBOUNCE_MS = 1500
const REDIRECT_DEBOUNCE_MS = 800

let lastPromptAt = 0
let redirectPending = false

/** 根据路径判断当前应使用会员还是管理员登录 */
export function resolveAuthContext(path?: string): AuthContext {
  const current = path ?? router.currentRoute.value.path
  return current.startsWith(ADMIN_BASE_PATH) ? 'admin' : 'member'
}

/** 是否具备有效登录态 */
export function hasValidLogin(context?: AuthContext): boolean {
  const ctx = context ?? resolveAuthContext()
  if (ctx === 'admin') {
    const userStore = useUserStore()
    return userStore.isLogin && !!userStore.accessToken
  }
  return useMemberStore().getIsLogin
}

/** 登录页路径 */
export function getLoginPath(context?: AuthContext): string {
  return (context ?? resolveAuthContext()) === 'admin' ? ADMIN_LOGIN_PATH : MEMBER_LOGIN_PATH
}

/** 弹出「请先登录」提示（防抖） */
export function showLoginRequiredMessage(): void {
  const now = Date.now()
  if (now - lastPromptAt < PROMPT_DEBOUNCE_MS) return
  lastPromptAt = now
  ElMessage.warning(LOGIN_REQUIRED_MESSAGE)
}

/**
 * 提示并跳转登录页
 * @returns 恒为 false，便于 `if (!promptLoginAndRedirect()) return`
 */
export function promptLoginAndRedirect(options?: {
  context?: AuthContext
  redirect?: string
  showMessage?: boolean
}): false {
  const { showMessage = true } = options ?? {}
  const context = options?.context ?? resolveAuthContext()
  const loginPath = getLoginPath(context)
  const redirect = options?.redirect ?? router.currentRoute.value.fullPath
  const currentPath = router.currentRoute.value.path

  if (showMessage) {
    showLoginRequiredMessage()
  }

  if (
    currentPath === loginPath ||
    currentPath === MEMBER_REGISTER_PATH ||
    currentPath === `${ADMIN_BASE_PATH}/register` ||
    currentPath === `${ADMIN_BASE_PATH}/forget-password`
  ) {
    return false
  }

  if (redirectPending) return false
  redirectPending = true
  setTimeout(() => {
    redirectPending = false
  }, REDIRECT_DEBOUNCE_MS)

  router.push({
    path: loginPath,
    query: redirect && redirect !== loginPath ? { redirect } : undefined
  })
  return false
}

/**
 * 校验登录态，未登录则提示并跳转
 * @returns 已登录 true；未登录 false（已触发跳转）
 */
export function requireLogin(options?: {
  context?: AuthContext
  redirect?: string
  showMessage?: boolean
}): boolean {
  if (hasValidLogin(options?.context)) return true
  promptLoginAndRedirect(options)
  return false
}

/** 会员端路径是否需要登录（/user 下除登录/注册外） */
export function memberPathRequiresAuth(path: string): boolean {
  if (!path.startsWith('/user')) return false
  return path !== MEMBER_LOGIN_PATH && path !== MEMBER_REGISTER_PATH
}

/** 无需 Token 的 API 路径（匿名或登录流程） */
export function isPublicApiRequest(url: string): boolean {
  const path = (url || '').split('?')[0]

  if (path.startsWith('/site/')) return true
  if (path.startsWith('/captcha/')) return true

  const publicAuthPaths = [
    '/auth/login',
    '/auth/register',
    '/auth/captcha',
    '/auth/checkCaptcha',
    '/auth/refresh',
    '/auth/forgetPassword',
    '/auth/resetPassword',
    '/auth/logout'
  ]

  if (path.startsWith('/admin/') || path.startsWith('/member/')) {
    if (publicAuthPaths.some((suffix) => path.endsWith(suffix))) return true
  }

  // 前台菜单：未登录也可拉取（驱动导航）
  if (path === '/member/user/menus') return true

  return false
}

/** 请求发起前校验 Token，缺失则拦截 */
export function ensureRequestAuthorized(url: string): boolean {
  if (isPublicApiRequest(url)) return true

  if (url.startsWith('/member/')) {
    if (useMemberStore().getToken()) return true
    promptLoginAndRedirect({ context: 'member' })
    return false
  }

  if (url.startsWith('/admin/')) {
    if (useUserStore().accessToken) return true
    promptLoginAndRedirect({ context: 'admin' })
    return false
  }

  return true
}
