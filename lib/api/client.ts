import type { ApiResponse } from '@/lib/api/envelope';

const MEMBER_TOKEN_KEY = 'member_token';
const MEMBER_REFRESH_KEY = 'member_refresh_token';

export function getMemberToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(MEMBER_TOKEN_KEY);
}

export function setMemberTokens(access: string, refresh?: string) {
  localStorage.setItem(MEMBER_TOKEN_KEY, access);
  if (refresh) localStorage.setItem(MEMBER_REFRESH_KEY, refresh);
}

export function clearMemberTokens() {
  localStorage.removeItem(MEMBER_TOKEN_KEY);
  localStorage.removeItem(MEMBER_REFRESH_KEY);
}

async function refreshMemberToken(): Promise<boolean> {
  const refresh = localStorage.getItem(MEMBER_REFRESH_KEY);
  if (!refresh) return false;
  const res = await fetch('/api/member/auth/refresh', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refreshToken: refresh }),
  });
  if (!res.ok) return false;
  const json = (await res.json()) as ApiResponse<{
    accessToken: string;
    refreshToken?: string;
  }>;
  if (json.code !== 0 || !json.data?.accessToken) return false;
  setMemberTokens(json.data.accessToken, json.data.refreshToken);
  return true;
}

export class ApiClientError extends Error {
  constructor(
    message: string,
    public code: number,
    public status: number,
  ) {
    super(message);
    this.name = 'ApiClientError';
  }
}

export async function memberFetch<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const headers = new Headers(init?.headers);
  if (!headers.has('Content-Type') && init?.body) {
    headers.set('Content-Type', 'application/json');
  }
  const token = getMemberToken();
  if (token) headers.set('Authorization', `Bearer ${token}`);

  let res = await fetch(path, { ...init, headers });

  if (res.status === 401) {
    const refreshed = await refreshMemberToken();
    if (refreshed) {
      const retryHeaders = new Headers(init?.headers);
      if (!retryHeaders.has('Content-Type') && init?.body) {
        retryHeaders.set('Content-Type', 'application/json');
      }
      const newToken = getMemberToken();
      if (newToken) retryHeaders.set('Authorization', `Bearer ${newToken}`);
      res = await fetch(path, { ...init, headers: retryHeaders });
    } else {
      clearMemberTokens();
      if (typeof window !== 'undefined') {
        window.location.href = '/user/login';
      }
      throw new ApiClientError('登录已过期，请重新登录', 1001, 401);
    }
  }

  const json = (await res.json()) as ApiResponse<T>;
  if (!res.ok || json.code !== 0) {
    throw new ApiClientError(json.message || '请求失败', json.code, res.status);
  }
  return json.data as T;
}

export async function publicFetch<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const headers = new Headers(init?.headers);
  if (!headers.has('Content-Type') && init?.body) {
    headers.set('Content-Type', 'application/json');
  }
  const res = await fetch(path, { ...init, headers });
  const json = (await res.json()) as ApiResponse<T>;
  if (!res.ok || json.code !== 0) {
    throw new ApiClientError(json.message || '请求失败', json.code, res.status);
  }
  return json.data as T;
}
