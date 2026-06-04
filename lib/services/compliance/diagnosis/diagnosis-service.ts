import { prisma, serializeBigInt } from '@/lib/db';
import {
  calculateTaxComparison,
  recommendPlan,
} from '@/lib/services/compliance/diagnosis/tax-calculator';

export async function createDiagnosis(input: {
  memberId?: bigint;
  platforms: string[];
  monthlyIncomeRange: string;
  annualCostEstimate: number;
  existingEntity: string;
  hasFiledTax: boolean;
  taxBureauContact: boolean;
}) {
  const rangeMap: Record<string, number> = {
    '0-2万': 120000,
    '2-5万': 420000,
    '5-15万': 1200000,
    '15万+': 2400000,
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

  return serializeBigInt({ ...diagnosis, taxComparison, recommendedPlan, annualRevenue });
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
