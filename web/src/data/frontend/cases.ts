export interface AppCase {
  id: string
  tag: string
  title: string
  profile: string
  income: string
  pain: string
  solution: string
  outcome: string
  highlight: string
}

/** 金税管家应用案例（示例数据，仅供展示） */
export const appCases: AppCase[] = [
  {
    id: 'content',
    tag: '内容创作',
    title: '内容创作者：多平台收入汇总申报',
    profile: '林女士 · 杭州 · 短视频+直播带货',
    income: '年综合收入约 86 万元',
    pain: '各平台代扣口径不一，年底汇算补税超预期，费用凭证零散难追溯。',
    solution: '通过金税管家完成收入结构诊断，建立月度流水台账与可扣除费用清单，由顾问协助按季核对申报数据。',
    outcome: '首年汇算清缴顺利完成，补税金额较自行估算下降约 18%，全年申报节点零逾期。',
    highlight: '多渠道收入一本账'
  },
  {
    id: 'designer',
    tag: '自由职业',
    title: '独立设计师：劳务与经营所得合规切换',
    profile: '陈先生 · 成都 · UI/品牌设计接单',
    income: '年接单收入约 42 万元',
    pain: '长期按劳务报酬申报，边际税率偏高；办公租赁、软件订阅等费用未能系统入账。',
    solution: '诊断后选择合规经营主体方案，开通收入/费用分类记账，按月生成利润表与申报提醒。',
    outcome: '合法享受小型微利企业优惠政策，综合税负较原方案下降，财务资料可应对抽查。',
    highlight: '税负对比可量化'
  },
  {
    id: 'consultant',
    tag: '知识服务',
    title: '咨询顾问：申报日历避免逾期风险',
    profile: '周先生 · 广州 · 企业培训+顾问服务',
    income: '年综合收入约 58 万元',
    pain: '收入类型混杂（课程、顾问费、合作分成），常错过申报截止日，产生滞纳金隐患。',
    solution: '签约金税管家标准套餐，启用申报日历与到期提醒，顾问协助整理分类型收入明细。',
    outcome: '连续 12 个月申报按时完成，滞纳金为零；月度对账单可用于与客户对账。',
    highlight: '申报节点自动提醒'
  },
  {
    id: 'ecommerce',
    tag: '电商个体',
    title: '网店店主：进销项与费用台账一体化',
    profile: '赵女士 · 义乌 · 跨境小商品零售',
    income: '年营业额约 120 万元',
    pain: '采购、物流、平台佣金分散在多张卡与多个店铺，年底整理耗时且易漏记可扣费用。',
    solution: '使用金税管家费用台账模块按类目归集成本，结合顾问复核后生成季度申报底稿。',
    outcome: '申报准备时间由 2 周缩短至 3 天，费用扣除有据可查，税务沟通效率明显提升。',
    highlight: '台账清晰可追溯'
  }
]
