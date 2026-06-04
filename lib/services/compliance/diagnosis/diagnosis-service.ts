import { prisma, serializeBigInt } from '@/lib/db';
import { isRedisConfigured } from '@/lib/redis';
import {
  calculateTaxComparison,
  recommendPlan,
} from '@/lib/services/compliance/diagnosis/tax-calculator';
import {
  getGuestDiagnosis,
  isGuestDiagnosisId,
  saveGuestDiagnosis,
} from '@/lib/services/compliance/diagnosis/guest-diagnosis-cache';

export type DiagnosisInput = {
  memberId?: bigint;
  guestSessionId?: string;
  platforms: string[];
  monthlyIncomeRange: string;
  annualCostEstimate: number;
  existingEntity: string;
  hasFiledTax: boolean;
  taxBureauContact: boolean;
};

function computeDiagnosisFields(input: DiagnosisInput) {
  const rangeMap: Record<string, number> = {
    '0-2万': 120000,
    '2-5万': 420000,
    '5-15万': 1200000,
    '15万+': 2400000,
    '0-2': 120000,
    '2-5': 420000,
    '5-15': 1200000,
    '15+': 2400000,
  };
  const annualRevenue = rangeMap[input.monthlyIncomeRange] ?? 600000;
  const taxComparison = calculateTaxComparison({
    annualRevenue,
    annualCost: input.annualCostEstimate,
    monthlyRevenue: annualRevenue / 12,
  });
  const recommendedPlan = recommendPlan(
    input.monthlyIncomeRange,
    annualRevenue,
    input.taxBureauContact,
  );
  return { annualRevenue, taxComparison, recommendedPlan };
}

export async function createDiagnosis(input: DiagnosisInput) {
  const { annualRevenue, taxComparison, recommendedPlan } = computeDiagnosisFields(input);

  if (input.memberId) {
    const diagnosis = await prisma.complianceDiagnosis.create({
      data: {
        memberId: input.memberId,
        platforms: input.platforms,
        monthlyIncomeRange: input.monthlyIncomeRange,
        annualCostEstimate: input.annualCostEstimate,
        existingEntity: input.existingEntity,
        hasFiledTax: input.hasFiledTax,
        taxBureauContact: input.taxBureauContact,
        recommendedPlan,
        taxComparison: taxComparison as object,
      },
    });
    return serializeBigInt({
      ...diagnosis,
      taxComparison,
      recommendedPlan,
      annualRevenue,
      persisted: true,
    });
  }

  if (!input.guestSessionId?.trim()) {
    throw new Error('GUEST_SESSION_REQUIRED');
  }
  if (!isRedisConfigured()) {
    throw new Error('REDIS_NOT_CONFIGURED');
  }

  const tempId = await saveGuestDiagnosis(input.guestSessionId.trim(), {
    platforms: input.platforms,
    monthlyIncomeRange: input.monthlyIncomeRange,
    annualCostEstimate: input.annualCostEstimate,
    existingEntity: input.existingEntity,
    hasFiledTax: input.hasFiledTax,
    taxBureauContact: input.taxBureauContact,
    recommendedPlan,
    taxComparison: taxComparison as object,
    annualRevenue,
  });

  return serializeBigInt({
    id: tempId,
    memberId: null,
    platforms: input.platforms,
    monthlyIncomeRange: input.monthlyIncomeRange,
    annualCostEstimate: input.annualCostEstimate,
    existingEntity: input.existingEntity,
    hasFiledTax: input.hasFiledTax,
    taxBureauContact: input.taxBureauContact,
    recommendedPlan,
    taxComparison,
    annualRevenue,
    persisted: false,
  });
}

export async function getDiagnosisHistory(memberId: bigint) {
  const list = await prisma.complianceDiagnosis.findMany({
    where: { memberId, deleted: false },
    orderBy: { createdAt: 'desc' },
  });
  return serializeBigInt(list);
}

export async function getDiagnosisById(id: bigint) {
  const d = await prisma.complianceDiagnosis.findFirst({ where: { id, deleted: false } });
  if (!d) return null;
  return serializeBigInt(d);
}

export async function getDiagnosisByIdOrGuest(id: string) {
  if (isGuestDiagnosisId(id)) {
    if (!isRedisConfigured()) return null;
    const payload = await getGuestDiagnosis(id);
    if (!payload) return null;
    return serializeBigInt({
      id,
      memberId: null,
      ...payload,
      persisted: false,
    });
  }

  if (!/^\d+$/.test(id)) return null;
  return getDiagnosisById(BigInt(id));
}
