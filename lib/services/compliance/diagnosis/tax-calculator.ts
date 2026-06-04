export interface TaxPlanResult {
  plan: 'none' | 'labor' | 'individual' | 'opc';
  label: string;
  taxAmount: number;
  effectiveRate: number;
  warning?: string;
}

export interface TaxComparisonInput {
  annualRevenue: number;
  annualCost: number;
  monthlyRevenue?: number;
}

function laborTax(income: number): number {
  if (income <= 0) return 0;
  const taxable = income * 0.8;
  if (taxable <= 36000) return taxable * 0.03;
  if (taxable <= 144000) return taxable * 0.1 - 2520;
  if (taxable <= 300000) return taxable * 0.2 - 16920;
  if (taxable <= 420000) return taxable * 0.25 - 31920;
  if (taxable <= 660000) return taxable * 0.3 - 52920;
  if (taxable <= 960000) return taxable * 0.35 - 85920;
  return taxable * 0.45 - 181920;
}

function individualBusinessTax(profit: number): number {
  if (profit <= 0) return 0;
  if (profit <= 30000) return profit * 0.05;
  if (profit <= 90000) return profit * 0.1 - 1500;
  if (profit <= 300000) return profit * 0.2 - 10500;
  if (profit <= 500000) return profit * 0.3 - 40500;
  return profit * 0.35 - 65500;
}

function opcTax(revenue: number, cost: number, monthlyRevenue?: number): number {
  const netRevenue = revenue / 1.01;
  const profit = Math.max(0, netRevenue - cost);
  const monthly = monthlyRevenue ?? revenue / 12;
  const monthlyNet = monthly / 1.01;
  let vat = 0;
  if (monthlyNet >= 100000) {
    vat = netRevenue * 0.01;
  }
  const surcharge = vat * 0.06;
  const cit = profit <= 3000000 ? profit * 0.05 : profit * 0.25;
  const dividendTax = profit * 0.2 * 0.3;
  return vat + surcharge + cit + dividendTax;
}

export function calculateTaxComparison(input: TaxComparisonInput): TaxPlanResult[] {
  const { annualRevenue, annualCost, monthlyRevenue } = input;
  const profit = Math.max(0, annualRevenue - annualCost);

  return [
    {
      plan: 'none',
      label: '不报税',
      taxAmount: 0,
      effectiveRate: 0,
      warning: '高危：平台已报送，存在补税+滞纳金+罚款风险',
    },
    {
      plan: 'labor',
      label: '纯劳务',
      taxAmount: Math.round(laborTax(annualRevenue)),
      effectiveRate: annualRevenue > 0 ? laborTax(annualRevenue) / annualRevenue : 0,
    },
    {
      plan: 'individual',
      label: '个体户',
      taxAmount: Math.round(individualBusinessTax(profit)),
      effectiveRate: annualRevenue > 0 ? individualBusinessTax(profit) / annualRevenue : 0,
    },
    {
      plan: 'opc',
      label: 'OPC',
      taxAmount: Math.round(opcTax(annualRevenue, annualCost, monthlyRevenue)),
      effectiveRate:
        annualRevenue > 0 ? opcTax(annualRevenue, annualCost, monthlyRevenue) / annualRevenue : 0,
    },
  ];
}

export function recommendPlan(
  monthlyIncomeRange: string,
  annualRevenue: number,
  taxBureauContact: boolean,
): 'opc' | 'individual' | 'labor' | 'transitional' {
  if (taxBureauContact) return 'opc';
  if (monthlyIncomeRange === '0-2万' || annualRevenue < 300000) return 'transitional';
  if (annualRevenue >= 1000000) return 'opc';
  if (annualRevenue >= 300000) return 'individual';
  return 'transitional';
}

export function calculateVat(monthlyGrossRevenue: number): number {
  const net = monthlyGrossRevenue / 1.01;
  if (net < 100000) return 0;
  return net * 0.01;
}

export function calculateSurcharge(vat: number): number {
  return vat * 0.06;
}

export function calculateCitQuarterly(cumulativeProfit: number): number {
  if (cumulativeProfit <= 0) return 0;
  if (cumulativeProfit <= 3000000) return cumulativeProfit * 0.05;
  return cumulativeProfit * 0.25;
}
