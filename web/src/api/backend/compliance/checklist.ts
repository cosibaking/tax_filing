/** 报税前自查清单九项（与后端 rules 一致） */
export const FILING_CHECKLIST_ITEMS = [
  { key: 'income_match', label: '收入流水是否与平台数据一致？' },
  { key: 'invoice_filed', label: '是否有该申报未申报的发票？' },
  { key: 'cost_booked', label: '成本费用发票是否已入账？' },
  { key: 'payroll_tax', label: '经营主体是否有员工要报工资个税？', naWhenNoEmployee: true },
  { key: 'vat_filed', label: '本季度增值税申报了吗？' },
  { key: 'cit_prepaid', label: '企业所得税预缴了吗？' },
  { key: 'social_insurance', label: '有员工的话，社保申报了吗？', naWhenNoEmployee: true },
  { key: 'other_platform', label: '是否还有其他渠道收入要合并计算？' },
  { key: 'prior_correction', label: '上一期申报是否有错误要更正？' }
] as const

export function isChecklistItemNA(key: string, employmentStatus?: string) {
  return employmentStatus === 'no_employee' && (key === 'payroll_tax' || key === 'social_insurance')
}
