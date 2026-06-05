import { stripBasePath, withBasePath } from '@/lib/base-path';

export const MEMBER_LOGIN_PATH = '/user/login';
export const ADMIN_LOGIN_PATH = '/admin/login';

const MEMBER_PUBLIC_PATHS = [MEMBER_LOGIN_PATH, '/user/register'];
const ADMIN_PUBLIC_PATHS = [ADMIN_LOGIN_PATH];

const MEMBER_DEFAULT = '/user/overview';
const ADMIN_DEFAULT = '/admin/compliance/customers';

export function isMemberPublicPath(pathname: string): boolean {
  return MEMBER_PUBLIC_PATHS.includes(pathname);
}

export function isAdminPublicPath(pathname: string): boolean {
  return ADMIN_PUBLIC_PATHS.includes(pathname);
}

export function buildLoginRedirectUrl(loginPath: string, returnTo: string): string {
  const params = new URLSearchParams();
  params.set('callbackUrl', returnTo);
  return `${withBasePath(loginPath)}?${params.toString()}`;
}

export function getSafeCallbackUrl(
  callbackUrl: string | null | undefined,
  loginPath: string,
  fallback: string,
): string {
  if (!callbackUrl) return fallback;
  const normalized = stripBasePath(callbackUrl.split('?')[0] ?? callbackUrl);
  const query = callbackUrl.includes('?') ? callbackUrl.slice(callbackUrl.indexOf('?')) : '';
  if (!normalized.startsWith('/') || normalized.startsWith('//')) return fallback;
  if (
    normalized === MEMBER_LOGIN_PATH ||
    normalized === ADMIN_LOGIN_PATH ||
    normalized.startsWith(`${MEMBER_LOGIN_PATH}?`) ||
    normalized.startsWith(`${ADMIN_LOGIN_PATH}?`) ||
    normalized.startsWith('/user/register')
  ) {
    return fallback;
  }
  if (loginPath === MEMBER_LOGIN_PATH && !normalized.startsWith('/user')) {
    return fallback;
  }
  if (loginPath === ADMIN_LOGIN_PATH && !normalized.startsWith('/admin')) {
    return fallback;
  }
  return normalized + query;
}

export function getMemberCallbackFromLocation(
  fallback = MEMBER_DEFAULT,
): string {
  if (typeof window === 'undefined') return fallback;
  const params = new URLSearchParams(window.location.search);
  return getSafeCallbackUrl(params.get('callbackUrl'), MEMBER_LOGIN_PATH, fallback);
}

export function getAdminCallbackFromLocation(fallback = ADMIN_DEFAULT): string {
  if (typeof window === 'undefined') return fallback;
  const params = new URLSearchParams(window.location.search);
  return getSafeCallbackUrl(params.get('callbackUrl'), ADMIN_LOGIN_PATH, fallback);
}

export function getCurrentReturnPath(): string {
  if (typeof window === 'undefined') return '/';
  return window.location.pathname + window.location.search;
}
