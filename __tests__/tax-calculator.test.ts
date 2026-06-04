import { calculateTaxComparison, calculateVat, calculateCitQuarterly } from '../lib/services/compliance/diagnosis/tax-calculator';

describe('tax-calculator', () => {
  it('calculates four plan comparison', () => {
    const result = calculateTaxComparison({ annualRevenue: 1000000, annualCost: 400000 });
    expect(result).toHaveLength(4);
    expect(result[0].plan).toBe('none');
    expect(result[0].warning).toBeDefined();
    expect(result[3].plan).toBe('opc');
  });

  it('applies vat exemption under 100k monthly net', () => {
    expect(calculateVat(100000)).toBe(0);
  });

  it('calculates vat above threshold', () => {
    const vat = calculateVat(1200000);
    expect(vat).toBeGreaterThan(0);
  });

  it('calculates cit quarterly at 5% micro rate', () => {
    expect(calculateCitQuarterly(100000)).toBe(5000);
  });
});
