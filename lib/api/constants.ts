export const ErrorCodes = {
  VALIDATION: 1002,
  UNAUTHORIZED: 1001,
  NOT_FOUND: 1004,
  AUTH_INVALID: 2001,
  AUTH_OTP_EXPIRED: 2002,
  DIAGNOSIS_INVALID: 3001,
  ORDER_INVALID: 4001,
  OPC_NOT_ACTIVE: 4002,
  COMPLIANCE_BLOCKED: 7001,
  COST_RATIO_WARNING: 5001,
} as const;

export const INCOME_CATEGORIES = [
  { value: 'tip', label: '打赏' },
  { value: 'commission', label: '带货佣金' },
  { value: 'ad', label: '广告' },
  { value: 'slot_fee', label: '坑位费' },
  { value: 'offline', label: '线下商单' },
  { value: 'other', label: '其他' },
] as const;

export const PLATFORMS = [
  { value: 'douyin', label: '抖音' },
  { value: 'kuaishou', label: '快手' },
  { value: 'bilibili', label: 'B站' },
  { value: 'channels', label: '视频号' },
  { value: 'xiaohongshu', label: '小红书' },
] as const;

export const INCOME_SORT_OPTIONS = [
  { value: 'createdAt_desc', label: '添加时间（新→旧）' },
  { value: 'createdAt_asc', label: '添加时间（旧→新）' },
  { value: 'occurredAt_desc', label: '收入日期（新→旧）' },
  { value: 'occurredAt_asc', label: '收入日期（旧→新）' },
] as const;

export type IncomeSort = (typeof INCOME_SORT_OPTIONS)[number]['value'];

export const DEFAULT_INCOME_SORT: IncomeSort = 'createdAt_desc';

const INCOME_SORT_SET = new Set<string>(INCOME_SORT_OPTIONS.map((o) => o.value));

export function parseIncomeSort(raw: string | null | undefined): IncomeSort {
  if (raw && INCOME_SORT_SET.has(raw)) return raw as IncomeSort;
  return DEFAULT_INCOME_SORT;
}

/** 收入台账默认筛选：当月 1 日 ~ 月末 */
export function defaultIncomeMonthRange(): { dateFrom: string; dateTo: string } {
  const now = new Date();
  const y = now.getFullYear();
  const m = now.getMonth();
  const pad = (n: number) => String(n).padStart(2, '0');
  const lastDay = new Date(y, m + 1, 0).getDate();
  return {
    dateFrom: `${y}-${pad(m + 1)}-01`,
    dateTo: `${y}-${pad(m + 1)}-${pad(lastDay)}`,
  };
}

const PLATFORM_SET = new Set<string>(PLATFORMS.map((p) => p.value));
const INCOME_CATEGORY_SET = new Set<string>(INCOME_CATEGORIES.map((c) => c.value));

export function isValidIncomePlatform(value: string): boolean {
  return PLATFORM_SET.has(value);
}

export function isValidIncomeCategory(value: string): boolean {
  return INCOME_CATEGORY_SET.has(value);
}

export const EXPENSE_CATEGORIES = [
  { value: 'equipment', label: '设备' },
  { value: 'network', label: '网费' },
  { value: 'venue', label: '场地' },
  { value: 'promotion', label: '投流' },
  { value: 'outsource', label: '外包' },
  { value: 'travel', label: '差旅' },
  { value: 'other', label: '其他' },
] as const;

export const INVOICE_TYPES = [
  { value: 'general', label: '增值税普通发票' },
  { value: 'special', label: '增值税专用发票' },
  { value: 'electronic', label: '电子发票' },
] as const;

export const EXPENSE_SORT_OPTIONS = [
  { value: 'createdAt_desc', label: '添加时间（新→旧）' },
  { value: 'createdAt_asc', label: '添加时间（旧→新）' },
  { value: 'occurredAt_desc', label: '费用日期（新→旧）' },
  { value: 'occurredAt_asc', label: '费用日期（旧→新）' },
  { value: 'amount_desc', label: '金额（高→低）' },
  { value: 'amount_asc', label: '金额（低→高）' },
] as const;

export type ExpenseSort = (typeof EXPENSE_SORT_OPTIONS)[number]['value'];

export const DEFAULT_EXPENSE_SORT: ExpenseSort = 'createdAt_desc';

const EXPENSE_SORT_SET = new Set<string>(EXPENSE_SORT_OPTIONS.map((o) => o.value));
const EXPENSE_CATEGORY_SET = new Set<string>(EXPENSE_CATEGORIES.map((c) => c.value));
const INVOICE_TYPE_SET = new Set<string>(INVOICE_TYPES.map((t) => t.value));

export function parseExpenseSort(raw: string | null | undefined): ExpenseSort {
  if (raw && EXPENSE_SORT_SET.has(raw)) return raw as ExpenseSort;
  return DEFAULT_EXPENSE_SORT;
}

export function isValidExpenseCategory(value: string): boolean {
  return EXPENSE_CATEGORY_SET.has(value);
}

export function isValidInvoiceType(value: string): boolean {
  return INVOICE_TYPE_SET.has(value);
}

export const LEDGER_ACCOUNT_LABELS: Record<string, string> = {
  '1001': '银行存款',
  '6001': '主营业务收入',
  '6401': '管理费/推广费',
};

export const PROFIT_VIEWS = ['monthly', 'quarterly', 'yearly'] as const;
export type ProfitView = (typeof PROFIT_VIEWS)[number];

export function parseProfitView(raw: string | null | undefined): ProfitView {
  if (raw && (PROFIT_VIEWS as readonly string[]).includes(raw)) return raw as ProfitView;
  return 'monthly';
}

export const RISK_DISCLOSURE_ITEMS = [
  '本服务为合规方案，非逃税方案',
  '不承诺「包不被查」',
  '不提供虚开发票、隐瞒收入服务',
  '税负取决于真实成本与收入',
] as const;

/** 与 .env 中 LEGAL_*_VERSION 保持一致，供前端提交同意记录时使用 */
export const LEGAL_DOCUMENT_VERSIONS = {
  terms: '1.0',
  privacy: '1.0',
  risk_disclosure: '1.0',
  plan_confirm: '1.0',
  contract_sign: '1.0',
} as const;

export const CONSENT_TYPES = [
  'terms',
  'privacy',
  'risk_disclosure',
  'plan_confirm',
  'contract_sign',
] as const;

export type ConsentType = (typeof CONSENT_TYPES)[number];

/** 需关联订单的同意类型 */
export const ORDER_LINKED_CONSENT_TYPES = [
  'risk_disclosure',
  'plan_confirm',
  'contract_sign',
] as const;

export const isDev = process.env.NODE_ENV === 'development';

export const TAX_CHECKLIST_ITEMS = [
  '收入流水是否与平台数据一致？',
  '是否有该申报未申报的发票？',
  '成本费用发票是否已入账？',
  'OPC 是否有员工要报工资个税？',
  '本季度增值税申报了吗？',
  '企业所得税预缴了吗？',
  '有员工的话，社保申报了吗？',
  '主播是否还有其他平台收入要合并计算？',
  '上一期申报是否有错误要更正？',
] as const;

export const TAX_FILING_STATUS_LABELS: Record<string, string> = {
  pending: '待申报',
  filed: '已申报',
  overdue: '已逾期',
};

export const TAX_TYPE_LABELS: Record<string, string> = {
  vat: '增值税（小规模）',
  surcharge: '附加税费',
  cit_quarterly: '企业所得税季度预缴',
  cit_annual: '企业所得税年度汇算',
};

export const EXPENSE_KEYWORD_BLACKLIST = ['咨询费', '服务费', '借款', '赠与'];
