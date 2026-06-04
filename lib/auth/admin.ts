import { SignJWT, jwtVerify } from 'jose';
import bcrypt from 'bcryptjs';
import { cookies } from 'next/headers';
import { prisma } from '@/lib/db';

const ADMIN_SECRET = new TextEncoder().encode(
  process.env.JWT_SECRET ?? 'dev-jwt-secret-change-in-production',
);

const ADMIN_COOKIE = 'admin_session';

export async function hashPassword(password: string): Promise<string> {
  return bcrypt.hash(password, 12);
}

export async function verifyAdminPassword(password: string, hash: string): Promise<boolean> {
  return bcrypt.compare(password, hash);
}

export async function signAdminToken(adminId: string, username: string, role: string) {
  return new SignJWT({ sub: adminId, username, role, type: 'admin' })
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setExpirationTime('8h')
    .sign(ADMIN_SECRET);
}

export async function setAdminSession(token: string) {
  const cookieStore = await cookies();
  cookieStore.set(ADMIN_COOKIE, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    path: '/',
    maxAge: 8 * 60 * 60,
  });
}

export async function clearAdminSession() {
  const cookieStore = await cookies();
  cookieStore.delete(ADMIN_COOKIE);
}

export async function getAdminFromRequest(request: Request) {
  const auth = request.headers.get('authorization');
  let token = auth?.startsWith('Bearer ') ? auth.slice(7) : null;
  if (!token) {
    token = request.headers.get('cookie')?.match(/admin_session=([^;]+)/)?.[1] ?? null;
  }
  if (!token) throw new Error('UNAUTHORIZED');
  const { payload } = await jwtVerify(token, ADMIN_SECRET);
  const admin = await prisma.adminUser.findFirst({
    where: { id: BigInt(payload.sub as string), deleted: false },
  });
  if (!admin) throw new Error('UNAUTHORIZED');
  return admin;
}

export async function getAdminFromCookies() {
  const cookieStore = await cookies();
  const token = cookieStore.get(ADMIN_COOKIE)?.value;
  if (!token) return null;
  try {
    const { payload } = await jwtVerify(token, ADMIN_SECRET);
    return prisma.adminUser.findFirst({
      where: { id: BigInt(payload.sub as string), deleted: false },
    });
  } catch {
    return null;
  }
}
