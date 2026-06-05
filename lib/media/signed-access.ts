import 'server-only';
import { createHmac, timingSafeEqual } from 'crypto';
import { withBasePath } from '@/lib/base-path';

const DEFAULT_TTL_SEC = 3600;

function getMediaSecret(): string {
  const secret = process.env.MEDIA_SIGN_SECRET ?? process.env.JWT_SECRET;
  if (!secret) throw new Error('MEDIA_SIGN_SECRET or JWT_SECRET is required');
  return secret;
}

function getAllowedOrigins(): string[] {
  return [
    process.env.NEXT_PUBLIC_APP_URL,
    process.env.APP_URL,
    'http://localhost:3000',
    'http://127.0.0.1:3000',
  ].filter(Boolean) as string[];
}

/** 生成带签名的媒体访问 URL（防盗链 + 时效） */
export function signMediaUrl(attachmentId: string, memberId: string, ttlSec = DEFAULT_TTL_SEC): string {
  const exp = Math.floor(Date.now() / 1000) + ttlSec;
  const sig = computeMediaSignature(attachmentId, memberId, exp);
  return withBasePath(`/api/media/${attachmentId}?uid=${memberId}&exp=${exp}&sig=${sig}`);
}

function computeMediaSignature(attachmentId: string, memberId: string, exp: number): string {
  const payload = `${attachmentId}:${memberId}:${exp}`;
  return createHmac('sha256', getMediaSecret()).update(payload).digest('hex');
}

export function verifyMediaSignature(
  attachmentId: string,
  memberId: string,
  expRaw: string,
  sig: string,
): boolean {
  const exp = Number.parseInt(expRaw, 10);
  if (!Number.isFinite(exp) || exp < Math.floor(Date.now() / 1000)) return false;
  const expected = computeMediaSignature(attachmentId, memberId, exp);
  try {
    const a = Buffer.from(sig, 'hex');
    const b = Buffer.from(expected, 'hex');
    return a.length === b.length && timingSafeEqual(a, b);
  } catch {
    return false;
  }
}

/** Referer / Sec-Fetch-Site 防盗链校验 */
export function assertAntiHotlink(request: Request): void {
  const referer = request.headers.get('referer');
  const secFetchSite = request.headers.get('sec-fetch-site');
  const allowed = getAllowedOrigins();

  if (secFetchSite === 'cross-site') {
    throw new Error('FORBIDDEN');
  }

  if (!referer) return;

  const ok = allowed.some((origin) => referer.startsWith(origin));
  if (!ok) throw new Error('FORBIDDEN');
}

export function isImageMime(mime: string): boolean {
  return mime.startsWith('image/');
}
