/** 政策与合规科普（F-04，面客端展示） */
export interface PolicyTopic {
  id: string
  title: string
  summary: string
  icon: string
}

export const policyTopics: PolicyTopic[] = [
  {
    id: 'platform-report',
    title: '平台收入报送',
    summary: '主流平台向税务机关报送经营者收入数据，个人流水与申报数据比对已成为常态。',
    icon: 'ri:database-2-line'
  },
  {
    id: 'golden-tax',
    title: '金税四期与数据比对',
    summary: '税务、银行、平台数据联网比对，隐瞒收入、公私账混用风险显著上升。',
    icon: 'ri:shield-keyhole-line'
  },
  {
    id: 'entity-choice',
    title: '主体类型怎么选',
    summary: '劳务报酬、个体户、小微公司税负与合规成本不同，需结合收入规模与凭证情况综合评估。',
    icon: 'ri:scales-3-line'
  },
  {
    id: 'small-micro',
    title: '小微企业所得税优惠',
    summary: '年利润不超过 300 万元的小微企业，可适用实际 5% 左右的企业所得税优惠（以现行政策为准）。',
    icon: 'ri:percent-line'
  },
  {
    id: 'vat-small',
    title: '小规模增值税规则',
    summary: '月销售额 10 万元以下（普票场景）可享免税；小规模纳税人适用优惠征收率，需按期申报。',
    icon: 'ri:calendar-check-line'
  },
  {
    id: 'compliance-bottom',
    title: '合规底线',
    summary: '合法合规方案≠逃税方案。禁止隐瞒收入、虚开发票、无真实业务的成本列支。',
    icon: 'ri:error-warning-line'
  }
]

export const serviceJourneySteps = [
  { key: 'aware', title: '认知触达', desc: '了解政策变化与税负对比，明确合规必要性' },
  { key: 'diagnosis', title: '合规诊断', desc: '填写问卷，获取劳务/个体户/公司方案对比' },
  { key: 'sign', title: '方案签约', desc: '阅读风险告知，确认服务套餐并电子签约' },
  { key: 'entity', title: '主体设立', desc: '工商注册、税务登记、银行开户一站式代办' },
  { key: 'ledger', title: '日常记账', desc: '收入费用台账、凭证上传、利润自动汇总' },
  { key: 'filing', title: '申报交付', desc: '增值税/企税申报提醒，月度对账单与回执存档' }
]
