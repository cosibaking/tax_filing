import { SignJWT, jwtVerify } from 'jose';
import bcrypt from 'bcryptjs';
import { prisma } from '@/lib/db';

const JWT_SECRET = new TextEncoder().encode(
  process.env.JWT_SECRET ?? 'dev-jwt-secret-change-in-production',
);

export interface MemberTokenPayload {
  sub: string;
  phone: string;
  type: 'access' | 'refresh';
}

export async function hashPassword(password: string): Promise<string> {
  return bcrypt.hash(password, 12);
}

export async function verifyPassword(password: string, hash: string): Promise<boolean> {
  return bcrypt.compare(password, hash);
}

export async function signMemberToken(
  payload: Omit<MemberTokenPayload, 'type'>,
  type: 'access' | 'refresh',
): Promise<string> {
  const expiresIn = type === 'access' ? '15m' : '7d';
  return new SignJWT({ ...payload, type })
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setExpirationTime(expiresIn)
    .sign(JWT_SECRET);
}

export async function verifyMemberToken(token: string): Promise<MemberTokenPayload> {
  const { payload } = await jwtVerify(token, JWT_SECRET);
  return payload as unknown as MemberTokenPayload;
}

export async function getMemberFromRequest(request: Request) {
  const auth = request.headers.get('authorization');
  if (!auth?.startsWith('Bearer ')) throw new Error('UNAUTHORIZED');
  const token = auth.slice(7);
  try {
    const payload = await verifyMemberToken(token);
    if (payload.type !== 'access') throw new Error('UNAUTHORIZED');
    const member = await prisma.member.findFirst({
      where: { id: BigInt(payload.sub), deleted: false },
    });
    if (!member) throw new Error('UNAUTHORIZED');
    return member;
  } catch (error) {
    if (error instanceof Error && error.message === 'UNAUTHORIZED') throw error;
    throw new Error('UNAUTHORIZED');
  }
}

export function getClientIp(request: Request): string | null {
  return request.headers.get('x-forwarded-for')?.split(',')[0]?.trim() ?? null;
}

export function getUserAgent(request: Request): string | null {
  return request.headers.get('user-agent');
}
