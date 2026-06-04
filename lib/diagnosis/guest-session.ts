const GUEST_SESSION_COOKIE = 'diagnosis_guest_session';

export function resolveGuestSessionId(
  bodySessionId: string | undefined,
  cookieHeader: string | null,
): string | undefined {
  if (bodySessionId?.trim()) return bodySessionId.trim();
  if (!cookieHeader) return undefined;
  const match = cookieHeader.match(new RegExp(`${GUEST_SESSION_COOKIE}=([^;]+)`));
  return match?.[1] ? decodeURIComponent(match[1]) : undefined;
}
