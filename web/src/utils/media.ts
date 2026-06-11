/**
 * 将后端返回的相对资源路径转为浏览器可访问的完整 URL。
 * 附件、头像等静态资源走站点根路径（VITE_API_URL），而非前端子路径（VITE_BASE_URL）。
 */
export function resolveMediaUrl(url?: string | null): string {
  if (!url) return ''
  const trimmed = url.trim()
  if (!trimmed) return ''
  if (/^https?:\/\//i.test(trimmed) || trimmed.startsWith('data:') || trimmed.startsWith('blob:')) {
    return trimmed
  }

  const apiBase = (import.meta.env.VITE_API_URL || '/').replace(/\/$/, '')
  const path = trimmed.startsWith('/') ? trimmed : `/${trimmed}`
  if (!apiBase || apiBase === '/') return path
  return `${apiBase}${path}`
}
