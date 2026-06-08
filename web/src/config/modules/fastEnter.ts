/**
 * 快速入口配置
 */
import type { FastEnterConfig } from '@/types/config'

const fastEnterConfig: FastEnterConfig = {
  minWidth: 1200,
  applications: [
    {
      name: '工作台',
      description: '系统概览与数据统计',
      icon: 'ri:pie-chart-line',
      iconColor: '#377dff',
      enabled: true,
      order: 1,
      routeName: 'Console'
    },
    {
      name: '合规客户',
      description: '主播合规服务客户列表',
      icon: 'ri:shield-check-line',
      iconColor: '#13DEB9',
      enabled: true,
      order: 2,
      routeName: 'ComplianceCustomers'
    },
    {
      name: 'OPC 任务',
      description: 'OPC 落地进度处理',
      icon: 'ri:file-list-3-line',
      iconColor: '#7A7FFF',
      enabled: true,
      order: 3,
      routeName: 'ComplianceOpcTasks'
    }
  ],
  quickLinks: []
}

export default fastEnterConfig
