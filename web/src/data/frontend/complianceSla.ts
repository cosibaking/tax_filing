/** 服务 SLA 承诺（对齐商业计划书第五章） */
export interface SlaItem {
  phase: string
  deadline: string
  note?: string
  icon: string
}

export const SERVICE_SLA_ITEMS: SlaItem[] = [
  { phase: '诊断出方案', deadline: '沟通后 24 小时内', note: '输出一页纸方案建议', icon: 'ri:file-list-3-line' },
  { phase: '工商注册（有限公司）', deadline: '资料齐全后 5—7 个工作日', icon: 'ri:building-2-line' },
  { phase: '税务登记', deadline: '执照下发后 2 个工作日内', note: '电子税务局激活', icon: 'ri:government-line' },
  { phase: '首次记账启动', deadline: '签约后 7 天内', note: '开账与历史流水补录', icon: 'ri:book-2-line' },
  { phase: '月度对账单', deadline: '次月 5 日前', icon: 'ri:file-chart-line' },
  { phase: '季度申报', deadline: '季度末前 3 天', note: '增值税/企税预缴', icon: 'ri:calendar-check-line' },
  { phase: '客户咨询响应', deadline: '工作时间 4 小时内', note: '9:00—21:00', icon: 'ri:customer-service-2-line' },
  { phase: '稽查/催缴紧急事项', deadline: '2 小时内响应', note: '全年无休', icon: 'ri:alarm-warning-line' },
]

/** 月度服务节奏（PRD 6.1） */
export const MONTHLY_SOP_ITEMS = [
  { period: '每月 1—5 日', action: '催收上月银行流水、平台截图、成本发票' },
  { period: '每月 6—10 日', action: '入账、对账、生成凭证' },
  { period: '每月 11—15 日', action: '完成增值税及企税预缴（季末月）' },
  { period: '每月 16—20 日', action: '发送月度对账单与服务完成通知' },
]
