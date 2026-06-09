export interface AboutValue {
  title: string
  desc: string
  icon: string
}

export interface AboutMilestone {
  year: string
  title: string
  desc: string
}

/** 金税管家 — 关于我们（面客端展示） */
export const aboutIntro = {
  title: '关于金税管家',
  subtitle: '专注个人税务合规与申报服务，让主播、自由职业者与个体经营者省心办税。',
  description:
    '金税管家面向个人创作者、自由职业者与电商个体经营者，提供从合规诊断、收入费用台账、申报提醒到月度对账单的一站式税务合规服务。我们以合法合规为前提，帮助用户把「收入分散、凭证零散、申报易漏」的痛点，转化为可执行、可追踪的合规流程。'
}

export const aboutMission = {
  mission: '让个人税务合规可理解、可执行、可追踪',
  vision: '成为个人经营者最信赖的税务合规服务伙伴'
}

export const aboutValues: AboutValue[] = [
  {
    title: '合法合规',
    desc: '基于现行税法与优惠政策提供方案，拒绝逃税避税承诺',
    icon: 'ri:shield-check-line'
  },
  {
    title: '专业可量化',
    desc: '税负对比、台账归集、申报节点全程可视',
    icon: 'ri:bar-chart-box-line'
  },
  {
    title: '省心省力',
    desc: '顾问协助申报准备，减少自行整理资料的时间成本',
    icon: 'ri:time-line'
  },
  {
    title: '全程陪伴',
    desc: '从免费诊断到签约服务，再到记账申报持续跟进',
    icon: 'ri:customer-service-2-line'
  }
]

export const aboutServices = [
  '免费合规诊断与税负对比',
  '经营主体设立代办（个体户 / OPC）',
  '收入与费用分类台账',
  '申报日历与到期提醒',
  '顾问协助申报与月度对账单'
]

export const aboutMilestones: AboutMilestone[] = [
  {
    year: '服务定位',
    title: '聚焦个人税务合规',
    desc: '面向直播带货、短视频、自由职业、电商个体等多元收入场景'
  },
  {
    year: '核心能力',
    title: '诊断 + 台账 + 申报',
    desc: '打通从风险识别到申报落地的完整服务链路'
  },
  {
    year: '服务承诺',
    title: '透明可预期',
    desc: '重要告知、服务协议与风险告知前置，不设零风险承诺'
  }
]

export const aboutContact = {
  serviceTime: '工作日 9:00 — 18:00',
  email: 'service@jingshuiguanjia.com',
  note: '如需商务合作或企业服务咨询，请通过平台内客服或注册账号后联系专属顾问。'
}
