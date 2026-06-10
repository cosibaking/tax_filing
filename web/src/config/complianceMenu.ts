/**
 * 会员合规侧栏菜单配置（对齐 M8 §5.2 动态可见规则）
 */

export type OpcStatus = 'none' | 'pending' | 'active'

/** 套餐/订单状态，驱动菜单可见性 */
export interface CompliancePlanState {
  /** 是否有已生效（已签约）的服务订单 */
  hasActiveOrder: boolean
  /** 是否存在待签约的 pending 订单 */
  hasPendingOrder?: boolean
  /** 是否已完成至少一次免费诊断 */
  hasDiagnosis?: boolean
  /** OPC 设立进度 */
  opcStatus: OpcStatus
}

export interface ComplianceMenuItem {
  id: string
  name: string
  icon: string
  path: string
  /** 侧栏分组：service=合规服务，ledger=台账与申报 */
  menuGroup: 'service' | 'ledger'
  /** 是否需要 OPC 已 active（台账类路由） */
  requiresOpcActive?: boolean
  visible: (plan: CompliancePlanState) => boolean
}

export interface ComplianceMenuGroup {
  id: string
  name: string
  items: ComplianceMenuItem[]
}

const always = () => true

/** 全部合规菜单项定义 */
export const complianceMenuItems: ComplianceMenuItem[] = [
  {
    id: 'diagnosis',
    name: '诊断历史',
    icon: 'ri:file-search-line',
    path: '/user/compliance/diagnosis',
    menuGroup: 'service',
    visible: always
  },
  {
    id: 'plan',
    name: '方案与签约',
    icon: 'ri:file-list-3-line',
    path: '/user/compliance/plan',
    menuGroup: 'service',
    visible: (plan) => !plan.hasActiveOrder
  },
  {
    id: 'opc',
    name: '主体设立',
    icon: 'ri:building-2-line',
    path: '/user/compliance/opc',
    menuGroup: 'service',
    visible: (plan) => plan.hasActiveOrder && plan.opcStatus !== 'active'
  },
  {
    id: 'social-consult',
    name: '社保咨询',
    icon: 'ri:question-answer-line',
    path: '/user/compliance/social-consult',
    menuGroup: 'service',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  },
  {
    id: 'income',
    name: '收入台账',
    icon: 'ri:money-cny-circle-line',
    path: '/user/compliance/income',
    menuGroup: 'ledger',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  },
  {
    id: 'expense',
    name: '费用台账',
    icon: 'ri:wallet-3-line',
    path: '/user/compliance/expense',
    menuGroup: 'ledger',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  },
  {
    id: 'ledger',
    name: '利润报表',
    icon: 'ri:line-chart-line',
    path: '/user/compliance/ledger',
    menuGroup: 'ledger',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  },
  {
    id: 'filing',
    name: '报税中心',
    icon: 'ri:file-edit-line',
    path: '/user/compliance/filing',
    menuGroup: 'ledger',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  },
  {
    id: 'tax',
    name: '申报管理',
    icon: 'ri:file-paper-2-line',
    path: '/user/compliance/tax',
    menuGroup: 'ledger',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  },
  {
    id: 'social-guide',
    name: '社保指引',
    icon: 'ri:heart-pulse-line',
    path: '/user/compliance/social-guide',
    menuGroup: 'ledger',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  },
  {
    id: 'statement',
    name: '对账单',
    icon: 'ri:file-chart-line',
    path: '/user/compliance/statement',
    menuGroup: 'ledger',
    requiresOpcActive: true,
    visible: (plan) => plan.opcStatus === 'active'
  }
]

/** 按分组组织可见菜单 */
export function buildComplianceMenuTree(plan: CompliancePlanState): ComplianceMenuGroup[] {
  const visible = complianceMenuItems.filter((item) => item.visible(plan))
  if (visible.length === 0) return []

  const serviceItems = visible.filter((item) => item.menuGroup === 'service')
  const ledgerItems = visible.filter((item) => item.menuGroup === 'ledger')

  const groups: ComplianceMenuGroup[] = []
  if (serviceItems.length > 0) {
    groups.push({ id: 'service', name: '合规服务', items: serviceItems })
  }
  if (ledgerItems.length > 0) {
    groups.push({ id: 'ledger', name: '台账与申报', items: ledgerItems })
  }
  return groups
}

/** 需要 OPC active 才能访问的路由前缀 */
export const opcActiveRoutePrefixes = complianceMenuItems
  .filter((item) => item.requiresOpcActive)
  .map((item) => item.path)
