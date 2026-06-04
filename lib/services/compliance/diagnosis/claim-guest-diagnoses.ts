import { prisma } from '@/lib/db';
import { isRedisConfigured } from '@/lib/redis';
import { claimGuestDiagnoses } from '@/lib/services/compliance/diagnosis/guest-diagnosis-cache';

export async function persistClaimedGuestDiagnoses(
  guestSessionId: string | undefined,
  memberId: bigint,
): Promise<string[]> {
  if (!isRedisConfigured() || !guestSessionId?.trim()) return [];

  const pending = await claimGuestDiagnoses(guestSessionId.trim());
  if (!pending.length) return [];

  const ids: string[] = [];
  for (const { payload } of pending) {
    const diagnosis = await prisma.complianceDiagnosis.create({
      data: {
        memberId,
        platforms: payload.platforms,
        monthlyIncomeRange: payload.monthlyIncomeRange,
        annualCostEstimate: payload.annualCostEstimate,
        existingEntity: payload.existingEntity,
        hasFiledTax: payload.hasFiledTax,
        taxBureauContact: payload.taxBureauContact,
        recommendedPlan: payload.recommendedPlan,
        taxComparison: payload.taxComparison,
        createdAt: new Date(payload.createdAt),
      },
    });
    ids.push(diagnosis.id.toString());
  }
  return ids;
}
