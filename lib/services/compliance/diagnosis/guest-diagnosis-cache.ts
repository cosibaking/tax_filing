import { randomUUID } from 'crypto';
import { getRedis } from '@/lib/redis';

const TTL_SECONDS = 7 * 24 * 60 * 60;
const GUEST_PREFIX = 'diagnosis:guest:';
const PENDING_PREFIX = 'diagnosis:pending:';
export const GUEST_DIAGNOSIS_ID_PREFIX = 'g_';

export type GuestDiagnosisPayload = {
  platforms: string[];
  monthlyIncomeRange: string;
  annualCostEstimate: number;
  existingEntity: string;
  hasFiledTax: boolean;
  taxBureauContact: boolean;
  recommendedPlan: string;
  taxComparison: object;
  annualRevenue: number;
  createdAt: string;
};

function guestSetKey(sessionId: string) {
  return `${GUEST_PREFIX}${sessionId}`;
}

function pendingKey(tempId: string) {
  return `${PENDING_PREFIX}${tempId}`;
}

export function isGuestDiagnosisId(id: string): boolean {
  return id.startsWith(GUEST_DIAGNOSIS_ID_PREFIX);
}

export async function saveGuestDiagnosis(
  guestSessionId: string,
  payload: Omit<GuestDiagnosisPayload, 'createdAt'>,
): Promise<string> {
  const redis = getRedis();
  const tempId = `${GUEST_DIAGNOSIS_ID_PREFIX}${randomUUID()}`;
  const record: GuestDiagnosisPayload = { ...payload, createdAt: new Date().toISOString() };
  const setKey = guestSetKey(guestSessionId);

  await redis
    .multi()
    .set(pendingKey(tempId), JSON.stringify(record), 'EX', TTL_SECONDS)
    .sadd(setKey, tempId)
    .expire(setKey, TTL_SECONDS)
    .exec();

  return tempId;
}

export async function getGuestDiagnosis(tempId: string): Promise<GuestDiagnosisPayload | null> {
  const redis = getRedis();
  const raw = await redis.get(pendingKey(tempId));
  if (!raw) return null;
  return JSON.parse(raw) as GuestDiagnosisPayload;
}

export async function claimGuestDiagnoses(
  guestSessionId: string,
): Promise<Array<{ tempId: string; payload: GuestDiagnosisPayload }>> {
  const redis = getRedis();
  const setKey = guestSetKey(guestSessionId);
  const tempIds = await redis.smembers(setKey);
  if (!tempIds.length) return [];

  const keys = tempIds.map((id) => pendingKey(id));
  const raws = await redis.mget(...keys);
  const claimed: Array<{ tempId: string; payload: GuestDiagnosisPayload }> = [];
  const pipeline = redis.multi();

  tempIds.forEach((tempId, index) => {
    const raw = raws[index];
    if (!raw) return;
    claimed.push({ tempId, payload: JSON.parse(raw) as GuestDiagnosisPayload });
    pipeline.del(pendingKey(tempId));
  });
  pipeline.del(setKey);
  await pipeline.exec();

  return claimed;
}
