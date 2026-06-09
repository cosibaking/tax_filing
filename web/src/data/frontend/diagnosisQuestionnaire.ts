import type { ExistingEntity, HasFiledTax, MonthlyIncomeRange } from '@/api/frontend/compliance/diagnosis'

/** 从业身份 */
export interface PersonaOption {
  value: string
  label: string
  desc: string
  icon: string
}

/** 收入来源 / 平台 */
export interface ChannelOption {
  value: string
  label: string
  icon: string
}

/** 收入区间 */
export interface IncomeRangeOption {
  value: MonthlyIncomeRange
  label: string
  desc: string
  hint: string
}

/** 经营主体 */
export interface EntityOption {
  value: ExistingEntity
  label: string
  desc: string
}

/** 报税状态 */
export interface TaxFiledOption {
  value: HasFiledTax
  label: string
  desc: string
}

/** 合规风险信号 */
export interface RiskOption {
  value: string
  label: string
  desc: string
  severity: 'high' | 'medium' | 'low'
}

/** 成本类目 */
export interface CostItemOption {
  key: string
  label: string
  placeholder: string
  hint: string
}

/** 核心诉求 */
export interface ConcernOption {
  value: string
  label: string
}

export const diagnosisSteps = [
  { key: 'profile', title: '身份与渠道', subtitle: '帮助我们了解您的从业类型与主要收入来源' },
  { key: 'income', title: '收入规模', subtitle: '用于估算年收入区间与税负对比基准' },
  { key: 'compliance', title: '合规现状', subtitle: '了解经营主体与历史申报情况' },
  { key: 'cost-risk', title: '成本与诉求', subtitle: '细化可扣费用，识别合规风险与您的核心需求' }
] as const

export const personaOptions: PersonaOption[] = [
  { value: '直播带货', label: '直播带货', desc: '抖音/快手等平台直播卖货', icon: 'ri:live-line' },
  { value: '短视频博主', label: '短视频博主', desc: '内容创作、广告合作、打赏收入', icon: 'ri:video-line' },
  { value: '自由职业', label: '自由职业接单', desc: '设计、咨询、翻译等项目制收入', icon: 'ri:palette-line' },
  { value: '电商个体', label: '电商网店', desc: '淘宝、拼多多、跨境小店等', icon: 'ri:store-2-line' },
  { value: '知识付费', label: '知识付费', desc: '课程、专栏、社群等付费内容', icon: 'ri:book-open-line' },
  { value: '其他', label: '其他个人收入', desc: '兼职、合作分成等', icon: 'ri:user-star-line' }
]

export const channelOptions: ChannelOption[] = [
  { value: '抖音', label: '抖音', icon: 'ri:tiktok-fill' },
  { value: '快手', label: '快手', icon: 'ri:live-line' },
  { value: 'B站', label: 'B站', icon: 'ri:bilibili-fill' },
  { value: '小红书', label: '小红书', icon: 'ri:book-mark-line' },
  { value: '视频号', label: '视频号', icon: 'ri:wechat-line' },
  { value: '淘宝/天猫', label: '淘宝/天猫', icon: 'ri:shopping-bag-line' },
  { value: '拼多多', label: '拼多多', icon: 'ri:store-line' },
  { value: '微信小店', label: '微信小店', icon: 'ri:wechat-2-line' },
  { value: '线下接单', label: '线下/私域接单', icon: 'ri:hand-coin-line' },
  { value: '其他平台', label: '其他', icon: 'ri:more-line' }
]

export const incomeTypeOptions: ChannelOption[] = [
  { value: '直播打赏/礼物', label: '直播打赏/礼物', icon: 'ri:gift-line' },
  { value: '带货佣金', label: '带货佣金/分销', icon: 'ri:shopping-cart-line' },
  { value: '广告合作', label: '广告/商单合作', icon: 'ri:megaphone-line' },
  { value: '知识付费', label: '课程/专栏收入', icon: 'ri:graduation-cap-line' },
  { value: '设计咨询接单', label: '设计/咨询接单', icon: 'ri:briefcase-line' },
  { value: '网店销售', label: '网店商品销售', icon: 'ri:store-2-line' },
  { value: '工资兼职', label: '工资/劳务兼职', icon: 'ri:wallet-line' }
]

export const incomeRangeOptions: IncomeRangeOption[] = [
  { value: '0-2万', label: '2 万元以下', desc: '月均不足 2 万', hint: '适合评估劳务报酬与过渡方案' },
  { value: '2-5万', label: '2 — 5 万元', desc: '月均 2—5 万', hint: '常见自由职业与博主区间' },
  { value: '5-15万', label: '5 — 15 万元', desc: '月均 5—15 万', hint: '建议重点对比个体户与 OPC' },
  { value: '15万+', label: '15 万元以上', desc: '月均 15 万以上', hint: '高收入需关注合规主体与台账' }
]

export const entityOptions: EntityOption[] = [
  { value: 'none', label: '暂无经营主体', desc: '收入以个人名义取得，未注册个体户或公司' },
  { value: 'individual', label: '已有个体工商户', desc: '已领取营业执照，有基本户或经营账户' },
  { value: 'company', label: '已有有限公司', desc: '已设立公司主体（含 OPC 等）' },
  { value: 'other', label: '其他/不确定', desc: '挂靠、代运营或主体情况待确认' }
]

export const taxFiledOptions: TaxFiledOption[] = [
  { value: 'yes', label: '按时申报', desc: '按月季年规律申报，暂无逾期' },
  { value: 'no', label: '存在漏报/逾期', desc: '有未申报期间或收到补税提醒' },
  { value: 'unsure', label: '不清楚', desc: '由平台代扣或代理处理，自己未系统申报' }
]

export const riskSignalOptions: RiskOption[] = [
  { value: 'tax_bureau', label: '收到税务机关联系', desc: '电话、短信或约谈通知', severity: 'high' },
  { value: 'platform_notice', label: '收到平台补税/合规通知', desc: '平台要求补税、升级资质等', severity: 'high' },
  { value: 'overdue', label: '曾有申报逾期或滞纳金', desc: '历史存在逾期记录', severity: 'medium' },
  { value: 'none', label: '暂无上述情况', desc: '目前未遇到明显合规压力', severity: 'low' }
]

export const concernOptions: ConcernOption[] = [
  { value: '税负对比', label: '想知道哪种申报方式更省税' },
  { value: '台账管理', label: '收入费用太散，想系统记账' },
  { value: '申报提醒', label: '怕错过申报截止日' },
  { value: '主体设立', label: '想设立个体户/公司主体' },
  { value: '应对抽查', label: '担心税务抽查，想规范留档' }
]

export const costItemOptions: CostItemOption[] = [
  { key: 'device', label: '设备器材', placeholder: '如相机、电脑、灯光等', hint: '可按购置年份折算年度折旧' },
  { key: 'software', label: '软件/素材订阅', placeholder: '设计软件、剪辑工具、素材会员', hint: '年费或月费合计' },
  { key: 'marketing', label: '投流/推广费', placeholder: '千川、随心推、平台推广', hint: '与经营直接相关的推广支出' },
  { key: 'purchase', label: '采购/进货', placeholder: '货品采购、打样、包装', hint: '电商/带货类可重点填写' },
  { key: 'venue', label: '场地/租赁', placeholder: '工作室、仓库、设备租赁', hint: '需有合同或支付凭证' },
  { key: 'service', label: '人力/物流/外包', placeholder: '助理、剪辑外包、快递物流', hint: '与业务相关的服务支出' },
  { key: 'other', label: '差旅/杂项', placeholder: '出差、通讯、其他经营费用', hint: '保留发票或支付记录' }
]
