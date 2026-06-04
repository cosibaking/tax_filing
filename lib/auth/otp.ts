import { prisma } from '@/lib/db';

const OTP_EXPIRY_MINUTES = 5;
const MOCK_CODE = '123456';

export async function sendOtp(phone: string, purpose: string): Promise<void> {
  const code =
    process.env.SMS_PROVIDER === 'mock' ? MOCK_CODE : String(Math.floor(100000 + Math.random() * 900000));
  const expiresAt = new Date(Date.now() + OTP_EXPIRY_MINUTES * 60 * 1000);
  await prisma.otpCode.create({
    data: { phone, code, purpose, expiresAt },
  });
  if (process.env.SMS_PROVIDER === 'mock') {
    console.info(`[OTP Mock] phone=${phone} code=${code} purpose=${purpose}`);
  }
}

export async function verifyOtp(phone: string, purpose: string, code: string): Promise<boolean> {
  const record = await prisma.otpCode.findFirst({
    where: { phone, purpose, used: false, expiresAt: { gt: new Date() } },
    orderBy: { createdAt: 'desc' },
  });
  if (!record || record.code !== code) return false;
  await prisma.otpCode.update({ where: { id: record.id }, data: { used: true } });
  return true;
}
