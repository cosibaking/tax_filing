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

export const EXPENSE_KEYWORD_BLACKLIST = ['咨询费', '服务费', '借款', '赠与'];
